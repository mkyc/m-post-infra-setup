package cmd

import (
	"encoding/json"
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
	err := enableCheckAndDiffMode()
	if err != nil {
		logger.Fatal().Err(err).Msg("enableCheckAndDiffMode failed")
	}
	options := gar.Options{
		AnsibleRunnerDir: ResourcesDirectoryPath,
		Playbook:         "entrypoint.yml",
		Ident:            time.Now().Format("20060102-150405"),
		LogsLevel:        makeAnsibleLogLevel(ansibleDebugLevel),
		Logger:           ZeroLogger{},
	}
	_, err = gar.Run(options)
	return checkPlayResults(options, err)
}

func ansibleRun() (string, error) {
	logger.Debug().Msg("ansibleRun")
	options := gar.Options{
		AnsibleRunnerDir: ResourcesDirectoryPath,
		Playbook:         "entrypoint.yml",
		Ident:            time.Now().Format("20060102-150405"),
		LogsLevel:        makeAnsibleLogLevel(ansibleDebugLevel),
		Logger:           ZeroLogger{},
	}
	_, err := gar.Run(options)
	return checkPlayResults(options, err)
}

func checkPlayResults(options gar.Options, err error) (string, error) {

	output, err2 := gar.GetOutput(options)
	if err2 != nil {
		return string(output), err2
	}
	rc, status, _ := gar.GetStatus(options)
	logger.Info().Msgf("RC: %d, Status: %s", rc, status)

	pr, _ := gar.GetPlayRecap(options)
	prBytes, _ := json.Marshal(pr)
	logger.Debug().Msgf("PlayRecap: %s", string(prBytes))

	tc := gar.Count(*pr)
	logger.Info().Msgf("Changed: %d, Failures: %d, Ignored: %d, Ok: %d, Processed: %d, Rescued: %d, Skipped: %d.",
		tc.Changed, tc.Failures, tc.Ignored, tc.Ok, tc.Processed, tc.Rescued, tc.Skipped)

	return string(output), err
}
