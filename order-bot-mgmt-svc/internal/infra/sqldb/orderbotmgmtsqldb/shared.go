package orderbotmgmtsqldb

import (
	"fmt"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

func resolveDB(base *gorm.DB, tx store.Tx) (*gorm.DB, error) {
	if tx == nil {
		return base, nil
	}
	gtx, ok := tx.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("orderbotmgmtsqldb.resolveDB: expected *gorm.DB, got %T", tx)
	}
	return gtx, nil
}
