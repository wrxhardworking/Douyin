package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Mysql    Mysql    `yaml:"mysql"`
	Redis    Redis    `yaml:"redis"`
	Resouece Resouece `yaml:"resouece"`
}

type Mysql struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Ipaddress string `yaml:"ipaddress"`
	Port      string `yaml:"port"`
	Dbname    string `yaml:"dbname"`
}

type Redis struct {
	Ipaddress    string `yaml:"ipaddress"`
	Port         string `yaml:"port"`
	Authpassword string `yaml:"authpassword"`
	Maxidle      int    `yaml:"maxidle"`
	Maxactive    int    `yaml:"maxactive"`
}

type Resouece struct {
	Ipaddress string `yaml:"ipaddress"`
	Port      string `yaml:"port"`
}

var C Config

func ConfInit() error {
	yamlFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// 将读取的yaml文件解析为响应的 struct
	err = yaml.Unmarshal(yamlFile, &C)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
