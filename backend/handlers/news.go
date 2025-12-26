package handlers

import (
	"myAPI/database"
	"myAPI/models"

	"github.com/gofiber/fiber/v2"
)

// GetNews returns list of news items
func GetNews(c *fiber.Ctx) error {
    rows, err := database.DB.Query(`SELECT id, title, description, image, created_at FROM news ORDER BY created_at DESC`)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to fetch news"})
    }
    defer rows.Close()

    var list []models.News
    for rows.Next() {
        var n models.News
        if err := rows.Scan(&n.ID, &n.Title, &n.Description, &n.Image, &n.CreatedAt); err != nil {
            continue
        }
        list = append(list, n)
    }

    return c.JSON(list)
}
