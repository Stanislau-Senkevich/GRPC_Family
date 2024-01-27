package models

import famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"

type Invite struct {
	ID       int64 `bson:"invite_id"`
	FamilyID int64 `bson:"family_id"`
	UserID   int64 `bson:"user_id"`
}

func ConvertToInviteModel(invite *Invite) *famv1.InviteModel {
	return &famv1.InviteModel{
		InviteId: invite.ID,
		FamilyId: invite.FamilyID,
		UserId:   invite.UserID,
	}
}
