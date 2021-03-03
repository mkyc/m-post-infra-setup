package cmd

import (
	"time"

	gar "github.com/mkyc/go-ansible-runner"
)

func makeAnsibleLogLevel(logLevel int) gar.LogsLevel {
	logger.Debug().Msgf("makeAnsibleLogLevel(%d)", logLevel)
	var l gar.LogsLevel
	switch logLevel {
	case 0:
		l = gar.L0
	case 1:
		l = gar.L1
	case 2:
		l = gar.L2
	case 3:
		l = gar.L3
	case 4:
		l = gar.L4
	case 5:
		l = gar.L5
	case 6:
		l = gar.L6
	default:
		l = gar.L1
	}
	logger.Trace().Msgf("will return: %v", l)
	return l
}

func ansiblePlan() (string, error) {
	logger.Debug().Msg("ansiblePlan")
	err := setCheckAndDiff()
	if err != nil {
		logger.Fatal().Err(err).Msg("setCheckAndDiff failed")
	}
	options := gar.Options{
		AnsibleRunnerDir: ResourcesDirectory,
		Playbook:         "entrypoint.yml",
		Ident:            time.Now().Format("20060102-150405"),
		LogsLevel:        makeAnsibleLogLevel(ansibleDebugLevel),
		Logger:           ZeroLogger{},
	}
	return gar.Run(options)
}

func ansibleRun() (string, error) {
	logger.Debug().Msg("ansibleRun")
	options := gar.Options{
		AnsibleRunnerDir: ResourcesDirectory,
		Playbook:         "entrypoint.yml",
		Ident:            time.Now().Format("20060102-150405"),
		LogsLevel:        makeAnsibleLogLevel(ansibleDebugLevel),
		Logger:           ZeroLogger{},
	}
	return gar.Run(options)
}
