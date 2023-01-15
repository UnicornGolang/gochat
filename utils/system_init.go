package utils

import (
	"fmt"
	"github.com/spf13/viper"
)


func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
  fmt.Println("##################################################")
	fmt.Println("###-----------   load app config     ----------###")
  fmt.Println("##################################################")
  fmt.Println("")
}


