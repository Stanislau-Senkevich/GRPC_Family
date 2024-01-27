package sso

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/client/sso/grpc"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/domain/models"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"google.golang.org/grpc/metadata"
	"log/slog"
)

type SSOService struct {
	client        *grpc.Client
	manager       *jwt.Manager
	adminEmail    string
	adminPassword string
}

func New(
	client *grpc.Client,
	manager *jwt.Manager,
	adminEmail string,
	adminPassword string,
) *SSOService {
	return &SSOService{
		client:        client,
		manager:       manager,
		adminEmail:    adminEmail,
		adminPassword: adminPassword,
	}
}

// GetUserInfo retrieves user information from the SSO service.
func (s *SSOService) GetUserInfo(userID int64) (*models.User, error) {
	const op = "sso.service.GetUserInfo"

	log := s.client.Log.With(
		slog.String("op", op),
	)

	ctx, err := s.signInAndGetContext()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	req, err := s.client.Userinfo.GetUserInfoByID(ctx, &ssov1.GetUserInfoByIDRequest{
		UserId: userID,
	})
	if err != nil {
		log.Warn("failed to find user", slog.Int64("user_id", userID))
		return nil, fmt.Errorf("%s (id: %d): %w", op, userID, err)
	}

	return &models.User{
		ID:           req.GetUserId(),
		Email:        req.GetEmail(),
		PhoneNumber:  req.GetPhoneNumber(),
		Name:         req.GetName(),
		Surname:      req.GetSurname(),
		RegisteredAt: req.GetRegisteredAt().AsTime(),
	}, nil
}

// AddFamilyToList adds a family to the user's family list in the SSO service.
func (s *SSOService) AddFamilyToList(ctx context.Context, familyID int64) error {
	const op = "sso.service.AddFamilyToList"

	log := s.client.Log.With(
		slog.String("op", op),
	)

	userID := s.manager.GetUserIDFromContext(ctx)

	adminCtx, err := s.signInAndGetContext()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.client.Userinfo.AddFamily(adminCtx, &ssov1.AddFamilyRequest{
		UserId:   userID,
		FamilyId: familyID,
	})
	if err != nil {
		log.Warn("failed to add family", sl.Err(err),
			slog.Int64("family_id", familyID), slog.Int64("user_id", userID))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// RemoveFamilyFromList removes a family from the user's family list in the SSO service.
func (s *SSOService) RemoveFamilyFromList(userID, familyID int64) error {
	const op = "sso.service.RemoveFamilyFromList"

	log := s.client.Log.With(
		slog.String("op", op),
	)

	adminCtx, err := s.signInAndGetContext()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.client.Userinfo.DeleteFamily(adminCtx, &ssov1.DeleteFamilyRequest{
		UserId:   userID,
		FamilyId: familyID,
	})
	if err != nil {
		log.Warn("failed to delete family", sl.Err(err),
			slog.Int64("family_id", familyID), slog.Int64("user_id", userID))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SSOService) signInAndGetContext() (context.Context, error) {
	const op = "sso.signInAndGetContext"

	log := s.client.Log.With(
		slog.String("op", op),
	)

	respSign, err := s.client.Auth.SignIn(context.Background(),
		&ssov1.SignInRequest{
			Email:    s.adminEmail,
			Password: s.adminPassword,
		})
	if err != nil {
		log.Error("failed to sign in to sso", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token := "Bearer " + respSign.GetToken()

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", token)

	return ctx, nil
}
