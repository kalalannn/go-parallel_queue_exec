# Parallel queue processing using Golang
### Overview
This project implements a set of parallel queue processing services using ```Golang```.
It provides ```REST API```, ```WebSocket API```, and ```Web UI``` interfaces for managing
and interacting with the queue processing system.
The system efficiently handles concurrent task execution by leveraging **goroutines** and **synchronization mechanisms** such as **mutexes** and **channels**
to ensure thread safety and proper coordination between workers.

The architecture enables scalable and efficient processing,
preventing race conditions and ensuring that queued tasks are executed in a controlled manner.

## Get started

> Note: refer to `Makefile` for more commands and information.

### 1. Local environment

These commands help you build and run the services locally.

* **`make local_build`** - Builds four binaries: `app_rest`, `app_ws`, `app_html`, and `app_html_ws`. These binaries represent different types of applications.

* **`make local_test`** - Runs local tests to ensure the codebase is working as expected.

* **`make local_run_app_(rest|ws|html|html_ws)`** - Runs the specified application locally.
    * `make local_run_app_rest` - Runs the application with **only REST API endpoints**.
    * `make local_run_app_ws` - Runs the application with **REST API endpoints and the WebSocket API**.
    * `make local_run_app_html` - Runs the application with the **REST API and the Web UI interface**.
    * `make local_run_app_html_ws` - Runs the application with the **REST API, WebSocket API, and Web UI interface**, providing the full application experience.

* > **Note:** `app_html_ws` provides the complete application with all features.  Use `make local_run_app_html_ws` to run this application locally.

### 2. Dockerized environment

These commands help you build and run the services within Docker containers.

* **`make build_images`** - Builds the necessary Docker images: `core`, `base`, `app_rest`, `app_ws`, `app_html`, and `app_html_ws`.

* **`make run_tests`** - Runs tests inside a Docker container to ensure the codebase functions correctly within the containerized environment.

* **`make run_app_(rest|ws|html|html_ws)_image`** - Runs the specified application image within a Docker container.

* **Stopping Containers:**
    * **`make stop_app_(rest|ws|html|html_ws)`** - Stops the Docker container associated with the specified application.
    *   **`make stop_all`** - Stops *all* running Docker containers associated with the project.

* **Stopping and Removing Containers:**
    * **`make stop_rm_app_(rest|ws|html|html_ws)`** - Stops and **removes** the Docker container associated with the specified application. This completely removes the container.
    *   **`make stop_rm_all`** - Stops and **removes** *all* running Docker containers associated with the project.

* > **Note:** To run the Docker image with all features (REST API, WebSocket API, and Web UI),
simply use `make run_app_html_ws_image`.
The Makefile will handle all the necessary steps to build and run the `app_html_ws` image.

## API Endpoints
### Web UI
* **`GET /`** - Serves the main page of the Web UI.
Access this endpoint in your web browser to interact with the application.

### REST API
* **`POST /schedule`** - Schedules tasks for processing.
This endpoint adds tasks to the queue to be processed.
    
    Example request:
    ```json
    {
        "task123": 1000,  // Task ID: "task123", Duration: 1000ms (1 second)
        "task456": 2500   // Task ID: "task789", Duration: 2500ms (2.5 seconds)
    }
    ```
* **`GET /tasks`** - Returns a map of `active` tasks (currently being processed) and list of `scheduled` tasks (waiting in the queue) at the moment.

    Example response:
    ```json
    {
        "active": {
            "task123": 1000,  // Task ID: "task123", Duration: 1000ms (1 second)
            "task456": 2500   // Task ID: "task456", Duration: 2500ms (2.5 seconds)
        },
        "scheduled": [
            { "task123": 1000 },  // Task ID: "task123", Duration: 1000ms (1 second)
            { "task456": 2500 },  // Task ID: "task456", Duration: 2500ms (2.5 seconds)
            { "task789": 5000 },  // Task ID: "task789", Duration: 5000ms (5 seconds)
        ]
    }
    ```

### WebSocket API
* **`ws://<host>:<port>/ws`** - Establishes a WebSocket connection for real-time communication.
After a connection is established, the server will push updates to the client as events occur.
Possible event types are:
    * **`scheduled`**: `[]` - A list of new scheduled tasks.
    * **`next`**: `{}` - The next task that will start processing when a worker becomes available.
    * **`start`**: `{}` - A worker has just started processing a task.
    * **`done`**: `{}` - A worker has finished processing a task.

    Example events:
    ```json
    // New tasks scheduled
    { "scheduled": [{ "task789": 2500 }, { "task012": 7000 }] }

    // Task is about to be processed
    { "next": { "task789": 2500 } }

    // Task started processing
    { "start": { "task789": 2500 } }

    // Task completed processing
    { "done": { "task789": 2500 } }
    ```

## Configuration
> Note: App configuration [folder](https://github.com/kalalannn/go-parallel_queue_exec/tree/main/config)

Yaml file:
```yaml
app:
  host: 0.0.0.0                         # Server host
  port: 8080                            # Server port
  static_endpoint: /static              # Endpoint for static files
  public_folder: ./public               # Server folder for static files
  views_folder: ./views                 # Server folder for HTML templates
  templates_ext: .html                  # HTML templates extension
  fiber_shutdown_timeout_ms: 1000       # Fiber shutdown timeout
  executor_shutdown_timeout_ms: 5000    # Executor shutdown timeout

executor_service:
  workers_limit: 10                     # Workers limit for parallel queue processing
```
