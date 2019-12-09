# sample-blog

It's a sample blog, written in Go.

## Repo

Project layout should follow <https://github.com/golang-standards/project-layout>.

## Libraries

### Framework

There are a lot of web/REST frameworks. Since it's a simple demo, any can be good, including the raw `net/http`. Selection criterias:

* Supports [OpenAPI](https://www.openapis.org/)
* Supports API versions
* Supports JWT
* Supports incoming test, similar to `net/http/httptest`
* Supports outgoing test, similar to faking `net/http.Client.Transport`
* Supports online metrics, if possible by Prometheus protocol.
* Supports [RFC7807](https://tools.ietf.org/html/rfc7807) (Problem Details for HTTP APIs)

### OpenAPI

There are several possibilities to generate Go code from OpenAPI spec or to generate OpenAPI spec from Go code. Finally, none of them was selected, but `swaggo/gin-swagger` (Go --> OpenAPI) can be the best alternative (not implemented).

**<https://github.com/OpenAPITools/openapi-generator>**

Example for generating Gin source code:

```sh
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate \
    -i /local/docs/petstore/petstore.yaml \
    -g go-gin-server \
    -o /local \
    --additional-properties packageName=petstoreserver \
    --additional-properties packageVersion=0.0.1 \
    --additional-properties hideGenerationTimestamp=true
```

Result: Non-conform project layout. Non-conform variable names. Gofmt should be run.

**<https://github.com/mikkeloscar/gin-swagger>**

Limited functionality. Generated interface should be implemented.

Uses:

* <https://github.com/go-swagger/go-swagger>
* <https://github.com/go-openapi>

Result: Non-conform project layout. Gofmt should be run.

**<https://github.com/swaggo/gin-swagger>**

Generates OpenAPI spec (Swagger 2.0) from Go annotations.

Uses:

* <https://github.com/swaggo/swag> (supports: Gin, Echo, Buffallo, net/http)

Examples:

* <https://martinheinz.dev/blog/9>

## API

A simple API **plan** for the servce (with `/api/v1` prefix):

`/login`

* `POST`: Login (IMPLEMENTED)

`/refresh_token`

* `GET`: Refresh auth token (IMPLEMENTED)

`/entry`

* `POST`: Write a new blog entry
* `GET`: Get posts by filter

`/entry/:entry`

* `GET`: Get a specific blog entry
* `DELETE`: Delete a specific blog entry

`/entry/:entry/comment`

* `POST`: Write a new comment (IMPLEMENTED)
* `GET`: Get comments of a entry by filter

`/entry/:entry/comment/:comment`

* `DELETE`: Delete a specific comment

`/stat/user-post-comment`

* `GET`: Get a specific statistics by filter (IMPLEMENTED)

Paths for O&M (without `/api/v1` prefix):

`/version`

* Build num and timestamp, Git tag (IMPLEMENTED)

`/metrics`

* Prometheus metrics (IMPLEMENTED)

## Database

It's simple: Sqlite. Sqlite supports in-memory storage, so faking is not needed during automatic tests.

If different DBs should be supported, Factory+Builder design pattern can be used, where in-memory Sqlite can be used for automatic tests. Currently, Postgres might be used, too, but not tested.

Unfortunatelly, Sqlite does not support foreign keys: <https://github.com/jinzhu/gorm/issues/635>. Example for workaround:

```go
UserID uint   `json:"userID" sql:"type:integer REFERENCES users(id)"`
```

There are a few other issues of Sqlite, because it's a very simple database, see TODO at `GetUserPostCommentStats` and Dockefile (binary cannot be compiled statically).

Running frontend, using filesystem:

```sh
mkdir -p tmp/sqlite
./sample-blog frontend --db-dsn tmp/sqlite/blog.db
```

### Authentication

A popular Gin JWT framework was selected: <https://github.com/appleboy/gin-jwt>. First example on the project page was adopted.

Example for getting token:

```text
curl -s -H "Content-Type: application/json" -X POST --data '{"username":"kovacsj","password":"kovacs12"}' localhost:8088/api/v1/login | jq

{
  "code": 200,
  "expire": "2019-12-08T20:55:16+01:00",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzU4MzQ5MTYsImlkIjoia292YWNzaiIsIm9yaWdfaWF0IjoxNTc1ODMxMzE2fQ.1iYMVvUqYZxgqIHS9WFsj34IZJCH6LcKgjUB2MHpY50"
}
```

Example for refreshing token:

```text
curl -s -H "Content-Type: application/json" -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzU4MzQ5MTYsImlkIjoia292YWNzaiIsIm9yaWdfaWF0IjoxNTc1ODMxMzE2fQ.1iYMVvUqYZxgqIHS9WFsj34IZJCH6LcKgjUB2MHpY50" localhost:8088/api/v1/refresh_token | jq

{
  "code": 200,
  "expire": "2019-12-08T20:56:19+01:00",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzU4MzQ5NzksImlkIjoia292YWNzaiIsIm9yaWdfaWF0IjoxNTc1ODMxMzc5fQ.dU9hDYT5-WksdkC7k82RN6XVklysTCvJyFTwjg1AwN4"
}
```

## Automatic test

Automatic tests can be executed by below commands:

```sh
go mod verify
go vet
golangci-lint run
go test -v ./...
```

### Unit tests

Below unit tests were implemented:

* `TestNewUser`
* `TestGetUser`
* `TestGetUserFailed`

### Function tests

E2E function test were written, because it makes better coverage than unit tests. The test cases imply a HTTP request to Gin, the middleware calls DB component, too. For making artificial environment, test case hacks a few things and starts Gin and handlers with different parameters (for example decorator function: `DecorHandlerDB`).

`stat/user-post-comment`

Positive test: `TestGetUserPostCommentStats`, see same with curl:

```text
curl -s localhost:8088/api/v1/stat/user-post-comment?days=4 | jq

{
  "1": {
    "userName": "kovacsj",
    "entries": 0,
    "comments": 1
  },
  "2": {
    "userName": "szabop",
    "entries": 0,
    "comments": 0
  },
  "3": {
    "userName": "kocsisi",
    "entries": 0,
    "comments": 0
  }
}
```

Negative test: `TestGetUserPostCommentStatsFailed`, see same with curl:

```text
curl -s localhost:8088/api/v1/stat/user-post-comment?days=négy | jq

{
  "type": "about:blank",
  "title": "Bad Request",
  "status": 400,
  "detail": "invalid interval: strconv.ParseUint: parsing \"négy\": invalid syntax",
  "details": {
    "error": "\"invalid interval: strconv.ParseUint: parsing \\\"négy\\\": invalid syntax\"",
    "time": "\"2019-12-08T16:02:20+01:00\""
  },
  "callstack": [
    "internal/dao.(*Handler).GetUserPostCommentStats() dao.go:210",
    "pkg/frontend.GetUserPostCommentStats() frontend.go:37",
    "internal/web.DecorHandlerDB.func1() gin.go:16",
    "github.com/gin-gonic/gin.(*Context).Next() context.go:124",
    "github.com/gin-gonic/gin.(*Engine).handleHTTPRequest() gin.go:389",
    "github.com/gin-gonic/gin.(*Engine).ServeHTTP() gin.go:351",
    "net/http.serverHandler.ServeHTTP() server.go:2802",
    "net/http.(*conn).serve() server.go:1890",
    "runtime.goexit() asm_amd64.s:1357"
  ]
}
```

`/entry/:entry/comment`

Positive test: `TestPostComment`, see same with curl:

```text
curl -s -H "Content-Type: application/json" -H "Authorization:Bearer $TOKEN" -X POST --data '{"text":"hello"}' localhost:8088/api/v1/entry/1/comment | jq

{
  "ID": 13,
  "CreatedAt": "2019-12-08T22:43:37.350267157+01:00",
  "UpdatedAt": "2019-12-08T22:43:37.350267157+01:00",
  "DeletedAt": null,
  "userID": 1,
  "entryID": 1,
  "text": "hello"
}
```

Negative test: `TestPostCommentFailed`, see similar with curl:

```text
curl -s -H "Content-Type: application/json" -H "Authorization:Bearer $TOKEN" -X POST --data '{"text":"hello"}' localhost:8088/api/v1/entry/2/comment | jq

{
  "type": "about:blank",
  "title": "Bad Request",
  "status": 400,
  "detail": "only own entry can be commented",
  "details": {
    "error": "\"only own entry can be commented\"",
    "time": "\"2019-12-08T23:59:23+01:00\""
  }
}
```

### Test coverage

Below commands can generate coverage metrics:

```sh
go test ./... -coverprofile tmp/cover.out
go tool cover -html=tmp/cover.out
```

Unfortunately, it does not detect E2E function test coverage.

## Prometheus

Slected library: <https://github.com/Depado/ginprom>. It looks better than <https://github.com/zsais/go-gin-prometheus> (no PushGateway).

Collecting metrics is not supported during tests (metrics will be registered multiple times).

Useful system metrics:

```text
curl -s localhost:8088/metrics | egrep 'routine|alloc_bytes |_inuse_bytes |open_fds'

# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 9
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.45784e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 1.45784e+06
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 3.579904e+06
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 6944
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 57664
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 557056
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 7
```

Gin-related metrics:

```text
curl -s localhost:8088/metrics | grep gin_

# HELP gin_gin_request_duration_seconds The HTTP request latencies in seconds.
# TYPE gin_gin_request_duration_seconds summary
gin_gin_request_duration_seconds{quantile="0.5"} 0.000736012
gin_gin_request_duration_seconds{quantile="0.9"} 0.001213785
gin_gin_request_duration_seconds{quantile="0.99"} 0.001429643
gin_gin_request_duration_seconds_sum 0.012871773
gin_gin_request_duration_seconds_count 17
# HELP gin_gin_request_size_bytes The HTTP request sizes in bytes.
# TYPE gin_gin_request_size_bytes summary
gin_gin_request_size_bytes{quantile="0.5"} 157
gin_gin_request_size_bytes{quantile="0.9"} 314
gin_gin_request_size_bytes{quantile="0.99"} 314
gin_gin_request_size_bytes_sum 3023
gin_gin_request_size_bytes_count 17
# HELP gin_gin_requests_total How many HTTP requests processed, partitioned by status code and HTTP method.
# TYPE gin_gin_requests_total counter
gin_gin_requests_total{code="200",handler="github.com/appleboy/gin-jwt/v2.(*GinJWTMiddleware).LoginHandler-fm",host="localhost:8088",method="POST",path="/api/v1/login"} 3
gin_gin_requests_total{code="200",handler="github.com/appleboy/gin-jwt/v2.(*GinJWTMiddleware).RefreshHandler-fm",host="localhost:8088",method="GET",path="/api/v1/refresh_token"} 2
gin_gin_requests_total{code="200",handler="github.com/pgillich/sample-blog/internal/web.DecorHandlerDB.func1",host="localhost:8088",method="GET",path="/api/v1/entry/:entry/comment"} 6
gin_gin_requests_total{code="200",handler="github.com/pgillich/sample-blog/internal/web.DecorHandlerDB.func1",host="localhost:8088",method="POST",path="/api/v1/entry/:entry/comment"} 3
gin_gin_requests_total{code="400",handler="github.com/pgillich/sample-blog/internal/web.DecorHandlerDB.func1",host="localhost:8088",method="GET",path="/api/v1/entry/:entry/comment"} 2
gin_gin_requests_total{code="400",handler="github.com/pgillich/sample-blog/internal/web.DecorHandlerDB.func1",host="localhost:8088",method="POST",path="/api/v1/entry/:entry/comment"} 1
# HELP gin_gin_response_size_bytes The HTTP response sizes in bytes.
# TYPE gin_gin_response_size_bytes summary
gin_gin_response_size_bytes{quantile="0.5"} 165
gin_gin_response_size_bytes{quantile="0.9"} 907
gin_gin_response_size_bytes{quantile="0.99"} 907
gin_gin_response_size_bytes_sum 4529
gin_gin_response_size_bytes_count 17
```

Example Docker compose deployment: <https://github.com/stefanprodan/dockprom>

## Build

Below script runs all tests and builds the binary and Docker image with build info:

```sh
scripts/build.sh
```

## Usage

Starting the service:

```sh
go build && ./sample-blog frontend
```

```sh
docker run --rm -it -p "8088:8088" github.com/pgillich/sample-blog:0.0.1
```

Parameters are written in help:

```text
./sample-blog frontend --help

Start blog frontend service.

Usage:
  sample-blog frontend [flags]

Flags:
      --db-dialect string   DB_DIALECT, DB dialect (Gorm driver name) (default "sqlite3")
      --db-dsn string       DB_DSN, DB connection info (default ":memory:")
      --db-sample           DB_SAMPLE, DB sample filling (default true)
  -h, --help                help for frontend
      --listen string       LISTEN, host:port listening on (default ":8088")

Global Flags:
      --log-level string   LOG_LEVEL, log level (default "DEBUG")
```

Examples for non-auth urls:

```text
curl -s localhost:8088/api/v1/stat/user-post-comment?days=20 | jq

{
  "1": {
    "userName": "kovacsj",
    "entries": 1,
    "comments": 8
  },
  "2": {
    "userName": "szabop",
    "entries": 4,
    "comments": 4
  },
  "3": {
    "userName": "kocsisi",
    "entries": 1,
    "comments": 0
  }
}
```

```text
curl -s  localhost:8088/version | jq

{
  "tag": "0.0.2",
  "commit": "a2b1095",
  "branch": "master",
  "buildTime": "2019-12-09T02:53:00+01:00"
}
```

Examples for auth urls:

```text
TOKEN=$(curl -s -H "Content-Type: application/json" -X POST --data '{"username":"kovacsj","password":"kovacs12"}' localhost:8088/api/v1/login | jq -r '.token'); echo $TOKEN

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzU4NDQ5MjEsImlkIjoia292YWNzaiIsIm9yaWdfaWF0IjoxNTc1ODQxMzIxfQ.PlWOBY_Hy1JCaKsebjfjdZWsR2IU8-RUP6at57MZOmU
```

```text
TOKEN=$(curl -s -H "Content-Type: application/json" -H "Authorization:Bearer $TOKEN" localhost:8088/api/v1/refresh_token | jq -r '.token'); echo $TOKEN

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzU4NDU1MTcsImlkIjoia292YWNzaiIsIm9yaWdfaWF0IjoxNTc1ODQxOTE3fQ.PFkD7sLxlw092849y1nKR_7QWJ3XAYc-Z2oaarQ_kuI
```

```text
curl -s -H "Content-Type: application/json" -H "Authorization:Bearer $TOKEN" -X POST --data '{"text":"hello"}' localhost:8088/api/v1/entry/1/comment | jq

{
  "ID": 13,
  "CreatedAt": "2019-12-08T22:43:37.350267157+01:00",
  "UpdatedAt": "2019-12-08T22:43:37.350267157+01:00",
  "DeletedAt": null,
  "userID": 1,
  "entryID": 1,
  "text": "hello"
}
```

## TODO

* Gin handlers should get and return more status codes.
* OpenAPI documentation.
* Docker compose deployment
* Kubernetes deployment
