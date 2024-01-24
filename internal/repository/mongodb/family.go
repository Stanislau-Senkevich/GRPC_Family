package mongodb

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/config"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
)

func (m *MongoRepository) CreateFamily(ctx context.Context, leaderID int64) (int64, error) {
	const op = "family.mongo.CreateFamily"

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.FamilyCollection])

	familyID, err := m.getNewUniqueFamilyID(ctx)
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

func (m *MongoRepository) GetFamilyMembersID(ctx context.Context, familyID int64) ([]int64, error) {
	const op = "family.mongo.GetFamilyMembersID"

	family, err := m.getFamily(ctx, familyID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return family.MembersID, nil
}

func (m *MongoRepository) RemoveUser(ctx context.Context, familyID int64, userID int64) error {
	const op = "family.mongo.LeaveFamily"

	log := m.log.With(
		slog.String("op", op),
	)

	inFamily, err := m.IsUserInFamily(ctx, familyID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !inFamily {
		log.Warn(grpcerror.ErrUserNotFound.Error())
		return grpcerror.ErrUserNotFound
	}

	family, err := m.getFamily(ctx, familyID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newMembers := removeUser(family.MembersID, userID)
	if newMembers == nil {
		log.Error("failed to remove user id from slice")
		return grpcerror.ErrInternalError
	}

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.FamilyCollection])

	filter := bson.D{
		{"family_id", familyID},
	}

	update := bson.D{
		{"$set", bson.D{
			{"members", newMembers},
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

func (m *MongoRepository) AddUser(ctx context.Context, familyID int64, userID int64) error {
	const op = "family.mongo.AddUser"

	log := m.log.With(
		slog.String("op", op),
	)

	inFamily, err := m.IsUserInFamily(ctx, familyID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !inFamily {
		log.Warn(grpcerror.ErrUserNotFound.Error())
		return grpcerror.ErrUserNotFound
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

func (m *MongoRepository) IsUserInFamily(ctx context.Context, familyID int64, userID int64) (bool, error) {
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

func (m *MongoRepository) GetFamilyLeaderID(ctx context.Context, familyID int64) (int64, error) {
	family, err := m.getFamily(ctx, familyID)
	if err != nil {
		return -1, err
	}

	return family.LeaderUserID, nil
}

func (m *MongoRepository) getNewUniqueFamilyID(ctx context.Context) (int64, error) {
	var seq models.Sequence

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.SequenceCollection])

	filter := bson.D{
		{"collection_name", m.Config.Collections[config.FamilyCollection]},
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

func removeUser(members []int64, userID int64) []int64 {
	if len(members) == 1 {
		return make([]int64, 0)
	}

	if members[0] == userID {
		return members[1:]
	}

	for i, id := range members {
		if id == userID {
			return append(members[:i], members[i+1:]...)
		}
	}
	return nil
}
