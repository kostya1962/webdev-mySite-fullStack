# Хостинг проекта (Docker Compose)

Краткая инструкция по управлению локальным хостингом проекта (запуск/остановка/пересборка). Предполагается, что вы используете `docker compose` из корня репозитория.

--

**Важно:** заменяйте `<HOST_IP>` на IP-адрес хоста (например `192.168.1.101`) при проверках в браузере или curl.

**Файлы:** основные точки внимания: [docker-compose.yml](docker-compose.yml), [frontend](frontend/), [backend](backend/).

## 1) Как отключать хостинг

- Остановить все контейнеры (без удаления):
```bash
docker compose stop
```
- Полностью остановить и удалить контейнеры/сеть (сохранённые тома остаются):
```bash
docker compose down
```
- Полностью удалить контейнеры и тома (ОСТОРОЖНО — удалит БД и прочие персистентные данные):
```bash
docker compose down -v
```

## 2) Как включать хостинг

- Запустить сервисы (использует существующие образы):
```bash
docker compose up -d
```
- Собрать заново образы и запустить (рекомендуется после изменений в Dockerfile / зависимостях):
```bash
docker compose up -d --build
```

## 3) Как вносить изменения и обновлять хостинг

Общий подход:

- Внести изменения локально в `frontend/` или `backend/` (редактируйте код в рабочем каталоге).
- Если изменяли только фронтенд (Vue/Nuxt): пересоберите только фронтенд и перезапустите сервис:
```bash
docker compose build frontend
docker compose up -d --no-deps --build frontend
```
- Если изменяли только бэкенд (Go):
```bash
docker compose build backend
docker compose up -d --no-deps --build backend
```
- Если вносились изменения в инфраструктуру (Dockerfile, nginx конфиг, compose): пересоберите весь стек:
```bash
docker compose down
docker compose up -d --build
```

Замечания:
- Флаг `--no-deps` перезапускает только указанный сервис без зависимых контейнеров.
- Если вы меняете только шаблоны/статические файлы внутри `frontend/dist`, можно копировать их в runtime nginx образ при ручной отладке, но рекомендуем пересборку.

## 4) Другие полезные команды

- Просмотр логов для всего стека (в реальном времени):
```bash
docker compose logs -f --tail=200
```
- Просмотр логов конкретного сервиса (например, backend):
```bash
docker compose logs -f backend
```
- Перезапустить один сервис:
```bash
docker compose restart backend
```
- Выполнить команду внутри запущенного контейнера (полезно для отладки):
```bash
docker compose exec backend sh
# или для frontend nginx
docker compose exec frontend sh
```
- Проверка доступности сайта и API с хоста:
```bash
curl -I http://<HOST_IP>/
curl -s http://<HOST_IP>/api/products | jq .
```
- Найти, какой процесс занимает порт 80 (если `docker compose` не может забиндить порт):
```bash
sudo ss -ltnp | grep ':80'
# или
sudo lsof -i :80
```
- Скопировать файл БД (SQLite) из контейнера на хост для бэкапа:
```bash
# узнайте имя контейнера: docker compose ps
docker compose exec backend sh -c "cp /app/app.db /tmp/app.db && chmod 644 /tmp/app.db"
docker cp $(docker compose ps -q backend):/tmp/app.db ./backup-app.db
```

## Отладка гидрационных ошибок (специфично для фронтенда)

- Если в браузере видите ошибки вида "Hydration completed but contains mismatches", попробуйте:
  - Обернуть динамически загружаемые части в `ClientOnly` (см. `frontend/pages/index.vue`).
  - Очистить кеш браузера и сделать Ctrl+F5.
  - Проверить, что `NUXT_PUBLIC_APIURL` и `NUXT_PUBLIC_IMAGEURL` в runtime указывают корректно (в собранном `dist` это обычно видно в `index.html` как конфиг).

## Рекомендации и предостережения

- Делайте бэкап базы перед `down -v` или перед экспериментами с таблицами.
- На продакшене не запускайте контейнеры, которые слушают порт 80, если на хосте уже есть веб-сервер — либо смените порт, либо останавливайте системный сервис.

Если нужно, могу добавить автоматические скрипты (`scripts/`) для однотипных операций или пример `make`-целей.

---

Если хотите, вставлю эти инструкции в другой файл (например, [frontend/README.md](frontend/README.md)) или дополню шагами для CI/CD.
