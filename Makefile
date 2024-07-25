BINARY_NAME=main.out

build:
	go build -o ./bin/${BINARY_NAME} ./cmd/main/main.go

run:
	go build -o ./bin/${BINARY_NAME} ./cmd/main/main.go
	./bin/${BINARY_NAME}

clean:
	go clean
	rm ./bin/${BINARY_NAME}