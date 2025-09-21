package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"fico/gol/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(gdb *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		// DB: fetch user
		rec, err := db.GetUserByEmail(gdb, req.Email)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		// verify password
		if bcrypt.CompareHashAndPassword([]byte(rec.PasswordHash), []byte(req.Password)) != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		// JWT
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "dev-secret"
		}
		claims := jwt.MapClaims{
			"sub":   rec.ID,
			"email": rec.Email,
			"exp":   time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(secret))
		if err != nil {
			http.Error(w, "could not sign token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(LoginResponse{Token: signed})
	}
}

func Logout(gdb *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Read Authorization: Bearer <token>
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusOK)
			return
		}
		var tokenStr string
		if _, err := fmt.Sscanf(auth, "Bearer %s", &tokenStr); err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Verify signature before revoking
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "dev-secret"
		}
		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil })
		if err != nil || !tok.Valid {
			w.WriteHeader(http.StatusOK)
			return
		}
		claims, _ := tok.Claims.(jwt.MapClaims)
		var expUnix int64
		if exp, ok := claims["exp"].(float64); ok {
			expUnix = int64(exp)
		}
		_ = db.RevokeToken(gdb, tokenStr, expUnix)
		_ = db.CleanupExpiredRevoked(gdb)
		w.WriteHeader(http.StatusOK)
	}
}

type EditUserRequest struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	ApellidoPaterno string `json:"apellidoPaterno"`
	ApellidoMaterno string `json:"apellidoMaterno"`
	Email           string `json:"email"`
}

func EditUser(gdb *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req EditUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == 0 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		err := db.EditUser(gdb, db.UserRecord{
			ID:              req.ID,
			Name:            req.Name,
			ApellidoPaterno: req.ApellidoPaterno,
			ApellidoMaterno: req.ApellidoMaterno,
			Email:           req.Email,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot edit: %v", err), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type DeleteUserRequest struct {
	ID int64 `json:"id"`
}

func DeleteUser(gdb *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req DeleteUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == 0 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if err := db.DeleteUser(gdb, req.ID); err != nil {
			http.Error(w, fmt.Sprintf("cannot delete: %v", err), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type CreateUserRequest struct {
	Name            string `json:"name"`
	ApellidoPaterno string `json:"apellidoPaterno"`
	ApellidoMaterno string `json:"apellidoMaterno"`
	Email           string `json:"email"`
	Password        string `json:"password"`
}

func CreateUser(gdb *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" || req.Name == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		err := db.CreateUser(gdb, db.UserRecord{
			Name:            req.Name,
			ApellidoPaterno: req.ApellidoPaterno,
			ApellidoMaterno: req.ApellidoMaterno,
			Email:           req.Email,
			PasswordHash:    string(hashed),
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot create: %v", err), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetUser(gdb *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)
		if id <= 0 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		rec, err := db.GetUser(gdb, id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		rec.PasswordHash = ""
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rec)
	}
}

func GetUsers(gdb *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		users, err := db.GetUsers(gdb, 100, 0)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		for i := range users {
			users[i].PasswordHash = ""
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(users)
	}
}
