# authenticationProject

Данный проект реализует часть сервиса аутентификации.

## Требования
- Генерация **Access** и **Refresh** токенов для пользователей.
- Использование **JWT** с алгоритмом **SHA512** для Access токенов.
- Хранение Refresh токенов в виде bcrypt-хэшей в базе данных PostgreSQL.
- Проверка изменения IP-адреса при обновлении токенов и отправка предупреждения на почту (моковые данные).
- Полностью контейнеризован для простого развертывания и тестирования через Docker.

## Технологии
- Go
- PostgreSQL
- JWT
- Docker

## Установка и запуск
### Клонирование репозитория
```bash
git clone https://github.com/molodoymaxim/authenticationProject.git
cd authenticationProject
```
### Настройка переменных окружения .env
Пример файла .env лежит внутри проекта. Необходимо подставить свои значения.
```
# База данных
POSTGRES_USER=user
POSTGRES_PASSWORD=1234
POSTGRES_DB=authdb
POSTGRES_HOST=localhost(для запуска через go run ./cmd/main.go) или db (для запуска в докере)
POSTGRES_PORT=5432
DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable

# Приложение
APP_PORT=8080
SECRET_KEY=ebe95a4ccc28f066bc9056334d07e0ef9e738e6fabc9c0a8d9d3530e515888ee
EMAIL_SERVICE_API_KEY=email_service_api_key (можно не изменять, используется моковые данные)
LOG_LEVEL=debug
```
### Запуск через Docker
1. Выполнить команду для сборки и запуска.
```bash
docker-compose up -d --build
```
2. Проверить работоспособность контейнера.
```bash
docker-compose ps
```
3. Проверить логи.
```bash
docker-compose logs -f auth-service
```
### Запуск тестов
Для запуска тестов используйте команду:
```go
go test ./...
```
## Описание API
### 1. Генерация токенов
**Эндпоинт**: `POST /auth/token`\
**URL**: `http://localhost:8080/auth/token` \
**Заголовки**: `Content-Type: application/json` \
**Тело**: `{ "user_id": "test-user-id" }`
### 2. Обновление токенов
**Эндпоинт**: `POST /auth/refresh` \
**URL**: `http://localhost:8080/auth/refresh` \
**Заголовки**: `Content-Type: application/json`,  `Authorization: Bearer <access_token>` \
**Тело**: `{ "refresh_token": "<refresh_token>" }`

## Дополнительно
Был добавлен swagger. Для просмотра необходимо запустить приложение и перейти по адресу:
```bash
 http://localhost:8080/swagger/index.html
```
Пример работы:
![image](https://github.com/user-attachments/assets/581b17e2-9ea1-4705-919c-57e22afca03f)

