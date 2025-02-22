#! Local variables
BIN_APP_REST=bin/app_rest
BIN_APP_WS=bin/app_ws
BIN_APP_HTML=bin/app_html
BIN_APP_HTML_WS=bin/app_html_ws
SRC_APP_REST=cmd/app_rest/main.go
SRC_APP_WS=cmd/app_ws/main.go
SRC_APP_HTML=cmd/app_html/main.go
SRC_APP_HTML_WS=cmd/app_html_ws/main.go

#! Local builds
_mkdir_bin:
	@mkdir -p bin
_local_build_app_rest:
	go build -o ${BIN_APP_REST} ${SRC_APP_REST}
_local_build_app_ws:
	go build -o ${BIN_APP_WS} ${SRC_APP_WS}
_local_build_app_html:
	go build -o ${BIN_APP_HTML} ${SRC_APP_HTML}
_local_build_app_html_ws:
	go build -o ${BIN_APP_HTML_WS} ${SRC_APP_HTML_WS}

local_build: _mkdir_bin _local_build_app_rest _local_build_app_ws _local_build_app_html _local_build_app_html_ws

local_build_app_rest:    _mkdir_bin _local_build_app_rest
local_build_app_ws:      _mkdir_bin _local_build_app_ws
local_build_app_html:    _mkdir_bin _local_build_app_html
local_build_app_html_ws: _mkdir_bin _local_build_app_html_ws

#! Local runs
_local_run_app_rest:
	${BIN_APP_REST}
_local_run_app_ws:
	${BIN_APP_WS}
_local_run_app_html:
	${BIN_APP_HTML}
_local_run_app_html_ws:
	${BIN_APP_HTML_WS}

local_run_app_rest:    _mkdir_bin _local_build_app_rest    _local_run_app_rest
local_run_app_ws:      _mkdir_bin _local_build_app_ws      _local_run_app_ws
local_run_app_html:    _mkdir_bin _local_build_app_html    _local_run_app_html
local_run_app_html_ws: _mkdir_bin _local_build_app_html_ws _local_run_app_html_ws

#! Local CompileDaemons
PATTERN_GO_HTML_CSS_JS ="(.+\\.go|.+\\.html|.+\\.css|.+\\.js)"
local_daemon_app_rest:
	CompileDaemon -build="make _local_build_app_rest" -command="make _local_run_app_rest"
local_daemon_app_ws:
	CompileDaemon -build="make _local_build_app_ws" -command="make _local_run_app_ws"
local_daemon_app_html:
	CompileDaemon -build="make _local_build_app_html" -command="make _local_run_app_html" -pattern ${PATTERN_GO_HTML_CSS_JS}
local_daemon_app_html_ws:
	CompileDaemon -build="make _local_build_app_html_ws" -command="make _local_run_app_html_ws" -pattern ${PATTERN_GO_HTML_CSS_JS}

#! Docker variables
CORE_CONTAINER_NAME=go-parallel_queue_exec-core
BASE_CONTAINER_NAME=go-parallel_queue_exec-base
APP_REST_CONTAINER_NAME=go-parallel_queue_exec-app_rest
APP_WS_CONTAINER_NAME=go-parallel_queue_exec-app_ws
APP_HTML_CONTAINER_NAME=go-parallel_queue_exec-app_html
APP_HTML_WS_CONTAINER_NAME=go-parallel_queue_exec-app_html_ws

CORE_IMAGE=$(CORE_CONTAINER_NAME):latest
BASE_IMAGE=$(BASE_CONTAINER_NAME):latest
APP_REST_IMAGE=$(APP_REST_CONTAINER_NAME):latest
APP_WS_IMAGE=$(APP_WS_CONTAINER_NAME):latest
APP_HTML_IMAGE=$(APP_HTML_CONTAINER_NAME):latest
APP_HTML_WS_IMAGE=$(APP_HTML_WS_CONTAINER_NAME):latest

show_images:
	-docker images | grep "parallel_queue_exec"

