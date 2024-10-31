swagger:
	swagger serve -F swagger --no-open ./swagger.json

swagger-generate:
	swagger generate spec -o ./swagger.json

run:
	go run main.go