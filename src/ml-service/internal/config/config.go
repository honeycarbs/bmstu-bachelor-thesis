package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"ml/pkg/logging"
	"sync"
)

type Config struct {
	Listen struct {
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	DB struct {
		Host     string `yaml:"host" env-default:"localhost"`
		Port     string `yaml:"port" env-default:"5432"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"db"`
	ClusterAmount int    `yaml:"cluster_amount" env-default:"500"`
	FilePath      string `yaml:"model_file_path" env-default:"_meta/model-%s.json"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("reading application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("_meta/config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
