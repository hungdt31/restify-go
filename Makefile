.PHONY: dev run build clean

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
