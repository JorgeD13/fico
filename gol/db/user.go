package db

import (
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRecord struct {
	ID              int64
	Name            string
	ApellidoPaterno string
	ApellidoMaterno string
	Email           string
	PasswordHash    string
}

// CreateUser inserts a new user row.
func CreateUser(gdb *gorm.DB, u UserRecord) error {
	return WithTx(gdb, func(tx *gorm.DB) error {
		const q = `INSERT INTO users (name, apellido_paterno, apellido_materno, email, password_hash, created_at)
                    VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`
		res := tx.Exec(q, u.Name, u.ApellidoPaterno, u.ApellidoMaterno, u.Email, u.PasswordHash)
		return res.Error
	})
}

// EditUser updates user core fields (not password).
func EditUser(gdb *gorm.DB, u UserRecord) error {
	if u.ID == 0 {
		return errors.New("missing id")
	}
	return WithTx(gdb, func(tx *gorm.DB) error {
		const q = `UPDATE users SET name = ?, apellido_paterno = ?, apellido_materno = ?, email = ? WHERE id = ?`
		res := tx.Exec(q, u.Name, u.ApellidoPaterno, u.ApellidoMaterno, u.Email, u.ID)
		return res.Error
	})
}

// DeleteUser removes a user by id.
func DeleteUser(gdb *gorm.DB, id int64) error {
	return WithTx(gdb, func(tx *gorm.DB) error {
		const q = `DELETE FROM users WHERE id = ?`
		return tx.Exec(q, id).Error
	})
}

// GetUserByEmail fetches a user by email (for auth).
func GetUserByEmail(gdb *gorm.DB, email string) (UserRecord, error) {
	const q = `SELECT id, name, apellido_paterno, apellido_materno, email, password_hash FROM users WHERE email = ? LIMIT 1`
	var rec UserRecord
	row := gdb.Raw(q, email).Row()
	if err := row.Scan(&rec.ID, &rec.Name, &rec.ApellidoPaterno, &rec.ApellidoMaterno, &rec.Email, &rec.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserRecord{}, fmt.Errorf("not found")
		}
		return UserRecord{}, err
	}
	return rec, nil
}

// GetUser fetches a user by id.
func GetUser(gdb *gorm.DB, id int64) (UserRecord, error) {
	const q = `SELECT id, name, apellido_paterno, apellido_materno, email, password_hash FROM users WHERE id = ? LIMIT 1`
	var rec UserRecord
	row := gdb.Raw(q, id).Row()
	if err := row.Scan(&rec.ID, &rec.Name, &rec.ApellidoPaterno, &rec.ApellidoMaterno, &rec.Email, &rec.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserRecord{}, fmt.Errorf("not found")
		}
		return UserRecord{}, err
	}
	return rec, nil
}

// ListUsers returns a page of users.
func GetUsers(gdb *gorm.DB, limit, offset int) ([]UserRecord, error) {
	if limit <= 0 {
		limit = 100
	}
	const q = `SELECT id, name, apellido_paterno, apellido_materno, email, password_hash FROM users ORDER BY id LIMIT ? OFFSET ?`
	var list []UserRecord
	if err := gdb.Raw(q, limit, offset).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
