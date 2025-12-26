package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"myAPI/database"
	"myAPI/models"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AdminBackup создает копию файла базы данных и возвращает путь
func AdminBackup(c *fiber.Ctx) error {
	src := "app.db"
	// Создаем папку backups
	backupDir := "backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create backup dir"})
	}

	dst := filepath.Join(backupDir, "backup-"+time.Now().Format("20060102-150405")+".db")

	in, err := os.Open(src)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to open source db"})
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create backup file"})
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to copy db"})
	}

	return c.JSON(fiber.Map{"success": true, "path": dst})
}

// AdminCreateProduct - простой эндпоинт для создания товара админом
func AdminCreateProduct(c *fiber.Ctx) error {
	var p models.Product
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	imagesJSON := "[]"
	if p.Images != nil {
		if v, err := p.Images.Value(); err == nil {
			if s, ok := v.(string); ok {
				imagesJSON = s
			}
		}
	}

	res, err := database.DB.Exec(`INSERT INTO products (name, price, short_description, long_description, sku, discount, images, category_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name, p.Price, p.ShortDescription, p.LongDescription, p.SKU, p.Discount, imagesJSON, p.CategoryID, time.Now(), time.Now(),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create product"})
	}
	id, _ := res.LastInsertId()
	p.ID = int(id)

	return c.Status(201).JSON(p)
}

// AdminGetResource - получить все записи ресурса
func AdminGetResource(c *fiber.Ctx) error {
	resource := c.Params("resource")

	var items []map[string]interface{}
	var err error

	switch resource {
	case "products":
		items, err = getProductsForAdmin()
	case "categories":
		items, err = getCategoriesForAdmin()
	case "orders":
		items, err = getOrdersForAdmin()
	case "news":
		items, err = getNewsForAdmin()
	case "banners":
		items, err = getBannersForAdmin()
	case "users":
		items, err = getUsersForAdmin()
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Unknown resource"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if items == nil {
		items = []map[string]interface{}{}
	}

	return c.JSON(fiber.Map{"data": items})
}

// AdminCreateResource - создать новую запись ресурса
func AdminCreateResource(c *fiber.Ctx) error {
	resource := c.Params("resource")
	var body map[string]interface{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var id int64
	var err error

	switch resource {
	case "products":
		id, err = createProduct(body)
	case "categories":
		id, err = createCategory(body)
	case "orders":
		id, err = createOrder(body)
	case "news":
		id, err = createNews(body)
	case "banners":
		id, err = createBanner(body)
	case "users":
		id, err = createUser(body)
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Unknown resource"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	body["id"] = id
	return c.Status(201).JSON(body)
}

// AdminUpdateResource - обновить запись ресурса
func AdminUpdateResource(c *fiber.Ctx) error {
	resource := c.Params("resource")
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	switch resource {
	case "products":
		err = updateProduct(id, body)
	case "categories":
		err = updateCategory(id, body)
	case "orders":
		err = updateOrder(id, body)
	case "news":
		err = updateNews(id, body)
	case "banners":
		err = updateBanner(id, body)
	case "users":
		err = updateUser(id, body)
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Unknown resource"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(body)
}

// AdminDeleteResource - удалить запись ресурса
func AdminDeleteResource(c *fiber.Ctx) error {
	resource := c.Params("resource")
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	switch resource {
	case "products":
		err = deleteProduct(id)
	case "categories":
		err = deleteCategory(id)
	case "orders":
		err = deleteOrder(id)
	case "news":
		err = deleteNews(id)
	case "banners":
		err = deleteBanner(id)
	case "users":
		err = deleteUser(id)
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Unknown resource"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

// Helper functions for Products
func getProductsForAdmin() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT id, name, price, short_description, long_description, sku, discount, images, category_id, created_at, updated_at 
		FROM products
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id, discount, categoryID int
		var name, shortDesc, longDesc, sku, images string
		var price float64
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &name, &price, &shortDesc, &longDesc, &sku, &discount, &images, &categoryID, &createdAt, &updatedAt); err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":                  id,
			"name":                name,
			"price":               price,
			"short_description":   shortDesc,
			"long_description":    longDesc,
			"sku":                 sku,
			"discount":            discount,
			"images":              images,
			"category_id":         categoryID,
			"created_at":          createdAt.Format(time.RFC3339),
			"updated_at":          updatedAt.Format(time.RFC3339),
		}
		items = append(items, item)
	}

	return items, nil
}

func createProduct(data map[string]interface{}) (int64, error) {
	name := toString(data["name"])
	price := toFloat64(data["price"])
	shortDesc := toString(data["short_description"])
	longDesc := toString(data["long_description"])
	sku := toString(data["sku"])
	discount := toInt(data["discount"])
	images := toString(data["images"])
	categoryID := toInt(data["category_id"])

	if images == "" {
		images = "[]"
	}

	result, err := database.DB.Exec(`
		INSERT INTO products (name, price, short_description, long_description, sku, discount, images, category_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, name, price, shortDesc, longDesc, sku, discount, images, categoryID, time.Now(), time.Now())

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func updateProduct(id int, data map[string]interface{}) error {
	name := toString(data["name"])
	price := toFloat64(data["price"])
	shortDesc := toString(data["short_description"])
	longDesc := toString(data["long_description"])
	sku := toString(data["sku"])
	discount := toInt(data["discount"])
	images := toString(data["images"])
	categoryID := toInt(data["category_id"])

	if images == "" {
		images = "[]"
	}

	_, err := database.DB.Exec(`
		UPDATE products 
		SET name = ?, price = ?, short_description = ?, long_description = ?, sku = ?, discount = ?, images = ?, category_id = ?, updated_at = ? 
		WHERE id = ?
	`, name, price, shortDesc, longDesc, sku, discount, images, categoryID, time.Now(), id)

	return err
}

func deleteProduct(id int) error {
	_, err := database.DB.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}

// Helper functions for Categories
func getCategoriesForAdmin() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`SELECT id, name, alias FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id int
		var name, alias string

		if err := rows.Scan(&id, &name, &alias); err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":    id,
			"name":  name,
			"alias": alias,
		}
		items = append(items, item)
	}

	return items, nil
}

func createCategory(data map[string]interface{}) (int64, error) {
	name := toString(data["name"])
	alias := toString(data["alias"])

	result, err := database.DB.Exec(`INSERT INTO categories (name, alias) VALUES (?, ?)`, name, alias)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func updateCategory(id int, data map[string]interface{}) error {
	name := toString(data["name"])
	alias := toString(data["alias"])

	_, err := database.DB.Exec(`UPDATE categories SET name = ?, alias = ? WHERE id = ?`, name, alias, id)
	return err
}

func deleteCategory(id int) error {
	_, err := database.DB.Exec("DELETE FROM categories WHERE id = ?", id)
	return err
}

// Helper functions for Orders
func getOrdersForAdmin() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`SELECT id, user_id, product_ids, status, created_at FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id, userID int
		var productIDsStr, status string
		var createdAt time.Time

		if err := rows.Scan(&id, &userID, &productIDsStr, &status, &createdAt); err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":           id,
			"user_id":      userID,
			"product_ids":  productIDsStr,
			"status":       status,
			"created_at":   createdAt.Format(time.RFC3339),
		}
		items = append(items, item)
	}

	return items, nil
}

