package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// StringArray для работы с массивом строк в SQLite
type StringArray []string

func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = nil
		return nil
	}

	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), sa)
	case []byte:
		return json.Unmarshal(v, sa)
	default:
		return errors.New("cannot scan into StringArray")
	}
}

func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return nil, nil
	}
	return json.Marshal(sa)
}

type Category struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Alias string `json:"alias" db:"alias"`
}

type Product struct {
	ID               int         `json:"id" db:"id"`
	Name             string      `json:"name" db:"name"`
	Price            float64     `json:"price" db:"price"`
	ShortDescription string      `json:"short_description" db:"short_description"`
	LongDescription  string      `json:"long_description" db:"long_description"`
	SKU              string      `json:"sku" db:"sku"`
	Discount         int         `json:"discount" db:"discount"`
	Images           StringArray `json:"images" db:"images"`
	CategoryID       int         `json:"category_id" db:"category_id"`
	Category         *Category   `json:"category,omitempty"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
}

type Review struct {
	ID        int       `json:"id" db:"id"`
	ProductID int       `json:"product_id" db:"product_id"`
	Name      string    `json:"name" db:"name"`
	Text      string    `json:"text" db:"text"`
	Rating    int       `json:"rating" db:"rating"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ProductListRequest struct {
	Limit       int      `query:"limit"`
	Offset      int      `query:"offset"`
	CategoryID  *int     `query:"category_id"`
	PriceFrom   *float64 `query:"price_from"`
	PriceTo     *float64 `query:"price_to"`
	HasDiscount *bool    `query:"has_discount"`
	Search      string   `query:"search"`
}

type ProductListResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
}
