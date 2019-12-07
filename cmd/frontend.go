package cmd

import (
	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pgillich/sample-blog/configs"
	"github.com/pgillich/sample-blog/internal/dao"
	"github.com/pgillich/sample-blog/internal/logger"
	"github.com/pgillich/sample-blog/pkg/frontend"
)

// nolint:gochecknoglobals
var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Frontend",
	Long:  `Start chat bot frontend service.`,
	Run: func(cmd *cobra.Command, args []string) {
		startFrontend()
	},
}

func init() { // nolint:gochecknoinits
	RootCmd.AddCommand(frontendCmd)

	registerStringOption(frontendCmd, configs.OptServiceHostPort, configs.DefaultServiceHostPort, "host:port listening on")

	registerStringOption(frontendCmd, configs.OptDbDsn, configs.DefaultDbDsn, "DB connection info")
}

func startFrontend() {
	hostPort := viper.GetString(configs.OptServiceHostPort)
	logger.Get().Infof("Start Frontend on %s", hostPort)

	db, err := dao.ConnectSqlite(viper.GetString(configs.OptDbDsn))
	if err != nil {
		logger.Get().Panic(err)
	}
	defer db.Close() //nolint:errcheck

	if err := frontend.SetupGin(gin.New(), db).Run(); err != nil {
		logger.Get().Panic(err)
	}

	logger.Get().Info("App closing...")
}
