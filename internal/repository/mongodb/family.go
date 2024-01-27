package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/config"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"slices"
)

// CreateFamily creates a new family in the database with the specified leader ID.
func (m *MongoRepository) CreateFamily(ctx context.Context, leaderID int64) (int64, error) {
	const op = "family.mongo.CreateFamily"

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.FamilyCollection])

	familyID, err := m.getNewID(ctx, m.Config.Collections[config.FamilyCollection])
	if err != nil {
		log.Error("failed to generate new family id", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	family := models.Family{
		ID:           familyID,
		LeaderUserID: leaderID,
		MembersID:    []int64{leaderID},
	}

	_, err = coll.InsertOne(ctx, family)
	if err != nil {
		log.Error("failed to insert family into db", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return familyID, nil
}

// GetFamilyMembersID retrieves the member IDs of the family with the specified ID from the database.
func (m *MongoRepository) GetFamilyMembersID(ctx context.Context, familyID int64) ([]int64, error) {
	const op = "family.mongo.GetFamilyMembersID"

	family, err := m.getFamily(ctx, familyID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return family.MembersID, nil
}

// GetFamilyLeaderID retrieves the leader's user ID of the family with the specified ID from the database.
func (m *MongoRepository) GetFamilyLeaderID(ctx context.Context, familyID int64) (int64, error) {
	family, err := m.getFamily(ctx, familyID)
	if err != nil {
		return -1, err
	}

	return family.LeaderUserID, nil
}

// IsUserInFamily checks whether the specified user is a member of the family with the given ID.
func (m *MongoRepository) IsUserInFamily(ctx context.Context, familyID, userID int64) (bool, error) {
	const op = "family.mongo.isUserInFamily"

	var family models.Family

	family, err := m.getFamily(ctx, familyID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	for _, id := range family.MembersID {
		if id == userID {
			return true, nil
		}
	}

	return false, nil
}

// AddUserToFamily adds a user to the specified family.
// It first checks if the user is already a member of the family using IsUserInFamily method.
// If the user is already in the family, it returns an error indicating that the user is already a member.
// Otherwise, it retrieves the current list of family members using GetFamilyMembersID method,
// appends the new user ID to the list, and updates the family document in the database.
func (m *MongoRepository) AddUserToFamily(ctx context.Context, familyID, userID int64) error {
	const op = "family.mongo.AddUserToFamily"

	log := m.log.With(
		slog.String("op", op),
	)

	inFamily, err := m.IsUserInFamily(ctx, familyID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if inFamily {
		log.Warn(grpcerror.ErrUserInFamily.Error())
		return grpcerror.ErrUserInFamily
	}

	members, err := m.GetFamilyMembersID(ctx, familyID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	members = append(members, userID)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.FamilyCollection])

	filter := bson.D{
		{"family_id", familyID},
	}

	update := bson.D{
		{"$set", bson.D{
			{"members", members},
		},
		},
	}

	_, err = coll.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error("failed to update family in db", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// RemoveUserFromFamily removes a user from the specified family.
// It retrieves the family document from the database using the familyID.
// If the user is the only member of the family, it deletes the entire family document from the database.
// If the user is not the only member, it removes the user from the family members list and updates the family document.
// If the user being removed is the leader of the family, it updates the leader to the next available member.
func (m *MongoRepository) RemoveUserFromFamily(ctx context.Context, familyID, userID int64) error {
	const op = "family.mongo.RemoveUserFromFamily"

	log := m.log.With(
		slog.String("op", op),
	)

	family, err := m.getFamily(ctx, familyID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newMembers := slices.DeleteFunc(family.MembersID, func(i int64) bool {
		return i == userID
	})

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.FamilyCollection])

	filter := bson.D{
		{"family_id", familyID},
	}

	if len(newMembers) == 0 {
		_, err = coll.DeleteOne(ctx, filter)
		if err != nil {
			log.Error("failed to delete family", sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}

	if family.LeaderUserID == userID {
		family.LeaderUserID = newMembers[0]
	}

	update := bson.D{
		{"$set", bson.D{
			{"members", newMembers},
			{"leader_id", family.LeaderUserID},
		},
		},
	}

	_, err = coll.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error("failed to update family in db", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// DeleteFamily deletes the family with the specified ID from the database.
// It first attempts to find the family document using the familyID.
// If the family is found, it deletes the document from the database and returns the IDs of its members.
// If the family is not found, it returns ErrFamilyNotFound.
// If any other error occurs during the process, it returns an error.
func (m *MongoRepository) DeleteFamily(ctx context.Context, familyID int64) ([]int64, error) {
	const op = "family.mongo.DeleteFamily"

	log := m.log.With(
		slog.String("op", op),
	)

	var family models.Family

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.FamilyCollection])

	filter := bson.D{
		{"family_id", familyID},
	}

	res := coll.FindOneAndDelete(ctx, filter)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		log.Warn(grpcerror.ErrFamilyNotFound.Error())
		return nil, fmt.Errorf("%s: %w", op, grpcerror.ErrFamilyNotFound)
	}
	if res.Err() != nil {
		log.Error("failed to find and delete family", sl.Err(res.Err()))
		return nil, fmt.Errorf("%s: %w", op, res.Err())
	}

	if err := res.Decode(&family); err != nil {
		log.Error("failed to decode family", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return family.MembersID, nil
}

func (m *MongoRepository) getFamily(ctx context.Context, familyID int64) (models.Family, error) {
	const op = "family.mongo.getFamily"

	var family models.Family

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.FamilyCollection])

	filter := bson.D{
		{"family_id", familyID},
	}

	res := coll.FindOne(ctx, filter)
	if res.Err() != nil {
		log.Warn(grpcerror.ErrFamilyNotFound.Error())
		return models.Family{}, fmt.Errorf("%s: %w", op, grpcerror.ErrFamilyNotFound)
	}

	if err := res.Decode(&family); err != nil {
		log.Error("failed to decode family", sl.Err(err))
		return models.Family{}, fmt.Errorf("%s: %w", op, err)
	}

	return family, nil
}

func (m *MongoRepository) getNewID(ctx context.Context, collectionName string) (int64, error) {
	var seq models.Sequence

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.SequenceCollection])

	filter := bson.D{
		{"collection_name", collectionName},
	}

	update := bson.D{
		{"$inc", bson.D{
			{"counter", 1},
		},
		},
	}

	res := coll.FindOneAndUpdate(ctx, filter, update)
	if res.Err() != nil {
		return -1, fmt.Errorf("failed to get id: %w", res.Err())
	}

	err := res.Decode(&seq)
	if err != nil {
		return -1, fmt.Errorf("failed to decode sequence: %w", err)
	}

	return seq.Counter, nil
}
