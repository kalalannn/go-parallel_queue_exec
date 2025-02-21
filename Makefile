_mkdir_bin:
	@mkdir -p bin
build: _mkdir_bin
	go build -o bin/app_html_ws cmd/app_html_ws/main.go
run: build
	./bin/app_html_ws
run_daemon:
	CompileDaemon -build="make build" -command="make run" -pattern "(.+\\.go|.+\\.html|.+\\.css|.+\\.js)"