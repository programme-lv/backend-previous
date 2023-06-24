package environment

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/viper"
)

type EnvConfig struct {
	SqlxConnString string `mapstructure:"SQLX_CONN_STRING"`
}

func ReadEnvConfig() *EnvConfig {
	config := getDefaultConfig()

	viper.SetConfigName(".env")
	viper.SetConfigType("env") // Ideally, this should be determined based on the config file extension
	viper.AddConfigPath(".")   // Look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found")
		} else {
			log.Fatalf("Fatal error config file: %s \n", err)
		}
	}

	viper.AutomaticEnv()
	viper.BindEnv("SQLX_CONN_STRING")

	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return config
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
		LENGTH_LIMIT := 30
		if len(strRepresentation) > LENGTH_LIMIT {
			strRepresentation = strRepresentation[:LENGTH_LIMIT] + "..."
		}
		t.AppendRow(table.Row{typeOfS.Field(i).Name, strRepresentation})
	}
	t.Render()
}
