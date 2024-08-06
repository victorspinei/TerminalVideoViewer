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
	echo "https://youtu.be/cNWKRMM_H3A?si=Xq6ltLjXOCNV5zBH" | ./bin/${BINARY_NAME}