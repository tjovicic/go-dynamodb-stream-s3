build: format
	env GOOS=linux go build -ldflags="-s -w" -o bin/main cmd/main.go

format:
	go fmt ./...

clean:
	rm -rf ./bin

deploy: clean format build
	sls deploy --verbose

remove:
	sls remove
