package user

import (
	"context"
	"database/sql"
	"errors"
	"order-bot-mgmt-svc/internal/models/entities"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	insertUserQuery     = `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3);`
	selectUserByEmail   = `SELECT id, email, password_hash FROM users WHERE email = $1;`
	uniqueViolationCode = "23505"
)

var (
	ErrNotFound   = errors.New("user not found")
	ErrUserExists = errors.New("user already exists")
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, user entities.User) error {
	record := RecordFromModel(user)
	_, err := s.db.ExecContext(ctx, insertUserQuery, record.ID, record.Email, record.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func (s *Store) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	var record Record
	err := s.db.QueryRowContext(ctx, selectUserByEmail, email).Scan(&record.ID, &record.Email, &record.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, ErrNotFound
		}
		return entities.User{}, err
	}
	return record.ToModel(), nil
}