func createOrder(data map[string]interface{}) (int64, error) {
	userID := toInt(data["user_id"])
	productIDs := toString(data["product_ids"])
	status := toString(data["status"])

	if status == "" {
		status = "новый"
	}

	result, err := database.DB.Exec(`INSERT INTO orders (user_id, product_ids, status, created_at) VALUES (?, ?, ?, ?)`,
		userID, productIDs, status, time.Now())

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func updateOrder(id int, data map[string]interface{}) error {
	userID := toInt(data["user_id"])
	productIDs := toString(data["product_ids"])
	status := toString(data["status"])

	_, err := database.DB.Exec(`UPDATE orders SET user_id = ?, product_ids = ?, status = ? WHERE id = ?`,
		userID, productIDs, status, id)

	return err
}

func deleteOrder(id int) error {
	_, err := database.DB.Exec("DELETE FROM orders WHERE id = ?", id)
	return err
}

// Helper functions for News
func getNewsForAdmin() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`SELECT id, title, description, image, created_at FROM news`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id int
		var title, description, image string
		var createdAt time.Time

		if err := rows.Scan(&id, &title, &description, &image, &createdAt); err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": description,
			"image":       image,
			"created_at":  createdAt.Format(time.RFC3339),
		}
		items = append(items, item)
	}

	return items, nil
}

