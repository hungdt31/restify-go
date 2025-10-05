.PHONY: dev run build clean staging-up staging-down prod-build prod-up prod-down docker-dev-up docker-dev-down

# Chạy với hot reload
dev:
	@echo 🔥 Starting development server with hot reload...
	@go run github.com/air-verse/air@latest

# Chạy trực tiếp không hot reload
run:
	@echo 🚀 Starting server...
	@go run main.go

# Build ứng dụng
build:
	@echo 🔨 Building application...
	@go build -o tmp/main.exe .

# Dọn dẹp
clean:
	@echo 🧹 Cleaning...
	@if exist tmp rmdir /s /q tmp
	@if exist build-errors.log del build-errors.log

# Docker compose - development
docker-dev-up:
	@echo 🔥 Starting DEVELOPMENT environment (Docker + hot reload)...
	@if not exist .env.development copy .env.development.example .env.development >nul 2>&1
	docker compose -f docker-compose.dev.yml --env-file .env.development up -d --build

docker-dev-down:
	@echo 🛑 Stopping DEVELOPMENT environment...
	docker compose -f docker-compose.dev.yml --env-file .env.development down -v

# Docker compose - staging
staging-up:
	@echo 🚀 Starting STAGING environment...
	@if not exist .env.staging copy .env.staging.example .env.staging >nul 2>&1
	docker compose -f docker-compose.staging.yml --env-file .env.staging up -d --build

staging-down:
	@echo 🛑 Stopping STAGING environment...
	docker compose -f docker-compose.staging.yml --env-file .env.staging down -v

# Docker compose - production
prod-build:
	@echo 🔨 Building PRODUCTION image...
	docker build -t gofirstapp:prod .

prod-up:
	@echo 🚀 Starting PRODUCTION environment...
	@if not exist .env.production echo APP_ENV=production>nul
	docker compose -f docker-compose.prod.yml --env-file .env.production up -d --build

prod-down:
	@echo 🛑 Stopping PRODUCTION environment...
	docker compose -f docker-compose.prod.yml --env-file .env.production down
