
build:
	go build -o bin/goq main.go

run: build
	./bin/goq

destroy: clean
	rm -f internal/server/handler*.go

clean:
	rm -f */**/*.gen.go

gen:
	go generate ./...
