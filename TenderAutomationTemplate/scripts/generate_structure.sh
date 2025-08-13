#!/bin/bash

# 🏗️ Скрипт для создания полной структуры файлов проекта автоматизации тендеров
# Создает все необходимые файлы с TODO комментариями

set -e

echo "🏗️ Создание структуры файлов Tender Automation Template..."

# Функция для создания файла с TODO комментарием
create_todo_file() {
    local file_path="$1"
    local description="$2"
    local package_name="$3"
    
    # Создать директорию если не существует
    mkdir -p "$(dirname "$file_path")"
    
    # Определить расширение файла
    extension="${file_path##*.}"
    
    if [[ "$extension" == "go" ]]; then
        cat > "$file_path" << EOF
package $package_name

// TODO: $description
// Этот файл должен содержать:
// - Реализацию интерфейсов
// - Бизнес-логику
// - Обработку ошибок
// - Логирование
// - Тесты

// Placeholder implementation
type Placeholder struct {
    // TODO: Добавить необходимые поля
}

// TODO: Реализовать методы согласно интерфейсам и требованиям архитектуры
EOF
    else
        cat > "$file_path" << EOF
# TODO: $description
# Этот файл должен содержать соответствующую конфигурацию/документацию
EOF
    fi
    
    echo "✅ Создан: $file_path"
}

# ===== DOMAIN LAYER =====
echo "📊 Создание Domain слоя..."

# Tender domain
create_todo_file "internal/domain/tender/repository.go" "Интерфейс репозитория для тендеров - определяет методы для работы с БД" "tender"

# Product domain  
create_todo_file "internal/domain/product/entity.go" "Доменная сущность медицинского оборудования с бизнес-логикой" "product"
create_todo_file "internal/domain/product/repository.go" "Интерфейс репозитория для товаров медицинского оборудования" "product"

# Supplier domain
create_todo_file "internal/domain/supplier/entity.go" "Доменная сущность поставщика медицинского оборудования" "supplier"
create_todo_file "internal/domain/supplier/repository.go" "Интерфейс репозитория для поставщиков" "supplier"

# Analysis domain
create_todo_file "internal/domain/analysis/entity.go" "Доменная сущность для AI анализа тендеров и ценовой оптимизации" "analysis"
create_todo_file "internal/domain/analysis/repository.go" "Интерфейс репозитория для результатов анализа" "analysis"

# Email Campaign domain
create_todo_file "internal/domain/email_campaign/entity.go" "Доменная сущность для email кампаний с поставщиками" "email_campaign"
create_todo_file "internal/domain/email_campaign/repository.go" "Интерфейс репозитория для email кампаний" "email_campaign"

# ===== USE CASES LAYER =====
echo "💼 Создание Use Cases слоя..."

# Tender Discovery
create_todo_file "internal/usecase/tender_discovery/discover_tenders.go" "Use case для поиска новых тендеров на закупочных площадках" "tender_discovery"
create_todo_file "internal/usecase/tender_discovery/interfaces.go" "Интерфейсы для web scraping и парсинга тендерных площадок" "tender_discovery"

# Document Processing
create_todo_file "internal/usecase/document_processing/download_documents.go" "Use case для скачивания документов тендеров" "document_processing"
create_todo_file "internal/usecase/document_processing/extract_text.go" "Use case для извлечения текста из DOC/PDF файлов" "document_processing"
create_todo_file "internal/usecase/document_processing/find_technical_task.go" "Use case для поиска технического задания в документах" "document_processing"
create_todo_file "internal/usecase/document_processing/interfaces.go" "Интерфейсы для обработки документов" "document_processing"

# AI Analysis  
create_todo_file "internal/usecase/ai_analysis/analyze_tender_relevance.go" "Use case для AI анализа релевантности тендера" "ai_analysis"
create_todo_file "internal/usecase/ai_analysis/extract_products.go" "Use case для извлечения товаров из техзадания с помощью AI" "ai_analysis"
create_todo_file "internal/usecase/ai_analysis/categorize_equipment.go" "Use case для категоризации медицинского оборудования" "ai_analysis"
create_todo_file "internal/usecase/ai_analysis/interfaces.go" "Интерфейсы для работы с Llama 4 Maviric и другими AI моделями" "ai_analysis"

# Price Optimization
create_todo_file "internal/usecase/price_optimization/analyze_market_prices.go" "Use case для анализа рыночных цен на медицинское оборудование" "price_optimization"
create_todo_file "internal/usecase/price_optimization/calculate_optimal_price.go" "Use case для расчета оптимальной цены для участия в тендере" "price_optimization"
create_todo_file "internal/usecase/price_optimization/predict_win_probability.go" "Use case для предсказания вероятности выигрыша тендера" "price_optimization"
create_todo_file "internal/usecase/price_optimization/interfaces.go" "Интерфейсы для ценовой аналитики и оптимизации" "price_optimization"

