package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"myAPI/database"
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

// AdminListBackups - получить список доступных бэкапов
func AdminListBackups(c *fiber.Ctx) error {
	backupDir := "backups"
	
	// Если папка не существует, создадим её
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create backup dir"})
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to read backups dir"})
	}

	type BackupInfo struct {
		Name    string `json:"name"`
		Path    string `json:"path"`
		ModTime string `json:"modTime"`
		Size    int64  `json:"size"`
	}

	var backups []BackupInfo

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".db" {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			backups = append(backups, BackupInfo{
				Name:    entry.Name(),
				Path:    filepath.Join(backupDir, entry.Name()),
				ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
				Size:    info.Size(),
			})
		}
	}

	return c.JSON(fiber.Map{"backups": backups})
}

// AdminRestoreBackup - восстановить БД из бэкапа
func AdminRestoreBackup(c *fiber.Ctx) error {
	var req struct {
		BackupName string `json:"backupName"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.BackupName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "backupName is required"})
	}

	backupPath := filepath.Join("backups", req.BackupName)

	// Проверим что файл существует и находится в папке backups
	absBackupPath, err := filepath.Abs(backupPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to get absolute path"})
	}

	absBackupDir, err := filepath.Abs("backups")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to get absolute path"})
	}

	// Защита от path traversal атак
	if !filepath.HasPrefix(absBackupPath, absBackupDir) {
		return c.Status(400).JSON(fiber.Map{"error": "invalid backup path"})
	}

	_, err = os.Stat(backupPath)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "backup file not found"})
	}

	// Создаем бэкап текущей БД перед восстановлением (как страховка)
	currentDBBackup := filepath.Join("backups", "backup-before-restore-"+time.Now().Format("20060102-150405")+".db")
	if in, err := os.Open("app.db"); err == nil {
		defer in.Close()
		if out, err := os.Create(currentDBBackup); err == nil {
			defer out.Close()
			io.Copy(out, in)
		}
	}

	// Читаем бэкап
	in, err := os.Open(backupPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to open backup file"})
	}
	defer in.Close()

	// Пишем в app.db
	out, err := os.Create("app.db")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create app.db"})
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to restore backup"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "backup restored successfully"})
}


func AdminCreateProduct(c *fiber.Ctx) error {
	// Parse incoming body as generic map to be tolerant to different client payload shapes
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	id, err := createProduct(body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create product"})
	}

	body["id"] = id
	return c.Status(201).JSON(body)
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

// AdminUploadFile загружает файл и сохраняет его в соответствующую папку images
func AdminUploadFile(c *fiber.Ctx) error {
	kind := c.Params("kind")

	// map resource kind to folder name
	folder := ""
	switch kind {
	case "products":
		folder = "jewelry"
	case "news":
		folder = "news"
	case "banners":
		folder = "banner"
	default:
		return c.Status(400).JSON(fiber.Map{"error": "unknown upload kind"})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file not provided"})
	}

	// ensure images/<folder> exists
	destDir := filepath.Join("images", folder)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create dir"})
	}

	// generate filename
	fname := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	dst := filepath.Join(destDir, fname)

	if err := c.SaveFile(fileHeader, dst); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to save file"})
	}

	// return public path
	publicPath := "/images/" + folder + "/" + fname
	return c.JSON(fiber.Map{"path": publicPath})
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

	// allow admin-provided timestamps
	createdAt := parseTimeFromMap(data, "created_at")
	updatedAt := parseTimeFromMap(data, "updated_at")

	result, err := database.DB.Exec(`
		INSERT INTO products (name, price, short_description, long_description, sku, discount, images, category_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, name, price, shortDesc, longDesc, sku, discount, images, categoryID, createdAt, updatedAt)

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

	updatedAt := parseTimeFromMap(data, "updated_at")

	_, err := database.DB.Exec(`
		UPDATE products 
		SET name = ?, price = ?, short_description = ?, long_description = ?, sku = ?, discount = ?, images = ?, category_id = ?, updated_at = ? 
		WHERE id = ?
	`, name, price, shortDesc, longDesc, sku, discount, images, categoryID, updatedAt, id)

	return err
}

func deleteProduct(id int) error {
	// Delete dependent records (reviews, banners) first to avoid foreign key constraint errors
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM reviews WHERE product_id = ?", id); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM banners WHERE product_id = ?", id); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM products WHERE id = ?", id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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
	rows, err := database.DB.Query(`SELECT id, user_id, product_ids, status, price, created_at FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id, userID int
		var productIDsStr, status string
		var price float64
		var createdAt time.Time

		if err := rows.Scan(&id, &userID, &productIDsStr, &status, &price, &createdAt); err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":           id,
			"user_id":      userID,
			"product_ids":  productIDsStr,
			"status":       status,
			"price":        price,
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
		status = "оплачен"
	}

	createdAt := parseTimeFromMap(data, "created_at")

	result, err := database.DB.Exec(`INSERT INTO orders (user_id, product_ids, status, created_at) VALUES (?, ?, ?, ?)`,
		userID, productIDs, status, createdAt)

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

	createdAt := parseTimeFromMap(data, "created_at")

	result, err := database.DB.Exec(`INSERT INTO news (title, description, image, created_at) VALUES (?, ?, ?, ?)`,
		title, description, image, createdAt)

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

	createdAt := parseTimeFromMap(data, "created_at")

	result, err := database.DB.Exec(`INSERT INTO banners (product_id, image, position, created_at) VALUES (?, ?, ?, ?)`,
		productID, image, position, createdAt)

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

// Helper function to convert sql.NullString to string
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
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
		var email, role string
		var name, phone, deliveryAddress sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &email, &role, &name, &phone, &deliveryAddress, &createdAt, &updatedAt); err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":                 id,
			"email":              email,
			"role":               role,
			"name":               nullStringToString(name),
			"phone":              nullStringToString(phone),
			"delivery_address":   nullStringToString(deliveryAddress),
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

	createdAt := parseTimeFromMap(data, "created_at")
	updatedAt := parseTimeFromMap(data, "updated_at")

	result, err := database.DB.Exec(`
		INSERT INTO users (email, password, role, name, phone, delivery_address, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, email, password, role, name, phone, deliveryAddress, createdAt, updatedAt)

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

	updatedAt := parseTimeFromMap(data, "updated_at")

	_, err := database.DB.Exec(`
		UPDATE users 
		SET email = ?, role = ?, name = ?, phone = ?, delivery_address = ?, updated_at = ? 
		WHERE id = ?
	`, email, role, name, phone, deliveryAddress, updatedAt, id)

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

// parseTimeFromMap parses a time string from the provided map at key and returns time.Time.
// It understands RFC3339 and the 'datetime-local' format '2006-01-02T15:04'.
func parseTimeFromMap(data map[string]interface{}, key string) time.Time {
	if data == nil {
		return time.Now()
	}
	s := toString(data[key])
	if s == "" {
		return time.Now()
	}

	// try RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}

	// try datetime-local format without seconds
	layouts := []string{"2006-01-02T15:04", "2006-01-02 15:04:05", "2006-01-02T15:04:05"}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t
		}
	}

	// fallback
	return time.Now()
}
