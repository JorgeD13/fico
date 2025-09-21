package db

import "gorm.io/gorm"

// WithTx runs fn within a transaction. If tx exists in db, it reuses it.
func WithTx(gdb *gorm.DB, fn func(tx *gorm.DB) error) error {
	return gdb.Transaction(func(tx *gorm.DB) error { return fn(tx) })
}
