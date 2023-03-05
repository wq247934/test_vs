package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func main() {
	config := viper.New()
	config.SetConfigFile("./config.yaml")
	config.Set("key.key", "1111111")
	if err := config.ReadInConfig(); err != nil {
		fmt.Println(os.IsNotExist(err))
		panic(err)
	}

	fmt.Println(config.GetString("luosimao.user"))
	fmt.Println(config.GetString("key.key"))
	config.WriteConfig()
}
