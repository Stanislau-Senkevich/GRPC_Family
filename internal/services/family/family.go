package family

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/repository"
	"log/slog"
)

type FamilyService struct {
	log     *slog.Logger
	repo    repository.FamilyRepository
	manager *jwt.Manager
}

func New(log *slog.Logger, repo repository.FamilyRepository, manager *jwt.Manager) *FamilyService {
	return &FamilyService{log: log, repo: repo, manager: manager}
}

func (s *FamilyService) CreateFamily(ctx context.Context) (int64, error) {
	//TODO: Add functionality to store families of the user in SSO
	const op = "family.service.CreateFamily"

	log := s.log.With(
		slog.String("op", op),
	)

	userID, err := s.manager.GetUserIDFromContext(ctx)
	if err != nil {
		log.Error("failed to get user id", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.repo.CreateFamily(ctx, userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *FamilyService) GetFamilyMembersInfo(ctx context.Context, familyID int64) ([]models.User, error) {
	const op = "family.service.GetFamilyMembersInfo"

	log := s.log.With(
		slog.String("op", op),
	)

	userID, err := s.manager.GetUserIDFromContext(ctx)
	if err != nil {
		log.Error("failed to get user_id from token")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	inFamily, err := s.repo.IsUserInFamily(ctx, familyID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !inFamily {
		log.Warn("forbidden access")
		return nil, grpcerror.ErrForbidden
	}

	_, err = s.repo.GetFamilyMembersID(ctx, familyID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: Add functionality to retrieve users data

	return []models.User{}, nil
}

func (s *FamilyService) LeaveFamily(ctx context.Context, familyID int64) error {
	const op = "family.service.LeaveFamily"

	userID, err := s.manager.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", err)
	}

	//TODO: Add functionality to delete family of the user in SSO

	err = s.repo.RemoveUser(ctx, familyID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
