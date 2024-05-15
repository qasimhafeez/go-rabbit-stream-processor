.DEFAULT_GOAL := app
app:
	docker compose down && docker compose up --build
lint:
	gofmt -w -s . && go mod tidy && go vet ./...
