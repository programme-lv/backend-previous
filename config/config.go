package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server struct {
		Port int `toml:"port"`
	} `toml:"server"`
	Postgres struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		DBName   string `toml:"dbname"`
		SSLMode  string `toml:"sslmode"`
	} `toml:"database"`
	Redis struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	} `toml:"redis"`
	S3 struct {
		Endpoint string `toml:"endpoint"`
		Bucket   string `toml:"bucket"`
		Key      string `toml:"key"`
		Secret   string `toml:"secret"`
	} `toml:"s3"`
	Director struct {
		Endpoint string `toml:"endpoint"`
		AuthKey  string `toml:"auth_key"`
	} `toml:"director"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (config *Config) Print() {
	fmt.Printf("Server Port: %d\n", config.Server.Port)
	fmt.Printf("Postgres Host: %s\n", config.Postgres.Host)
	fmt.Printf("Postgres Port: %d\n", config.Postgres.Port)
	fmt.Printf("Postgres User: %s\n", config.Postgres.User)
	fmt.Printf("Postgres Password: %s\n", config.Postgres.Password)
	fmt.Printf("Postgres Name: %s\n", config.Postgres.DBName)
	fmt.Printf("Postgres SSLMode: %s\n", config.Postgres.SSLMode)
	fmt.Printf("Redis Host: %s\n", config.Redis.Host)
	fmt.Printf("Redis Port: %d\n", config.Redis.Port)
	fmt.Printf("S2 Endpoint: %s\n", config.S3.Endpoint)
	fmt.Printf("S2 Bucket: %s\n", config.S3.Bucket)
	fmt.Printf("S2 Key: %s\n", config.S3.Key)
	fmt.Printf("S2 Secret: %s\n", config.S3.Secret)
	fmt.Printf("Director Endpoint: %s\n", config.Director.Endpoint)
	fmt.Printf("Director Auth Key: %s\n", config.Director.AuthKey)
}
