# Insomnia Test Collection для Admin API

## Setup

1. Получите токен администратора через LOGIN запрос
2. Скопируйте токен и используйте в переменной `{{token}}`
3. Используйте примеры ниже для тестирования

---

## 1. LOGIN - Получить токен администратора

**Method:** POST
**URL:** `http://localhost:3000/api/auth/login`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "role": "admin",
    "name": "Admin User",
    "created_at": "2025-12-26T10:30:00Z"
  }
}
```

**Сохраните токен в переменную:**
- Environment Variable: `token` = `{{response.body.token}}`

---

## PRODUCTS

### GET - Получить все товары

**Method:** GET
**URL:** `http://localhost:3000/api/admin/products`

**Headers:**
```
Authorization: Bearer {{token}}
```

**Expected Status:** 200

---

### POST - Создать новый товар

**Method:** POST
**URL:** `http://localhost:3000/api/admin/products`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Premium Bracelet",
  "price": 199.99,
  "short_description": "Elegant gold bracelet",
  "long_description": "A beautiful and elegant gold bracelet with precious stones",
  "sku": "BR-GOLD-001",
  "discount": 15,
  "images": "[\"bracelet1.jpg\", \"bracelet2.jpg\"]",
  "category_id": 1
}
```

**Expected Status:** 201

---

### PUT - Обновить товар

**Method:** PUT
**URL:** `http://localhost:3000/api/admin/products/1`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Premium Gold Bracelet",
  "price": 249.99,
  "discount": 20
}
```

**Expected Status:** 200

---

### DELETE - Удалить товар

**Method:** DELETE
**URL:** `http://localhost:3000/api/admin/products/1`

**Headers:**
```
Authorization: Bearer {{token}}
```

**Expected Status:** 200

---

## CATEGORIES

### GET - Получить все категории

**Method:** GET
**URL:** `http://localhost:3000/api/admin/categories`

**Headers:**
```
Authorization: Bearer {{token}}
```

---

### POST - Создать категорию

**Method:** POST
**URL:** `http://localhost:3000/api/admin/categories`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Rings",
  "alias": "rings"
}
```

---

### PUT - Обновить категорию

**Method:** PUT
**URL:** `http://localhost:3000/api/admin/categories/1`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Gold Rings",
  "alias": "gold-rings"
}
```

---

### DELETE - Удалить категорию

**Method:** DELETE
**URL:** `http://localhost:3000/api/admin/categories/1`

**Headers:**
```
Authorization: Bearer {{token}}
```

---

## NEWS

### GET - Получить все новости

**Method:** GET
**URL:** `http://localhost:3000/api/admin/news`

**Headers:**
```
Authorization: Bearer {{token}}
```

---

### POST - Создать новость

**Method:** POST
**URL:** `http://localhost:3000/api/admin/news`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "title": "New Collection Launch",
  "description": "We are excited to announce our new winter collection of jewelry",
  "image": "news-image.jpg"
}
```

---

### PUT - Обновить новость

**Method:** PUT
**URL:** `http://localhost:3000/api/admin/news/1`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "title": "Winter Collection Now Available",
  "description": "Discover our exclusive winter jewelry collection with premium materials"
}
```

---

### DELETE - Удалить новость

**Method:** DELETE
**URL:** `http://localhost:3000/api/admin/news/1`

**Headers:**
```
Authorization: Bearer {{token}}
```

---

## BANNERS

### GET - Получить все баннеры

**Method:** GET
**URL:** `http://localhost:3000/api/admin/banners`

**Headers:**
```
Authorization: Bearer {{token}}
```

---

### POST - Создать баннер

**Method:** POST
**URL:** `http://localhost:3000/api/admin/banners`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "product_id": 5,
  "image": "banner-product5.jpg",
  "position": 1
}
```

---

### PUT - Обновить баннер

**Method:** PUT
**URL:** `http://localhost:3000/api/admin/banners/1`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "product_id": 5,
  "image": "banner-product5-updated.jpg",
  "position": 0
}
```

---

### DELETE - Удалить баннер

**Method:** DELETE
**URL:** `http://localhost:3000/api/admin/banners/1`

**Headers:**
```
Authorization: Bearer {{token}}
```

---

## ORDERS

### GET - Получить все заказы

**Method:** GET
**URL:** `http://localhost:3000/api/admin/orders`

**Headers:**
```
Authorization: Bearer {{token}}
```

---

### POST - Создать заказ (админ)

**Method:** POST
**URL:** `http://localhost:3000/api/admin/orders`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "user_id": 2,
  "product_ids": "[1, 2, 3]",
  "status": "обработан"
}
```

---

### PUT - Обновить статус заказа

**Method:** PUT
**URL:** `http://localhost:3000/api/admin/orders/1`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "status": "отправлен"
}
```

---

### DELETE - Удалить заказ

**Method:** DELETE
**URL:** `http://localhost:3000/api/admin/orders/1`

**Headers:**
```
Authorization: Bearer {{token}}
```

---

## USERS

### GET - Получить всех пользователей

**Method:** GET
**URL:** `http://localhost:3000/api/admin/users`

**Headers:**
```
Authorization: Bearer {{token}}
```

---

### POST - Создать пользователя (админ)

**Method:** POST
**URL:** `http://localhost:3000/api/admin/users`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "email": "newuser@example.com",
  "name": "John Doe",
  "phone": "+1234567890",
  "delivery_address": "123 Main St",
  "role": "user"
}
```

**Note:** Пароль автоматически установится как "changeme"

---

### PUT - Обновить пользователя

**Method:** PUT
**URL:** `http://localhost:3000/api/admin/users/2`

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Jane Doe",
  "phone": "+9876543210",
  "delivery_address": "456 Oak Ave"
}
```

---

### DELETE - Удалить пользователя

**Method:** DELETE
**URL:** `http://localhost:3000/api/admin/users/2`

**Headers:**
```
Authorization: Bearer {{token}}
```

---

## ERROR HANDLING

### Test 401 - Без токена

**Method:** GET
**URL:** `http://localhost:3000/api/admin/products`

**Expected Response:**
```json
{
  "error": "Unauthorized"
}
```

**Status:** 401

---

### Test 403 - Обычный пользователь

Залогинитесь как обычный пользователь (role = 'user') и попробуйте админ endpoint.

**Expected Response:**
```json
{
  "error": "Forbidden"
}
```

**Status:** 403

---

### Test 400 - Неверный ресурс

**Method:** GET
**URL:** `http://localhost:3000/api/admin/invalid_resource`

**Headers:**
```
Authorization: Bearer {{token}}
```

**Expected Response:**
```json
{
  "error": "Unknown resource"
}
```

**Status:** 400

---

## TIPS

1. **Переменные:** Сохраняйте токен в переменную окружения для удобства
2. **Pre-request Script:** Используйте pre-request скрипты для подготовки данных
3. **Tests:** Настройте автоматические проверки ответов
4. **Documentation:** Экспортируйте коллекцию для команды

### Пример Pre-request Script:

```javascript
const token = pm.environment.get("token");
if (!token) {
  console.log("Warning: token not set");
}
pm.request.headers.upsert({
  key: "Authorization",
  value: `Bearer ${token}`
});
```

---

## CHECKLIST

- [ ] Получен токен администратора
- [ ] Все GET запросы работают (200)
- [ ] Все POST запросы работают (201)
- [ ] Все PUT запросы работают (200)
- [ ] Все DELETE запросы работают (200)
- [ ] Проверены ошибки (401, 403, 400)
- [ ] Тестирование завершено
