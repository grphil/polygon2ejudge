package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"polygon2ejudge/lib/console"
)

const UserConfigVersion = "user.1"

type TUserConfig struct {
	Version string `yaml:"version"`

	ApiKey    string `yaml:"api_key"`
	ApiSecret string `yaml:"api_secret"`

	EjudgeLogin    string `yaml:"ejudge_login"`
	EjudgePassword string `yaml:"ejudge_password"`

	NolintString string `yaml:"nolint_string"`
}

var UserConfig *TUserConfig

func init() {
	confPath := getUserConfigPath()
	_, err := os.Stat(confPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ResetUserConfig()
		} else {
			panic(err)
		}
	}

	confFileData, err := os.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	userConfig := &TUserConfig{}
	err = yaml.Unmarshal(confFileData, userConfig)
	if err != nil {
		ResetUserConfig()
	}

	if userConfig.Version != UserConfigVersion {
		ResetUserConfig()
	}

	UserConfig = userConfig
}

func ResetUserConfig() {
	err := os.MkdirAll(getConfigDir(), 0700)
	if err != nil {
		panic(err)
	}

	UserConfig = &TUserConfig{
		Version: UserConfigVersion,
	}

	UserConfig.ApiKey = console.ReadValue("Enter polygon api key:")
	UserConfig.ApiSecret = console.ReadSecret("Enter polygon api secret:")
	UserConfig.EjudgeLogin = console.ReadValue("Enter ejudge login:")
	UserConfig.EjudgePassword = console.ReadValue("Enter ejudge password:")
	UserConfig.NolintString = console.ReadValue("Enter nolint string (will be added before all submitted solutions). Leave empty for no nolint string")

	outputFile, err := os.OpenFile(getUserConfigPath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	encoder := yaml.NewEncoder(outputFile)
	err = encoder.Encode(UserConfig)
	if err != nil {
		panic(err)
	}
}

func getConfigDir() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homePath, ".config", "polygon2ejudge")
}

func getUserConfigPath() string {
	return filepath.Join(getConfigDir(), "user_config.yaml")
}
