package pqsql

import (
	"order-bot-mgmt-svc/internal/models/entities"
)

type UserRecord struct {
	ID           string
	Email        string
	PasswordHash string
}

func UserRecordFromModel(user entities.User) UserRecord {
	return UserRecord{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}
}

func (r UserRecord) ToModel() entities.User {
	return entities.User{
		ID:           r.ID,
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
	}
}
