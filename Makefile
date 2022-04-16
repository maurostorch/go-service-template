BINARY_NAME=app

clean:
	go clean
	rm ${BINARY_NAME}

build:
	GOARCH=amd64 go build -o ${BINARY_NAME} main.go

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

pkg:
	docker build -t go-service-template:latest .

deploy:
	envsubst < k8s/deployment.yml | kubectl apply -f -

all: dep vet test test_coverage build
