package familyleader

import (
	"context"
	"fmt"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/repository"
	"log/slog"
)

type FamilyLeaderService struct {
	log        *slog.Logger
	familyRepo repository.FamilyRepository
	inviteRepo repository.InviteRepository
	manager    *jwt.Manager
}

func New(
	log *slog.Logger,
	familyRepo repository.FamilyRepository,
	inviteRepo repository.InviteRepository,
	manager *jwt.Manager,
) *FamilyLeaderService {
	return &FamilyLeaderService{
		log:        log,
		familyRepo: familyRepo,
		inviteRepo: inviteRepo,
		manager:    manager,
	}
}

// RemoveUserFromFamily allows a family leader to remove a user from the family.
// It first checks if the caller has the rights to remove users from the family.
// If the caller does not have the rights (is not the family leader or admin), it returns a forbidden error.
// If the caller has the rights, it checks if the specified user is a member of the family.
// If the user is not a member of the family, it returns a user not in family error.
// If the user is a member of the family, it removes the user from the family.
func (s *FamilyLeaderService) RemoveUserFromFamily(
	ctx context.Context,
	familyID, userID int64) error {
	const op = "familyleader.service.RemoveUserFromFamily"

	log := s.log.With(
		slog.String("op", op),
	)

	isLeader, err := s.hasRightsToRemove(ctx, familyID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !isLeader {
		log.Warn("failed to remove user",
			sl.Err(grpcerror.ErrForbidden),
			slog.Int64("user_id", userID))
		return fmt.Errorf("%s: %w", op, grpcerror.ErrForbidden)
	}

	inFamily, err := s.familyRepo.IsUserInFamily(ctx, familyID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !inFamily {
		log.Warn(grpcerror.ErrUserNotInFamily.Error(),
			slog.Int64("family_id", familyID),
			slog.Int64("user_id", userID))
		return fmt.Errorf("%s: %w", op, grpcerror.ErrUserNotInFamily)
	}

	err = s.familyRepo.RemoveUserFromFamily(ctx, familyID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// DeleteFamily allows a family leader to delete the specified family.
// It first checks if the caller has the rights to delete the family.
// If the caller does not have the rights (is not the family leader or admin), it returns a forbidden error.
// If the caller has the rights, it deletes the family and removes its ID from users' family lists.
func (s *FamilyLeaderService) DeleteFamily(
	ctx context.Context,
	familyID int64,
) ([]int64, error) {
	const op = "familyleader.service.DeleteFamily"

	log := s.log.With(
		slog.String("op", op),
	)

	isLeader, err := s.hasRightsToRemove(ctx, familyID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if isLeader {
		log.Warn("failed to delete family", sl.Err(grpcerror.ErrForbidden))
		return nil, fmt.Errorf("%s: %w", op, grpcerror.ErrForbidden)
	}

	return s.familyRepo.DeleteFamily(ctx, familyID)
}

func (s *FamilyLeaderService) hasRightsToRemove(ctx context.Context, familyID int64) (bool, error) {
	const op = "familyleader.service.hasRightsToRemove"

	userID := s.manager.GetUserIDFromContext(ctx)

	leaderID, err := s.familyRepo.GetFamilyLeaderID(ctx, familyID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return userID == leaderID || s.manager.IsAdmin(ctx), nil
}
