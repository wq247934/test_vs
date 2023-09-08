package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	conf := viper.New()
	conf.SetConfigFile("/Users/wangqian/Documents/go-project/test_vs/viper_config/config.yaml")

	if err := conf.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(conf.GetString("expire_2time"))
}
