// 🚀 УНИВЕРСАЛЬНЫЙ ШАБЛОН MAIN.GO для Clean Architecture
//
// ✅ ВСЕГДА одинаково:
// - Загрузка конфигурации
// - Создание DI контейнера  
// - Graceful shutdown
// - Общая структура main функции
//
// 🔀 МЕНЯЕТСЯ в зависимости от проекта:
// - Импорты (название модуля)
// - Дополнительная инициализация (логгеры, трейсинг, etc.)
// - Тип сервера (HTTP, gRPC, GraphQL)

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	// 🔀 ИЗМЕНИТЕ: импорты под ваш проект
	"your-project-name/configs"
	"your-project-name/pkg/di"

	// 🔀 ДОПОЛНИТЕЛЬНЫЕ ИМПОРТЫ (добавьте при необходимости):
	
	// Для логирования
	// "github.com/sirupsen/logrus"
	// "go.uber.org/zap"
	
	// Для трейсинга
	// "github.com/opentracing/opentracing-go"
	
	// Для метрик
	// "github.com/prometheus/client_golang/prometheus"
)

func main() {
	// 🔀 ИЗМЕНИТЕ: настройте логирование под ваши нужды
	fmt.Println("🚀 Starting Your Application...")  // 🔀 ИЗМЕНИТЕ: название приложения
	
	// 📝 ОПЦИОНАЛЬНЫЕ шаги инициализации (раскомментируйте при необходимости):
	
	// ✅ Инициализация логгера
	// initLogger()
	
	// ✅ Инициализация трейсинга
	// initTracing()
	
	// ✅ Инициализация метрик
	// initMetrics()

	// ШАГ 1: Загружаем конфигурацию
	// ✅ ВСЕГДА: эта часть не меняется
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	
	// 📝 ОПЦИОНАЛЬНО: валидация конфигурации
	// if err := config.Validate(); err != nil {
	//     log.Fatalf("❌ Invalid config: %v", err)
	// }

	// ШАГ 2: Создаем контейнер зависимостей
	// ✅ ВСЕГДА: эта часть не меняется
	container, err := di.NewContainer(config)
	if err != nil {
		log.Fatalf("❌ Failed to create DI container: %v", err)
	}
	defer container.Close() // ✅ ВСЕГДА: освобождаем ресурсы

	// ШАГ 3: Настраиваем graceful shutdown
	// ✅ ВСЕГДА: эта часть не меняется
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для ошибок сервера
	serverErrors := make(chan error, 1)

	// Канал для системных сигналов
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// ШАГ 4: Запускаем сервер(ы)
	// 🔀 МЕНЯЕТСЯ: тип сервера зависит от проекта
	
	// HTTP Server (самый частый случай)
	go func() {
		fmt.Printf("🌐 HTTP Server starting on :%s\n", config.Server.Port)
		if err := container.Server.Start(); err != nil {
			serverErrors <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()
	
	// 📝 ДОПОЛНИТЕЛЬНЫЕ СЕРВЕРЫ (раскомментируйте при необходимости):
	
	// gRPC Server
	// go func() {
	//     fmt.Printf("🔗 gRPC Server starting on :%s\n", config.GRPC.Port)
	//     if err := container.GRPCServer.Start(); err != nil {
	//         serverErrors <- fmt.Errorf("gRPC server error: %w", err)
	//     }
	// }()
	
	// GraphQL Server
	// go func() {
	//     fmt.Printf("📊 GraphQL Server starting on :%s\n", config.GraphQL.Port)
	//     if err := container.GraphQLServer.Start(); err != nil {
	//         serverErrors <- fmt.Errorf("GraphQL server error: %w", err)
	//     }
	// }()
	
	// Metrics Server (Prometheus)
	// go func() {
	//     fmt.Printf("📈 Metrics Server starting on :%s\n", config.Metrics.Port)
	//     if err := container.MetricsServer.Start(); err != nil {
	//         serverErrors <- fmt.Errorf("metrics server error: %w", err)
	//     }
	// }()

	// ШАГ 5: Ожидаем завершения
	// ✅ ВСЕГДА: эта часть не меняется
	select {
	case err := <-serverErrors:
		log.Printf("❌ Server error: %v", err)
		cancel()

	case sig := <-shutdown:
		log.Printf("🛑 Shutdown signal received: %v", sig)
		cancel()
	}

	// ШАГ 6: Graceful shutdown
	// ✅ ВСЕГДА: эта часть не меняется, но таймаут можете изменить
	
	// 🔀 ИЗМЕНИТЕ: таймаут можно настроить под ваши нужды
	shutdownTimeout := 30 * time.Second
	
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	fmt.Println("🔄 Shutting down gracefully...")

	// Останавливаем серверы
	if err := container.Server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️  Force shutdown: %v", err)
	}
	
	// 📝 ДОПОЛНИТЕЛЬНЫЕ SHUTDOWNS (раскомментируйте при необходимости):
	// if err := container.GRPCServer.Shutdown(shutdownCtx); err != nil {
	//     log.Printf("⚠️  gRPC force shutdown: %v", err)
	// }

	// Закрываем соединения с БД, кешем, etc.
	// ✅ ВСЕГДА: container.Close() вызывается в defer выше

	fmt.Println("✅ Application stopped gracefully")
}

// 📝 ДОПОЛНИТЕЛЬНЫЕ ФУНКЦИИ ИНИЦИАЛИЗАЦИИ
// Раскомментируйте и адаптируйте под ваш проект:

// func initLogger() {
//     // 🔀 ПРИМЕР: настройка Logrus
//     logrus.SetFormatter(&logrus.JSONFormatter{})
//     logrus.SetLevel(logrus.InfoLevel)
//     
//     // 🔀 ПРИМЕР: настройка Zap
//     // logger, _ := zap.NewProduction()
//     // defer logger.Sync()
//     // zap.ReplaceGlobals(logger)
// }

// func initTracing() {
//     // 🔀 ПРИМЕР: настройка Jaeger
//     // tracer, closer := jaeger.NewTracer(
//     //     "your-service-name",
//     //     jaeger.NewConstSampler(true),
//     //     jaeger.NewRemoteReporter(...),
//     // )
//     // opentracing.SetGlobalTracer(tracer)
//     // defer closer.Close()
// }

// func initMetrics() {
//     // 🔀 ПРИМЕР: настройка Prometheus
//     // prometheus.MustRegister(prometheus.NewGoCollector())
//     // prometheus.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
// }

// 📚 ИНСТРУКЦИЯ ПО АДАПТАЦИИ:
//
// 1. 🔀 ИЗМЕНИТЕ импорты:
//    Замените "your-project-name" на название вашего модуля
//
// 2. 🔀 НАСТРОЙТЕ логирование:
//    Раскомментируйте и настройте initLogger() под ваши нужды
//
// 3. 🔀 ВЫБЕРИТЕ тип сервера:
//    HTTP/gRPC/GraphQL - раскомментируйте нужные серверы
//
// 4. 🔀 ДОБАВЬТЕ дополнительные сервисы:
//    Метрики, трейсинг, health checks и т.д.
//
// 5. 🔀 НАСТРОЙТЕ таймауты:
//    Измените shutdownTimeout под ваши требования
//
// 6. ✅ ВСЕГДА сохраняйте:
//    - Порядок инициализации
//    - Graceful shutdown логику
//    - Обработку ошибок