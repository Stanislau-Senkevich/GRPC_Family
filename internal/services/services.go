package services

import (
	"context"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/domain/models"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
)

type SSO interface {
	GetUserInfo(userID int64) (*models.User, error)
	AddFamilyToList(ctx context.Context, familyID int64) error
	RemoveFamilyFromList(userID, familyID int64) error
}

type Family interface {
	CreateFamily(ctx context.Context) (int64, error)
	GetFamilyMembersIDs(ctx context.Context, familyID int64) ([]int64, error)
	LeaveFamily(ctx context.Context, familyID int64) (int64, error)
}

type FamilyLeader interface {
	RemoveUserFromFamily(ctx context.Context, familyID, userID int64) error
	DeleteFamily(ctx context.Context, familyID int64) ([]int64, error)
}

type Invite interface {
	SendInvite(ctx context.Context, familyID, userID int64) (int64, error)
	GetInvites(ctx context.Context) ([]*famv1.InviteModel, error)
	AcceptInvite(ctx context.Context, inviteID int64) (int64, error)
	DenyInvite(ctx context.Context, inviteID int64) error
	DeleteUserInvites(ctx context.Context, userID int64) error
}
