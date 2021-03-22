package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

// applyCmd represents the run command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies planned changes to machines",
	Long: `Applies planned changes to machines. 

Using provided configuration file this command applies expected ansible logic. This command performs following steps: 
 - validates presence of config and module state files
 - performs 'ansible run' operation
 - updates module state file with applied config.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("run called")
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
		_, err = ansibleRun()
		if err != nil {
			logger.Fatal().Err(err).Msg("ansible run failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
