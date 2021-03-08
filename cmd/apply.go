package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

// applyCmd represents the run command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("run called")
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		config, _, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("checkAndLoad failed")
		}

		err = templateInventory(config)
		if err != nil {
			logger.Fatal().Err(err).Msg("templateInventory failed")
		}
		err = prepareSshKey(config)
		if err != nil {
			logger.Fatal().Err(err).Msg("prepareSshKey failed")
		}
		_, err = ansibleRun()
		if err != nil {
			logger.Fatal().Err(err).Msg("ansible run failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}