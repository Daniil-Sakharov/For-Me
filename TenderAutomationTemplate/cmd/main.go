package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tender-automation/configs"
	"tender-automation/pkg/di"
)

// TODO: Основная точка входа приложения
// Эта функция должна:
// 1. Загружать конфигурацию из переменных окружения
// 2. Инициализировать DI контейнер
// 3. Запускать все необходимые сервисы (HTTP сервер, планировщик, воркеры)
// 4. Настраивать graceful shutdown
// 5. Обрабатывать сигналы операционной системы

func main() {
	// TODO: Загрузка конфигурации
	// Использовать viper или аналогичную библиотеку для чтения:
	// - .env файлов
	// - переменных окружения
	// - конфигурационных файлов (config.yml)
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// TODO: Инициализация логгера
	// Настроить structured logging с:
	// - Разными уровнями логирования (debug, info, warn, error)
	// - Форматированием JSON для production
	// - Ротацией логов
	// - Отправкой критических ошибок в Slack/email

	// TODO: Инициализация DI контейнера
	// DI контейнер должен инициализировать:
	// - Подключения к базам данных (PostgreSQL, MongoDB, Redis)
	// - Clients для внешних сервисов (Ollama для ИИ, SMTP для email)
	// - Все сервисы и репозитории
	// - HTTP сервер и middleware
	container, err := di.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize DI container: %v", err)
	}
	defer container.Close()

	// TODO: Запуск компонентов системы
	// 1. HTTP сервер для API и веб-интерфейса
	// 2. Pipeline контроллер для обработки тендеров
	// 3. Планировщик задач (cron jobs)
	// 4. Воркеры для обработки очередей
	// 5. Метрики и мониторинг (Prometheus)

	// TODO: Graceful shutdown
	// Настроить корректное завершение работы при получении сигналов:
	// - SIGTERM, SIGINT
	// - Завершение текущих задач
	// - Закрытие соединений с БД
	// - Сохранение состояния
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, gracefully shutting down...")
		cancel()
	}()

	// TODO: Запуск приложения
	// Здесь должен быть запуск всех сервисов и ожидание завершения
	if err := container.Run(ctx); err != nil {
		log.Fatalf("Application failed: %v", err)
	}

	log.Println("Application stopped")
}