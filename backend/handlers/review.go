package handlers

import (
	"database/sql"
	"myAPI/database"
	"myAPI/models"
	"myAPI/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CreateReview handles creating a new review for a product
func CreateReview(c *fiber.Ctx) error {
    id := c.Params("id")
    productID, err := strconv.Atoi(id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid product id"})
    }

    var body struct {
        Name          string `json:"name"`
        Text          string `json:"text"`
        Rating        int    `json:"rating"`
        SaveToAccount bool   `json:"save_to_account"`
    }

    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    if body.Name == "" || body.Text == "" || body.Rating < 1 || body.Rating > 5 {
        return c.Status(400).JSON(fiber.Map{"error": "invalid review data"})
    }

    res, err := database.DB.Exec(`INSERT INTO reviews (product_id, name, text, rating) VALUES (?, ?, ?, ?)`, productID, body.Name, body.Text, body.Rating)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to insert review"})
    }

    lastID, _ := res.LastInsertId()

    // If user provided authorization and requested to save name, update user
    authHeader := c.Get("Authorization")
    if authHeader != "" && body.SaveToAccount {
        // validate token
        tokenPart := ""
        // using utils.ValidateToken requires just token string without Bearer
        // extract token (expected format: "Bearer <token>")
        if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
            tokenPart = authHeader[7:]
        }
        if tokenPart != "" {
            claims, err := utils.ValidateToken(tokenPart)
            if err == nil {
                // update user name
                _, _ = database.DB.Exec(`UPDATE users SET name = ? WHERE id = ?`, body.Name, claims.UserID)
            }
        }
    }

    var review models.Review
    row := database.DB.QueryRow(`SELECT id, product_id, name, text, rating, created_at FROM reviews WHERE id = ?`, lastID)
    err = row.Scan(&review.ID, &review.ProductID, &review.Name, &review.Text, &review.Rating, &review.CreatedAt)
    if err != nil && err != sql.ErrNoRows {
        return c.Status(500).JSON(fiber.Map{"error": "failed to fetch created review"})
    }

    return c.Status(201).JSON(review)
}
