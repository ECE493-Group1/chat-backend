package config

import (
	"os"
	"fmt"
)

type Config struct {
	Port            string
	CORSOrigin      string
	KeywordEndpoint string
	KeywordUpdates  bool
}

func GetConfig() *Config {
	kw_host, exists := os.LookupEnv("KEYWORD_HOST")
	if !exists {
		kw_host = "localhost"
	}
	return &Config{
		Port:            ":8000",
		CORSOrigin:      "http://localhost:8080",
		KeywordEndpoint: fmt.Sprintf("http://%s:5000/keywords/message", kw_host),
		KeywordUpdates:  true,
	}
}
