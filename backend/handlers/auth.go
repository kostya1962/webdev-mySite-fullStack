package handlers

import (
	"database/sql"
	"myAPI/database"
	"myAPI/models"
	"myAPI/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
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
		"INSERT INTO users (email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		req.Email, string(hashedPassword), "user", time.Now(), time.Now(),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	userID, _ := result.LastInsertId()

	// Генерируем токен
	token, err := utils.GenerateToken(int(userID), req.Email, "user")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	user := models.User{
		ID:        int(userID),
		Email:     req.Email,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return c.JSON(models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ищем пользователя
	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, email, password, COALESCE(role, 'user'), COALESCE(name, ''), COALESCE(phone, ''), COALESCE(delivery_address, ''), created_at, updated_at FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.Name, &user.Phone, &user.DeliveryAddress, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Генерируем токен
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Убираем пароль из ответа
	user.Password = ""

	return c.JSON(models.AuthResponse{
		Token: token,
		User:  user,
	})
}
