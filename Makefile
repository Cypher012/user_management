# ======================
# Config
# ======================

MIGRATIONS_DIR := ./internal/db/migrations
GOOSE := goose

# ======================
# Goose commands
# ======================

.PHONY: goose-status
goose-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

.PHONY: goose-up
goose-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

.PHONY: goose-down
goose-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

.PHONY: goose-reset
goose-reset:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" reset

.PHONY: goose-redo
goose-redo:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" redo

.PHONY: goose-create
goose-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make goose-create name=init_schema"; \
		exit 1; \
	fi
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql
