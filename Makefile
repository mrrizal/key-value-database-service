run:
	go run main.go

test:
	go clean -testcache && go test -v -coverprofile=coverage.out ./...

see-coverage:
	go tool cover -html=coverage.out

build:
	CGO_ENABLED=0 GOOS=linux go build -a -o kvs
