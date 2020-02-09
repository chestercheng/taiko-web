package config

import (
	"github.com/spf13/viper"
	"log"
)

type Version struct {
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	CommitShort string `json:"commit_short"`
	Url         string `json:"url"`
}

type WebConfig struct {
	AssetsBaseurl string  `json:"assets_baseurl"`
	SongsBaseurl  string  `json:"songs_baseurl"`
	Version       Version `json:"_version"`
}

type DBConfig struct {
	Dialect, ConnString string
}

type TaikoWebConfig struct {
	Mode, Port string
	Web        WebConfig
	DB         DBConfig
}

var config TaikoWebConfig

func Init(mode string) {
	var err error
	configFile := viper.New()
	configFile.SetConfigType("yaml")
	configFile.SetConfigName(mode)
	configFile.AddConfigPath("../config/")
	configFile.AddConfigPath("config/")
	err = configFile.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}

	config.Mode = mode
	config.Port = configFile.GetString("port")
	config.Web = WebConfig{
		configFile.GetString("web.assets_baseurl"),
		configFile.GetString("web.songs_baseurl"),
		Version{Url: configFile.GetString("url")},
	}
	config.DB = DBConfig{
		configFile.GetString("db.dialect"),
		configFile.GetString("db.conn_string"),
	}
}

func GetConfig() *TaikoWebConfig {
	return &config
}
