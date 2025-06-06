run:	build
	@./bin/app
build: 
	@go build -o bin/app cmd/app/main.go

css:
	npx @tailwindcss/cli -i ./views/css/app.css -o ./cmd/app/public/assets/app.css --watch

# tailwindcss -i views/css/app.css -o /cmd/app/public/assets/app.css --watch 
# npx @tailwindcss/cli -i ./views/css/app.css -o ./cmd/app/public/assets/app.css --watch
# templ generate --watch --proxy="http://localhost:3000"


# css:
# 	npx @tailwindcss/cli -i ./views/css/app.css -o ./cmd/app/public/assets/app.css --watch

# dev:
# 	trap 'kill 0' SIGINT; \
# 	npx @tailwindcss/cli -i ./views/css/app.css -o ./cmd/app/public/assets/app.css --watch & \
# 	air & \
# 	templ generate --watch --proxy="http://localhost:3000" & \
# 	wait