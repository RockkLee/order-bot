package store

// Tx is an abstract transaction handle passed across service/store layers.
// Concrete implementations are infrastructure-specific (e.g. *sql.Tx, *gorm.DB).
type Tx interface{}
