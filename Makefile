BINARY_NAME=main.out

build:
	go build -o ./bin/${BINARY_NAME} ./cmd/main/main.go

run:
	go build -o ./bin/${BINARY_NAME} ./cmd/main/main.go
	./bin/${BINARY_NAME}

clean:
	go clean
	rm ./bin/${BINARY_NAME}
	rm -rf ./temp/*

test:
	go build -o ./bin/${BINARY_NAME} ./cmd/main/main.go
	echo "https://www.youtube.com/watch?v=ucZl6vQ_8Uo" | ./bin/${BINARY_NAME}