# =====================================================================
# 🐳 DOCKERFILE ДЛЯ TENDER AUTOMATION MVP
# =====================================================================
#
# Multi-stage Dockerfile для оптимизации размера и безопасности:
# 1. Builder stage - компиляция Go приложения
# 2. Development stage - для разработки с hot reload
# 3. Production stage - минимальный образ для production
#
# TODO: При расширении функциональности добавить:
# - Установка дополнительных утилит для обработки документов
# - Настройка CA сертификатов для HTTPS клиентов
# - Пользователь с ограниченными правами для безопасности
# - Health check скрипт

# =====================================================================
# 🔧 BUILDER STAGE - Компиляция приложения
# =====================================================================
FROM golang:1.21-alpine AS builder

# 📝 Метаданные образа
LABEL maintainer="tender-automation-team"
LABEL description="Tender Automation MVP - Builder Stage"
LABEL version="1.0.0"

# 🔧 Установка необходимых пакетов для сборки
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    curl \
    bash

# 📊 Создание рабочей директории
WORKDIR /build

# 📦 Копирование go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# ⬇️ Загрузка зависимостей (будет закешировано если go.mod не изменился)
RUN go mod download && go mod verify

# 📋 Копирование исходного кода
COPY . .

# 🏗️ Компиляция приложения с оптимизациями
# Флаги:
# -ldflags="-s -w" - удаление debug информации для уменьшения размера
# -trimpath - удаление путей файлов из executable
# -tags netgo - статическая линковка сетевых библиотек
# -a - принудительная пересборка всех пакетов
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -trimpath \
    -tags netgo \
    -a -installsuffix cgo \
    -o app cmd/api/main.go

# 🧪 Компиляция дополнительных утилит (scraper, cli)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -trimpath \
    -tags netgo \
    -a -installsuffix cgo \
    -o scraper cmd/scraper/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -trimpath \
    -tags netgo \
    -a -installsuffix cgo \
    -o cli cmd/cli/main.go

# ✅ Проверка что executable файлы созданы корректно
RUN ./app --version || echo "Version command not implemented yet"
RUN file ./app
RUN ls -la ./app

# =====================================================================
# 🛠️ DEVELOPMENT STAGE - Для разработки
# =====================================================================
FROM golang:1.21-alpine AS development

# 📝 Метаданные
LABEL stage="development"
LABEL description="Tender Automation MVP - Development Stage"

# 🔧 Установка инструментов для разработки
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    curl \
    bash \
    make \
    gcc \
    musl-dev \
    sqlite \
    postgresql-client

# 🔧 Установка дополнительных Go инструментов для разработки
RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 👤 Создание пользователя для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 📊 Настройка рабочей директории
WORKDIR /app

# 🔧 Настройка переменных окружения для разработки
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIN_MODE=debug

# 📦 Копирование go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# 📋 Копирование исходного кода
COPY . .

# 👤 Изменение владельца файлов
RUN chown -R appuser:appgroup /app

# 👤 Переключение на пользователя
USER appuser

# 🌐 Открытие порта
EXPOSE 8080

# 🔍 Health check для development
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 🚀 Команда по умолчанию (с hot reload через air)
CMD ["air", "-c", ".air.toml"]

# =====================================================================
# 🚀 PRODUCTION STAGE - Минимальный образ для production
# =====================================================================
FROM alpine:3.18 AS production

# 📝 Метаданные для production
LABEL stage="production"
LABEL description="Tender Automation MVP - Production Stage"
LABEL version="1.0.0"
LABEL vendor="Tender Automation Team"

# 🔧 Установка только необходимых runtime зависимостей
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    && update-ca-certificates

# 👤 Создание пользователя без shell для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup -s /sbin/nologin

# 📊 Создание директорий с правильными разрешениями
RUN mkdir -p /app/logs /app/tmp /app/uploads && \
    chown -R appuser:appgroup /app

# 📊 Настройка рабочей директории
WORKDIR /app

# 📦 Копирование скомпилированных приложений из builder stage
COPY --from=builder --chown=appuser:appgroup /build/app ./
COPY --from=builder --chown=appuser:appgroup /build/scraper ./
COPY --from=builder --chown=appuser:appgroup /build/cli ./

# 📋 Копирование конфигурационных файлов
COPY --chown=appuser:appgroup .env.example ./
COPY --chown=appuser:appgroup migrations/ ./migrations/

# 📋 Создание конфигурационного файла для production
RUN echo '#!/bin/sh\n\
echo "🚀 Starting Tender Automation in Production Mode"\n\
echo "📊 Application: $APP_NAME"\n\
echo "🌍 Environment: $APP_ENV"\n\
echo "🌐 Server Address: $SERVER_HOST:$SERVER_PORT"\n\
echo "🗃️ Database: $DB_HOST:$DB_PORT/$DB_NAME"\n\
echo "🤖 AI Service: $AI_URL"\n\
exec ./app "$@"' > entrypoint.sh && \
    chmod +x entrypoint.sh && \
    chown appuser:appgroup entrypoint.sh

