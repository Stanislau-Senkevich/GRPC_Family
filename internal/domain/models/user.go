package models

import (
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	ID           int64     `bson:"user_id"`
	Email        string    `bson:"email"`
	PhoneNumber  string    `bson:"phone_number"`
	Name         string    `bson:"name"`
	Surname      string    `bson:"surname"`
	PassHash     string    `bson:"pass_hash"`
	RegisteredAt time.Time `bson:"registered_at"`
	Role         Role      `bson:"role"`
}

func ConvertToInfo(user *User) *famv1.UserInfo {
	return &famv1.UserInfo{
		UserId:       user.ID,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Name:         user.Name,
		Surname:      user.Surname,
		RegisteredAt: timestamppb.New(user.RegisteredAt),
	}
}
