#!/bin/bash

# ===============================
# 🧠 Go Project Management Script
# ===============================

# Load environment variables from .env if it exists
if [ -f .env ]; then
  echo "📥 Loading environment variables from .env..."
  export $(grep -v '^#' .env | xargs)
else
  echo "⚠️  No .env file found. Using default DB configuration."
fi

# Compose DB_URL from individual .env variables (with fallback)
DB_USER=${DB_USER:-myuser}
DB_PASS=${DB_PASS:-mypass}
DB_HOST=${DB_HOST:-127.0.0.1}
DB_PORT=${DB_PORT:-3306}
DB_NAME=${DB_NAME:-mydb}

DB_URL="mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

# Function: main menu
show_menu() {
  echo ""
  echo "==============================="
  echo "🚀 Go Project Manager"
  echo "==============================="
  echo "1. Run program (go run main.go)"
  echo "2. Run with Air (auto reload)"
  echo "3. Run migration (up/down/goto)"
  echo "4. Create new migration file"
  echo "5. Exit"
  echo "==============================="
  read -p "👉 Choose an option: " choice
}

# Option 1
run_program() {
  echo "▶ Running main.go ..."
  go run main.go
}

# Option 2
run_with_air() {
  echo "🔥 Running with Air (auto reload)..."
  go run github.com/air-verse/air@latest
}

# Option 3
run_migrate() {
  echo ""
  echo "📦 Migration options:"
  echo "1. Up"
  echo "2. Down"
  echo "3. Goto"
  read -p "👉 Choose: " mopt

  case $mopt in
    1)
      echo "🚀 Running migrations up ..."
      go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate \
        -database "$DB_URL" -path database/migrations up
      ;;
    2)
      read -p "How many steps to go down? (default=1): " steps
      steps=${steps:-1}
      echo "🔁 Rolling back $steps step(s)..."
      go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate \
        -database "$DB_URL" -path database/migrations down $steps
      ;;
    3)
      read -p "Enter migration version number to go to: " version
      echo "➡️ Migrating to version $version ..."
      go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate \
        -database "$DB_URL" -path database/migrations goto $version
      ;;
    *)
      echo "❌ Invalid migration option."
      ;;
  esac
}

# Option 4
create_migration() {
  read -p "Enter migration name (e.g. add_fullname_to_users): " name
  echo "📄 Creating new migration: $name ..."
  go run github.com/golang-migrate/migrate/v4/cmd/migrate create \
    -ext sql -dir database/migrations -seq "$name"
}

# ===============================
# Main loop
# ===============================
while true; do
  show_menu
  case $choice in
    1) run_program ;;
    2) run_with_air ;;
    3) run_migrate ;;
    4) create_migration ;;
    5) echo "👋 Bye!"; exit 0 ;;
    *) echo "❌ Invalid choice. Please try again." ;;
  esac
done
