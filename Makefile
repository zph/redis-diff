binary := bin/redis-diff

build: deps
	go build -o $(binary)

deps:
	go get gopkg.in/redis.v4
	go get github.com/stretchr/testify/assert

test:
	go test

clean:
	rm -f $(binary)

binary: build