# 🔧 Настройка переменных окружения для production
ENV APP_ENV=production \
    GIN_MODE=release \
    LOG_LEVEL=info \
    LOG_FORMAT=json \
    LOG_OUTPUT=stdout

# 👤 Переключение на пользователя
USER appuser

# 🌐 Открытие порта
EXPOSE 8080

# 🔍 Health check для production
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 📊 Настройка volume для логов (опционально)
VOLUME ["/app/logs"]

# 🚀 Точка входа
ENTRYPOINT ["./entrypoint.sh"]

# 🔧 Команда по умолчанию
CMD []

# =====================================================================
# 🧪 TEST STAGE - Для запуска тестов
# =====================================================================
FROM golang:1.21-alpine AS test

# 🔧 Установка инструментов для тестирования
RUN apk add --no-cache \
    git \
    ca-certificates \
    curl \
    bash \
    gcc \
    musl-dev

# 📊 Рабочая директория
WORKDIR /app

# 📦 Копирование исходного кода
COPY . .

# ⬇️ Загрузка зависимостей
RUN go mod download

# 🧪 Запуск тестов
RUN go test -v -race -coverprofile=coverage.out ./...

# 📊 Генерация отчета о покрытии
RUN go tool cover -html=coverage.out -o coverage.html

# =====================================================================
# 🔧 ИНСТРУКЦИИ ПО ИСПОЛЬЗОВАНИЮ
# =====================================================================

# TODO: Команды для сборки и запуска:
#
# 🏗️ Сборка для разработки:
#   docker build --target development -t tender-automation:dev .
#
# 🚀 Сборка для production:
#   docker build --target production -t tender-automation:prod .
#
# 🧪 Запуск тестов:
#   docker build --target test -t tender-automation:test .
#
# 🔧 Запуск development контейнера:
#   docker run -p 8080:8080 -v $(pwd):/app tender-automation:dev
#
# 🚀 Запуск production контейнера:
#   docker run -p 8080:8080 \
#     -e DB_HOST=localhost \
#     -e DB_PASSWORD=password \
#     tender-automation:prod
#
# 📊 Проверка размера образов:
#   docker images | grep tender-automation
#
# 🔍 Инспекция образа:
#   docker inspect tender-automation:prod
#
# 📋 Просмотр логов:
#   docker logs container_name
#
# 🐚 Подключение к контейнеру:
#   docker exec -it container_name sh
#
# 🧪 Копирование отчета о покрытии:
#   docker run --rm -v $(pwd):/output tender-automation:test \
#     cp coverage.html /output/

# =====================================================================
# 🔒 БЕЗОПАСНОСТЬ
# =====================================================================

# TODO: Рекомендации по безопасности:
#
# 1. 👤 Используется непривилегированный пользователь
# 2. 🔧 Минимальный base образ (Alpine Linux)
# 3. 📦 Статическая компиляция без CGO
# 4. 🗑️ Удаление build зависимостей в production
# 5. 🔍 Health checks для monitoring
# 6. 📊 Read-only файловая система (можно добавить)
# 7. 🌐 Только необходимые порты открыты
#
# Дополнительно для production:
# - Используйте секреты вместо переменных окружения
# - Настройте AppArmor/SELinux профили
# - Сканируйте образы на уязвимости
# - Используйте signed образы

# =====================================================================
# 📊 ОПТИМИЗАЦИЯ
# =====================================================================

# TODO: Оптимизации образа:
#
# 📏 Размеры образов (примерно):
# - Builder: ~800MB (только для сборки)
# - Development: ~500MB (с инструментами)
# - Production: ~20MB (минимальный)
# - Test: ~300MB (с coverage tools)
#
# 🚀 Performance оптимизации:
# - Multi-stage build для минимального размера
# - Кеширование слоев Docker
# - Статическая компиляция
# - Удаление debug информации
#
# 💾 Storage оптимизации:
# - .dockerignore для исключения ненужных файлов
# - Минимальный base образ
# - Очистка package cache

# =====================================================================
# ✅ ЗАКЛЮЧЕНИЕ
# =====================================================================

# Этот Dockerfile обеспечивает:
# ✅ Multi-stage сборку для оптимизации
# ✅ Отдельные образы для dev/test/prod
# ✅ Безопасность через непривилегированного пользователя
# ✅ Минимальный размер production образа
# ✅ Health checks для мониторинга
# ✅ Переменные окружения для конфигурации
# ✅ Статическую компиляцию для портабельности
#
# Поддерживаемые targets:
# - development: для разработки с hot reload
# - production: минимальный образ для prod
# - test: для запуска тестов и генерации coverage