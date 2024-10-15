package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

type AlibabaCanal struct {
	Address     string `yaml:"address"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Destination string `yaml:"destination"`

	Database string `yaml:"database"`
}
type Config struct {
	EsAddress    []string     `yaml:"esAddress"`
	AlibabaCanal AlibabaCanal `yaml:"alibabaCanal"`
	//Ab           string       `yaml:"ab"`
}

var AppConfig Config

func Ab() {

}
func init() {
	abs, _ := filepath.Abs("./")
	fmt.Println(abs)
	file, _ := ioutil.ReadFile(abs + "/.env.yaml")
	yaml.Unmarshal(file, &AppConfig)
	fmt.Println(AppConfig)
}