# Supplier Communication
create_todo_file "internal/usecase/supplier_communication/find_suppliers.go" "Use case для поиска поставщиков медицинского оборудования в интернете" "supplier_communication"
create_todo_file "internal/usecase/supplier_communication/generate_emails.go" "Use case для генерации персонализированных писем поставщикам с помощью AI" "supplier_communication"
create_todo_file "internal/usecase/supplier_communication/send_email_campaign.go" "Use case для отправки email кампании поставщикам" "supplier_communication"
create_todo_file "internal/usecase/supplier_communication/process_email_responses.go" "Use case для обработки ответов от поставщиков" "supplier_communication"
create_todo_file "internal/usecase/supplier_communication/interfaces.go" "Интерфейсы для email автоматизации и связи с поставщиками" "supplier_communication"

# Data Collection
create_todo_file "internal/usecase/data_collection/collect_tender_results.go" "Use case для сбора результатов завершенных тендеров" "data_collection"
create_todo_file "internal/usecase/data_collection/analyze_competitors.go" "Use case для анализа конкурентов и их стратегий" "data_collection"
create_todo_file "internal/usecase/data_collection/build_supplier_database.go" "Use case для построения базы данных поставщиков" "data_collection"
create_todo_file "internal/usecase/data_collection/interfaces.go" "Интерфейсы для сбора и анализа данных" "data_collection"

# ===== INFRASTRUCTURE LAYER =====
echo "🌐 Создание Infrastructure слоя..."

# Scraping
create_todo_file "internal/infrastructure/scraping/zakupki_scraper.go" "Скрапер для zakupki.gov.ru с использованием Colly/ChromeDP" "scraping"
create_todo_file "internal/infrastructure/scraping/szvo_scraper.go" "Скрапер для szvo.gov35.ru" "scraping"
create_todo_file "internal/infrastructure/scraping/spb_scraper.go" "Скрапер для gz-spb.ru" "scraping"
create_todo_file "internal/infrastructure/scraping/base_scraper.go" "Базовый класс для всех скраперов с общей функциональностью" "scraping"

# AI Integration
create_todo_file "internal/infrastructure/ai/ollama_client.go" "Клиент для работы с локальной Llama 4 Maviric через Ollama" "ai"
create_todo_file "internal/infrastructure/ai/openai_client.go" "Клиент для работы с OpenAI API (альтернатива)" "ai"
create_todo_file "internal/infrastructure/ai/prompt_templates.go" "Шаблоны промптов для различных AI задач" "ai"

# Document Processing
create_todo_file "internal/infrastructure/document/pdf_processor.go" "Обработчик PDF документов для извлечения текста и таблиц" "document"
create_todo_file "internal/infrastructure/document/office_processor.go" "Обработчик DOC/DOCX файлов с использованием UniOffice" "document"
create_todo_file "internal/infrastructure/document/excel_processor.go" "Обработчик Excel файлов для извлечения данных" "document"
create_todo_file "internal/infrastructure/document/text_extractor.go" "Универсальный извлекатель текста из различных форматов" "document"

# Database
create_todo_file "internal/infrastructure/database/postgres_connection.go" "Подключение к PostgreSQL и миграции" "database"
create_todo_file "internal/infrastructure/database/tender_repository.go" "Реализация репозитория тендеров для PostgreSQL" "database"
create_todo_file "internal/infrastructure/database/product_repository.go" "Реализация репозитория товаров для PostgreSQL" "database"
create_todo_file "internal/infrastructure/database/supplier_repository.go" "Реализация репозитория поставщиков для PostgreSQL" "database"
create_todo_file "internal/infrastructure/database/mongodb_connection.go" "Подключение к MongoDB для документов" "database"
create_todo_file "internal/infrastructure/database/redis_connection.go" "Подключение к Redis для кеширования и очередей" "database"

# Email
create_todo_file "internal/infrastructure/email/smtp_client.go" "SMTP клиент для отправки писем поставщикам" "email"
create_todo_file "internal/infrastructure/email/imap_client.go" "IMAP клиент для чтения ответов от поставщиков" "email"
create_todo_file "internal/infrastructure/email/email_parser.go" "Парсер email сообщений и вложений" "email"

# External Services
create_todo_file "internal/infrastructure/external/file_storage.go" "Клиент для работы с файловым хранилищем (S3/MinIO/Local)" "external"
create_todo_file "internal/infrastructure/external/search_client.go" "Клиент для поисковых движков (Elasticsearch/Bleve)" "external"
create_todo_file "internal/infrastructure/external/notification_client.go" "Клиент для уведомлений (Slack/Email)" "external"

# ===== INTERFACE ADAPTERS =====
echo "🔌 Создание Interface Adapters..."

# API Controllers
create_todo_file "internal/interfaces/api/tender_controller.go" "HTTP контроллер для работы с тендерами" "api"
create_todo_file "internal/interfaces/api/product_controller.go" "HTTP контроллер для работы с товарами" "api"
create_todo_file "internal/interfaces/api/supplier_controller.go" "HTTP контроллер для работы с поставщиками" "api"
create_todo_file "internal/interfaces/api/analysis_controller.go" "HTTP контроллер для AI анализа и аналитики" "api"
create_todo_file "internal/interfaces/api/middleware.go" "HTTP middleware для аутентификации, логирования, CORS" "api"
create_todo_file "internal/interfaces/api/routes.go" "Настройка маршрутов API" "api"

