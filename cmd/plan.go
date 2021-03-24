package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "performs module plan operation",
	Long: `Performs module plan operation. 

To illustrate proposed changes this module runs expected logic with '--diff' and '--check' flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("plan called")
		configFilePath := filepath.Join(SharedDirectoryPath, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectoryPath, stateFileName)
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
		_, err = ansiblePlan()
		if err != nil {
			logger.Fatal().Err(err).Msg("ansible run failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
