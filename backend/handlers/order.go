package handlers

import (
	"database/sql"
	"encoding/json"
	"myAPI/database"
	"myAPI/models"
	"myAPI/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// CreateOrder - создание заказа с регистрацией пользователя
func CreateOrder(c *fiber.Ctx) error {
	var req models.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем существование товаров
	if len(req.ProductIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Product IDs are required",
		})
	}

	// Проверяем существует ли пользователь
	var existingUser models.User
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&existingUser.ID)
	if err != sql.ErrNoRows {
		return c.Status(400).JSON(fiber.Map{
			"error": "User already exists",
		})
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Создаем пользователя
	result, err := database.DB.Exec(
		"INSERT INTO users (email, password, name, phone, delivery_address, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		req.Email, string(hashedPassword), req.Name, req.Phone, req.DeliveryAddress, time.Now(), time.Now(),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	userID, _ := result.LastInsertId()

	// Создаем заказ
	productIDsJSON, _ := json.Marshal(req.ProductIDs)
	orderResult, err := database.DB.Exec(
		"INSERT INTO orders (user_id, product_ids, status, created_at) VALUES (?, ?, ?, ?)",
		userID, string(productIDsJSON), "новый", time.Now(),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	orderID, _ := orderResult.LastInsertId()

	// Получаем созданный заказ с товарами
	order, err := getOrderWithProducts(int(orderID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch order",
		})
	}

	// Генерируем токен
	token, err := utils.GenerateToken(int(userID), req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(models.OrderResponse{
		Order: *order,
		Token: token,
	})
}

// CreateOrderAuth - создание заказа для авторизованного пользователя
func CreateOrderAuth(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	var req models.CreateOrderAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем существование товаров
	if len(req.ProductIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Product IDs are required",
		})
	}

	// Обновляем профиль пользователя
	_, err := database.DB.Exec(
		"UPDATE users SET name = ?, phone = ?, delivery_address = ?, updated_at = ? WHERE id = ?",
		req.Name, req.Phone, req.DeliveryAddress, time.Now(), userID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update user profile",
		})
	}

	// Создаем заказ
	productIDsJSON, _ := json.Marshal(req.ProductIDs)
	orderResult, err := database.DB.Exec(
		"INSERT INTO orders (user_id, product_ids, status, created_at) VALUES (?, ?, ?, ?)",
		userID, string(productIDsJSON), "новый", time.Now(),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	orderID, _ := orderResult.LastInsertId()

	// Получаем созданный заказ с товарами
	order, err := getOrderWithProducts(int(orderID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch order",
		})
	}

	return c.JSON(models.OrderResponse{
		Order: *order,
	})
}

// GetUserOrders - получение списка заказов пользователя
func GetUserOrders(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	query := `
		SELECT id, user_id, product_ids, status, created_at
		FROM orders
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.ProductIDs, &order.Status, &order.CreatedAt)
		if err != nil {
			continue
		}

		// Получаем товары для каждого заказа
		products, err := getProductsByIDs(order.ProductIDs)
		if err == nil {
			order.Products = products
		}

		orders = append(orders, order)
	}

	return c.JSON(fiber.Map{
		"orders": orders,
	})
}

// Вспомогательные функции

func getOrderWithProducts(orderID int) (*models.Order, error) {
	var order models.Order
	var user models.User

	query := `
		SELECT o.id, o.user_id, o.product_ids, o.status, o.created_at,
			   u.id, u.email, COALESCE(u.name, ''), COALESCE(u.phone, ''), COALESCE(u.delivery_address, ''), u.created_at, u.updated_at
		FROM orders o
		JOIN users u ON o.user_id = u.id
		WHERE o.id = ?
	`

	err := database.DB.QueryRow(query, orderID).Scan(
		&order.ID, &order.UserID, &order.ProductIDs, &order.Status, &order.CreatedAt,
		&user.ID, &user.Email, &user.Name, &user.Phone, &user.DeliveryAddress, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	order.User = &user

	// Получаем товары
	products, err := getProductsByIDs(order.ProductIDs)
	if err == nil {
		order.Products = products
	}

	return &order, nil
}

func getProductsByIDs(productIDs models.IntArray) ([]models.Product, error) {
	if len(productIDs) == 0 {
		return []models.Product{}, nil
	}

	// Создаем плейсхолдеры для IN запроса
	placeholders := ""
	args := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = id
	}

	query := `
		SELECT p.id, p.name, p.price, p.short_description, p.long_description,
		       p.sku, p.discount, p.images, p.category_id, p.created_at, p.updated_at,
		       c.id, c.name, c.alias
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.id IN (` + placeholders + `)
	`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var category models.Category

		err := rows.Scan(
			&product.ID, &product.Name, &product.Price, &product.ShortDescription,
			&product.LongDescription, &product.SKU, &product.Discount, &product.Images,
			&product.CategoryID, &product.CreatedAt, &product.UpdatedAt,
			&category.ID, &category.Name, &category.Alias,
		)
		if err != nil {
			continue
		}

		product.Category = &category
		products = append(products, product)
	}

	return products, nil
}
