package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Debug bool

var Web struct {
	BindIP string
	Port   int
}

var Database struct {
	Driver   string
	Host     string
	Protocol string
	Name     string
	User     string
	Password string
	Charset  string
	SSL      bool
}

func init() {
	viper.SetDefault("debug", false)

	viper.SetDefault("web.port", 8080)
	viper.SetDefault("web.bind_ip", "0.0.0.0")

	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.protocol", "tcp")
	viper.SetDefault("database.name", "plumbie")
	viper.SetDefault("database.user", "plumbie")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.charset", "utf8")
	viper.SetDefault("database.ssl", false)
}

func LoadConfig() error {
	viper.SetConfigFile("./config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if err := loadMainConfig(); err != nil {
		return err
	}
	if err := loadWebConfig(); err != nil {
		return err
	}
	if err := loadDatabaseConfig(); err != nil {
		return err
	}
	return nil
}

func loadMainConfig() error {
	Debug = viper.GetBool("debug")
	return nil
}

func loadWebConfig() error {
	Web.BindIP = viper.GetString("web.bind_ip")
	Web.Port = viper.GetInt("web.port")
	return nil
}

func loadDatabaseConfig() error {
	Database.Driver = viper.GetString("database.driver")
	Database.Host = viper.GetString("database.host")
	Database.Protocol = viper.GetString("database.protocol")
	Database.Name = viper.GetString("database.name")
	Database.User = viper.GetString("database.user")
	Database.Password = viper.GetString("database.password")
	Database.Charset = viper.GetString("database.charset")
	Database.SSL = viper.GetBool("database.ssl")

	if Database.Driver != "mysql" {
		return fmt.Errorf("Unsupported database: %s", Database.Driver)
	}
	if Database.Driver != "tcp" || Database.Driver != "unix" {
		return fmt.Errorf("Unsupported protocol: %s. Only tcp and unix are supported.", Database.Driver)
	}
	if Database.Charset != "utf8" || Database.Charset != "utf8mb4" {
		return fmt.Errorf("Unsupported database charset: %s. Please use utf8 or utf8mb4", Database.Charset)
	}
	return nil
}