func createNews(data map[string]interface{}) (int64, error) {
	title := toString(data["title"])
	description := toString(data["description"])
	image := toString(data["image"])

	result, err := database.DB.Exec(`INSERT INTO news (title, description, image, created_at) VALUES (?, ?, ?, ?)`,
		title, description, image, time.Now())

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func updateNews(id int, data map[string]interface{}) error {
	title := toString(data["title"])
	description := toString(data["description"])
	image := toString(data["image"])

	_, err := database.DB.Exec(`UPDATE news SET title = ?, description = ?, image = ? WHERE id = ?`,
		title, description, image, id)

	return err
}

func deleteNews(id int) error {
	_, err := database.DB.Exec("DELETE FROM news WHERE id = ?", id)
	return err
}

// Helper functions for Banners
func getBannersForAdmin() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`SELECT id, product_id, image, position, created_at FROM banners`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id, productID, position int
		var image string
		var createdAt time.Time

		if err := rows.Scan(&id, &productID, &image, &position, &createdAt); err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":         id,
			"product_id": productID,
			"image":      image,
			"position":   position,
			"created_at": createdAt.Format(time.RFC3339),
		}
		items = append(items, item)
	}

	return items, nil
}

func createBanner(data map[string]interface{}) (int64, error) {
	productID := toInt(data["product_id"])
	image := toString(data["image"])
	position := toInt(data["position"])

	result, err := database.DB.Exec(`INSERT INTO banners (product_id, image, position, created_at) VALUES (?, ?, ?, ?)`,
		productID, image, position, time.Now())

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func updateBanner(id int, data map[string]interface{}) error {
	productID := toInt(data["product_id"])
	image := toString(data["image"])
	position := toInt(data["position"])

	_, err := database.DB.Exec(`UPDATE banners SET product_id = ?, image = ?, position = ? WHERE id = ?`,
		productID, image, position, id)

	return err
}

func deleteBanner(id int) error {
	_, err := database.DB.Exec("DELETE FROM banners WHERE id = ?", id)
	return err
}

// Helper functions for Users
func getUsersForAdmin() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT id, email, role, name, phone, delivery_address, created_at, updated_at 
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id int
		var email, role, name, phone, deliveryAddress string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &email, &role, &name, &phone, &deliveryAddress, &createdAt, &updatedAt); err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":                 id,
			"email":              email,
			"role":               role,
			"name":               name,
			"phone":              phone,
			"delivery_address":   deliveryAddress,
			"created_at":         createdAt.Format(time.RFC3339),
			"updated_at":         updatedAt.Format(time.RFC3339),
		}
		items = append(items, item)
	}

	return items, nil
}

func createUser(data map[string]interface{}) (int64, error) {
	email := toString(data["email"])
	role := toString(data["role"])
	name := toString(data["name"])
	phone := toString(data["phone"])
	deliveryAddress := toString(data["delivery_address"])

	if role == "" {
		role = "user"
	}

	// Используем пустой пароль или временный
	password := "changeme"

	result, err := database.DB.Exec(`
		INSERT INTO users (email, password, role, name, phone, delivery_address, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, email, password, role, name, phone, deliveryAddress, time.Now(), time.Now())

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func updateUser(id int, data map[string]interface{}) error {
	email := toString(data["email"])
	role := toString(data["role"])
	name := toString(data["name"])
	phone := toString(data["phone"])
	deliveryAddress := toString(data["delivery_address"])

	_, err := database.DB.Exec(`
		UPDATE users 
		SET email = ?, role = ?, name = ?, phone = ?, delivery_address = ?, updated_at = ? 
		WHERE id = ?
	`, email, role, name, phone, deliveryAddress, time.Now(), id)

	return err
}

func deleteUser(id int) error {
	_, err := database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// Utility conversion functions
func toString(val interface{}) string {
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%v", int64(v))
	case int:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%v", v)
	default:
		if b, err := json.Marshal(v); err == nil {
			return string(b)
		}
		return fmt.Sprintf("%v", v)
	}
}

func toInt(val interface{}) int {
	if val == nil {
		return 0
	}

	switch v := val.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	return 0
}

func toFloat64(val interface{}) float64 {
	if val == nil {
		return 0
	}

	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}

	return 0
}

func toIntArray(val interface{}) []int {
	if val == nil {
		return []int{}
	}

	var arr []int
	switch v := val.(type) {
	case string:
		var jsonArr []int
		if err := json.Unmarshal([]byte(v), &jsonArr); err == nil {
			return jsonArr
		}
	case []interface{}:
		for _, item := range v {
			if intVal, ok := item.(float64); ok {
				arr = append(arr, int(intVal))
			}
		}
		return arr
	}

	return arr
}
