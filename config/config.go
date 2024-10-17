package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

type AlibabaCanal struct {
	DbHost      string `yaml:"dbHost"`
	DbPort      string `yaml:"dbPort"`
	Address     string `yaml:"address"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Destination string `yaml:"destination"`

	Database string `yaml:"database"`
}
type IndexCondition struct {
	Day     string `yaml:"day"`
	DocsNum string `yaml:"docsNum"`
	MaxSize string `yaml:"maxSize"`
}
type Config struct {
	DatabaseName   string         `yaml:"databaseName"`
	EsAddress      []string       `yaml:"esAddress"`
	AlibabaCanal   AlibabaCanal   `yaml:"alibabaCanal"`
	IndexCondition IndexCondition `yaml:"indexCondition"`
}

var AppConfig Config

func InitConfig(RootPath string) {
	abs, _ := filepath.Abs(RootPath)
	envYamlFilePath := abs + "/.env.yaml"
	fmt.Println(envYamlFilePath)
	file, _ := ioutil.ReadFile(envYamlFilePath)
	err := yaml.Unmarshal(file, &AppConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(AppConfig)
}
