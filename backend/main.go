package main

import (
	"log"
	"myAPI/database"
	"myAPI/handlers"
	"myAPI/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Проверяем, существовала ли БД до инициализации
	needSeed := false
	if _, err := os.Stat("app.db"); os.IsNotExist(err) {
		needSeed = true
	} else if err != nil {
		log.Printf("Warning checking app.db: %v", err)
	}

	// Инициализация базы данных
	database.InitDatabase()

	// Заполнение тестовыми данными (только при первом запуске, если БД не существовала до инициализации)
	if needSeed {
		database.SeedData()
	}
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
	products.Post(":id/reviews", handlers.CreateReview)

	// Категории
	categories := api.Group("/categories")
	categories.Get("/", handlers.GetCategories)

	// Баннеры
	banners := api.Group("/banners")
	banners.Get("/", handlers.GetBanners)

	// Новости
	news := api.Group("/news")
	news.Get("/", handlers.GetNews)

	// Админ-панель (требует авторизацию и роль admin)
	admin := api.Group("/admin", utils.AuthMiddleware, utils.AdminMiddleware)
	admin.Post("/backup", handlers.AdminBackup)
	admin.Get("/backups", handlers.AdminListBackups)
	admin.Post("/restore", handlers.AdminRestoreBackup)
	admin.Post("/products", handlers.AdminCreateProduct)
	admin.Post("/upload/:kind", handlers.AdminUploadFile)
	
	// Generic admin CRUD endpoints for all resources
	admin.Get("/:resource", handlers.AdminGetResource)           // GET /api/admin/{resource}
	admin.Post("/:resource", handlers.AdminCreateResource)       // POST /api/admin/{resource}
	admin.Put("/:resource/:id", handlers.AdminUpdateResource)    // PUT /api/admin/{resource}/{id}
	admin.Delete("/:resource/:id", handlers.AdminDeleteResource) // DELETE /api/admin/{resource}/{id}

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

	// Избранное (favorites)
	favorites := api.Group("/favorites")
	favorites.Post("/", handlers.SaveFavorites)
	favorites.Get("/", handlers.GetFavorites)

	// Запуск сервера
	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
