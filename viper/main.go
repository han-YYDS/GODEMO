package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("app")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetDefault("redis.port", 6381)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed: %v", err)
	}
	fmt.Println(viper.Get("app_name"))
	fmt.Println(viper.Get("log_level"))
	fmt.Println("mysql ip: ", viper.Get("mysql.ip"))
	viper.Set("mysql.ip", "192.168")
	fmt.Println("mysql port: ", viper.Get("mysql.port"))
	fmt.Println("mysql user: ", viper.Get("mysql.user"))
	fmt.Println("mysql password: ", viper.Get("mysql.password"))
	fmt.Println("mysql database: ", viper.Get("mysql.database"))
	fmt.Println("redis ip: ", viper.Get("redis.ip"))
	fmt.Println("redis port: ", viper.Get("redis.port"))

	err = viper.WriteConfig()
	if err != nil {
		log.Fatal("write config failed: ", err)
	}
}