# CLI Interface
create_todo_file "internal/interfaces/cli/discover_command.go" "CLI команда для запуска поиска тендеров" "cli"
create_todo_file "internal/interfaces/cli/analyze_command.go" "CLI команда для анализа конкретного тендера" "cli"
create_todo_file "internal/interfaces/cli/email_command.go" "CLI команда для управления email кампаниями" "cli"
create_todo_file "internal/interfaces/cli/stats_command.go" "CLI команда для просмотра статистики" "cli"

# Scheduler
create_todo_file "internal/interfaces/scheduler/tender_discovery_job.go" "Планировщик для автоматического поиска тендеров" "scheduler"
create_todo_file "internal/interfaces/scheduler/document_processing_job.go" "Планировщик для обработки документов" "scheduler"
create_todo_file "internal/interfaces/scheduler/email_campaign_job.go" "Планировщик для email кампаний" "scheduler"
create_todo_file "internal/interfaces/scheduler/analytics_job.go" "Планировщик для аналитических задач" "scheduler"

# ===== PKG UTILITIES =====
echo "🧩 Создание PKG утилит..."

# Dependency Injection
create_todo_file "pkg/di/container.go" "DI контейнер для инициализации всех компонентов системы" "di"

# AI Utilities
create_todo_file "pkg/ai/prompt_builder.go" "Утилиты для построения промптов для AI" "ai"
create_todo_file "pkg/ai/response_parser.go" "Парсер ответов от AI моделей" "ai"
create_todo_file "pkg/ai/model_manager.go" "Менеджер для работы с различными AI моделями" "ai"

# Document Parser
create_todo_file "pkg/parser/table_extractor.go" "Извлекатель таблиц из документов" "parser"
create_todo_file "pkg/parser/text_cleaner.go" "Очистка и нормализация извлеченного текста" "parser"
create_todo_file "pkg/parser/format_detector.go" "Детектор формата документов" "parser"

# ===== ДОПОЛНИТЕЛЬНЫЕ КОМПОНЕНТЫ =====
echo "🔧 Создание дополнительных компонентов..."

# CMD Components
create_todo_file "cmd/pipeline/main.go" "Pipeline контроллер для координации всех процессов" "main"
create_todo_file "cmd/scraper/main.go" "Standalone scraper для автономного парсинга" "main"
create_todo_file "cmd/processor/main.go" "Standalone обработчик документов" "main"
create_todo_file "cmd/cli/main.go" "CLI интерфейс для управления системой" "main"

# Additional Pipelines
create_todo_file "pipelines/email-campaign.yml" "Конфигурация пайплайна email кампаний" ""
create_todo_file "pipelines/price-optimization.yml" "Конфигурация пайплайна ценовой оптимизации" ""

# Deployment
create_todo_file "deployments/docker-compose.yml" "Docker Compose для локальной разработки" ""
create_todo_file "deployments/Dockerfile" "Dockerfile для контейнеризации приложения" ""
create_todo_file "Makefile" "Makefile с командами для разработки и развертывания" ""

# Documentation
create_todo_file "docs/api.md" "Документация API endpoints" ""
create_todo_file "docs/deployment.md" "Инструкции по развертыванию" ""
create_todo_file "docs/architecture.md" "Описание архитектуры системы" ""
create_todo_file "docs/ai_prompts.md" "Документация AI промптов и моделей" ""

# Scripts
create_todo_file "scripts/setup.sh" "Скрипт первоначальной настройки окружения" ""
create_todo_file "scripts/deploy.sh" "Скрипт развертывания в production" ""
create_todo_file "scripts/backup.sh" "Скрипт резервного копирования данных" ""

# Tests
create_todo_file "internal/domain/tender/entity_test.go" "Тесты для доменной сущности тендера" "tender"
create_todo_file "internal/usecase/ai_analysis/analyze_tender_relevance_test.go" "Тесты для AI анализа тендеров" "ai_analysis"
create_todo_file "internal/infrastructure/scraping/zakupki_scraper_test.go" "Тесты для скрапера zakupki.gov.ru" "scraping"

echo ""
echo "🎉 Структура файлов создана успешно!"
echo ""
echo "📊 Статистика:"
echo "- Domain entities: 5 модулей"
echo "- Use cases: 6 основных процессов"
echo "- Infrastructure: 4 категории интеграций"
echo "- Interface adapters: API, CLI, Scheduler"
echo "- Utilities: DI, AI, Parser"
echo "- Documentation: API, Architecture, Deployment"
echo "- Scripts: Setup, Deploy, Backup"
echo ""
echo "🚀 Следующие шаги:"
echo "1. Заполните TODO файлы согласно Clean Architecture"
echo "2. Реализуйте интерфейсы в каждом слое"
echo "3. Настройте DI контейнер"
echo "4. Напишите тесты для критической функциональности"
echo "5. Настройте CI/CD pipeline"
echo ""
echo "💡 Помните: это шаблон для ЛЮБОЙ программы с Clean Architecture!"
echo "   Адаптируйте под конкретные требования вашего проекта."