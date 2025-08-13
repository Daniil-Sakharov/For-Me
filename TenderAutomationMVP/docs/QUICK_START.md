# 🚀 Быстрый старт - Tender Automation MVP

**Пошаговое руководство для запуска системы автоматизации тендеров**

## 📋 Требования

### Системные требования
- **Go 1.21+** для локальной разработки
- **Docker & Docker Compose** для контейнеризации
- **PostgreSQL 15+** для основной БД
- **Ollama** для AI функциональности (Llama модели)
- **4GB RAM** минимум (8GB рекомендуется)
- **10GB свободного места** для Docker образов и данных

### Операционные системы
- ✅ Linux (Ubuntu 20.04+, CentOS 8+)
- ✅ macOS (10.15+)
- ✅ Windows 10+ (с WSL2)

## 🔧 Метод 1: Быстрый запуск с Docker

### Шаг 1: Клонирование проекта
```bash
# Скачайте проект в удобную директорию
cd ~/projects
git clone <repository-url> TenderAutomationMVP
cd TenderAutomationMVP
```

### Шаг 2: Настройка окружения
```bash
# Скопируйте пример конфигурации
cp .env.example .env

# Отредактируйте конфигурацию под ваши нужды
nano .env
```

**Основные настройки в .env:**
```env
# Database
DB_PASSWORD=secure_password_123
DB_PORT=5432

# AI
AI_MODEL=llama2
OLLAMA_PORT=11434

# Server
SERVER_PORT=8080
LOG_LEVEL=debug
```

### Шаг 3: Запуск сервисов
```bash
# Запуск всех сервисов в фоне
docker-compose up -d

# Просмотр логов запуска
docker-compose logs -f app
```

### Шаг 4: Инициализация AI модели
```bash
# Загрузка Llama модели (может занять время)
docker-compose exec ollama ollama pull llama2

# Проверка доступности модели
docker-compose exec ollama ollama list
```

### Шаг 5: Применение миграций БД
```bash
# TODO: После реализации миграционного инструмента
# docker-compose exec app migrate -path migrations -database "postgres://..." up
echo "⚠️ Миграции будут добавлены при реализации"
```

### Шаг 6: Проверка работы
```bash
# Проверка health check
curl http://localhost:8080/health

# Ожидаемый ответ:
# {"status":"ok","service":"tender-automation-api","timestamp":"..."}

# Проверка основного endpoint
curl http://localhost:8080/
```

## 🛠️ Метод 2: Локальная разработка

### Шаг 1: Установка зависимостей

#### Go разработка
```bash
# Проверка версии Go
go version  # Должно быть 1.21+

# Установка зависимостей
go mod tidy
```

#### PostgreSQL
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS
brew install postgresql
brew services start postgresql

# Docker (альтернатива)
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=tender_automation \
  -p 5432:5432 postgres:15
```

#### Ollama для AI
```bash
# Linux/macOS
curl -fsSL https://ollama.ai/install.sh | sh

# Или через Docker
docker run -d --name ollama \
  -p 11434:11434 \
  -v ollama:/root/.ollama \
  ollama/ollama:latest

# Загрузка модели
ollama pull llama2
```

### Шаг 2: Настройка базы данных
```bash
# Создание БД (если не используете Docker)
createdb tender_automation

# TODO: Запуск миграций
# migrate -path migrations -database "postgres://..." up
```

### Шаг 3: Конфигурация
```bash
# Копирование и настройка .env
cp .env.example .env

# Настройка для локальной разработки:
# DB_HOST=localhost
# AI_URL=http://localhost:11434
# LOG_LEVEL=debug
# LOG_FORMAT=console
```

### Шаг 4: Запуск приложения
```bash
# Запуск API сервера
go run cmd/api/main.go

# В отдельном терминале - scraper (когда будет реализован)
# go run cmd/scraper/main.go

# CLI утилита (когда будет реализована)
# go run cmd/cli/main.go --help
```

## 🧪 Тестирование системы

### Проверка компонентов
```bash
# 1. Health check API
curl http://localhost:8080/health

# 2. Проверка БД подключения
# TODO: После реализации добавить endpoint
# curl http://localhost:8080/health/db

# 3. Проверка AI сервиса
# TODO: После реализации добавить endpoint  
# curl http://localhost:8080/health/ai

# 4. Проверка основного API
curl http://localhost:8080/api/v1/tenders
```

### Создание тестового тендера
```bash
# TODO: После реализации API endpoints
# curl -X POST http://localhost:8080/api/v1/tenders \
#   -H "Content-Type: application/json" \
#   -d '{
#     "external_id": "test-123",
#     "title": "Тестовый тендер",
#     "platform": "zakupki",
#     "url": "http://example.com"
#   }'
```

### Запуск тестов
```bash
# Unit тесты
go test ./...

# Тесты с покрытием
go test -cover ./...

# Интеграционные тесты
go test -tags=integration ./test/integration/...

