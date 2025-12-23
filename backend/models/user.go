package models

import (
	"time"
)

type User struct {
	ID              int       `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	Password        string    `json:"-" db:"password"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	Phone           string    `json:"phone" db:"phone"`
	DeliveryAddress string    `json:"delivery_address" db:"delivery_address"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
