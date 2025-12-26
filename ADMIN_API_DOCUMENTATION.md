# Admin Panel API Documentation

## Overview

Backend теперь полностью поддерживает работу с админ-панелью frontend'а. Все CRUD операции для управления базой данных доступны через унифицированные API endpoints.

## Authentication

Все администраторские endpoints требуют:
1. **Authorization Header**: `Authorization: Bearer {token}`
2. **Admin Role**: Пользователь должен иметь роль `admin` в базе данных

## API Endpoints

### Получить все записи ресурса
```
GET /api/admin/{resource}
```

**Параметры:**
- `resource` - тип ресурса: `products`, `categories`, `orders`, `news`, `banners`, `users`

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Product Name",
      "price": 99.99,
      ...
    }
  ]
}
```

**Example:**
```bash
curl -H "Authorization: Bearer {token}" \
  http://localhost:3000/api/admin/products
```

---

### Создать новую запись
```
POST /api/admin/{resource}
```

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Body:**
```json
{
  "name": "New Product",
  "price": 99.99,
  "description": "Product description"
}
```

**Response (201 Created):**
```json
{
  "id": 5,
  "name": "New Product",
  "price": 99.99,
  ...
}
```

**Example:**
```bash
curl -X POST \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Product","price":99.99,"category_id":1}' \
  http://localhost:3000/api/admin/products
```

---

### Обновить запись
```
PUT /api/admin/{resource}/{id}
```

**Parameters:**
- `resource` - тип ресурса
- `id` - ID записи для обновления

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Updated Name",
  "price": 149.99
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Updated Name",
  "price": 149.99,
  ...
}
```

**Example:**
```bash
curl -X PUT \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Product","price":149.99}' \
  http://localhost:3000/api/admin/products/1
```

---

### Удалить запись
```
DELETE /api/admin/{resource}/{id}
```

**Parameters:**
- `resource` - тип ресурса
- `id` - ID записи для удаления

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true
}
```

**Example:**
```bash
curl -X DELETE \
  -H "Authorization: Bearer {token}" \
  http://localhost:3000/api/admin/products/1
```

---

## Resource Schemas

### Products
```json
{
  "id": 1,
  "name": "Product Name",
  "price": 99.99,
  "short_description": "Short desc",
  "long_description": "Long description",
  "sku": "SKU-001",
  "discount": 10,
  "images": "[\"img1.jpg\", \"img2.jpg\"]",
  "category_id": 1,
  "created_at": "2025-12-26T10:30:00Z",
  "updated_at": "2025-12-26T10:30:00Z"
}
```

### Categories
```json
{
  "id": 1,
  "name": "Category Name",
  "alias": "category-alias"
}
```

### Orders
```json
{
  "id": 1,
  "user_id": 1,
  "product_ids": "[1,2,3]",
  "status": "новый",
  "created_at": "2025-12-26T10:30:00Z"
}
```

### News
```json
{
  "id": 1,
  "title": "News Title",
  "description": "News description",
  "image": "news.jpg",
  "created_at": "2025-12-26T10:30:00Z"
}
```

### Banners
```json
{
  "id": 1,
  "product_id": 1,
  "image": "banner.jpg",
  "position": 0,
  "created_at": "2025-12-26T10:30:00Z"
}
```

### Users
```json
{
  "id": 1,
  "email": "user@example.com",
  "role": "admin",
  "name": "User Name",
  "phone": "+1234567890",
  "delivery_address": "Street 123",
  "created_at": "2025-12-26T10:30:00Z",
  "updated_at": "2025-12-26T10:30:00Z"
}
```

---

## Frontend Integration

### ResourceTable.vue Component

Компонент автоматически работает со следующими API вызовами:

**1. Получить все записи:**
```javascript
// GET /api/admin/products
const data = await $fetch('/api/admin/products', {
  method: 'GET',
  headers: { Authorization: `Bearer ${token}` }
})
```

**2. Создать запись:**
```javascript
// POST /api/admin/products
await $fetch('/api/admin/products', {
  method: 'POST',
  headers: { Authorization: `Bearer ${token}` },
  body: { name: 'New Product', price: 99.99, ... }
})
```

**3. Обновить запись:**
```javascript
// PUT /api/admin/products/1
await $fetch('/api/admin/products/1', {
  method: 'PUT',
  headers: { Authorization: `Bearer ${token}` },
  body: { name: 'Updated', price: 149.99, ... }
})
```

**4. Удалить запись:**
```javascript
// DELETE /api/admin/products/1
await $fetch('/api/admin/products/1', {
  method: 'DELETE',
  headers: { Authorization: `Bearer ${token}` }
})
```

---

## Error Handling

### Invalid Resource Type
```
Status: 400
{
  "error": "Unknown resource"
}
```

### Invalid ID
```
Status: 400
{
  "error": "Invalid ID"
}
```

### Database Error
```
Status: 500
{
  "error": "Database error message"
}
```

### Unauthorized
```
Status: 401
{
  "error": "Unauthorized"
}
```

### Forbidden (Not Admin)
```
Status: 403
{
  "error": "Forbidden"
}
```

---

## Testing with Insomnia

### 1. Create Admin Token
POST `/api/auth/login`
```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

### 2. Get Products
GET `/api/admin/products`
Header: `Authorization: Bearer {token}`

### 3. Create Product
POST `/api/admin/products`
```json
{
  "name": "New Product",
  "price": 99.99,
  "short_description": "A great product",
  "long_description": "Detailed description",
  "sku": "SKU-123",
  "discount": 0,
  "images": "[]",
  "category_id": 1
}
```

### 4. Update Product
PUT `/api/admin/products/1`
```json
{
  "name": "Updated Product",
  "price": 149.99
}
```

### 5. Delete Product
DELETE `/api/admin/products/1`

---

## Type Conversion

Backend автоматически конвертирует типы данных:

| Frontend | Backend | Conversion |
|----------|---------|-----------|
| `string` | `string` | As is |
| `number` | `int` | `int(value)` |
| `number` | `float64` | As is |
| `boolean` | `string` | `"true"` or `"false"` |
| `object` | `string` | JSON string |

---

## Features

✅ Generic CRUD operations for all resources
✅ Automatic type conversion
✅ Error handling and validation
✅ Authentication and authorization checks
✅ Database transaction support
✅ Timestamps (created_at, updated_at)
✅ Support for arrays (product_ids, images)

---

## Notes

- Все временные метки в формате RFC3339
- ID генерируются автоматически базой данных (AUTOINCREMENT)
- Пароль пользователя не возвращается в ответе
- Пароли хешируются с bcrypt
- Для удаления пользователя требуется особые права
