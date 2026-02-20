package config

import "os"

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// ExoplanetCalculationConfig содержит конфигурацию системы расчета экзопланет
type ExoplanetCalculationConfig struct {
	// PostgreSQL
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	
	// MinIO для хранения изображений телескопов
	MinIOEndpoint        string
	MinIOAccessKeyID     string
	MinIOSecretAccessKey string
	MinIOUseSSL          bool
	MinIOBucketName      string
	
	// Server
	ServerPort string

	// Lab 8: асинхронный сервис
	AsyncServiceURL string
	ServiceToken    string // 8 байт, для приёма результатов от async
}

// GetExoplanetCalculationConfig возвращает конфигурацию по умолчанию
func GetExoplanetCalculationConfig() *ExoplanetCalculationConfig {
	// Проверяем, запущено ли приложение в Docker
	isDocker := os.Getenv("DOCKER_ENV") == "true"
	
	config := &ExoplanetCalculationConfig{
		// PostgreSQL
		DBUser:     "postgres",
		DBPassword: "postgres123",
		DBName:     "exoplanet_calculations",
		DBPort:     5432,
		
		// MinIO для изображений телескопов
		MinIOAccessKeyID:     "minioadmin123",
		MinIOSecretAccessKey: "minioadmin123456",
		MinIOUseSSL:          false,
		MinIOBucketName:      "telescope-images",
		
		// Server
		ServerPort: "8081",

		// Lab 8: async service
		AsyncServiceURL: getEnv("ASYNC_SERVICE_URL", "http://localhost:8001"),
		ServiceToken:    getEnv("SERVICE_TOKEN", "Lab8Token"),
	}
	
	if isDocker {
		// Конфигурация для Docker
		config.DBHost = "postgres"
		config.MinIOEndpoint = "minio:9000"
	} else {
		// Конфигурация для локального запуска
		config.DBHost = "localhost"
		config.MinIOEndpoint = "localhost:9000"
	}
	
	return config
}
