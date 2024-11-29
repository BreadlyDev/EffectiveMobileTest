package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local" env-required:"true"`
	PageSize   int        `yaml:"page_size" env-default:"10"`
	Server     HTTPServer `yaml:"http_server" env-required:"true"`
	Storage    DBStorage  `yaml:"storage" env-required:"true"`
	Client     APIClient  `yaml:"api_client"`
	Migrations Migrations `yanl:"migrations"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost" env-required:"true"`
	Port        int           `yaml:"port" env-default:"5432" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DBStorage struct {
	Host    string `yaml:"host" env-default:"localhost" env-required:"true"`
	Port    int    `yaml:"port" env-default:"5432" env-required:"true"`
	DBName  string `yaml:"db_name" env-required:"true"`
	User    string `yaml:"user" env-required:"true"`
	Pass    string `yaml:"pass" env-required:"true"`
	SSLMode string `yaml:"ssl_mode" env-default:"disable" env-required:"true"`
}

type APIClient struct {
	Address  string `yaml:"address" env-default:"localhost:3000"`
	Protocol string `yaml:"protocol" env-default:"http"`
	Url      string `yaml:"url" env-default:"/info"`
}

type Migrations struct {
	Path  string `yaml:"path" env-default:"./migrations"`
	Table string `yaml:"table" env-default:"migrations"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
