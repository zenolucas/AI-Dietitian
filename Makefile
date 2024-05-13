run: build
	@./bin/AI-Dietitian

templ:
	@templ generate --watch --proxy=http://localhost:3000

css:
	@tailwindcss -i view/css/app.css -o public/styles.css --watch

install: 
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D tailwindcss
	@npm install -D daisyui@latest

build: 
	tailwindcss -i view/css/app.css -o public/styles.css
	@templ generate view
	@go build -o bin/AI-Dietitian main.go