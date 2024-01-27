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
)

// RegisterInvite registers a new invite in the database.
// It creates a new invite document with the specified familyID and userID,
// inserts it into the database, and returns the ID of the newly created invite.
func (m *MongoRepository) RegisterInvite(ctx context.Context, familyID, userID int64) (int64, error) {
	const op = "invite.mongo.RegisterInvite"

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.InviteCollection])

	id, err := m.getNewID(ctx, m.Config.Collections[config.InviteCollection])
	if err != nil {
		log.Error("failed to get new id for invite", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	invite := models.Invite{
		ID:       id,
		FamilyID: familyID,
		UserID:   userID,
	}

	_, err = coll.InsertOne(ctx, invite)
	if err != nil {
		log.Error("failed to insert new invite into db", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetInvites retrieves invites for a specific user from the database.
// It searches the database for invites associated with the specified userID,
// retrieves them, and returns a slice of models.Invite.
func (m *MongoRepository) GetInvites(ctx context.Context, userID int64) ([]models.Invite, error) {
	const op = "invite.mongo.GetInvites"

	var invites []models.Invite

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.InviteCollection])

	filter := bson.D{
		{"user_id", userID},
	}

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		log.Error("failed to search in db:", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = cur.All(ctx, &invites); err != nil {
		log.Error("failed to decode invites:", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return invites, nil
}

// IsUserInvited checks if a user with a specific ID is invited to join a family with a specific ID.
// It searches the database for an invitation matching the provided familyID and userID.
// If an invitation is found, it returns true, indicating that the user is invited.
func (m *MongoRepository) IsUserInvited(ctx context.Context, familyID, userID int64) (bool, error) {
	const op = "invite.mongo.IsUserInvited"

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.InviteCollection])

	filter := bson.D{
		{"family_id", familyID},
		{"user_id", userID},
	}

	res := coll.FindOne(ctx, filter)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return false, nil
	}
	if res.Err() != nil {
		log.Error("failed to search in mongo", sl.Err(res.Err()))
		return false, fmt.Errorf("%s: %w", op, res.Err())
	}

	return true, nil
}

// AcceptInvite accepts an invitation for a specific user.
// It searches for an invitation in the database with the provided userID and inviteID.
// If an invitation is found, it removes the invite from the database and returns the ID of the family associated with the invite.
func (m *MongoRepository) AcceptInvite(ctx context.Context, userID, inviteID int64) (int64, error) {
	const op = "invite.mongo.AcceptInvite"

	var invite models.Invite

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.InviteCollection])

	filter := bson.D{
		{"invite_id", inviteID},
		{"user_id", userID},
	}

	res := coll.FindOneAndDelete(ctx, filter)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		log.Warn(grpcerror.ErrInviteNotFound.Error(),
			slog.Int64("user_id", userID),
			slog.Int64("invite_id", inviteID))
		return -1, grpcerror.ErrInviteNotFound
	}
	if res.Err() != nil {
		log.Error("failed to find and delete invite", sl.Err(res.Err()))
		return -1, fmt.Errorf("%s: %w", op, res.Err())
	}

	if err := res.Decode(&invite); err != nil {
		log.Error("failed to decode invite", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return invite.FamilyID, nil
}

// DenyInvite denies an invitation for a specific user.
// It searches for an invitation in the database with the provided userID and inviteID.
// If an invitation is found, it removes the invite from the database.
func (m *MongoRepository) DenyInvite(ctx context.Context, userID, inviteID int64) error {
	const op = "invite.mongo.DenyInvite"

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.InviteCollection])

	filter := bson.D{
		{"invite_id", inviteID},
		{"user_id", userID},
	}

	res := coll.FindOneAndDelete(ctx, filter)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		log.Warn(grpcerror.ErrInviteNotFound.Error(),
			slog.Int64("user_id", userID),
			slog.Int64("invite_id", inviteID))
		return grpcerror.ErrInviteNotFound
	}
	if res.Err() != nil {
		log.Error("failed to find and delete invite", sl.Err(res.Err()))
		return fmt.Errorf("%s: %w", op, res.Err())
	}

	return nil
}

// DeleteUserInvites deletes all invites associated with a specific user.
// It searches for invites in the database with the provided userID and deletes them.
func (m *MongoRepository) DeleteUserInvites(ctx context.Context, userID int64) error {
	const op = "invite.mongo.DeleteUserInvites"

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.InviteCollection])

	filter := bson.D{
		{"user_id", userID},
	}

	_, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		log.Error("failed to delete user's invites", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
