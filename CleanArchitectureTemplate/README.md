# 🏗️ Универсальный шаблон Clean Architecture для Go

Этот шаблон предоставляет **готовую структуру** для создания любого Go приложения с Clean Architecture. Просто скопируйте, адаптируйте под ваш проект и начинайте разработку!

## 🎯 Что это?

**Универсальный стартер** для Go проектов с:
- ✅ Правильной структурой папок
- ✅ Готовыми шаблонами всех слоев
- ✅ Подробными комментариями
- ✅ Примерами для разных типов приложений
- ✅ Настроенным DI контейнером
- ✅ Конфигурацией для популярных сервисов

## 📁 Структура проекта

```
your-project/
├── cmd/
│   └── main.go                           # 🚀 Точка входа приложения
├── configs/
│   └── config.go                         # ⚙️ Конфигурация
├── internal/
│   ├── domain/                           # 🏛️ DOMAIN СЛОЙ
│   │   └── example_entity/
│   │       ├── entity.go                 # Доменная модель
│   │       └── repository.go             # Интерфейс репозитория
│   ├── usecase/                          # 💼 USE CASES СЛОЙ
│   │   └── example_feature/
│   │       └── interfaces.go             # Интерфейсы зависимостей
│   ├── infrastructure/                   # 🌐 INFRASTRUCTURE СЛОЙ
│   │   ├── database/                     # Реализации репозиториев
│   │   ├── external/                     # Внешние сервисы
│   │   └── web/                          # HTTP сервер
│   └── interfaces/                       # 🔌 INTERFACE ADAPTERS
│       └── controllers/                  # HTTP контроллеры
├── pkg/
│   └── di/
│       └── container.go                  # 🧩 DI контейнер
├── go.mod                                # 📦 Зависимости
├── .env.example                          # 📝 Пример конфигурации
└── README.md                             # 📚 Документация
```

## 🚀 Быстрый старт

### 1. Скопируйте шаблон

```bash
cp -r CleanArchitectureTemplate your-new-project
cd your-new-project
```

### 2. Адаптируйте под ваш проект

#### 📝 Обновите go.mod

```go
// Замените "your-project-name" на название вашего проекта
module my-awesome-app

go 1.21

// Раскомментируйте нужные зависимости
require (
    github.com/go-playground/validator/v10 v10.16.0
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/lib/pq v1.10.9  // PostgreSQL
    // github.com/gin-gonic/gin v1.9.1  // Если предпочитаете Gin
)
```

#### 🔄 Обновите импорты

Замените во всех файлах:
```
your-project-name → my-awesome-app
```

### 3. Определите ваши сущности

#### 🏛️ Создайте доменные модели

```bash
# Переименуйте example_entity в ваши сущности
mv internal/domain/example_entity internal/domain/user
mv internal/domain/example_entity internal/domain/product
# и т.д.
```

#### ✏️ Адаптируйте entity.go

```go
// internal/domain/user/entity.go
package user

type User struct {
    ID        uint      `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(email, name string) (*User, error) {
    // Ваша бизнес-логика валидации
}
```

### 4. Настройте Use Cases

#### 💼 Создайте ваши фичи

```bash
# Переименуйте example_feature в ваши фичи
mv internal/usecase/example_feature internal/usecase/auth
mv internal/usecase/example_feature internal/usecase/payment
# и т.д.
```

#### ⚙️ Определите интерфейсы

```go
// internal/usecase/auth/interfaces.go
package auth

type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(hashedPassword, password string) error
}

type TokenGenerator interface {
    Generate(userID uint, email string) (string, error)
    Validate(token string) (*TokenData, error)
}
```

### 5. Обновите DI контейнер

```go
// pkg/di/container.go
type Container struct {
    // Ваши репозитории
    UserRepo    user.Repository
    ProductRepo product.Repository
    
    // Ваши сервисы
    PasswordHasher   auth.PasswordHasher
    TokenGenerator   auth.TokenGenerator
    
    // Ваши Use Cases
    RegisterUC *auth.RegisterUseCase
    LoginUC    *auth.LoginUseCase
    
    // Ваши контроллеры
    AuthController *controllers.AuthController
}
```

### 6. Настройте конфигурацию

```bash
cp .env.example .env
```

```env
# .env
APP_NAME=My Awesome App
APP_VERSION=1.0.0

DB_HOST=localhost
DB_PORT=5432
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=myapp_db

JWT_SECRET=your-super-secret-key
```

### 7. Запустите приложение

```bash
go mod tidy
go run cmd/main.go
```

## 🎨 Примеры адаптации

### 🛒 E-commerce приложение

```go
// Сущности
internal/domain/user/
internal/domain/product/
internal/domain/order/
internal/domain/category/

// Фичи
internal/usecase/auth/        // Регистрация, авторизация
internal/usecase/catalog/     // Каталог товаров
internal/usecase/cart/        // Корзина
internal/usecase/payment/     // Платежи
internal/usecase/shipping/    // Доставка

// Сервисы
PasswordHasher, TokenGenerator
PaymentProcessor (Stripe/PayPal)
EmailSender (SMTP/SendGrid)
FileStorage (S3/Local)
```

### 📝 Blog/CMS

```go
// Сущности
internal/domain/user/
internal/domain/article/
internal/domain/comment/
internal/domain/category/

// Фичи
internal/usecase/auth/
internal/usecase/blog/
internal/usecase/moderation/
internal/usecase/analytics/

