.PHONY: dev css templ build run

# Generate Templ code
templ:
	templ generate

# Build Tailwind CSS (watch mode)
css:
	npx tailwindcss -i ./static/input.css -o ./static/styles.css --watch

# Build Go binary
build: templ
	go build -o tmp/main.exe ./cmd/main.go

# Run the server (no live reload)
run: build
	./tmp/main.exe

# Live reload with Air (needs .air.toml configured)
dev: templ
	air