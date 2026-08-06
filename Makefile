.PHONY: build run up down tidy swag  help migrate-up migrate-down migrate-create migrate-force migrate-status seed seed-crop seed-farm seed-growing

build:
	go build -o bin/agroapp ./cmd/agroapp

run:
	go run ./cmd/agroapp

up:
	docker compose up -d

down:
	docker compose down

tidy:
	go mod tidy


DB_DSN ?= postgres://tiun:tiun@lab.raspi:5432/tiun?sslmode=disable


# Запуск всех seed-данных
seed:
	@echo "Seeding all data..."
	go run cmd/seed/main.go -data=./data -module=all

# Только crop
seed-crop:
	@echo "Seeding crop data..."
	go run cmd/seed/main.go -data=./data -module=crop

# Только farm
seed-farm:
	@echo "Seeding farm data..."
	go run cmd/seed/main.go -data=./data -module=farm

# Только growing
seed-growing:
	@echo "Seeding growing data..."
	go run cmd/seed/main.go -data=./data -module=growing

# Проверка seed-данных (dry-run)
seed-dry-run:
	@echo "Validating seed data..."
	go run cmd/seed/main.go -data=./data -dry-run -module=all



help:
	@echo "Commands:"
	@echo "  make migrate-up MODULE=farm           - Run up migrations"
	@echo "  make migrate-down MODULE=farm STEPS=1 - Rollback"
	@echo "  make migrate-create NAME=xxx MODULE=farm - Create migration"
	@echo "  make migrate-force VERSION=1 MODULE=farm - Force version"
	@echo "  make migrate-status MODULE=farm       - Show status"

migrate-up:
	@if [ -z "$(MODULE)" ]; then \
		echo "Usage: make migrate-up MODULE=farm"; \
		exit 1; \
	fi
	 go run cmd/migrate/main.go -module=$(MODULE) -direction=up -verbose=true

migrate-down:
	@if [ -z "$(MODULE)" ]; then \
		echo "Usage: make migrate-down MODULE=farm STEPS=1"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go \
		-module=$(MODULE) \
		-direction=down \
		-steps=$(STEPS) \
		-verbose=true

migrate-create:
	@if [ -z "$(NAME)" ] || [ -z "$(MODULE)" ]; then \
		echo "Usage: make migrate-create NAME=migration_name MODULE=farm"; \
		exit 1; \
	fi
	@mkdir -p ./migrations/$(MODULE)
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	up_file="./migrations/$(MODULE)/$${timestamp}_$(NAME).up.sql"; \
	down_file="./migrations/$(MODULE)/$${timestamp}_$(NAME).down.sql"; \
	echo "-- Up migration: $(NAME)" > $$up_file; \
	echo "-- Created: $$(date)" >> $$up_file; \
	echo "" >> $$up_file; \
	echo "CREATE TABLE IF NOT EXISTS $(NAME) (" >> $$up_file; \
	echo "    id UUID PRIMARY KEY DEFAULT gen_random_uuid()" >> $$up_file; \
	echo ");" >> $$up_file; \
	echo "" >> $$up_file; \
	echo "-- Down migration: $(NAME)" > $$down_file; \
	echo "DROP TABLE IF EXISTS $(NAME);" >> $$down_file; \
	echo "Created: $$up_file"; \
	echo "Created: $$down_file"

migrate-force:
	@if [ -z "$(VERSION)" ] || [ -z "$(MODULE)" ]; then \
		echo "Usage: make migrate-force VERSION=1 MODULE=farm"; \
		exit 1; \
	fi
	cd migrations && go run cmd/migrate/main.go \
		-module=$(MODULE) \
		-force=$(VERSION)

migrate-status:
	@if [ -z "$(MODULE)" ]; then \
		echo "Usage: make migrate-status MODULE=farm"; \
		exit 1; \
	fi
	@psql "$(DB_DSN)" -c "SELECT version, applied_at FROM schema_migrations WHERE module = '$(MODULE)' ORDER BY version;"


desc:
	cd ./myproject && wails dev

swagger:
	swag init \
		-g main.go \
		-d cmd/api,internal/application/commands,internal/application/queries \
		--parseDependency \
		--parseInternal \
        --parseDependencyLevel 3 \
		--parseDepth 3 \
		--output docs/swagger \
		--outputTypes json,yaml