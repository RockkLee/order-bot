package sqldbold

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/infra/sqldbold/sqldbexecutor"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"github.com/jackc/pgx/v5/pgconn"
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

func NewUserStore(db *sqldb.DB) *UserStore {
	if db == nil {
		panic("sqldb.NewUserStore(), the db ptr is nil")
	}
	return &UserStore{db: db.Conn()}
}

func (s *UserStore) Create(ctx context.Context, tx store.Tx, user entities.User) error {
	record := UserRecordFromModel(user)
	exec, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.UserStore.Create: %w", err)
	}
	_, err = exec.ExecContext(ctx, insertUserQuery, record.ID, record.Email, record.PasswordHash, record.AccessToken, record.RefreshToken)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return fmt.Errorf("sqldb.UserStore.Create: %w", store.ErrUserExists)
		}
		return fmt.Errorf("sqldb.UserStore.Create: %w", err)
	}
	return nil
}

func (s *UserStore) FindByEmail(ctx context.Context, tx store.Tx, email string) (entities.User, error) {
	var record UserRecord
	exec, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByEmail: %w", err)
	}
	err = exec.QueryRowContext(ctx, selectUserByEmail, email).Scan(&record.ID, &record.Email, &record.PasswordHash, &record.AccessToken, &record.RefreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByEmail: %w", store.ErrNotFound)
		}
		return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByEmail: %w", err)
	}
	return record.ToModel(), nil
}

func (s *UserStore) FindByID(ctx context.Context, tx store.Tx, id string) (entities.User, error) {
	var record UserRecord
	exec, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByBotID: %w", err)
	}
	err = exec.QueryRowContext(ctx, selectUserByID, id).Scan(&record.ID, &record.Email, &record.PasswordHash, &record.AccessToken, &record.RefreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByBotID: %w", store.ErrNotFound)
		}
		return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByBotID: %w", err)
	}
	return record.ToModel(), nil
}

func (s *UserStore) UpdateTokens(ctx context.Context, tx store.Tx, id string, accessToken string, refreshToken string) error {
	exec, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.UserStore.UpdateTokens: %w", err)
	}
	result, err := exec.ExecContext(ctx, updateTokensQuery, accessToken, refreshToken, id)
	if err != nil {
		return fmt.Errorf("sqldb.UserStore.UpdateTokens: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqldb.UserStore.UpdateTokens: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("sqldb.UserStore.UpdateTokens: %w", store.ErrNotFound)
	}
	return nil
}
