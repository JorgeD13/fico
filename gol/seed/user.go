package seed

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"fico/gol/internal/models"

	"gorm.io/gorm"
)

// AutoMigrate applies schema changes for core models.
func AutoMigrate(db *gorm.DB) error {
	mig := db.Migrator()
	if !mig.HasTable(&models.User{}) {
		if err := db.AutoMigrate(&models.User{}, &models.RevokedToken{}); err != nil {
			return err
		}
		return nil
	}
	addCol := func(name, ddl string) error {
		if !mig.HasColumn(&models.User{}, name) {
			if err := db.Exec(fmt.Sprintf("ALTER TABLE users ADD COLUMN %s", ddl)).Error; err != nil {
				if !strings.Contains(strings.ToLower(err.Error()), "duplicate column") {
					return err
				}
			}
		}
		return nil
	}
	if err := addCol("apellido_paterno", "apellido_paterno TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := addCol("apellido_materno", "apellido_materno TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := addCol("password_hash", "password_hash TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	// SQLite no permite DEFAULT no constante en ALTER TABLE ADD COLUMN.
	// Usamos DEFAULT '' y luego seteamos valores.
	if err := addCol("created_at", "created_at TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	// Rellenar created_at vacío con now()
	_ = db.Exec("UPDATE users SET created_at = datetime('now') WHERE created_at IS NULL OR created_at = ''").Error
	if err := addCol("email", "email TEXT"); err != nil {
		return err
	}
	if !mig.HasIndex(&models.User{}, "idx_users_email") {
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error
	}
	// Ensure revoked_tokens table exists
	if !mig.HasTable(&models.RevokedToken{}) {
		if err := db.AutoMigrate(&models.RevokedToken{}); err != nil {
			return err
		}
	}
	return nil
}

// SeedDemo inserts minimal demo data.
func SeedDemo(db *gorm.DB) error {
	// default password: demo123
	hashed, _ := bcrypt.GenerateFromPassword([]byte("demo123"), bcrypt.DefaultCost)
	var count int64
	if err := db.Model(&models.User{}).Where("id = ?", 1).Count(&count).Error; err != nil {
		return fmt.Errorf("seed demo user count: %w", err)
	}
	if count == 0 {
		// Insert minimal row via Exec to avoid scanning issues
		if err := db.Exec("INSERT INTO users (id, name, apellido_paterno, apellido_materno, email, password_hash, created_at) VALUES (1, ?, '', '', ?, ?, datetime('now'))", "Demo", "demo@example.com", string(hashed)).Error; err != nil {
			return fmt.Errorf("seed demo user insert: %w", err)
		}
	} else {
		// Ensure fields exist
		if err := db.Exec("UPDATE users SET password_hash = CASE WHEN password_hash IS NULL OR password_hash = '' THEN ? ELSE password_hash END, email = CASE WHEN email IS NULL OR email = '' THEN ? ELSE email END, name = CASE WHEN name IS NULL OR name = '' THEN ? ELSE name END, created_at = CASE WHEN created_at IS NULL OR created_at = '' THEN datetime('now') ELSE created_at END WHERE id = 1", string(hashed), "demo@example.com", "Demo").Error; err != nil {
			return fmt.Errorf("seed demo user update: %w", err)
		}
	}
	return nil
}
