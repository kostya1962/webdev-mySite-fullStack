package handlers

import (
	"database/sql"
	"encoding/json"
	"myAPI/database"

	"github.com/gofiber/fiber/v2"
)

// SaveFavorites сохраняет список избранных товаров для пользователя
func SaveFavorites(c *fiber.Ctx) error {
    var req struct {
        Email      string `json:"email"`
        ProductIDs []int  `json:"productIDs"`
    }

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // Получаем ID пользователя
    var userID int
    err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&userID)
    if err == sql.ErrNoRows {
        return c.Status(404).JSON(fiber.Map{"error": "User not found"})
    }
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Database error"})
    }

    // Сохраняем productIDs как JSON
    data, _ := json.Marshal(req.ProductIDs)

    // Используем INSERT OR REPLACE для обновления
    _, err = database.DB.Exec(`INSERT OR REPLACE INTO favorites (user_id, product_ids, updated_at) VALUES (?, ?, datetime('now'))`, userID, string(data))
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to save favorites"})
    }

    return c.JSON(fiber.Map{"success": true})
}

// GetFavorites возвращает массив product IDs для пользователя
func GetFavorites(c *fiber.Ctx) error {
    email := c.Query("email")
    if email == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Email is required"})
    }

    var userID int
    err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
    if err == sql.ErrNoRows {
        return c.Status(404).JSON(fiber.Map{"error": "User not found"})
    }
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Database error"})
    }

    var raw string
    err = database.DB.QueryRow("SELECT product_ids FROM favorites WHERE user_id = ?", userID).Scan(&raw)
    if err == sql.ErrNoRows {
        // возвращаем пустой массив, если нет записи
        return c.JSON([]int{})
    }
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Database error"})
    }

    var ids []int
    if err := json.Unmarshal([]byte(raw), &ids); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to parse favorites"})
    }

    return c.JSON(ids)
}
