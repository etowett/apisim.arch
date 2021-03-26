TAG=v0.1.0
BINARY=apisim
NAME=ektowett/$(BINARY)
IMAGE=$(NAME):$(TAG)
LATEST=$(NAME):latest
LDFLAGS := -ldflags ""


run:
	@revel run

build:
	@docker-compose build

up:
	@docker-compose up -d

build_live:
	@echo "Building the image $(IMAGE)"
	@docker build -t $(IMAGE) . -f Dockerfile
	@echo "Tagging the image $(IMAGE) to latest"
	@docker tag $(IMAGE) $(LATEST)
	@echo "Remove the go binary"
	@rm $(BINARY)
	@echo "Done!"

logs:
	docker-compose logs -f

ps:
	@docker-compose ps

stop:
	@docker-compose stop

rm: stop
	@docker-compose rm

build_cli:
	@echo "Building cli to /tmp/apisim-cli"
	@go build -o /tmp/apisim-cli ./scripts/cli/main.go
	@echo "Done!"

# make migration name=create_users
migration:
	@echo "Creating migration $(name)!"
	@goose -dir app/migrations create $(name) sql
	@echo "Done!"

migrate_up:
	@echo "Migrating up!"
	@goose -dir app/migrations postgres 'postgres://apisim:apisim@127.0.0.1/apisim?sslmode=disable' up
	@echo "Done!"

migrate_down:
	@echo "Migrating down!"
	@goose -dir app/migrations postgres 'postgres://apisim:apisim@127.0.0.1/apisim?sslmode=disable' down
	@echo "Done!"

migrate_status:
	@echo "Getting migration status!"
	@goose -dir app/migrations postgres 'postgres://apisim:apisim@127.0.0.1/apisim?sslmode=disable' status
	@echo "Done!"

migrate_reset:
	@echo "Resetting migrations!"
	@goose -dir app/migrations postgres 'postgres://apisim:apisim@127.0.0.1/apisim?sslmode=disable' reset
	@echo "Done!"

migrate_version:
	@echo "Getting migration version!"
	@goose -dir app/migrations postgres 'postgres://apisim:apisim@127.0.0.1/apisim?sslmode=disable' version
	@echo "Done!"

migrate_redo: migrate_reset migrate_up
