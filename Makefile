
build:
	go build --tags "json1" -o bin/goq main.go

run: build
	mkdir -p /tmp/goq-session-store
	GOQ_MODE=DEVELOPMENT \
		GOQ_API_KEY=letmein  \
		GOQ_ROOT_USERNAME=root \
		GOQ_ROOT_PASSWORD=letmein \
		GOQ_SESSION_KEY=itsasecret \
		GOQ_SESSION_STORE_PATH=/tmp/goq-session-store \
		./bin/goq

run-no-build:
	mkdir -p /tmp/goq-session-store
	GOQ_MODE=DEVELOPMENT \
		GOQ_API_KEY=letmein  \
		GOQ_ROOT_USERNAME=root \
		GOQ_ROOT_PASSWORD=letmein \
		GOQ_SESSION_KEY=itsasecret \
		GOQ_SESSION_STORE_PATH=/tmp/goq-session-store \
		./bin/goq

destroy: clean
	rm -f internal/server/handler*.go

clean:
	rm -f */**/*.gen.go

gen:
	go generate ./...


live:
	CompileDaemon -build "make build" -command "make run-no-build" -graceful-kill -color -polling -verbose

