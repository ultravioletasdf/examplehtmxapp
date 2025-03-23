include .env
export $(shell sed 's/=.*//' .env)

dev:
	make -j2 dev/air dev/templ
dev/air:
	air
dev/templ:
	templ generate --watch --proxy="http://localhost:3000"
migrate:
	atlas schema apply --url ${DATABASE_MIGRATION_URL} --dev-url "sqlite://dev.db" --to "file://sql/schema.sql" --auto-approve
