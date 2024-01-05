
uses Go api starter kit for boilerplate -> https://github.com/qiangxue/go-rest-api 

Adds following features
* JWT-based authentication
* Rate limiting for each user
* Share notes with users
* Search using Postgres full text search
* Tests
 
TODO: 
- salt and hash passwords
- add indexes to tables
- add more tests

## Getting Started

Install Go [the instructions](https://golang.org/doc/install). Requires**Go 1.13 or above**.

[Install Docker](https://www.docker.com/get-started) **Docker 17.05 or higher** 

After installing Go and Docker, run the following commands to start experiencing this starter kit:

```shell
git clone https://github.com/niranjan92/notes-api

cd notes-api

# sync dependencies
go mod tidy

# start a PostgreSQL database server in a Docker container
make db-start

# seed the database with some test data
make testdata

make run

# or run the API server with live reloading, which is useful during development
# requires fswatch (https://github.com/emcrisostomo/fswatch)
make run-live
```

At this time, you have a RESTful API server running at `http://127.0.0.1:8080`. It provides the following endpoints:

* `GET /healthcheck`: a healthcheck service provided for health checking purpose (needed when implementing a server cluster)
* `POST /api/auth/signup`: authenticates a user and generates a JWT
* `POST /api/auth/login`: authenticates a user and generates a JWT
* `GET /api/notes`: returns a list of notes for the user (includes notes shared with the user)
* `GET /api/notes/:id`: returns the detailed information of an note
* `POST /api/notes`: creates a new note
* `PUT /api/notes/:id`: updates an existing note
* `DELETE /api/notes/:id`: deletes a note
* `POST /api/notes/:id/shares/:user_id`: shares a note with another user id
* `GET /api/search?q=<query>`: searches for matching word 

Try the URL `http://localhost:8080/healthcheck` in a browser, and you should see something like `"OK v1.0.0"` displayed.

```shell
# authenticate the user via: POST /v1/login
curl -X POST -H "Content-Type: application/json" -d '{"username": "demo", "password": "pass"}' http://localhost:8080/api/login
# should return a JWT token like: {"token":"...JWT token here..."}

# with the above JWT token, access the note resources, such as: GET /api/notes
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:8080/api/notes
```

To use the starter kit as a starting point of a real project whose package name is `github.com/abc/xyz`, do a global 
replacement of the string `github.com/qiangxue/go-rest-api` in all of project files with the string `github.com/abc/xyz`.


## Project Layout

The starter kit uses the following project layout:
 
```
.
├── cmd                  main applications of the project
│   └── server           the API server application
├── config               configuration files for different environments
├── internal             private application and library code
│   ├── notes            notes-related features
│   ├── auth             authentication feature
│   ├── config           configuration library
│   ├── entity           entity definitions and domain logic
│   ├── errors           error types and handling
│   ├── healthcheck      healthcheck feature
│   └── test             helpers for testing purpose
├── migrations           database migrations
├── pkg                  public library code
│   ├── accesslog        access log middleware
│   ├── graceful         graceful shutdown of HTTP server
│   ├── log              structured and context-aware logger
│   └── pagination       paginated list
└── testdata             test data scripts
```

The top level directories `cmd`, `internal`, `pkg` are commonly found in other popular Go projects, as explained in
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

For example, the `notes` directory contains the application logic related with the note feature. 

Within each feature package, code are organized in layers (API, service, repository), following the dependency guidelines
as described in the [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).


## Common Development Tasks

schema changes:

```shell
make migrate

make migrate-new

make migrate-down

make migrate-reset
```

### Managing Configurations

The `config` directory contains the configuration files named after different environments. For example,
`config/local.yml` corresponds to the local development environment and is used when running the application 
via `make run`.

```shell
./server -config=./config/prod.yml
```

```
