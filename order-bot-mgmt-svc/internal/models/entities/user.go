package entities

type User struct {
	ID           string
	Email        string
	PasswordHash string
	AccessToken  string
	RefreshToken string
}
