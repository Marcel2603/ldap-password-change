package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port   string
	Domain string
}

func Get() Config {
	portEnv, exists := os.LookupEnv("PORT")
	if !exists {
		portEnv = "3333"
	}
	_, err := strconv.Atoi(portEnv)
	fmt.Println(os.Getenv("DOMAIN"))
	if err != nil {
		panic(err)
	}
	return Config{
		Port:   portEnv,
		Domain: os.Getenv("DOMAIN"),
	}
}
