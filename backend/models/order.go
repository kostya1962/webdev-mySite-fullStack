package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// IntArray для работы с массивом int в SQLite
type IntArray []int

func (ia *IntArray) Scan(value interface{}) error {
	if value == nil {
		*ia = nil
		return nil
	}

	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), ia)
	case []byte:
		return json.Unmarshal(v, ia)
	default:
		return errors.New("cannot scan into IntArray")
	}
}

func (ia IntArray) Value() (driver.Value, error) {
	if ia == nil {
		return nil, nil
	}
	return json.Marshal(ia)
}

type Order struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	ProductIDs IntArray  `json:"product_ids" db:"product_ids"`
	Status     string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	User       *User     `json:"user,omitempty"`
	Products   []Product `json:"products,omitempty"`
}

type CreateOrderRequest struct {
	ProductIDs      []int  `json:"product_ids" validate:"required,min=1"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	Name            string `json:"name" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
	DeliveryAddress string `json:"delivery_address" validate:"required"`
}

type CreateOrderAuthRequest struct {
	ProductIDs      []int  `json:"product_ids" validate:"required,min=1"`
	Name            string `json:"name" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
	DeliveryAddress string `json:"delivery_address" validate:"required"`
}

type OrderResponse struct {
	Order Order  `json:"order"`
	Token string `json:"token,omitempty"`
}