// Сервисы
PasswordHasher, TokenGenerator
EmailSender
FileStorage (для изображений)
SearchEngine (Elasticsearch)
```

### 🏦 Финансовое приложение

```go
// Сущности
internal/domain/user/
internal/domain/account/
internal/domain/transaction/
internal/domain/card/

// Фичи
internal/usecase/auth/
internal/usecase/banking/
internal/usecase/payment/
internal/usecase/analytics/

// Сервисы
PasswordHasher, TokenGenerator
PaymentProcessor
EmailSender, SMSSender
AuditLogger
EncryptionService
```

## 🛠️ Доступные конфигурации

### 🗄️ Базы данных

```go
// PostgreSQL (по умолчанию)
github.com/lib/pq v1.10.9

// MySQL
github.com/go-sql-driver/mysql v1.7.1

// MongoDB
go.mongodb.org/mongo-driver v1.13.1

// Redis (кеширование)
github.com/redis/go-redis/v9 v9.3.0
```

### 🌐 Web фреймворки

```go
// net/http (по умолчанию) - уже в стандартной библиотеке

// Gin
github.com/gin-gonic/gin v1.9.1

// Fiber
github.com/gofiber/fiber/v2 v2.52.0
```

### 📧 Внешние сервисы

```go
// Email
github.com/go-mail/mail v2.3.1+incompatible

// AWS
github.com/aws/aws-sdk-go-v2 v1.24.0

// Message queues
github.com/segmentio/kafka-go v0.4.46

// Logging
github.com/sirupsen/logrus v1.9.3
go.uber.org/zap v1.26.0
```

## 📋 Чек-лист адаптации

### ✅ Обязательные шаги

- [ ] Переименовать модуль в `go.mod`
- [ ] Обновить импорты во всех файлах
- [ ] Переименовать `example_entity` в ваши сущности
- [ ] Переименовать `example_feature` в ваши фичи
- [ ] Адаптировать структуры в `configs/config.go`
- [ ] Обновить поля в DI контейнере
- [ ] Создать `.env` файл с вашими настройками

### 🔀 Опциональные шаги

- [ ] Выбрать и настроить базу данных
- [ ] Раскомментировать нужные зависимости
- [ ] Добавить дополнительные внешние сервисы
- [ ] Настроить дополнительные серверы (gRPC, GraphQL)
- [ ] Добавить метрики и логирование
- [ ] Настроить Docker и развертывание

## 🏗️ Принципы Clean Architecture

### 📏 Правила зависимостей

```
Infrastructure → Interfaces → Use Cases → Domain
     ↑              ↑            ↑         ↑
(PostgreSQL,    (HTTP,       (Business  (Business
 bcrypt,        JSON,        Logic,     Rules,
 JWT)           REST)        Orchestr.) Entities)
```

**Главное правило:** Зависимости всегда направлены внутрь!

### 🎯 Ответственности слоев

- **🏛️ Domain:** Чистые бизнес-правила, никаких внешних зависимостей
- **💼 Use Cases:** Оркестрация, координация между компонентами
- **🔌 Interfaces:** Адаптация HTTP/gRPC/GraphQL в доменные типы
- **🌐 Infrastructure:** Конкретные реализации (БД, внешние API)

### 🧩 Dependency Injection

```go
// ✅ Правильно: зависимости через интерфейсы
type RegisterUseCase struct {
    userRepo       user.Repository      // интерфейс
    passwordHasher PasswordHasher       // интерфейс
    tokenGenerator TokenGenerator       // интерфейс
}

// ❌ Неправильно: прямые зависимости
type RegisterUseCase struct {
    userRepo *PostgreSQLUserRepository  // конкретный тип
    bcrypt   *BcryptService            // конкретный тип
}
```

## 🔧 Расширение шаблона

### 📝 Добавление новой сущности

1. Создайте папку `internal/domain/new_entity/`
2. Добавьте `entity.go` и `repository.go`
3. Создайте реализацию в `internal/infrastructure/database/`
4. Добавьте в DI контейнер

### 💼 Добавление новой фичи

1. Создайте папку `internal/usecase/new_feature/`
2. Добавьте `interfaces.go` для зависимостей
3. Создайте Use Cases файлы
4. Реализуйте интерфейсы в `infrastructure/external/`
5. Добавьте контроллер в `interfaces/controllers/`
6. Обновите DI контейнер

### 🌐 Добавление нового протокола

1. Создайте обработчики в `interfaces/`
2. Добавьте сервер в `infrastructure/web/`
3. Обновите DI контейнер
4. Зарегистрируйте маршруты

## 🆘 Помощь и поддержка

### 📚 Полезные ресурсы

- [Clean Architecture (книга)](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164)
- [SOLID принципы](https://en.wikipedia.org/wiki/SOLID)
- [Go Clean Architecture примеры](https://github.com/bxcodec/go-clean-arch)

### 🐛 Частые проблемы

**Проблема:** Циклические импорты
**Решение:** Проверьте направление зависимостей, используйте интерфейсы

**Проблема:** "Interface not implemented" ошибки
**Решение:** Убедитесь, что все методы интерфейса реализованы

**Проблема:** DI контейнер не компилируется
**Решение:** Проверьте импорты и порядок инициализации

## 🎉 Заключение

Этот шаблон поможет вам:

- ⚡ **Быстро стартовать** новые Go проекты
- 🏗️ **Следовать правильной архитектуре** с самого начала
- 🔄 **Легко масштабировать** и изменять приложение
- 🧪 **Эффективно тестировать** каждый слой изолированно
- 👥 **Работать в команде** с понятной структурой

**Запомните главное:** структура и принципы остаются неизменными, меняется только содержимое под ваш проект!

Happy coding! 🚀