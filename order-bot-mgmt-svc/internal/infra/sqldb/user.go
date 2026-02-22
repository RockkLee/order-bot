package sqldb

import (
	"context"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserRecord struct {
	Base         BaseRecord `gorm:"embedded"`
	ID           string     `gorm:"column:id;primaryKey"`
	Email        string     `gorm:"column:email"`
	PasswordHash string     `gorm:"column:password_hash"`
	AccessToken  string     `gorm:"column:access_token"`
	RefreshToken string     `gorm:"column:refresh_token"`
}

func (UserRecord) TableName() string { return "users" }

func UserRecordFromModel(user entities.User) UserRecord {
	return UserRecord{ID: user.ID, Email: user.Email, PasswordHash: user.PasswordHash, AccessToken: user.AccessToken, RefreshToken: user.RefreshToken}
}
func (r UserRecord) ToModel() entities.User {
	return entities.User{ID: r.ID, Email: r.Email, PasswordHash: r.PasswordHash, AccessToken: r.AccessToken, RefreshToken: r.RefreshToken}
}

type UserStore struct{ db *gorm.DB }

func NewUserStore(db *DB) *UserStore {
	if db == nil {
		panic("sqldb.NewUserStore(), the db ptr is nil")
	}
	return &UserStore{db: db.Gorm()}
}

func (s *UserStore) Create(ctx context.Context, tx store.Tx, user entities.User) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.UserStore.Create: %w", err)
	}
	record := UserRecordFromModel(user)
	if err := db.WithContext(ctx).Create(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("sqldb.UserStore.Create: %w", store.ErrUserExists)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // or pgerrcode.UniqueViolation
				return fmt.Errorf("sqldb.UserStore.Create: %w", store.ErrUserExists)
			}
		}
		return fmt.Errorf("sqldb.UserStore.Create: %w", err)
	}
	return nil
}
func (s *UserStore) FindByEmail(ctx context.Context, tx store.Tx, email string) (entities.User, error) {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByEmail: %w", err)
	}
	var record UserRecord
	if err := db.WithContext(ctx).Where("email = ?", email).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByEmail: %w", store.ErrNotFound)
		}
		return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByEmail: %w", err)
	}
	return record.ToModel(), nil
}
func (s *UserStore) FindByID(ctx context.Context, tx store.Tx, id string) (entities.User, error) {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByID: %w", err)
	}
	var record UserRecord
	if err := db.WithContext(ctx).Where("id = ?", id).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByID: %w", store.ErrNotFound)
		}
		return entities.User{}, fmt.Errorf("sqldb.UserStore.FindByID: %w", err)
	}
	return record.ToModel(), nil
}
func (s *UserStore) UpdateTokens(ctx context.Context, tx store.Tx, id, accessToken, refreshToken string) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.UserStore.UpdateTokens: %w", err)
	}
	res := db.WithContext(ctx).Model(&UserRecord{}).Where("id = ?", id).Updates(map[string]any{"access_token": accessToken, "refresh_token": refreshToken})
	if res.Error != nil {
		return fmt.Errorf("sqldb.UserStore.UpdateTokens: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("sqldb.UserStore.UpdateTokens: %w", store.ErrNotFound)
	}
	return nil
}
