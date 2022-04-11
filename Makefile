BINARY_NAME=app

clean:
	go clean
	rm ${BINARY_NAME}

build:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME} main.go

run:
	go run main.go

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all

all: dep vet test test_coverage build
