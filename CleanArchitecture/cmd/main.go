package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"clean-url-shortener/configs"
	"clean-url-shortener/pkg/di"
)

// main - точка входа в приложение
// ЭТО ГЛАВНАЯ ФУНКЦИЯ, КОТОРАЯ ЗАПУСКАЕТ ВСЕ ПРИЛОЖЕНИЕ
func main() {
	// ШАГ 1: Загружаем конфигурацию
	fmt.Println("🚀 Starting URL Shortener Service...")
	
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}
	
	fmt.Printf("📋 Configuration loaded for environment: %s\n", config.App.Environment)
	
	// ШАГ 2: Создаем контейнер зависимостей
	fmt.Println("🔧 Initializing dependencies...")
	
	container, err := di.NewContainer(config)
	if err != nil {
		log.Fatalf("❌ Failed to initialize dependencies: %v", err)
	}
	
	// Обеспечиваем очистку ресурсов при завершении
	defer func() {
		fmt.Println("🧹 Cleaning up resources...")
		if err := container.Cleanup(); err != nil {
			log.Printf("⚠️ Error during cleanup: %v", err)
		}
	}()
	
	fmt.Println("✅ Dependencies initialized successfully")
	
	// ШАГ 3: Настраиваем graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Канал для получения сигналов ОС
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	
	// ШАГ 4: Запускаем HTTP сервер в отдельной горутине
	serverErrors := make(chan error, 1)
	
	go func() {
		fmt.Printf("🌐 HTTP server starting on %s:%d\n", 
			config.Server.Host, config.Server.Port)
		
		if err := container.Server.Start(); err != nil {
			serverErrors <- fmt.Errorf("server failed to start: %w", err)
		}
	}()
	
	// ШАГ 5: Ждем сигнала завершения или ошибки
	select {
	case err := <-serverErrors:
		log.Fatalf("❌ Server error: %v", err)
		
	case sig := <-signalChan:
		fmt.Printf("\n🛑 Received signal: %s\n", sig)
		fmt.Println("🔄 Initiating graceful shutdown...")
		
		// Создаем контекст с таймаутом для graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(), 
			30*time.Second,
		)
		defer shutdownCancel()
		
		// Останавливаем HTTP сервер
		if err := container.Server.Shutdown(shutdownCtx); err != nil {
			log.Printf("⚠️ Error during server shutdown: %v", err)
		} else {
			fmt.Println("✅ HTTP server stopped gracefully")
		}
		
		// Отменяем контекст приложения
		cancel()
	}
	
	fmt.Println("👋 Application terminated")
}

// ПРИНЦИПЫ ГЛАВНОЙ ФУНКЦИИ:
// 1. Загрузка конфигурации на старте
// 2. Инициализация всех зависимостей через DI container
// 3. Graceful shutdown при получении сигналов ОС
// 4. Proper error handling и логирование
// 5. Cleanup ресурсов при завершении
// 6. Информативные сообщения о состоянии приложения

// ЧТО ПРОИСХОДИТ ПРИ ЗАПУСКЕ:
// 1. Загружается конфигурация из переменных окружения
// 2. Создается DI container, который:
//    - Подключается к базе данных
//    - Выполняет миграции
//    - Создает все репозитории
//    - Инициализирует Use Cases
//    - Настраивает контроллеры и middleware
//    - Создает HTTP сервер
// 3. Запускается HTTP сервер
// 4. Приложение ждет сигналов завершения
// 5. При получении сигнала выполняется graceful shutdown

// АРХИТЕКТУРНЫЕ ПРЕИМУЩЕСТВА:
// 1. Все зависимости явно объявлены
// 2. Легко тестировать (можно заменить любой компонент)
// 3. Соблюдены принципы SOLID
// 4. Инверсия зависимостей на всех уровнях
// 5. Четкое разделение слоев и ответственности