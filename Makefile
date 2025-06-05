run:	build
	@./bin/app
build: 
	@go build -o bin/app cmd/app/main.go

css:
	npx @tailwindcss/cli -i ./views/css/app.css -o ./cmd/app/public/app.css --watch


# templ generate --watch --proxy="http://localhost:3000"