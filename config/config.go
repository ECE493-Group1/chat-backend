package config

type Config struct {
	Port            string
	CORSOrigin      string
	KeywordEndpoint string
	KeywordUpdates  bool
}

func GetConfig() *Config {
	return &Config{
		Port:            ":8000",
		CORSOrigin:      "http://localhost:8080",
		KeywordEndpoint: "http://localhost:5000/keywords/message",
		KeywordUpdates:  true,
	}
}
