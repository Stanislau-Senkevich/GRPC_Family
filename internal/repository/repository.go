package repository

import "context"

type FamilyRepository interface {
	CreateFamily(ctx context.Context, leaderID int64) (int64, error)
	GetFamilyMembersID(ctx context.Context, familyID int64) ([]int64, error)
	RemoveUser(ctx context.Context, familyID int64, userID int64) error
	AddUser(ctx context.Context, familyID int64, userID int64) error
	IsUserInFamily(ctx context.Context, familyID int64, userID int64) (bool, error)
	GetFamilyLeaderID(ctx context.Context, familyID int64) (int64, error)
}

type InviteRepository interface {
	IsUserInvited(ctx context.Context, familyID int64, userID int64) (bool, error)
	AcceptInvite(ctx context.Context, inviteID int64) error
	DenyInvite(ctx context.Context, inviteID int64) error
}

type FamilyLeaderRepository interface {
	InviteUser(ctx context.Context, familyID int64, userID int64) error
	RemoveUser(ctx context.Context, familyID int64, userID int64) error
}
