module pipeline-clean-architecture

go 1.21

require (
	// Web Framework
	github.com/gin-gonic/gin v1.9.1

	// Database
	github.com/lib/pq v1.10.9
	github.com/go-redis/redis/v8 v8.11.5

	// Configuration
	github.com/spf13/viper v1.18.2
	github.com/joho/godotenv v1.5.1

	// Logging
	go.uber.org/zap v1.26.0

	// Validation
	github.com/go-playground/validator/v10 v10.16.0

	// Background Jobs
	github.com/hibiken/asynq v0.24.1

	// HTTP Client
	github.com/go-resty/resty/v2 v2.11.0

	// Utilities
	github.com/google/uuid v1.5.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/time v0.5.0
)