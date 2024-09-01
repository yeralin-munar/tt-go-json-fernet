package config

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/spf13/viper"
)

//go:embed config.yaml
var embeddedConfigBytes []byte
var Cfg Config

type Config struct {
	Server Server `mapstructure:"server"`
	Data   Data   `mapstructure:"data"`
}

type Server struct {
	Mode string `mapstructure:"mode"`
	HTTP HTTP   `mapstructure:"http"`
}

type HTTP struct {
	Addr string `mapstructure:"addr"`
}

type Data struct {
	Folder   string   `mapstructure:"folder"`
	Crypto   Crypto   `mapstructure:"crypto"`
	Database Database `mapstructure:"database"`
}

type Crypto struct {
	Key string `mapstructure:"key"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Sslmode  string `mapstructure:"sslmode"`
	Name     string `mapstructure:"name"`
}

func NewConfig(conf string) {
	v := viper.New()

	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if len(embeddedConfigBytes) > 0 {
		if err := v.ReadConfig(bytes.NewReader(embeddedConfigBytes)); err != nil {
			panic(err)
		}
	}

	var appConfig Config
	if err := v.Unmarshal(&appConfig); err != nil {
		panic(err)
	}

	// return appConfig
	Cfg = appConfig

}
