package config

import "os"

type Config struct {
    Port     string
    RedisURL string
    MongoURL string
}

func Load() *Config {
    return &Config{
        Port:     getEnv("PORT", "8080"),
        RedisURL: getEnv("REDIS_URL", "localhost:6379"),
        MongoURL: getEnv("MONGO_URL", "localhost:27017"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
