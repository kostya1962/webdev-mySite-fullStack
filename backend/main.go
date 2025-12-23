package main

import (
	"log"
	"myAPI/database"
	"myAPI/handlers"
	"myAPI/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Инициализация базы данных
	database.InitDatabase()

	// Заполнение тестовыми данными (только при первом запуске)
	database.SeedData()

	// Создание Fiber приложения
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Статические файлы для изображений
	app.Static("/images", "./images")

	// Роуты
	api := app.Group("/api")

	// Аутентификация
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)

	// Товары
	products := api.Group("/products")
	products.Get("/", handlers.GetProducts)
	products.Get("/:id", handlers.GetProduct)

	// Категории
	categories := api.Group("/categories")
	categories.Get("/", handlers.GetCategories)

	// Баннеры
	banners := api.Group("/banners")
	banners.Get("/", handlers.GetBanners)

	// Заказы
	orders := api.Group("/orders")
	orders.Post("/", handlers.CreateOrder)                               // Создание заказа с регистрацией
	orders.Post("/auth", utils.AuthMiddleware, handlers.CreateOrderAuth) // Создание заказа для авторизованного пользователя
	orders.Get("/", utils.AuthMiddleware, handlers.GetUserOrders)        // Получение заказов пользователя

	// Корзина
	cart := api.Group("/cart")
	cart.Post("/", handlers.AddToCart)           // Добавить товар в корзину
	cart.Get("/", handlers.GetCart)              // Получить корзину пользователя
	cart.Delete("/", handlers.RemoveFromCart)    // Удалить товар из корзины

	// Запуск сервера
	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
