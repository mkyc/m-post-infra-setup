package cmd

import (
	"time"

	gar "github.com/mkyc/go-ansible-runner"
)

func ansibleRun() (string, error) {
	logger.Debug().Msg("ansibleRun")
	options := gar.Options{
		AnsibleRunnerDir: ResourcesDirectory,
		Playbook:         "entrypoint.yml",
		Ident:            time.Now().Format("20060102-150405"),
		Logger:           ZeroLogger{},
	}
	return gar.Run(options)
}
