package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Server struct {
		Port int `toml:"port"`
	} `toml:"server"`
	Database struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		DBName   string `toml:"dbname"`
		SSLMode  string `toml:"sslmode"`
	} `toml:"database"`
	AMQP struct {
		User     string `toml:"user"`
		Password string `toml:"password"`
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
	} `toml:"amqp"`
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

func loadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, err := loadConfig("config.toml")
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	fmt.Printf("Server Port: %d\n", config.Server.Port)
	fmt.Printf("Database Host: %s\n", config.Database.Host)
	fmt.Printf("Database Port: %d\n", config.Database.Port)
	fmt.Printf("Database User: %s\n", config.Database.User)
	fmt.Printf("Database Password: %s\n", config.Database.Password)
	fmt.Printf("Database Name: %s\n", config.Database.DBName)
	fmt.Printf("Database SSLMode: %s\n", config.Database.SSLMode)
	fmt.Printf("AMQP User: %s\n", config.AMQP.User)
	fmt.Printf("AMQP Password: %s\n", config.AMQP.Password)
	fmt.Printf("AMQP Host: %s\n", config.AMQP.Host)
	fmt.Printf("AMQP Port: %d\n", config.AMQP.Port)
	fmt.Printf("Redis Host: %s\n", config.Redis.Host)
	fmt.Printf("Redis Port: %d\n", config.Redis.Port)
	fmt.Printf("S3 Endpoint: %s\n", config.S3.Endpoint)
	fmt.Printf("S3 Bucket: %s\n", config.S3.Bucket)
	fmt.Printf("S3 Key: %s\n", config.S3.Key)
	fmt.Printf("S3 Secret: %s\n", config.S3.Secret)
	fmt.Printf("Director Endpoint: %s\n", config.Director.Endpoint)
	fmt.Printf("Director Auth Key: %s\n", config.Director.AuthKey)
}
