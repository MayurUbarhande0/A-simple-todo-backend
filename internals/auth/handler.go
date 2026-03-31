package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	if input.Name == "" || input.Email == "" || input.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing required fields")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("internal server error")
		return
	}

	_, err = h.db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", input.Name, input.Email, string(hashedPassword))
	if err != nil {
		fmt.Println("REGISTER ERROR:", err) // add this line
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("user created")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	if input.Email == "" || input.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing required fields")
		return
	}

	var user User
	err := h.db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", input.Email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("invalid email or password")
		return
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
