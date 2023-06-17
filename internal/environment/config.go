package environment

import (
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
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.BindEnv("SQLX_CONN_STRING")

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func getDefaultConfig() *EnvConfig {
	return &EnvConfig{
		SqlxConnString: "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
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
		t.AppendRow(table.Row{typeOfS.Field(i).Name, v.Field(i).Interface()})
	}
	t.Render()
}
