up:
	docker-compose -f deploy/docker-compose.yml up -d postgres redis
ps:
	docker-compose -f deploy/docker-compose.yml ps -a
down:
	docker-compose -f deploy/docker-compose.yml down
rebuild:
	docker-compose -f deploy/docker-compose.yml up --build --force-recreate
reload:
	cd web && air -c air.toml
test:
	cd web && go test ./internal/service
migrate_up_new:
	cd web && migrate -source file://./internal/db/migrations -database postgres://postgres:secret@localhost:5432/postgres_sa?sslmode=disable up
migrate_down_all:
	cd web && migrate -source file://./internal/db/migrations -database postgres://postgres:secret@localhost:5432/postgres_sa?sslmode=disable down
generate_swagger:
	cd web && swag init --parseInternal -d ./cmd/subsaggregator,./internal/router,./internal/service,./internal/model,./internal/utils --parseDependency
