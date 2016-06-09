package common

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func ReadConfig(path string) (Config, error) {

	var config Config

	if path == "" {

		config.setFromEnv()

		return config, nil
	}

	if _, err := os.Open(path); err != nil {

		if os.IsNotExist(err) {

			return config, fmt.Errorf("No such configuration file '%s'", path)
		}

		return config, fmt.Errorf("Could not open configuration file '%s': %v", path, err)
	}

	data, err := ioutil.ReadFile(path)

	if err != nil {

		return config, nil
	}

	if err := yaml.Unmarshal(data, &config); err != nil {

		return config, err
	}

	return config, nil
}

type Config struct {
	AppAddress string  `yaml:"app_address"`
	Loglevel   string  `yaml:"loglevel"`
	Connect    connect `yaml:"connect"`
	SMTP       smtp    `yaml:"smtp"`
	Telegram   struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
}

func (c *Config) setFromEnv() {

	var (
		dbPort   uint16 = 5432
		smptPort uint16 = 25
	)

	if value, err := strconv.ParseUint(os.Getenv("DB_PORT"), 10, 32); err == nil {

		dbPort = uint16(value)
	}

	if value, err := strconv.ParseUint(os.Getenv("SMTP_PORT"), 10, 32); err == nil {

		smptPort = uint16(value)
	}
	c.AppAddress = os.Getenv("APP_ADDRESS")
	c.Loglevel = os.Getenv("LOG_LEVEL")
	c.Connect = connect{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
	}
	c.SMTP = smtp{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     smptPort,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	c.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")
}

func (c *Config) LogLevel() log.Level {

	switch c.Loglevel {
	case "info":
		return log.InfoLevel
	case "warning":
		return log.WarnLevel
	}

	return log.ErrorLevel
}

type connect struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (c *connect) DSN() string {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Username, c.Password, c.Database)

	if c.Port != 0 {

		dsn += fmt.Sprintf(" port=%d", c.Port)
	}

	return dsn
}

type smtp struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
