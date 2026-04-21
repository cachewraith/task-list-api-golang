package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents the domain entity for a user
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose in JSON
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterInput represents the input for user registration
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginInput represents the input for user login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate performs validation on RegisterInput
func (r RegisterInput) Validate() error {
	if r.Email == "" {
		return ErrValidationFailed("email is required")
	}
	if r.Password == "" {
		return ErrValidationFailed("password is required")
	}
	if len(r.Password) < 6 {
		return ErrValidationFailed("password must be at least 6 characters")
	}
	if r.Name == "" {
		return ErrValidationFailed("name is required")
	}
	return nil
}

// Validate performs validation on LoginInput
func (l LoginInput) Validate() error {
	if l.Email == "" {
		return ErrValidationFailed("email is required")
	}
	if l.Password == "" {
		return ErrValidationFailed("password is required")
	}
	return nil
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// CheckPassword verifies if the provided password matches the hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// NewUser creates a new User from RegisterInput
func NewUser(input RegisterInput) *User {
	now := time.Now().UTC()
	return &User{
		ID:        uuid.New(),
		Email:     input.Email,
		Password:  input.Password,
		Name:      input.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SafeUser returns user data without sensitive fields
func (u *User) SafeUser() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"email":      u.Email,
		"name":       u.Name,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}