# Через Docker
docker build --target test -t tender-automation:test .
```

## 🔧 Разработка и отладка

### Полезные команды
```bash
# Просмотр логов всех сервисов
docker-compose logs -f

# Просмотр логов конкретного сервиса
docker-compose logs -f app
docker-compose logs -f postgres
docker-compose logs -f ollama

# Перезапуск сервиса
docker-compose restart app

# Подключение к БД
docker-compose exec postgres psql -U postgres -d tender_automation

# Выполнение команд в контейнере приложения
docker-compose exec app sh
docker-compose exec app go mod tidy
```

### Мониторинг ресурсов
```bash
# Использование ресурсов контейнерами
docker stats $(docker-compose ps -q)

# Размер образов
docker images | grep tender

# Использование дискового пространства
docker system df
```

### Hot reload для разработки
```bash
# Установка air для hot reload
go install github.com/cosmtrek/air@latest

# Создание .air.toml конфигурации
cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main cmd/api/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
EOF

# Запуск с hot reload
air
```

## 🐛 Решение проблем

### Частые проблемы

#### 1. Ошибка подключения к БД
```bash
# Проверка доступности PostgreSQL
docker-compose ps postgres

# Просмотр логов БД
docker-compose logs postgres

# Проверка подключения
docker-compose exec postgres pg_isready -U postgres
```

#### 2. Ollama не отвечает
```bash
# Проверка статуса Ollama
curl http://localhost:11434/api/version

# Просмотр логов
docker-compose logs ollama

# Перезапуск сервиса
docker-compose restart ollama

# Проверка загруженных моделей
docker-compose exec ollama ollama list
```

#### 3. Порты заняты
```bash
# Проверка использования портов
netstat -tulpn | grep :8080
netstat -tulpn | grep :5432
netstat -tulpn | grep :11434

# Остановка конфликтующих сервисов
sudo systemctl stop postgresql  # если локальный postgres
```

#### 4. Проблемы с Docker
```bash
# Полная очистка
docker-compose down -v --remove-orphans

# Пересборка образов
docker-compose build --no-cache

# Очистка системы Docker
docker system prune -f
```

### Логи и отладка

#### Настройка уровня логирования
```env
# В .env файле
LOG_LEVEL=debug    # Подробные логи
LOG_FORMAT=console # Удобочитаемый формат для разработки
```

#### Просмотр детальных логов
```bash
# Логи приложения с временными метками
docker-compose logs -f --timestamps app

# Только ошибки
docker-compose logs app 2>&1 | grep ERROR

# Поиск по логам
docker-compose logs app 2>&1 | grep "конкретное сообщение"
```

## 📊 Мониторинг

### Health checks
```bash
# Основной health check
curl http://localhost:8080/health

# TODO: После реализации дополнительных checks
# curl http://localhost:8080/health/db      # Статус БД
# curl http://localhost:8080/health/ai      # Статус AI
# curl http://localhost:8080/health/scraper # Статус scraper
```

### Метрики (в планах)
```bash
# TODO: После реализации Prometheus интеграции
# curl http://localhost:8080/metrics
```

## 🚀 Следующие шаги

После успешного запуска MVP:

### 1. Реализация компонентов
- [ ] **Domain layer** - доменная логика уже создана
- [ ] **Repository** - реализация для PostgreSQL  
- [ ] **Use cases** - бизнес-логика приложения
- [ ] **HTTP handlers** - REST API endpoints
- [ ] **Scraper** - парсинг zakupki.gov.ru
- [ ] **AI integration** - интеграция с Ollama

### 2. Тестирование
- [ ] Unit тесты для всех компонентов
- [ ] Интеграционные тесты
- [ ] End-to-end тесты
- [ ] Performance тесты

### 3. Расширение функциональности
- [ ] Обработка документов (PDF/DOC)
- [ ] Email автоматизация
- [ ] Дополнительные тендерные площадки
- [ ] Web UI интерфейс
- [ ] Продвинутая аналитика

### 4. Production готовность
- [ ] Мониторинг и метрики
- [ ] Логирование и алертинг
- [ ] Backup и восстановление
- [ ] Безопасность и аудит
- [ ] CI/CD pipeline

## 📚 Дополнительные ресурсы

- **[README.md](../README.md)** - Полная документация системы
- **[API Documentation](api.md)** - Описание REST API
- **[Deployment Guide](deployment.md)** - Инструкции по развертыванию
- **[Architecture Overview](architecture.md)** - Архитектура системы

## 🆘 Получение помощи

Если возникли проблемы:

1. **Проверьте логи** - `docker-compose logs -f app`
2. **Проверьте health checks** - `curl http://localhost:8080/health`
3. **Просмотрите документацию** - все инструкции в этом файле
4. **Проверьте issues** - возможно проблема уже решена
5. **Создайте issue** - с детальным описанием проблемы

---

**🎉 Поздравляем! Вы запустили Tender Automation MVP**

Система готова к разработке и расширению функциональности. Следуйте Clean Architecture принципам и используйте подробные TODO комментарии в коде для реализации компонентов.