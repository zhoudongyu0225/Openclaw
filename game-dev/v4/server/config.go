package config

type Config struct {
    ServerPort   string
    RedisAddr   string
    MongoAddr   string
    LogLevel    string
    EnableDebug bool
}

func Load() *Config {
    return &Config{
        ServerPort:  getEnv("PORT", "8080"),
        RedisAddr:   getEnv("REDIS", "localhost:6379"),
        MongoAddr:   getEnv("MONGO", "localhost:27017"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
        EnableDebug: getEnv("DEBUG", "false") == "true",
    }
}

func getEnv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}
