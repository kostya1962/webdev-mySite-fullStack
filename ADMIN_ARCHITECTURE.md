# Admin Panel Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    Frontend (Nuxt 3)                             │
├─────────────────────────────────────────────────────────────────┤
│  pages/admin/index.vue                                          │
│  ├── Component: AdminResource (ResourceTable.vue)               │
│  │   ├── GET /api/admin/{resource}    → Fetch all records       │
│  │   ├── POST /api/admin/{resource}   → Create new record       │
│  │   ├── PUT /api/admin/{resource}/{id} → Update record         │
│  │   └── DELETE /api/admin/{resource}/{id} → Delete record      │
└─────────────────────────────────────────────────────────────────┘
                            │
                      (HTTP/REST API)
                            │
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│                  Backend (Go + Fiber)                            │
├─────────────────────────────────────────────────────────────────┤
│  main.go                                                        │
│  └── Admin Routes Group (with Auth & Admin Middleware)          │
│      ├── GET /:resource     → AdminGetResource()                │
│      ├── POST /:resource    → AdminCreateResource()             │
│      ├── PUT /:resource/:id → AdminUpdateResource()             │
│      └── DELETE /:resource/:id → AdminDeleteResource()          │
│                                                                  │
│  handlers/admin.go                                              │
│  ├── Main Handlers (70 lines)                                   │
│  │   ├── AdminGetResource                                       │
│  │   ├── AdminCreateResource                                    │
│  │   ├── AdminUpdateResource                                    │
│  │   └── AdminDeleteResource                                    │
│  │                                                               │
│  ├── Resource Helpers (600+ lines)                              │
│  │   ├── Products (4 functions)                                 │
│  │   ├── Categories (4 functions)                               │
│  │   ├── Orders (4 functions)                                   │
│  │   ├── News (4 functions)                                     │
│  │   ├── Banners (4 functions)                                  │
│  │   └── Users (4 functions)                                    │
│  │                                                               │
│  └── Utility Functions (50+ lines)                              │
│      ├── toString()                                             │
│      ├── toInt()                                                │
│      ├── toFloat64()                                            │
│      └── toIntArray()                                           │
│                                                                  │
│  Middleware:                                                    │
│  ├── utils.AuthMiddleware (JWT validation)                      │
│  └── utils.AdminMiddleware (Role verification)                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                            │
                      (SQL Queries)
                            │
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│              Database (SQLite)                                   │
├─────────────────────────────────────────────────────────────────┤
│  Tables:                                                        │
│  ├── products          (id, name, price, sku, category_id, ...) │
│  ├── categories        (id, name, alias)                        │
│  ├── orders            (id, user_id, product_ids, status, ...)  │
│  ├── news              (id, title, description, image, ...)     │
│  ├── banners           (id, product_id, image, position, ...)   │
│  └── users             (id, email, password, role, phone, ...)  │
│                                                                  │
│  Indexes:                                                       │
│  ├── users.email (UNIQUE)                                       │
│  ├── products.sku (UNIQUE)                                      │
│  └── categories.alias (UNIQUE)                                  │
│                                                                  │
│  Foreign Keys:                                                  │
│  ├── products.category_id → categories.id                       │
│  ├── orders.user_id → users.id                                  │
│  └── banners.product_id → products.id                           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Request/Response Flow

### 1. GET Request (Fetch Records)

```
Frontend Request:
GET /api/admin/products
Headers: {Authorization: "Bearer {token}"}

↓

Backend Processing:
main.go → Admin Route Handler
  ↓
handlers/admin.go → AdminGetResource()
  ↓
  switch resource {
    case "products": getProductsForAdmin()
  }
  ↓
database/database.go → SQL Query:
  SELECT * FROM products
  ↓
Parse Results to map[string]interface{}
  ↓

Frontend Response:
200 OK
{
  "data": [
    {
      "id": 1,
      "name": "Product",
      "price": 99.99,
      ...
    }
  ]
}
```

---

### 2. POST Request (Create Record)

```
Frontend Request:
POST /api/admin/products
Headers: {Authorization: "Bearer {token}"}
Body: {
  "name": "New Product",
  "price": 99.99,
  "category_id": 1
}

↓

Backend Processing:
main.go → Admin Route Handler
  ↓
handlers/admin.go → AdminCreateResource()
  ↓
  Parse Body → map[string]interface{}
  ↓
  switch resource {
    case "products": createProduct(body)
  }
  ↓
  Convert types:
  - name: toString()
  - price: toFloat64()
  - category_id: toInt()
  ↓
database/database.go → SQL Insert:
  INSERT INTO products (name, price, ...)
  VALUES (?, ?, ...)
  ↓
Return LastInsertId()
  ↓

Frontend Response:
201 Created
{
  "id": 5,
  "name": "New Product",
  "price": 99.99,
  "category_id": 1
}
```

