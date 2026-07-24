package config

import(
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port string 
	RedisAddr string
	RedisPasseord string
	RedisDB int
	LockTTL time.Duration
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", ":8080"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPasseord: getEnv("REDIS_PASSWORD", ""),
		RedisDB: getEnvInt("REDIS_DB", 0),
		LockTTL: getEnvDuration("LOCK_TTL", 60*time.Second),
	}
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed;
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value , exists := os.LookupEnv(key); exists {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed;
		}
	}
	return fallback
}