#! Build docker images
_build_core_image:
	docker build --file=Dockerfile.core --tag=${CORE_IMAGE} .
_build_base_image: _build_core_image
	docker build --file=Dockerfile.base --tag=${BASE_IMAGE} .

_build_app_rest_image:
	docker build --file=Dockerfile.app_rest --tag=${APP_REST_IMAGE} .
_build_app_ws_image:
	docker build --file=Dockerfile.app_ws --tag=${APP_WS_IMAGE} .
_build_app_html_image:
	docker build --file=Dockerfile.app_html --tag=${APP_HTML_IMAGE} .
_build_app_html_ws_image:
	docker build --file=Dockerfile.app_html_ws --tag=${APP_HTML_WS_IMAGE} .

build_images: _build_core_image _build_base_image _build_app_rest_image _build_app_ws_image _build_app_html_image _build_app_html_ws_image show_images

build_app_rest_image:    _build_core_image _build_base_image _build_app_rest_image
build_app_ws_image:      _build_core_image _build_base_image _build_app_ws_image
build_app_html_image:    _build_core_image _build_base_image _build_app_html_image
build_app_html_ws_image: _build_core_image _build_base_image _build_app_html_ws_image

#! Run docker images
run_app_rest_image: build_app_rest_image
	docker run -d -p 8080:8080 --name ${APP_REST_CONTAINER_NAME} ${APP_REST_IMAGE}
run_app_ws_image: build_app_ws_image
	docker run -d -p 8080:8080 --name ${APP_WS_CONTAINER_NAME} ${APP_WS_IMAGE}
run_app_html_image: build_app_html_image
	docker run -d -p 8080:8080 --name ${APP_HTML_CONTAINER_NAME} ${APP_HTML_IMAGE}
run_app_html_ws_image: build_app_html_ws_image
	docker run -d -p 8080:8080 --name ${APP_HTML_WS_CONTAINER_NAME} ${APP_HTML_WS_IMAGE}

#! Stop docker images
stop_app_rest_image:
	@docker kill --signal=SIGINT ${APP_REST_CONTAINER_NAME}
stop_app_ws_image:
	@docker kill --signal=SIGINT ${APP_WS_CONTAINER_NAME}
stop_app_html_image:
	@docker kill --signal=SIGINT ${APP_HTML_CONTAINER_NAME}
stop_app_html_ws_image:
	@docker kill --signal=SIGINT ${APP_HTML_WS_CONTAINER_NAME}

stop_all:
	@make stop_app_rest    2> /dev/null || true
	@make stop_app_ws      2> /dev/null || true
	@make stop_app_html    2> /dev/null || true
	@make stop_app_html_ws 2> /dev/null || true

#! Stop and remove docker images
stop_rm_app_rest: stop_app_rest_image
	@docker wait ${APP_REST_CONTAINER_NAME}
	@docker logs ${APP_REST_CONTAINER_NAME}
	@docker rm ${APP_REST_CONTAINER_NAME}

stop_rm_app_ws: stop_app_ws_image
	@docker wait ${APP_WS_CONTAINER_NAME}
	@docker logs ${APP_WS_CONTAINER_NAME}
	@docker rm ${APP_WS_CONTAINER_NAME}

stop_rm_app_html: stop_app_html_image
	@docker wait ${APP_HTML_CONTAINER_NAME}
	@docker logs ${APP_HTML_CONTAINER_NAME}
	@docker rm ${APP_HTML_CONTAINER_NAME}

stop_rm_app_html_ws: stop_app_html_ws_image
	@docker wait ${APP_HTML_WS_CONTAINER_NAME}
	@docker logs ${APP_HTML_WS_CONTAINER_NAME}
	@docker rm ${APP_HTML_WS_CONTAINER_NAME}

stop_rm_all:
	@make stop_rm_app_rest    2> /dev/null || true
	@make stop_rm_app_ws      2> /dev/null || true
	@make stop_rm_app_html    2> /dev/null || true
	@make stop_rm_app_html_ws 2> /dev/null || true
