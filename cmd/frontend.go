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
	Long:  `Start blog frontend service.`,
	Run: func(cmd *cobra.Command, args []string) {
		startFrontend()
	},
}

func init() { // nolint:gochecknoinits
	RootCmd.AddCommand(frontendCmd)

	registerStringOption(frontendCmd, configs.OptServiceHostPort, configs.DefaultServiceHostPort, "host:port listening on")

	registerStringOption(frontendCmd, configs.OptDbDialect, configs.DefaultDbDialect, "DB dialect (Gorm driver name)")
	registerStringOption(frontendCmd, configs.OptDbDsn, configs.DefaultDbDsn, "DB connection info")
	registerBoolOption(frontendCmd, configs.OptDbSample, configs.DefaultDbSample, "DB sample filling")
}

func startFrontend() {
	hostPort := viper.GetString(configs.OptServiceHostPort)
	logger.Get().Infof("Start Frontend on %s", hostPort)

	samples := []dao.CompactSample{}
	if viper.GetBool(configs.OptDbSample) {
		samples = dao.GetDefaultSampleFill()
	}

	dbHandler, err := dao.NewHandler(
		viper.GetString(configs.OptDbDialect), viper.GetString(configs.OptDbDsn), samples)
	if err != nil {
		logger.Get().Panic(err)
	}
	defer dbHandler.Close()

	if err := frontend.SetupGin(gin.New(), dbHandler, true).Run(hostPort); err != nil {
		logger.Get().Panic(err)
	}

	logger.Get().Info("App closing...")
}
