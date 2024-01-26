package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"polygon2ejudge/lib/console"
)

const GlobalConfigVersion = "global.1"

type TGlobalConfig struct {
	Version string `yaml:"version"`

	JudgesDir   string `yaml:"judges_dir"`
	GvaluerPath string `yaml:"gvaluer_path"`

	EjudgeUrl  string `yaml:"ejudge_url"`
	PolygonUrl string `yaml:"polygon_url"`

	LangIds map[string][]int `yaml:"lang_ids"`
}

var GlobalConfig *TGlobalConfig

var defaultGlobalConfig = &TGlobalConfig{
	JudgesDir:   "/home/judges",
	GvaluerPath: "/home/judges/001501/problems/gvaluer",
	EjudgeUrl:   "https://ejudge.algocode.ru/",
	PolygonUrl:  "https://polygon.codeforces.com/api/",

	LangIds: map[string][]int{
		"py":   {23, 64}, // python3 and pypy3
		"cpp":  {3, 52},  // g++ and clang
		"java": {18},     // java
		"pas":  {1},      // Free Pascal
	},
}

func init() {
	if TryLoadGlobalConfig() {
		return
	}

	GlobalConfig = defaultGlobalConfig
}

func TryLoadGlobalConfig() bool {
	path := getGlobalConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	globalConfig := &TGlobalConfig{}
	err = yaml.Unmarshal(data, globalConfig)
	if err != nil {
		ResetGlobalConfig()
	}

	if globalConfig.Version != GlobalConfigVersion {
		ResetGlobalConfig()
	}

	GlobalConfig = globalConfig

	fmt.Printf("Initialized global config from %s\n", path)
	return true
}

func ResetGlobalConfig() {
	err := os.MkdirAll(getConfigDir(), 0700)
	if err != nil {
		panic(err)
	}

	GlobalConfig = &TGlobalConfig{
		Version: GlobalConfigVersion,
		LangIds: make(map[string][]int),
	}

	GlobalConfig.JudgesDir = console.ReadValue(fmt.Sprintf("Enter judges dir (%s)", defaultGlobalConfig.JudgesDir))
	GlobalConfig.GvaluerPath = console.ReadValue(fmt.Sprintf("Enter gvaluer path (%s)", defaultGlobalConfig.GvaluerPath))
	GlobalConfig.EjudgeUrl = console.ReadValue(fmt.Sprintf("Enter ejudge url (%s)", defaultGlobalConfig.EjudgeUrl))
	GlobalConfig.PolygonUrl = console.ReadValue(fmt.Sprintf("Enter polygon url (%s)", defaultGlobalConfig.PolygonUrl))

	for lang, ids := range defaultGlobalConfig.LangIds {
		GlobalConfig.LangIds[lang] = console.ReadInts(fmt.Sprintf("Enter space ceparated lang ids for %s %v", lang, ids))
	}

	outputFile, err := os.OpenFile(getGlobalConfigPath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	encoder := yaml.NewEncoder(outputFile)
	err = encoder.Encode(GlobalConfig)
	if err != nil {
		panic(err)
	}
}

func getGlobalConfigPath() string {
	return filepath.Join(getConfigDir(), "polygon2ejudge.yaml")
}
