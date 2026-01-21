package pqsql

import (
	"context"
	"database/sql"
	"errors"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	insertUserQuery     = `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3);`
	selectUserByEmail   = `SELECT id, email, password_hash FROM users WHERE email = $1;`
	uniqueViolationCode = "23505"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, user entities.User) error {
	record := UserRecordFromModel(user)
	_, err := s.db.ExecContext(ctx, insertUserQuery, record.ID, record.Email, record.PasswordHash)
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
	err := s.db.QueryRowContext(ctx, selectUserByEmail, email).Scan(&record.ID, &record.Email, &record.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, store.ErrNotFound
		}
		return entities.User{}, err
	}
	return record.ToModel(), nil
}
