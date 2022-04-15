hello:
	echo "Hello"
serve:
	swag init -g core/app.go
	go build
	./io-service