.PHONY: build build-darwin build-linux

build: build-darwin build-linux

build-darwin: test
	env GOOS=darwin GOARCH=amd64 go build -o alternator-darwin -ldflags '-s -w -X github.com/kota65535/alternator/cmd.Version=0.0.0' main.go

build-linux: test
	env GOOS=linux GOARCH=amd64 go build -o alternator-linux -ldflags '-s -w -X github.com/kota65535/alternator/cmd.Version=0.0.0' main.go

test: yacc compose-up
	go test -v -cover -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html
	docker-compose down

compose-up:
	docker-compose up -d
	while ! (mysqladmin ping -h 127.0.0.1 -P 13306 -u root --silent); do sleep 5; done
	while ! (mysqladmin ping -h 127.0.0.1 -P 13307 -u root --silent); do sleep 5; done

yacc:
	goyacc -o parser/parser.go parser/parser.go.y

clean:
	rm -f alternator-darwin alternator-linux y.output cover.out cover.html