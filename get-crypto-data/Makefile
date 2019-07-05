build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/crypto pull/main.go

clean:
	rm -rf ./bin

prod: clean build
	sls deploy --verbose --stage prod

dev: clean build
	sls deploy --verbose --stage dev

test: clean build
	sls deploy --verbose --stage test