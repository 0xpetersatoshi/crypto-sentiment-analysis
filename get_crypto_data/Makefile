.PHONY: build clean deploy

all: clean build deploy clean

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/crypto pull/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
