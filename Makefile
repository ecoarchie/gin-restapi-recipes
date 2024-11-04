swagger:
	swagger serve -F swagger --no-open ./swagger.json

swagger-generate:
	swagger generate spec -o ./swagger.json

fmt:
	go fmt ./...

run: fmt
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
	MONGO_DATABASE=demo \
	go run main.go