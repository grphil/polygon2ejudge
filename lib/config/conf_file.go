package config

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

type ConfFile struct {
	ApiKey         string `yaml:"api_key"`
	ApiSecret      string `yaml:"api_secret"`
	EjudgeLogin    string `yaml:"ejudge_login"`
	EjudgePassword string `yaml:"ejudge_password"`
}

func GetConfigFile(resetCredentials bool) (*ConfFile, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	confPath := filepath.Join(homePath, ".config", "polygon2ejudge", "conf.yaml")

	if resetCredentials {
		return createConfFile(confPath)
	}

	_, err = os.Stat(confPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return createConfFile(confPath)
		} else {
			return nil, err
		}
	}

	confFileData, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	confFile := &ConfFile{}
	err = yaml.Unmarshal(confFileData, confFile)
	if err != nil {
		return nil, err
	}
	return confFile, nil
}

func createConfFile(confPath string) (*ConfFile, error) {
	err := os.MkdirAll(filepath.Dir(confPath), 0700)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(os.Stdin)

	confFile := &ConfFile{}

	fmt.Println("Enter polygon api key:")
	confFile.ApiKey, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	confFile.ApiKey = strings.TrimSpace(confFile.ApiKey)

	fmt.Println("Enter polygon api secret:")
	apiSecretB, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	confFile.ApiSecret = string(apiSecretB)

	fmt.Println("Enter ejudge login:")
	confFile.EjudgeLogin, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	confFile.EjudgeLogin = strings.TrimSpace(confFile.EjudgeLogin)

	fmt.Println("Enter polygon api secret:")
	ejudgePasswordB, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	confFile.EjudgePassword = string(ejudgePasswordB)

	outputFile, err := os.OpenFile(confPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	encoder := yaml.NewEncoder(outputFile)
	err = encoder.Encode(confFile)
	if err != nil {
		return nil, err
	} else {
		return confFile, nil
	}
}
