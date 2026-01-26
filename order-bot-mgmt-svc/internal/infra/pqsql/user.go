package pqsql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type UserRecord struct {
	ID           string
	Email        string
	PasswordHash string
	AccessToken  string
	RefreshToken string
}

func UserRecordFromModel(user entities.User) UserRecord {
	return UserRecord{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}
}

func (r UserRecord) ToModel() entities.User {
	return entities.User{
		ID:           r.ID,
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
	}
}

type UserStore struct {
	db *sql.DB
}

const (
	insertUserQuery     = `INSERT INTO users (id, email, password_hash, access_token, refresh_token) VALUES ($1, $2, $3, $4, $5);`
	selectUserByEmail   = `SELECT id, email, password_hash, access_token, refresh_token FROM users WHERE email = $1;`
	selectUserByID      = `SELECT id, email, password_hash, access_token, refresh_token FROM users WHERE id = $1;`
	updateTokensQuery   = `UPDATE users SET access_token = $1, refresh_token = $2 WHERE id = $3;`
	uniqueViolationCode = "23505"
)

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, user entities.User) error {
	record := UserRecordFromModel(user)
	_, err := s.db.ExecContext(ctx, insertUserQuery, record.ID, record.Email, record.PasswordHash, record.AccessToken, record.RefreshToken)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return store.ErrUserExists
		}
		return err
	}
	return nil
}

func (s *UserStore) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	var record UserRecord
	err := s.db.QueryRowContext(ctx, selectUserByEmail, email).Scan(&record.ID, &record.Email, &record.PasswordHash, &record.AccessToken, &record.RefreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, store.ErrNotFound
		}
		return entities.User{}, err
	}
	return record.ToModel(), nil
}

func (s *UserStore) FindByID(ctx context.Context, id string) (entities.User, error) {
	var record UserRecord
	err := s.db.QueryRowContext(ctx, selectUserByID, id).Scan(&record.ID, &record.Email, &record.PasswordHash, &record.AccessToken, &record.RefreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, store.ErrNotFound
		}
		return entities.User{}, err
	}
	return record.ToModel(), nil
}

func (s *UserStore) UpdateTokens(ctx context.Context, id string, accessToken string, refreshToken string) error {
	result, err := s.db.ExecContext(ctx, updateTokensQuery, accessToken, refreshToken, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return store.ErrNotFound
	}
	return nil
}