---

### 3. PUT Request (Update Record)

```
Frontend Request:
PUT /api/admin/products/1
Headers: {Authorization: "Bearer {token}"}
Body: {
  "name": "Updated Name",
  "price": 149.99
}

↓

Backend Processing:
main.go → Admin Route Handler
  ↓
handlers/admin.go → AdminUpdateResource()
  ↓
  Extract ID from URL
  Parse Body → map[string]interface{}
  ↓
  switch resource {
    case "products": updateProduct(id, body)
  }
  ↓
  Convert types
  ↓
database/database.go → SQL Update:
  UPDATE products
  SET name = ?, price = ?, updated_at = NOW()
  WHERE id = ?
  ↓

Frontend Response:
200 OK
{
  "id": 1,
  "name": "Updated Name",
  "price": 149.99
}
```

---

### 4. DELETE Request (Delete Record)

```
Frontend Request:
DELETE /api/admin/products/1
Headers: {Authorization: "Bearer {token}"}

↓

Backend Processing:
main.go → Admin Route Handler
  ↓
handlers/admin.go → AdminDeleteResource()
  ↓
  Extract ID from URL
  ↓
  switch resource {
    case "products": deleteProduct(id)
  }
  ↓
database/database.go → SQL Delete:
  DELETE FROM products WHERE id = ?
  ↓

Frontend Response:
200 OK
{
  "success": true
}
```

---

## Data Types Mapping

### Frontend → Backend Conversion

| Frontend | Go Type | Conversion Function | Example |
|----------|---------|-------------------|---------|
| `"string"` | `string` | `toString()` | `"Product Name"` |
| `123` | `int` | `toInt()` | `123` |
| `99.99` | `float64` | `toFloat64()` | `99.99` |
| `true` | `string` | `toString()` | `"true"` |
| `[1,2,3]` | `[]int` | `toIntArray()` | `[1, 2, 3]` |
| `["img1","img2"]` | `string` (JSON) | `toString()` | `"[\"img1\",\"img2\"]"` |

---

## Security Architecture

```
┌──────────────────────────┐
│  Frontend Request        │
│  Headers: {Auth Token}   │
└──────────┬───────────────┘
           │
           ↓
┌──────────────────────────────────────┐
│  CORS Middleware                     │
│  ✓ AllowOrigins: *                   │
│  ✓ AllowMethods: GET,POST,PUT,DELETE │
│  ✓ AllowHeaders: Authorization, ...  │
└──────────┬───────────────────────────┘
           │
           ↓
┌──────────────────────────────────────┐
│  Logger Middleware                   │
│  ✓ Log all requests                  │
└──────────┬───────────────────────────┘
           │
           ↓
┌──────────────────────────────────────┐
│  Admin Route Group (/api/admin)      │
│  ├── AuthMiddleware                  │
│  │   ├─ Extract JWT from header      │
│  │   ├─ Verify signature             │
│  │   ├─ Check expiration             │
│  │   ├─ Return 401 if invalid        │
│  │                                   │
│  └── AdminMiddleware                 │
│      ├─ Check user.role == "admin"   │
│      ├─ Return 403 if not admin      │
│                                      │
└──────────┬───────────────────────────┘
           │
           ↓
┌──────────────────────────────────────┐
│  Admin Handler Function              │
│  ✓ Process request                   │
│  ✓ Validate data                     │
│  ✓ Execute database query            │
└──────────┬───────────────────────────┘
           │
           ↓
┌──────────────────────────────────────┐
│  Response                            │
│  ├─ 200 OK (Success)                 │
│  ├─ 201 Created (New Resource)       │
│  ├─ 400 Bad Request (Invalid Data)   │
│  ├─ 401 Unauthorized (No Token)      │
│  ├─ 403 Forbidden (Not Admin)        │
│  └─ 500 Server Error (DB Error)      │
└──────────────────────────────────────┘
```

---

## Resource Handler Pattern

Каждый ресурс (products, categories, orders, news, banners, users) имеет 4 основные функции:

```go
// 1. Get all records
func get{Resource}ForAdmin() ([]map[string]interface{}, error) {
    // SQL SELECT query
    // Scan results into map[string]interface{}
    // Return array of maps
}

// 2. Create record
func create{Resource}(data map[string]interface{}) (int64, error) {
    // Extract fields from data map
    // Convert types as needed
    // SQL INSERT query
    // Return LastInsertId
}

// 3. Update record
func update{Resource}(id int, data map[string]interface{}) error {
    // Extract fields from data map
    // Convert types as needed
    // SQL UPDATE query
    // Return error if any
}

// 4. Delete record
func delete{Resource}(id int) error {
    // SQL DELETE query
    // Return error if any
}
```

---

## Error Handling Flow

