package db

import (
	"time"

	"gorm.io/gorm"
)

// Store both created_at and expires_at in UTC ISO8601 (SQLite-compatible) to avoid TZ drifts.
func RevokeToken(gdb *gorm.DB, token string, expUnix int64) error {
	return WithTx(gdb, func(tx *gorm.DB) error {
		expiresUTC := time.Unix(expUnix, 0).UTC().Format("2006-01-02 15:04:05")
		const q = `INSERT OR REPLACE INTO revoked_tokens (token, expires_at, created_at) VALUES (?, ?, strftime('%Y-%m-%d %H:%M:%S','now'))`
		return tx.Exec(q, token, expiresUTC).Error
	})
}

func IsTokenRevoked(gdb *gorm.DB, token string) (bool, error) {
	const q = `SELECT 1 FROM revoked_tokens WHERE token = ? LIMIT 1`
	var x int
	row := gdb.Raw(q, token).Row()
	if err := row.Scan(&x); err != nil {
		return false, nil
	}
	return true, nil
}

func CleanupExpiredRevoked(gdb *gorm.DB) error {
	const q = `DELETE FROM revoked_tokens WHERE expires_at < datetime('now')`
	return gdb.Exec(q).Error
}
