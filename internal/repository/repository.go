package repository

import (
	"context"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/domain/models"
)

type FamilyRepository interface {
	CreateFamily(ctx context.Context, leaderID int64) (int64, error)
	GetFamilyMembersID(ctx context.Context, familyID int64) ([]int64, error)
	GetFamilyLeaderID(ctx context.Context, familyID int64) (int64, error)
	IsUserInFamily(ctx context.Context, familyID, userID int64) (bool, error)
	AddUserToFamily(ctx context.Context, familyID, userID int64) error
	RemoveUserFromFamily(ctx context.Context, familyID, userID int64) error
	DeleteFamily(ctx context.Context, familyID int64) ([]int64, error)
}

type InviteRepository interface {
	RegisterInvite(ctx context.Context, familyID, userID int64) (int64, error)
	GetInvites(ctx context.Context, userID int64) ([]models.Invite, error)
	IsUserInvited(ctx context.Context, familyID, userID int64) (bool, error)
	AcceptInvite(ctx context.Context, userID, inviteID int64) (int64, error)
	DenyInvite(ctx context.Context, userID, inviteID int64) error
	DeleteUserInvites(ctx context.Context, userID int64) error
}
