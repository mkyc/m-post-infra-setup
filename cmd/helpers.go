package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	hi "github.com/epiphany-platform/e-structures/hi/v0"
	st "github.com/epiphany-platform/e-structures/state/v0"
)

func ensureDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

//TODO handle new states validation somehow
func loadState(path string) (*st.State, error) {
	logger.Debug().Msgf("loadState(%s)", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return st.NewState(), nil
	} else {
		state := &st.State{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = state.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return state, nil
	}
}

func saveState(path string, state *st.State) error {
	bytes, err := state.Marshal()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig(path string) (*hi.Config, error) {
	logger.Debug().Msgf("loadConfig(%s)", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return hi.NewConfig(), nil
	} else {
		config := &hi.Config{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = config.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
}

func saveConfig(path string, config *hi.Config) error {
	bytes, err := config.Marshal()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func backupFile(path string) error {
	logger.Debug().Msgf("backupFile(%s)", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	} else {
		backupPath := path + ".backup"

		input, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(backupPath, input, 0644)
		if err != nil {
			return err
		}
		return nil
	}
}

func checkAndLoad(stateFilePath string, configFilePath string) (*hi.Config, *st.State, error) {
	logger.Debug().Msgf("checkAndLoad(%s, %s)", stateFilePath, configFilePath)
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("state file does not exist, please run init first")
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("config file does not exist, please run init first")
	}

	state, err := loadState(stateFilePath)
	if err != nil {
		return nil, nil, err
	}

	config, err := loadConfig(configFilePath)
	if err != nil {
		return nil, nil, err
	}

	return config, state, nil
}

func templateInventory(config *hi.Config) error {
	logger.Debug().Msg("templateInventory")
	inventoryFilePath := filepath.Join(ResourcesDirectory, inventoryDir, inventoryFile)
	err := ensureDirectory(filepath.Join(ResourcesDirectory, inventoryDir))
	if err != nil {
		return err
	}

	groups := make(map[string]interface{})
	for _, vmg := range config.Params.VmGroups {
		hosts := make(map[string]interface{})
		vars := make(map[string]interface{})
		for _, vm := range vmg.Hosts {
			hosts[*vm.Name] = map[string]string{"ansible_host": *vm.Ip}
		}
		vars["ansible_user"] = *vmg.AdminUser
		mountPoints := make([]map[string]string, 0, 0)
		for _, mp := range vmg.MountPoints {
			mountPoints = append(mountPoints, map[string]string{
				"lun":        strconv.Itoa(*mp.Lun),
				"mountpoint": *mp.Path,
			})
		}
		vars["mountpoints"] = mountPoints
		groups[*vmg.Name] = map[string]interface{}{"hosts": hosts, "vars": vars}
	}

	data := map[string]interface{}{
		"all": map[string]interface{}{
			"children": groups,
		},
	}

	bytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	logger.Info().Msg(string(bytes))
	err = ioutil.WriteFile(inventoryFilePath, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func prepareSshKey(config *hi.Config) error {
	logger.Debug().Msg("prepareSshKey")
	sshKeyFilePath := filepath.Join(ResourcesDirectory, envDir, sshKeyFile)
	err := ensureDirectory(filepath.Join(ResourcesDirectory, envDir))
	if err != nil {
		return err
	}
	input, err := ioutil.ReadFile(*config.Params.RsaPrivateKeyPath)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(sshKeyFilePath, input, 0600)
	if err != nil {
		return err
	}
	logger.Debug().Msgf("file %s copied to %s", *config.Params.RsaPrivateKeyPath, sshKeyFilePath)
	return nil
}

func setCheckAndDiff() error {
	logger.Debug().Msg("setCheckAndDiff")
	cmdlineFilePath := filepath.Join(ResourcesDirectory, envDir, cmdlineFile)
	err := ensureDirectory(filepath.Join(ResourcesDirectory, envDir))
	if err != nil {
		return err
	}
	content := []byte("--check --diff")
	err = ioutil.WriteFile(cmdlineFilePath, content, 0644)
	if err != nil {
		return err
	}
	logger.Debug().Msgf("file %s created with content: [%s]", cmdlineFilePath, string(content))
	return nil
}
