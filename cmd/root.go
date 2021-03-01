package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

const (
	moduleShortName = "hi"
	configFileName  = "hi-config.json"
	stateFileName   = "state.json"

	defaultSharedDirectory    = "/shared"
	defaultResourcesDirectory = "/resources"
)

var (
	enableDebug bool

	Version string

	SharedDirectory    string
	ResourcesDirectory string

	logger zerolog.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "m-host-init",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PersistentPreRun")

		err := viper.BindPFlags(cmd.PersistentFlags())
		if err != nil {
			logger.Fatal().Err(err).Msg("BindPFlags failed")
		}

		SharedDirectory = viper.GetString("shared")
		ResourcesDirectory = viper.GetString("resources")
	},
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal().Err(err).Msg("rootCmd.Execute failed")
	}
}

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger = zerolog.New(output).With().Caller().Timestamp().Logger()

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&enableDebug, "debug", "d", false, "print debug information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if enableDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
	viper.AutomaticEnv() // read in environment variables that match
}
