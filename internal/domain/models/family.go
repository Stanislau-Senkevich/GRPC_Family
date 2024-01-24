package models

type Family struct {
	ID           int64   `bson:"family_id"`
	LeaderUserID int64   `bson:"leader_id"`
	MembersID    []int64 `bson:"members"`
}
