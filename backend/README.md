# Nuxt API

API на Golang Fiber с SQLite базой данных для аутентификации пользователей, каталога товаров и системы заказов.

## Установка и запуск

1. Установите зависимости:
```bash
go mod tidy
```

2. Запустите приложение:
```bash
go run main.go
```

Сервер запустится на порту 3000.

## Тестирование API

### Insomnia
Для удобного тестирования API импортируйте файл `insomnia-api.json` в Insomnia:
1. Откройте Insomnia
2. Нажмите "Import/Export" → "Import Data" → "From File"
3. Выберите файл `insomnia-api.json`
4. Все endpoints будут доступны с примерами запросов

### Переменные окружения
В Insomnia настроены переменные:
- `base_url` - базовый URL API (по умолчанию http://localhost:3000)
- `auth_token` - JWT токен для авторизованных запросов

## API Endpoints

### Аутентификация

#### Регистрация
```
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Логин
```
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Товары

#### Получить список товаров
```
GET /api/products?limit=20&offset=0&category_id=1&price_from=1000&price_to=50000&has_discount=true&search=iPhone
```

Параметры запроса:
- `limit` - количество товаров (по умолчанию 20, максимум 100)
- `offset` - смещение для пагинации (по умолчанию 0)
- `category_id` - ID категории для фильтрации
- `price_from` - минимальная цена
- `price_to` - максимальная цена
- `has_discount` - только товары со скидкой (true/false)
- `search` - поиск по названию и описанию

#### Получить товар по ID
```
GET /api/products/1
```

### Заказы

#### Создание заказа с регистрацией пользователя
```
POST /api/orders
Content-Type: application/json

{
  "product_ids": [1, 2, 3],
  "email": "user@example.com",
  "password": "password123",
  "delivery_address": "ул. Пушкина, д. 10, кв. 5",
  "name": "Иван Иванов",
  "phone": "+7 (999) 123-45-67"
}
```

Ответ:
```json
{
  "order": {
    "id": 1,
    "user_id": 1,
    "product_ids": [1, 2, 3],
    "status": "новый",
    "created_at": "2024-01-01T00:00:00Z",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "Иван Иванов",
      "phone": "+7 (999) 123-45-67",
      "delivery_address": "ул. Пушкина, д. 10, кв. 5"
    },
    "products": [...]
  },
  "token": "jwt_token_here"
}
```

#### Создание заказа для авторизованного пользователя
```
POST /api/orders/auth
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_ids": [1, 2],
  "delivery_address": "ул. Ленина, д. 20, кв. 15",
  "name": "Иван Петров",
  "phone": "+7 (999) 987-65-43"
}
```

#### Получение списка заказов пользователя
```
GET /api/orders
Authorization: Bearer <token>
```

Ответ:
```json
{
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "product_ids": [1, 2],
      "status": "новый",
      "created_at": "2024-01-01T00:00:00Z",
      "products": [
        {
          "id": 1,
          "name": "iPhone 15 Pro",
          "price": 99999.99,
          "category": {
            "id": 1,
            "name": "Электроника",
            "alias": "electronics"
          }
        }
      ]
    }
  ]
}
```

## Структура проекта

```
├── main.go              # Основной файл приложения
├── go.mod               # Зависимости Go
├── insomnia-api.json    # Конфигурация Insomnia для тестирования
├── models/
│   ├── user.go          # Модели пользователей
│   ├── product.go       # Модели товаров, категорий, отзывов
│   └── order.go         # Модели заказов
├── database/
│   ├── database.go      # Настройка базы данных
│   └── seed.go          # Тестовые данные
├── handlers/
│   ├── auth.go          # Хендлеры аутентификации
│   ├── product.go       # Хендлеры товаров
│   └── order.go         # Хендлеры заказов
├── utils/
│   ├── jwt.go           # Утилиты для JWT
│   └── middleware.go    # Middleware для аутентификации
└── app.db               # SQLite база данных (создается автоматически)
```

## База данных

### Таблицы:
- **users** - пользователи (с полями для доставки)
- **categories** - категории товаров
- **products** - товары
- **reviews** - отзывы на товары
- **orders** - заказы

### Тестовые данные:
При первом запуске автоматически создаются:
- 4 категории (Электроника, Одежда, Книги, Спорт)
- 5 товаров с разными параметрами
- Отзывы на товары

## Аутентификация

Для защищенных endpoints используется JWT токен в заголовке:
```
Authorization: Bearer <your_jwt_token>
```

## Технологии

- **Fiber** - веб-фреймворк
- **SQLite** - база данных
- **ncruces/go-sqlite3** - драйвер SQLite без CGO
- **JWT** - токены аутентификации
- **bcrypt** - хеширование паролей 