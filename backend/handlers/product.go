package handlers

import (
	"database/sql"
	"fmt"
	"myAPI/database"
	"myAPI/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var product models.Product
	var category models.Category

	query := `
		SELECT p.id, p.name, p.price, p.short_description, p.long_description, 
		       p.sku, p.discount, p.images, p.category_id, p.created_at, p.updated_at,
		       c.id, c.name, c.alias
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?
	`

	err = database.DB.QueryRow(query, productID).Scan(
		&product.ID, &product.Name, &product.Price, &product.ShortDescription,
		&product.LongDescription, &product.SKU, &product.Discount, &product.Images,
		&product.CategoryID, &product.CreatedAt, &product.UpdatedAt,
		&category.ID, &category.Name, &category.Alias,
	)

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

	product.Category = &category

	// Получаем отзывы для товара
	reviewsQuery := `
		SELECT id, product_id, name, text, rating, created_at
		FROM reviews
		WHERE product_id = ?
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(reviewsQuery, productID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch reviews",
		})
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(&review.ID, &review.ProductID, &review.Name,
			&review.Text, &review.Rating, &review.CreatedAt)
		if err != nil {
			continue
		}
		reviews = append(reviews, review)
	}

	return c.JSON(fiber.Map{
		"product": product,
		"reviews": reviews,
	})
}

func GetProducts(c *fiber.Ctx) error {
	var req models.ProductListRequest

	// Парсим query параметры
	if err := c.QueryParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	// Устанавливаем значения по умолчанию
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Строим запрос
	var conditions []string
	var args []interface{}
	// special case: ids list (comma separated)
	idsParam := c.Query("ids")
	if idsParam != "" {
		// parse ids
		parts := strings.Split(idsParam, ",")
		var placeholders []string
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			idVal, err := strconv.Atoi(p)
			if err != nil {
				continue
			}
			placeholders = append(placeholders, "?")
			args = append(args, idVal)
		}
		if len(placeholders) > 0 {
			conditions = append(conditions, fmt.Sprintf("p.id IN (%s)", strings.Join(placeholders, ",")))
		}
	}

	if req.CategoryID != nil {
		conditions = append(conditions, "p.category_id = ?")
		args = append(args, *req.CategoryID)
	}

	if req.PriceFrom != nil {
		conditions = append(conditions, "p.price >= ?")
		args = append(args, *req.PriceFrom)
	}

	if req.PriceTo != nil {
		conditions = append(conditions, "p.price <= ?")
		args = append(args, *req.PriceTo)
	}

	if req.HasDiscount != nil && *req.HasDiscount {
		conditions = append(conditions, "p.discount > 0")
	}

	if req.Search != "" {
		conditions = append(conditions, "(p.name LIKE ? OR p.short_description LIKE ? OR p.long_description LIKE ?)")
		searchTerm := "%" + req.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Запрос для подсчета общего количества
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM products p
		JOIN categories c ON p.category_id = c.id
		%s
	`, whereClause)

	var total int
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to count products",
		})
	}

	// Основной запрос
	query := fmt.Sprintf(`
		SELECT p.id, p.name, p.price, p.short_description, p.long_description,
		       p.sku, p.discount, p.images, p.category_id, p.created_at, p.updated_at,
		       c.id, c.name, c.alias
		FROM products p
		JOIN categories c ON p.category_id = c.id
		%s
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, req.Limit, req.Offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
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

	response := models.ProductListResponse{
		Products: products,
		Total:    total,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}

	return c.JSON(response)
}

// GetCategories возвращает список категорий
func GetCategories(c *fiber.Ctx) error {
	query := `SELECT id, name, alias FROM categories ORDER BY name`

	rows, err := database.DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Alias); err != nil {
			continue
		}
		categories = append(categories, category)
	}

	// Возвращаем объект с полем categories для совместимости с клиентом
	return c.JSON(fiber.Map{"categories": categories})
}
