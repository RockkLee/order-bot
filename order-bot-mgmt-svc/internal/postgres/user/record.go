package user

import "order-bot-mgmt-svc/internal/models"

type Record struct {
	ID           string
	Email        string
	PasswordHash string
}

func RecordFromModel(user models.User) Record {
	return Record{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}
}

func (r Record) ToModel() models.User {
	return models.User{
		ID:           r.ID,
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
	}
}
