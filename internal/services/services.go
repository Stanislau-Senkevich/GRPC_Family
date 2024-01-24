package services

import (
	"context"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/domain/models"
)

type Family interface {
	CreateFamily(ctx context.Context) (int64, error)
	GetFamilyMembersInfo(ctx context.Context, familyID int64) ([]models.User, error)
	LeaveFamily(ctx context.Context, familyID int64) error
}

type FamilyLeader interface {
}

type Invite interface {
}
