_mkdir_bin:
	@mkdir -p bin
build: _mkdir_bin
	go build -o bin/app cmd/app/main.go
run: build
	./bin/app
run_daemon:
	CompileDaemon -build="go build -o bin/app cmd/app/main.go" -command="./bin/app"