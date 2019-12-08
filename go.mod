module github.com/pgillich/sample-blog

go 1.13

// replace github.com/pgillich/errfmt => /home/peter/work/src/github.com/pgillich_errorformatter

require (
	emperror.dev/errors v0.4.3
	github.com/appleboy/gin-jwt/v2 v2.6.2
	github.com/gin-gonic/gin v1.5.0
	github.com/jinzhu/gorm v1.9.11
	github.com/pgillich/errfmt v0.1.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
)
