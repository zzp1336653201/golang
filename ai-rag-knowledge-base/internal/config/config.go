package config

import "github.com/spf13/viper"

func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// 默认配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("ollama.endpoint", "http://localhost:11434")
	viper.SetDefault("ollama.model", "llama3.2")
	viper.SetDefault("chroma.endpoint", "http://localhost:8000")

	return viper.ReadInConfig()
}
