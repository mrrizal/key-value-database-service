run:
	go run main.go

test:
	go clean -testcache && go test -v -coverprofile=coverage.out ./...

view-coverage:
	go tool cover -html=coverage.out
