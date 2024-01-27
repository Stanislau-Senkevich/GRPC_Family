package family

import (
	"context"
	"fmt"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/repository"
	"log/slog"
)

type FamilyService struct {
	log     *slog.Logger
	repo    repository.FamilyRepository
	manager *jwt.Manager
}

func New(
	log *slog.Logger,
	repo repository.FamilyRepository,
	manager *jwt.Manager,
) *FamilyService {
	return &FamilyService{log: log, repo: repo, manager: manager}
}

// CreateFamily creates a new family with the user making the request as the leader.
// It retrieves the user ID from the context and uses it as the leader ID when creating the family.
func (s *FamilyService) CreateFamily(ctx context.Context) (int64, error) {
	const op = "family.service.CreateFamily"

	userID := s.manager.GetUserIDFromContext(ctx)

	id, err := s.repo.CreateFamily(ctx, userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetFamilyMembersIDs retrieves the IDs of members belonging to the specified family.
// It first checks if the user making the request is a member of the family.
// If the user is not a member of the family, it returns a permission denied error.
func (s *FamilyService) GetFamilyMembersIDs(ctx context.Context, familyID int64) ([]int64, error) {
	const op = "family.service.GetFamilyMembersIDs"

	log := s.log.With(
		slog.String("op", op),
	)

	userID := s.manager.GetUserIDFromContext(ctx)

	inFamily, err := s.repo.IsUserInFamily(ctx, familyID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !inFamily {
		log.Warn(grpcerror.ErrForbidden.Error())
		return nil, grpcerror.ErrForbidden
	}

	IDs, err := s.repo.GetFamilyMembersID(ctx, familyID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return IDs, nil
}

// LeaveFamily allows a user to leave a family.
// It first checks if the user making the request is a member of the specified family.
// If the user is not a member of the family, it returns a user not in family error.
// If the user is a member of the family, it removes the user from the family.
func (s *FamilyService) LeaveFamily(ctx context.Context, familyID int64) (int64, error) {
	const op = "family.service.LeaveFamily"

	log := s.log.With(
		slog.String("op", op),
	)

	userID := s.manager.GetUserIDFromContext(ctx)

	inFamily, err := s.repo.IsUserInFamily(ctx, familyID, userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	if !inFamily {
		log.Warn(grpcerror.ErrUserNotInFamily.Error(),
			slog.Int64("family_id", familyID),
			slog.Int64("user_id", userID))
		return -1, fmt.Errorf("%s: %w", op, grpcerror.ErrUserNotInFamily)
	}

	err = s.repo.RemoveUserFromFamily(ctx, familyID, userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}
