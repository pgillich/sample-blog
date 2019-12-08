module github.com/pgillich/sample-blog

go 1.13

require (
	emperror.dev/emperror v0.21.3
	emperror.dev/errors v0.4.3
	github.com/gin-contrib/pprof v0.0.0-20180827024024-a27513940d36
	github.com/gin-gonic/gin v1.4.0
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/runtime v0.19.4
	github.com/go-openapi/strfmt v0.19.2
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.2
	github.com/jinzhu/gorm v1.9.11
	github.com/mikkeloscar/gin-swagger v0.0.0-20190909064159-056d7535fc0e
	github.com/opentracing/opentracing-go v1.0.2
	github.com/pgillich/errfmt v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/swaggo/swag v1.6.3 // indirect
	github.com/zalando/gin-oauth2 v1.5.0
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

replace github.com/pgillich/errfmt => /home/peter/work/src/github.com/pgillich_errorformatter
