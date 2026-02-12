package sqldb

import (
	"fmt"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

func resolveDB(base *gorm.DB, tx store.Tx) (*gorm.DB, error) {
	gtx, err := gormTx(tx)
	if err != nil {
		return nil, fmt.Errorf("sqldb.resolveDB: %w", err)
	}
	if gtx != nil {
		return gtx, nil
	}
	return base, nil
}
