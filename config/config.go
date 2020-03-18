package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Debug bool

var Web struct {
	BindIP string
	Port   int
}

var Database struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

var Plugins struct {
	Path string
}

func init() {
	viper.SetDefault("debug", false)

	viper.SetDefault("web.port", 8080)
	viper.SetDefault("web.bind_ip", "0.0.0.0")

	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "plumbie")
	viper.SetDefault("database.user", "plumbie")
	viper.SetDefault("database.password", "")
}

func LoadConfig() error {
	configPath := "./config.toml"
	args := os.Args[1:]
	if len(args) > 0 {
		configPath = args[0]
	}
	viper.SetConfigFile(configPath)
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
	if err := loadPluginsConfig(); err != nil {
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
	Database.Host = viper.GetString("database.host")
	Database.Port = viper.GetInt("database.port")
	Database.Name = viper.GetString("database.name")
	Database.User = viper.GetString("database.user")
	Database.Password = viper.GetString("database.password")
	return nil
}

func loadPluginsConfig() error {
	Plugins.Path = viper.GetString("plugins.path")
	fmt.Println("path", Plugins.Path)
	return nil
}