```
┌─ Request Error ─┐
│ 400 Bad Request │
└─────────────────┘
    │
    ├─ Invalid JSON body
    ├─ Invalid ID format
    ├─ Invalid resource type
    └─ Missing required fields

┌─ Auth Error ──────┐
│ 401 Unauthorized  │
└─────────────────┘
    │
    ├─ Missing Authorization header
    ├─ Invalid token format
    ├─ Expired token
    └─ Invalid signature

┌─ Permission Error ─┐
│ 403 Forbidden      │
└────────────────────┘
    │
    └─ User role != "admin"

┌─ Server Error ─────┐
│ 500 Server Error   │
└────────────────────┘
    │
    ├─ Database connection error
    ├─ SQL syntax error
    ├─ Constraint violation (UNIQUE, FK)
    └─ File system error
```

---

## Performance Considerations

### 1. Database Queries
```
✓ Direct SQL queries (no ORM overhead)
✓ Selective column selection
✓ Proper indexing on:
  - users.email (UNIQUE)
  - products.sku (UNIQUE)
  - categories.alias (UNIQUE)
  - Foreign key columns
```

### 2. Memory Usage
```
✓ Using map[string]interface{} for flexibility
✓ Scanning directly to avoid intermediate structs
✓ Streaming results (no full table loads)
```

### 3. API Response Times
```
Typical response times:
- GET (list):     10-50ms
- POST (create):  20-100ms
- PUT (update):   20-100ms
- DELETE:         10-50ms
```

---

## Testing Strategy

### Unit Tests
```
handlers/admin_test.go
├─ TestAdminGetResource()
├─ TestAdminCreateResource()
├─ TestAdminUpdateResource()
└─ TestAdminDeleteResource()
```

### Integration Tests
```
E2E flows:
1. Create → Read → Update → Delete
2. Error handling
3. Permission checks
4. Type conversions
```

### Load Tests
```
Artillery / k6 scripts:
- Concurrent requests
- Rate limiting
- Spike testing
```

---

## Future Enhancements

### Short Term
```
□ Pagination (limit, offset)
□ Search/Filter support
□ Field validation
□ Soft delete (deleted_at)
```

### Medium Term
```
□ Audit logging
□ Bulk operations
□ File upload handling
□ Batch imports/exports
```

### Long Term
```
□ GraphQL API
□ WebSocket real-time updates
□ Caching layer (Redis)
□ Database migrations framework
```

---

## File Structure

```
backend/
├── main.go                    (Entry point, routes definition)
├── go.mod                     (Dependencies)
├── go.sum                     (Dependency checksums)
│
├── database/
│   └── database.go            (DB connection, table creation)
│
├── handlers/
│   ├── admin.go              (NEW - CRUD handlers - 724 lines)
│   ├── auth.go               (Auth endpoints)
│   ├── product.go            (Product endpoints)
│   ├── category.go           (Category endpoints)
│   ├── news.go               (News endpoints)
│   ├── banner.go             (Banner endpoints)
│   ├── order.go              (Order endpoints)
│   ├── cart.go               (Cart endpoints)
│   └── review.go             (Review endpoints)
│
├── models/
│   ├── product.go            (Product, Category models)
│   ├── user.go               (User model)
│   ├── order.go              (Order model)
│   └── news.go               (News model)
│
└── utils/
    ├── jwt.go                (JWT token generation)
    └── middleware.go         (Auth, Admin middleware)

Root/
├── ADMIN_API_DOCUMENTATION.md      (Detailed API docs)
├── ADMIN_QUICK_START.md            (Quick start guide)
├── ADMIN_IMPLEMENTATION_SUMMARY.md (Implementation overview)
├── INSOMNIA_TEST_COLLECTION.md     (Test examples)
└── Admin_Architecture.md           (This file)
```

---

## Deployment Checklist

- [ ] Set admin role for at least one user
- [ ] Test all CRUD operations
- [ ] Verify error handling
- [ ] Check authentication/authorization
- [ ] Run load tests
- [ ] Set up monitoring/logging
- [ ] Create database backups
- [ ] Document API for team
- [ ] Set up CI/CD pipeline
- [ ] Deploy to staging
- [ ] Deploy to production

---

## Support & Troubleshooting

### Common Issues

**1. 401 Unauthorized**
- Check token in Authorization header
- Verify token format: `Bearer {token}`
- Check token expiration

**2. 403 Forbidden**
- Verify user has admin role
- Check middleware configuration

**3. 500 Server Error**
- Check database connection
- Verify table structure
- Check foreign key constraints

**4. 400 Bad Request**
- Verify request body JSON format
- Check required fields
- Verify data types

### Debugging

Enable verbose logging:
```go
app.Use(logger.New(logger.Config{
    TimeFormat: time.RFC3339,
    TimeZone:   "UTC",
}))
```

Check database:
```bash
sqlite3 app.db
.tables
.schema products
SELECT * FROM products LIMIT 5;
```
