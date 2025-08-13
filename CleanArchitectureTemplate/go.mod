// 🏗️ УНИВЕРСАЛЬНЫЙ ШАБЛОН GO.MOD для Clean Architecture
// 
// ✅ ВСЕГДА одинаково:
// - Название модуля нужно поменять на ваш проект
// - Версия Go обычно последняя стабильная
// - Базовые зависимости остаются теми же
//
// 🔀 МЕНЯЕТСЯ в зависимости от проекта:
// - Название модуля (your-project-name)
// - Дополнительные библиотеки (redis, kafka, grpc, etc.)
// - Версии зависимостей

module your-project-name  // 🔀 ИЗМЕНИТЕ: название вашего проекта

go 1.21  // 🔀 ИЗМЕНИТЕ: актуальная версия Go

require (
	// ✅ БАЗОВЫЕ ЗАВИСИМОСТИ (почти всегда нужны):
	
	// Валидация входных данных
	github.com/go-playground/validator/v10 v10.16.0
	
	// Работа с JWT токенами (для аутентификации)
	github.com/golang-jwt/jwt/v5 v5.2.0
	
	// Хеширование паролей
	golang.org/x/crypto v0.17.0
	
	// 🔀 БАЗЫ ДАННЫХ (выберите нужную):
	
	// PostgreSQL
	github.com/lib/pq v1.10.9
	
	// MySQL (раскомментируйте если нужен)
	// github.com/go-sql-driver/mysql v1.7.1
	
	// MongoDB (раскомментируйте если нужен)
	// go.mongodb.org/mongo-driver v1.13.1
	
	// Redis (раскомментируйте если нужен)
	// github.com/redis/go-redis/v9 v9.3.0
	
	// 🔀 WEB ФРЕЙМВОРКИ (выберите нужный):
	
	// Стандартный net/http (уже в стандартной библиотеке)
	
	// Gin Web Framework (раскомментируйте если предпочитаете Gin)
	// github.com/gin-gonic/gin v1.9.1
	
	// Fiber (раскомментируйте если предпочитаете Fiber)
	// github.com/gofiber/fiber/v2 v2.52.0
	
	// 🔀 ДОПОЛНИТЕЛЬНЫЕ СЕРВИСЫ (добавьте при необходимости):
	
	// Email отправка
	// github.com/go-mail/mail v2.3.1+incompatible
	
	// AWS SDK
	// github.com/aws/aws-sdk-go-v2 v1.24.0
	
	// Kafka
	// github.com/segmentio/kafka-go v0.4.46
	
	// gRPC
	// google.golang.org/grpc v1.60.1
	// google.golang.org/protobuf v1.31.0
	
	// Логирование
	// github.com/sirupsen/logrus v1.9.3
	// go.uber.org/zap v1.26.0
	
	// Конфигурация из ENV
	// github.com/joho/godotenv v1.4.0
)

// 📝 ИНСТРУКЦИЯ ПО ИСПОЛЬЗОВАНИЮ:
//
// 1. Переименуйте модуль:
//    Замените "your-project-name" на название вашего проекта
//
// 2. Выберите базу данных:
//    Раскомментируйте нужную БД, закомментируйте ненужные
//
// 3. Выберите веб-фреймворк:
//    Стандартный net/http уже включен, или раскомментируйте Gin/Fiber
//
// 4. Добавьте дополнительные зависимости:
//    Раскомментируйте нужные сервисы (email, AWS, kafka, etc.)
//
// 5. Обновите зависимости:
//    go mod tidy
//    go mod download