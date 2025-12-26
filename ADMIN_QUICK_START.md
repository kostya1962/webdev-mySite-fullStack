# Admin Panel API - Quick Start Guide

## Что было реализовано

Backend теперь полностью поддерживает админ-панель frontend'а. Добавлены CRUD операции для всех ресурсов.

## Новые Endpoints

### Структура маршрутов (все требуют авторизацию и роль admin):

```
GET    /api/admin/{resource}         - Получить все записи
POST   /api/admin/{resource}         - Создать новую запись
PUT    /api/admin/{resource}/{id}    - Обновить запись
DELETE /api/admin/{resource}/{id}    - Удалить запись
```

### Поддерживаемые ресурсы:
- `products` - Товары
- `categories` - Категории
- `orders` - Заказы
- `news` - Новости
- `banners` - Баннеры
- `users` - Пользователи

## Тестирование

### 1. Запуск backend сервера

```bash
cd backend
go run main.go
```

Сервер запустится на `http://localhost:3000`

### 2. Получение admin токена

Предварительно убедитесь, что в БД есть пользователь с ролью `admin`:

```bash
# Логин как админ
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'
```

Из ответа скопируйте `token`.

### 3. Получить список товаров

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:3000/api/admin/products
```

### 4. Создать новый товар

```bash
curl -X POST http://localhost:3000/api/admin/products \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Product",
    "price": 99.99,
    "short_description": "Short desc",
    "long_description": "Long description",
    "sku": "SKU-001",
    "discount": 10,
    "images": "[]",
    "category_id": 1
  }'
```

### 5. Обновить товар

```bash
curl -X PUT http://localhost:3000/api/admin/products/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product",
    "price": 149.99
  }'
```

### 6. Удалить товар

```bash
curl -X DELETE http://localhost:3000/api/admin/products/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Frontend Integration

Frontend ResourceTable.vue компонент уже настроен для работы с этими endpoints.

В админ-панели (`/admin`) вы можете:
1. ✅ Выбрать тип ресурса (Товары, Категории, Заказы и т.д.)
2. ✅ Просмотреть все записи в таблице
3. ✅ Нажать "Создать" для добавления новой записи
4. ✅ Нажать "Ред." для редактирования записи
5. ✅ Нажать "Удал." для удаления записи

## Важно

### Требования для работы:

1. **Авторизация**: Нужен валидный JWT токен в заголовке `Authorization: Bearer {token}`
2. **Роль admin**: Пользователь должен иметь `role = 'admin'` в таблице `users`
3. **Валидные данные**: Отправляемые данные должны соответствовать схеме таблицы

### Проверка авторизации:

Если вы получили ошибку `401 Unauthorized`:
- Проверьте, что токен передан в заголовке
- Проверьте, что токен не истек

Если вы получили ошибку `403 Forbidden`:
- Проверьте, что пользователь имеет роль `admin` в БД

### Обновление администратора

Если нужно сделать пользователя администратором:

```bash
# Через SQLite CLI
sqlite3 app.db
UPDATE users SET role = 'admin' WHERE email = 'user@example.com';
```

## Примеры использования

### JavaScript/Nuxt (как в ResourceTable.vue)

```typescript
// Получить записи
const data = await $fetch('/api/admin/products', {
  method: 'GET',
  headers: { Authorization: `Bearer ${auth.token}` }
})

// Создать
await $fetch('/api/admin/products', {
  method: 'POST',
  headers: { Authorization: `Bearer ${auth.token}` },
  body: { name: 'Product', price: 99 }
})

// Обновить
await $fetch('/api/admin/products/1', {
  method: 'PUT',
  headers: { Authorization: `Bearer ${auth.token}` },
  body: { name: 'Updated', price: 149 }
})

// Удалить
await $fetch('/api/admin/products/1', {
  method: 'DELETE',
  headers: { Authorization: `Bearer ${auth.token}` }
})
```

## Тестирование в Insomnia

1. Импортируйте `backend/insomnia-api.json` в Insomnia
2. Создайте новую папку "Admin API"
3. Используйте примеры выше для создания запросов
4. Установите переменную окружения `token` из响应 логина

## Файлы, которые были изменены

- **backend/handlers/admin.go** - Полная реализация CRUD handlers
- **backend/main.go** - Добавлены новые маршруты
- **ADMIN_API_DOCUMENTATION.md** - Подробная документация

## Заметки

- Все ошибки обрабатываются и возвращают соответствующие HTTP коды
- Типы данных автоматически конвертируются
- Поддерживается JSON для сложных типов (images, product_ids)
- Timestamps автоматически добавляются/обновляются

## Поддержка

Если возникнут проблемы:

1. Проверьте логи сервера (`go run main.go`)
2. Убедитесь, что база данных инициализирована
3. Проверьте, что пользователь имеет роль admin
4. Попробуйте тестовый запрос через curl
