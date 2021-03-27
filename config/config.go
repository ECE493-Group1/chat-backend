package config

type Config struct {
	Port       string
	CORSOrigin string
}

func GetConfig() *Config {
	return &Config{
		Port:       ":8000",
		CORSOrigin: "http://localhost:8080",
	}
}
