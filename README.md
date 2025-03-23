# Examplehtmxapp

Boilerplate for a HTMX app with templ, daisyui, user authentication sqlite/turso database, sqlc, and email verification

## Developing

###  Dependencies:

```sh
go install github.com/air-verse/air@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/a-h/templ/cmd/templ@latest

# Database migration tool
curl -sSf https://atlasgo.sh | sh

# Bun (or node)
curl -fsSL https://bun.sh/install | bash
bun add
```

### Running

Run `make dev` for live reloading

### Building

Run
```sh
bunx @tailwindcss/cli -i ./frontend/input.css -o ./assets/tailwind.css
sqlc generate
templ generate
go build
```
