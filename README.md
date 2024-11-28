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
```
APP_PORT=8080
DATABASE_URL=postgres://user:password@db:5432/authdb?sslmode=disable
SECRET_KEY=<ваш_секретный_ключ>
EMAIL_SERVICE_API_KEY=<можно_не_использовать_тк_используется_мок>
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
