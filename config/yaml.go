package config

import (
	"fmt"

	"github.com/toolkits/pkg/file"
)

type Config struct {
	Logger   loggerSection   `yaml:"logger"`
	Smtp     smtpSection     `yaml:"smtp"`
	Consumer consumerSection `yaml:"queue"`
	Redis    redisSection    `yaml:"redis"`
}

type loggerSection struct {
	Dir       string `yaml:"dir"`
	Level     string `yaml:"level"`
	KeepHours uint   `yaml:"keepHours"`
}

type redisSection struct {
	Addr    string         `yaml:"addr"`
	Pass    string         `yaml:"pass"`
	Idle    int            `yaml:"idle"`
	Timeout timeoutSection `yaml:"timeout"`
}

type timeoutSection struct {
	Conn  int `yaml:"conn"`
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
}

type consumerSection struct {
	Queue  string `yaml:"queue"`
	Worker int    `yaml:"worker"`
}

type smtpSection struct {
	FromName   string `yaml:"from_name"`
	FromMail   string `yaml:"from_mail"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	ServerHost string `yaml:"server_host"`
	ServerPort int    `yaml:"server_port"`
	UseSSL     bool   `yaml:"use_ssl"`
	StartTLS   bool   `yaml:"start_tls"`
}

var yaml Config

func Get() Config {
	return yaml
}

func ParseConfig(yf string) error {
	err := file.ReadYaml(yf, &yaml)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", yf, err)
	}
	return nil
}
