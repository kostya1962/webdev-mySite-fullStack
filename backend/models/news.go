package models

import "time"

type News struct {
    ID          int       `json:"id" db:"id"`
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    Image       string    `json:"image" db:"image"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
