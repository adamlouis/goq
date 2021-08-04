
build:
	go build --tags "json1" -o bin/goq main.go

run: build
	GOQ_MODE=LOCAL ./bin/goq

destroy: clean
	rm -f internal/server/handler*.go

clean:
	rm -f */**/*.gen.go

gen:
	go generate ./...
