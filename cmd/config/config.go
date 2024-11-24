package config

import (
  "os"
  "strconv"
)

type Config struct {
  Port string
  Host string
}

func Get() Config {
  portEnv, exists := os.LookupEnv("PORT")
  if !exists {
    portEnv = "3333"
  }
  _, err := strconv.Atoi(portEnv)
  if err != nil {
    panic(err)
  }
  return Config{
    Port: portEnv,
    Host: os.Getenv("HOST"),
  }
}
