# Admin Panel Deployment & Setup Guide

## Содержание
1. [Локальная разработка](#локальная-разработка)
2. [Тестирование](#тестирование)
3. [Продакшн развертывание](#продакшн-развертывание)
4. [Управление администраторами](#управление-администраторами)
5. [Мониторинг](#мониторинг)
6. [Решение проблем](#решение-проблем)

---

## Локальная разработка

### Предварительные требования
```bash
# Проверить Go установлен
go version  # должно быть >= 1.18

# Проверить Node.js установлен (для frontend)
node --version && npm --version

# Проверить SQLite установлен
sqlite3 --version
```

### Setup процесс

#### 1. Клонировать репозиторий
```bash
cd "simle shopper"
```

#### 2. Запустить backend
```bash
cd backend

# Установить зависимости (если нужно)
go mod download
go mod tidy

# Запустить сервер
go run main.go

# Или собрать бинарный файл
go build -o server main.go
./server
```

**Ожидаемый вывод:**
```
Server starting on :3000
Database connected successfully
```

#### 3. Запустить frontend (в другом терминале)
```bash
cd frontend

# Установить зависимости
npm install

# Запустить dev сервер
npm run dev

# Или собрать для продакшна
npm run build
```

#### 4. Проверить работоспособность

**Frontend:**
- Откройте `http://localhost:3000`
- Страницы должны загружаться корректно

**Backend:**
- Тестируйте API через curl или Insomnia
- ```bash
  curl http://localhost:3000/api/news
  ```

#### 5. Проверить админ-панель

1. Залогинитесь как обычный пользователь
2. Проверьте, что админ-панель недоступна (требует роль admin)
3. Сделайте пользователя администратором (смотри раздел ниже)
4. Перезагрузитесь
5. Админ-панель должна быть доступна на `/admin`

---

## Тестирование

### Unit тесты backend

```bash
cd backend

# Запустить все тесты
go test ./...

# С подробным выводом
go test -v ./...

# Только admin handlers
go test -v ./handlers -run Admin
```

### Integration тесты

```bash
# Создать тестовую базу данных
cp app.db app.test.db

# Запустить интеграционные тесты
go test -v -tags=integration ./...
```

### E2E тесты (Frontend)

```bash
cd frontend

# Запустить тесты Cypress
npm run test

# Или Playwright
npm run test:e2e
```

### API тесты через Insomnia

1. Откройте Insomnia
2. Импортируйте коллекцию из `INSOMNIA_TEST_COLLECTION.md`
3. Установите переменную окружения `token`
4. Запустите тесты для всех endpoints
5. Проверьте статусы ответов (200, 201, 400, 401, 403, 500)

### Load тесты

```bash
# Использовать Artillery.io
npm install -g artillery

# Создать файл artillery-load-test.yml
artillery run artillery-load-test.yml

# Или использовать k6
k6 run load-test.js
```

---

## Продакшн развертывание

### Подготовка

#### 1. Подготовить окружение
```bash
# Создать production database
sqlite3 app.prod.db

# Создать backup директорию
mkdir -p backups
```

#### 2. Конфигурировать переменные окружения

Создать файл `.env.production`:
```env
# Backend
BACKEND_PORT=3000
BACKEND_HOST=0.0.0.0
DATABASE_PATH=app.db
LOG_LEVEL=info

# Frontend
API_URL=https://your-domain.com/api

# Security
JWT_SECRET=your-super-secret-key-change-this
JWT_EXPIRATION=24h
```

#### 3. Собрать приложение

```bash
# Backend
cd backend
go build -o server main.go

# Frontend
cd ../frontend
npm run build
# Выходной файл в .nuxt/dist
```

### Docker развертывание

#### Dockerfile для Backend

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/app.db .
COPY --from=builder /app/backups ./backups

EXPOSE 3000

CMD ["./server"]
```

#### docker-compose.yml

```yaml
version: '3.8'

services:
  backend:
    build: ./backend
    ports:
      - "3000:3000"
    environment:
      - DATABASE_PATH=/root/app.db
    volumes:
      - ./backups:/root/backups
      - ./images:/root/images
    restart: unless-stopped

  frontend:
    build: ./frontend
    ports:
      - "3001:3000"
    environment:
      - API_URL=http://localhost:3000/api
    depends_on:
      - backend
    restart: unless-stopped
```

#### Запустить с Docker Compose

```bash
docker-compose up -d
```

### Nginx конфигурация

```nginx
upstream backend {
    server localhost:3000;
}

upstream frontend {
    server localhost:3001;
}

server {
    listen 80;
    server_name your-domain.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # Backend API
    location /api/ {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Frontend
    location / {
        proxy_pass http://frontend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # Статические файлы
    location /images/ {
        alias /root/images/;
    }
}
```

### Системд сервис (Linux)

```ini
# /etc/systemd/system/simle-shopper.service

[Unit]
Description=Simple Shopper Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/simle-shopper
ExecStart=/var/www/simle-shopper/backend/server
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

```bash
# Установить сервис
sudo systemctl daemon-reload
sudo systemctl enable simle-shopper
sudo systemctl start simle-shopper

# Проверить статус
sudo systemctl status simle-shopper
```

---

## Управление администраторами

### Создать нового администратора

#### Через API (требует существующего админа)

```bash
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}' \
  | jq -r '.token')

curl -X POST http://localhost:3000/api/admin/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newadmin@example.com",
    "role": "admin",
    "name": "New Admin"
  }'
```

#### Через SQLite (для первого администратора)

```bash
sqlite3 app.db

-- Просмотреть пользователей
SELECT id, email, role FROM users;

-- Сделать пользователя администратором
UPDATE users SET role = 'admin' WHERE email = 'user@example.com';

-- Или создать нового администратора
INSERT INTO users (email, password, role, created_at, updated_at)
VALUES ('admin@example.com', 'hashed_password', 'admin', datetime('now'), datetime('now'));

.exit
```

### Изменить роль пользователя

```bash
# Администратор
sqlite3 app.db
UPDATE users SET role = 'admin' WHERE id = 5;

# Обычный пользователь
UPDATE users SET role = 'user' WHERE id = 5;

.exit
```

### Сбросить пароль администратора

```bash
# Через админ-панель (обновить поле "password")
# Это должно содержать хешированный пароль!

# Лучше: использовать API
curl -X PUT http://localhost:3000/api/admin/users/1 \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"password": "new_hashed_password"}'
```

### Удалить администратора

```bash
# Через API
curl -X DELETE http://localhost:3000/api/admin/users/5 \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Или через SQLite
sqlite3 app.db
DELETE FROM users WHERE id = 5;
.exit
```

---

## Мониторинг

### Логирование

#### Backend логи

```bash
# Просмотреть логи в реальном времени
go run main.go 2>&1 | tee server.log

# Просмотреть логи из файла
tail -f server.log

# Фильтровать по уровню
grep "error" server.log
```

#### Структурированное логирование

```go
// Использовать логгер в коде
log.Println("Info message")
log.Fatal("Error message")
```

### Мониторинг производительности

#### Метрики базы данных

```bash
sqlite3 app.db

-- Размер базы данных
.dump | wc -c

-- Размер таблиц
SELECT name, sum(pgsize) as size
FROM dbstat('main')
GROUP BY name
ORDER BY size DESC;
```

#### Размер файлов

```bash
# Размер базы данных
du -sh app.db

# Размер backup'ов
du -sh backups/

# Размер изображений
du -sh images/
```

### Backup'ы

#### Автоматические backup'ы

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="./backups"
DB_FILE="app.db"

mkdir -p "$BACKUP_DIR"

# Создать backup
cp "$DB_FILE" "$BACKUP_DIR/backup-$(date +%Y%m%d-%H%M%S).db"

# Удалить старые backup'ы (старше 30 дней)
find "$BACKUP_DIR" -name "backup-*.db" -mtime +30 -delete

echo "Backup completed: $BACKUP_DIR"
```

```bash
# Добавить в crontab
0 2 * * * cd /path/to/project && bash backup.sh
```

#### Восстановление из backup'а

```bash
# Просмотреть доступные backup'ы
ls -la backups/

# Восстановить из backup'а
cp backups/backup-20251226-020000.db app.db

# Перезагрузить сервер
sudo systemctl restart simle-shopper
```

### Мониторинг здоровья

```bash
# Проверить, что сервер работает
curl -s http://localhost:3000/api/news | jq . && echo "✅ Backend OK" || echo "❌ Backend ERROR"

# Проверить доступ к БД
sqlite3 app.db "SELECT COUNT(*) FROM users;"

# Проверить свободное место на диске
df -h | grep "/" | head -1
```

---

## Решение проблем

### 401 Unauthorized

**Проблема:** API возвращает 401

**Решения:**
1. Проверить, что токен передан в заголовке
   ```bash
   curl -H "Authorization: Bearer $TOKEN" ...
   ```

2. Проверить, что токен не истек
   ```bash
   # Получить новый токен
   curl -X POST http://localhost:3000/api/auth/login ...
   ```

3. Проверить JWT_SECRET совпадает
   ```bash
   grep JWT_SECRET .env
   ```

### 403 Forbidden

**Проблема:** API возвращает 403

**Решение:** Пользователь не администратор

```bash
# Проверить роль
sqlite3 app.db "SELECT email, role FROM users WHERE email='user@example.com';"

# Сделать администратором
sqlite3 app.db "UPDATE users SET role = 'admin' WHERE email='user@example.com';"
```

### 500 Server Error

**Проблема:** API возвращает 500

**Решения:**
1. Проверить логи сервера
   ```bash
   tail -f server.log
   ```

2. Проверить подключение к БД
   ```bash
   sqlite3 app.db ".tables"
   ```

3. Проверить диск полный
   ```bash
   df -h /
   ```

4. Перезагрузить сервер
   ```bash
   sudo systemctl restart simle-shopper
   ```

### База данных повреждена

**Проблема:** "database disk image malformed"

**Решение:**
```bash
# Восстановить из backup'а
cp backups/backup-latest.db app.db

# Или пересоздать
rm app.db
go run main.go  # пересоздаст и заполнит тестовыми данными
```

### Высокое использование памяти

**Проблема:** Процесс использует много памяти

**Решения:**
1. Проверить, нет ли утечек памяти в коде
2. Ограничить процесс
   ```bash
   ulimit -m 512000  # 512MB limit
   ```
3. Перезагрузить сервер
   ```bash
   sudo systemctl restart simle-shopper
   ```

### Медленные запросы

**Проблема:** API работает медленно

**Решения:**
1. Проверить индексы БД
   ```bash
   sqlite3 app.db ".indices"
   ```

2. Проверить количество записей
   ```bash
   sqlite3 app.db "SELECT COUNT(*) FROM products;"
   ```

3. Добавить индексы если нужно
   ```sql
   CREATE INDEX idx_products_category ON products(category_id);
   CREATE INDEX idx_orders_user ON orders(user_id);
   ```

4. Оптимизировать запросы в коде

---

## Чек-лист развертывания

### Pre-deployment
- [ ] Все тесты прошли успешно
- [ ] Код пересмотрен
- [ ] Конфигурация проверена
- [ ] Backup'ы готовы
- [ ] SSL сертификат установлен

### Deployment
- [ ] Собрано приложение
- [ ] Копирование файлов на сервер
- [ ] Перезагрузка сервиса
- [ ] Проверка здоровья приложения

### Post-deployment
- [ ] API endpoints работают
- [ ] Админ-панель доступна
- [ ] Логи чистые (нет ошибок)
- [ ] База данных целая
- [ ] Резервные копии актуальны

---

## Поддержка и связь

### Где найти помощь
- 📖 Документация: `ADMIN_API_DOCUMENTATION.md`
- 🚀 Quick Start: `ADMIN_QUICK_START.md`
- 🏗️ Архитектура: `ADMIN_ARCHITECTURE.md`
- 🧪 Тесты: `INSOMNIA_TEST_COLLECTION.md`

### Сообщить об ошибке
1. Проверьте логи сервера
2. Проверьте консоль браузера (F12)
3. Создайте Issue с деталями
4. Приложите логи и версию Go/Node

---

**Последнее обновление:** 26 декабря 2025  
**Версия:** 1.0  
**Статус:** Production Ready ✅
