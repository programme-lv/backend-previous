package environment

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type EnvConfig struct {
	SqlxConnString   string `mapstructure:"SQLX_CONN_STRING"`
	AMQPConnString   string `mapstructure:"AMQP_CONN_STRING"`
	RedisConnString  string `mapstructure:"REDIS_CONN_STRING"`
	S3Endpoint       string `mapstructure:"S3_ENDPOINT"`
	S3Bucket         string `mapstructure:"S3_BUCKET"`
	DOSpacesKey      string `mapstructure:"DO_SPACES_KEY"`
	DOSpacesSecret   string `mapstructure:"DO_SPACES_SECRET"`
	DirectorEndpoint string `mapstructure:"DIRECTOR_ENDPOINT"`
	DirectorAuthKey  string `mapstructure:"DIRECTOR_AUTH_KEY"`
}

func ReadEnvConfig(log *slog.Logger) *EnvConfig {
	config := getDefaultConfig()

	viper.SetConfigName(".env")
	viper.SetConfigType("env") // Ideally, this should be determined based on the config file extension
	viper.AddConfigPath(".")   // Look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("No config file found")
		} else {
			log.Error("Error reading config file", err)
		}
	}

	viper.AutomaticEnv()
	viper.BindEnv("SQLX_CONN_STRING")
	viper.BindEnv("AMQP_CONN_STRING")
	viper.BindEnv("REDIS_CONN_STRING")
	viper.BindEnv("S3_ENDPOINT")
	viper.BindEnv("S3_BUCKET")
	viper.BindEnv("DO_SPACES_KEY")
	viper.BindEnv("DO_SPACES_SECRET")
	viper.BindEnv("DIRECTOR_ENDPOINT")
	viper.BindEnv("DIRECTOR_AUTH_KEY")

	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func ConnectToPostgresByEnvConf() (*sqlx.DB, error) {
	config := ReadEnvConfig(slog.Default())
	return sqlx.Connect("postgres", config.SqlxConnString)
}

func getDefaultConfig() *EnvConfig {
	return &EnvConfig{
		SqlxConnString: "host=localhost port=5432 user=proglv password=proglv dbname=proglv sslmode=disable",
	}
}

func (c *EnvConfig) Print() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.SetTitle("Environment Configuration")
	t.AppendHeader(table.Row{"field name", "value"})
	v := reflect.ValueOf(*c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		strRepresentation := fmt.Sprintf("%v", v.Field(i).Interface())
		LENGTH_LIMIT := 50
		if len(strRepresentation) > LENGTH_LIMIT {
			strRepresentation = strRepresentation[:LENGTH_LIMIT] + "..."
		}
		t.AppendRow(table.Row{typeOfS.Field(i).Name, strRepresentation})
	}
	t.Render()
}
