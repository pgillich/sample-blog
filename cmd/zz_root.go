//	Package cmd is the CLI handler
package cmd

import (
	goflag "flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pgillich/sample-blog/configs"
	"github.com/pgillich/sample-blog/internal/logger"
)

// RootCmd is the root command
var RootCmd = &cobra.Command{ // nolint:gochecknoglobals
	Use:   "sample-blog",
	Short: "Sample blog",
	Long:  `Sample blog`,
}

// Execute is the main function
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("Runtime error: %s\n", err)
		os.Exit(1)
	}
}

func getEnvReplacer() *strings.Replacer {
	return strings.NewReplacer("-", "_", ".", "_")
}

func init() { // nolint:gochecknoinits
	cobra.OnInitialize(initConfig)

	cobra.OnInitialize()

	registerStringOption(RootCmd, configs.OptLogLevel, configs.DefaultLogLevel, "log level")

	goflag.CommandLine.Usage = func() {
		RootCmd.Usage() // nolint:gosec,errcheck
	}
	goflag.Parse()

	logger.Init(viper.GetString(configs.OptLogLevel))
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(getEnvReplacer())
}

func registerStringOption(command *cobra.Command, name string, value string, usage string) {
	envName := getEnvReplacer().Replace(name)
	command.PersistentFlags().String(name, value, strings.ToUpper(envName)+", "+usage)
	viper.BindPFlag(name, command.PersistentFlags().Lookup(name)) // nolint:errcheck,gosec
}

func registerBoolOption(command *cobra.Command, name string, value bool, usage string) {
	envName := getEnvReplacer().Replace(name)
	command.PersistentFlags().Bool(name, value, strings.ToUpper(envName)+", "+usage)
	viper.BindPFlag(name, command.PersistentFlags().Lookup(name)) // nolint:errcheck,gosec
}
