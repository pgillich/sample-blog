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

There are several possibilities to generate Go code from OpenAPI spec or to generate OpenAPI spec from Go code. Finally, none of them was selected, but `swaggo/gin-swagger` (Go --> OpenAPI) can be the best alternative.

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

A simple API plan for the servce (with `/api/v1` prefix):

`/login`

* `POST`: Login

`/entry`

* `POST`: Write a new blog entry
* `GET`: Get posts by filter

`/entry/:entry`

* `GET`: Get a specific blog entry
* `DELETE`: Delete a specific blog entry

`/entry/:entry/comment`

* `POST`: Write a new comment
* `GET`: Get comments of a entry by filter

`/entry/:entry/comment/:comment`

* `DELETE`: Delete a specific comment

`/stat/user-post-comment`

* `GET`: Get a specific statistics by filter

Paths for O&M (without `/api/v1` prefix):

`/version`

* Build num and timestamp, Git tag

`/metrics`

* Prometheus metrics

## Database

It's simple: Sqlite. Sqlite supports in-memory storage, so faking is not needed during automatic tests.

If different DBs should be supported, Factory+Builder design pattern can be used, where in-memory Sqlite can be used for automatic tests. Currently, Postgres might be used, too, but not tested.

Unfortunatelly, Sqlite does not support foreign keys: <https://github.com/jinzhu/gorm/issues/635>. Example for workaround:

```go
UserID uint   `json:"userID" sql:"type:integer REFERENCES users(id)"`
```

There are a few other issues of Sqlite, because it's a very simple database, see TODO at `GetUserPostCommentStats`.

Running SQLite using filesystem:

```sh
mkdir -p tmp/sqlite
./sample-blog frontend --db-dsn tmp/sqlite/blog.db
```

## Automatic test

Automatic tests can be executed by below commands:

```sh
go vet
go mod verify
golangci-lint run
go test -v ./...
```

### Unit tests

Below unit tests were implemented:

* TestNewUser
* TestGetUser
* TestGetUserFailed

### Function tests

E2E function test were written, because it makes better coverage than unit tests. The test cases imply a HTTP request to Gin, the middleware calls DB component, too. For making artificial environment, test case hacks a few things and starts Gin and handlers with different parameters (for example decorator function: `DecorHandlerDB`).

`stat/user-post-comment`

Positive test: `TestGetUserPostCommentStats`, see same with curl:

```text
$ curl -s localhost:8088/api/v1/stat/user-post-comment?days=4 | jq
{
  "1": {
    "userName": "Kovács János",
    "entries": 0,
    "comments": 1
  },
  "2": {
    "userName": "Szabó Pál",
    "entries": 0,
    "comments": 0
  },
  "3": {
    "userName": "Kocsis Irma",
    "entries": 0,
    "comments": 0
  }
}
```

Negative test: `TestGetUserPostCommentStatsFailed`, see same with curl:

```text
$ curl -s localhost:8088/api/v1/stat/user-post-comment?days=négy | jq
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

## Prometheus

<https://github.com/Depado/ginprom>

## TODO

* Gin handlers should get and return more status codes.
