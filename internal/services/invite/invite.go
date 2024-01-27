package invite

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/repository"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"log/slog"
)

type InviteService struct {
	log        *slog.Logger
	inviteRepo repository.InviteRepository
	familyRepo repository.FamilyRepository
	manager    *jwt.Manager
}

func New(
	log *slog.Logger,
	inviteRepo repository.InviteRepository,
	familyRepo repository.FamilyRepository,
	manager *jwt.Manager) *InviteService {
	return &InviteService{
		log:        log,
		inviteRepo: inviteRepo,
		familyRepo: familyRepo,
		manager:    manager,
	}
}

// SendInvite allows the leader of a family to send an invitation to a user to join the family.
// It first checks if the caller has the necessary rights to send invites (must be the leader of the family).
// If the caller does not have the rights, it returns a forbidden error.
// If the user is already invited to the family, it returns an error indicating that the invite already exists.
// If the user is already a member of the family, it returns an error indicating that the user is already in the family.
// If all checks pass, it registers the invite in the repository and returns the invite ID.
func (s *InviteService) SendInvite(
	ctx context.Context,
	familyID, userID int64,
) (int64, error) {
	const op = "invite.service.SendInvite"

	log := s.log.With(
		slog.String("op", op),
	)

	clientID := s.manager.GetUserIDFromContext(ctx)

	leaderID, err := s.familyRepo.GetFamilyLeaderID(ctx, familyID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	if clientID != leaderID {
		log.Warn(grpcerror.ErrForbidden.Error())
		return -1, grpcerror.ErrForbidden
	}

	isInvited, err := s.inviteRepo.IsUserInvited(ctx, familyID, userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	if isInvited {
		log.Warn(grpcerror.ErrInviteExist.Error())
		return -1, grpcerror.ErrInviteExist
	}

	inFamily, err := s.familyRepo.IsUserInFamily(ctx, familyID, userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	if inFamily {
		log.Warn(grpcerror.ErrUserInFamily.Error())
		return -1, grpcerror.ErrUserInFamily
	}

	inviteID, err := s.inviteRepo.RegisterInvite(ctx, familyID, userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return inviteID, nil
}

// GetInvites retrieves the invites for the current user.
// It first retrieves the user ID from the context using the manager.
// Then, it calls the GetInvites method of the invite repository to fetch the invites associated with the user ID.
func (s *InviteService) GetInvites(ctx context.Context) ([]*famv1.InviteModel, error) {
	const op = "invite.service.GetInvites"

	var res []*famv1.InviteModel

	userID := s.manager.GetUserIDFromContext(ctx)

	invites, err := s.inviteRepo.GetInvites(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res = make([]*famv1.InviteModel, 0, len(invites))

	for _, invite := range invites {
		res = append(res, models.ConvertToInviteModel(&invite))
	}

	return res, nil
}

// AcceptInvite accepts the invite with the given inviteID for the current user.
// It retrieves the user ID from the context using the manager, then accepts the invite
// using the invite repository. If successful, it adds the user to the family associated
// with the accepted invite using the family repository.
func (s *InviteService) AcceptInvite(
	ctx context.Context,
	inviteID int64,
) (int64, error) {
	const op = "invite.service.AcceptInvite"

	userID := s.manager.GetUserIDFromContext(ctx)

	familyID, err := s.inviteRepo.AcceptInvite(ctx, userID, inviteID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	err = s.familyRepo.AddUserToFamily(ctx, familyID, userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return familyID, nil
}

// DenyInvite denies the invite with the given inviteID for the current user.
func (s *InviteService) DenyInvite(ctx context.Context, inviteID int64) error {
	userID := s.manager.GetUserIDFromContext(ctx)

	return s.inviteRepo.DenyInvite(ctx, userID, inviteID)
}

// DeleteUserInvites deletes all invites associated with the specified userID.
func (s *InviteService) DeleteUserInvites(ctx context.Context, userID int64) error {
	return s.inviteRepo.DeleteUserInvites(ctx, userID)
}
