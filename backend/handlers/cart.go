package handlers

import (
	"database/sql"
	"myAPI/database"
	"myAPI/models"

	"github.com/gofiber/fiber/v2"
)

// CartItem представляет элемент корзины
type CartItem struct {
	ProductID int `json:"product_id" db:"product_id"`
	Quantity  int `json:"quantity" db:"quantity"`
}

// AddToCart - добавить товар в корзину пользователя
func AddToCart(c *fiber.Ctx) error {
	var req struct {
		Email     string `json:"email"`
		ProductID int    `json:"productID"`
		Quantity  int    `json:"quantity"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем существование пользователя
	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Проверяем существование товара
	var productID int
	err = database.DB.QueryRow("SELECT id FROM products WHERE id = ?", req.ProductID).Scan(&productID)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error": "Product not found",
		})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Проверяем, есть ли уже товар в корзине
	var existingQuantity int
	err = database.DB.QueryRow(
		"SELECT quantity FROM cart_items WHERE user_id = ? AND product_id = ?",
		userID, req.ProductID,
	).Scan(&existingQuantity)

	if err == sql.ErrNoRows {
		// Добавляем новый товар в корзину
		_, err = database.DB.Exec(
			"INSERT INTO cart_items (user_id, product_id, quantity) VALUES (?, ?, ?)",
			userID, req.ProductID, req.Quantity,
		)
	} else if err == nil {
		// Обновляем количество
		_, err = database.DB.Exec(
			"UPDATE cart_items SET quantity = ? WHERE user_id = ? AND product_id = ?",
			req.Quantity, userID, req.ProductID,
		)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to add to cart",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Item added to cart",
	})
}

// GetCart - получить корзину пользователя
func GetCart(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	// Получаем ID пользователя
	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Получаем товары из корзины с их данными
	rows, err := database.DB.Query(`
		SELECT 
			p.id, p.name, p.price, p.short_description, p.long_description,
			p.sku, p.discount, p.category_id, p.created_at, p.updated_at,
			ci.quantity
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = ?
	`, userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get cart",
		})
	}
	defer rows.Close()

	type CartItemResponse struct {
		Product  models.Product `json:"product"`
		Quantity int            `json:"quantity"`
	}

	var cartItems []CartItemResponse

	for rows.Next() {
		var product models.Product
		var quantity int

		err := rows.Scan(
			&product.ID, &product.Name, &product.Price, &product.ShortDescription,
			&product.LongDescription, &product.SKU, &product.Discount, &product.CategoryID,
			&product.CreatedAt, &product.UpdatedAt, &quantity,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to parse cart items",
			})
		}

		cartItems = append(cartItems, CartItemResponse{
			Product:  product,
			Quantity: quantity,
		})
	}

	return c.JSON(cartItems)
}

// RemoveFromCart - удалить товар из корзины
func RemoveFromCart(c *fiber.Ctx) error {
	var req struct {
		Email     string `json:"email"`
		ProductID int    `json:"productID"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Получаем ID пользователя
	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Удаляем товар из корзины
	_, err = database.DB.Exec(
		"DELETE FROM cart_items WHERE user_id = ? AND product_id = ?",
		userID, req.ProductID,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to remove from cart",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Item removed from cart",
	})
}
