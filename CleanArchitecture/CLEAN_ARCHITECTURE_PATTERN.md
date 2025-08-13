# 🏗️ Clean Architecture как паттерн: что ВСЕГДА повторяется

Этот файл объясняет **универсальные принципы и структуры** Clean Architecture, которые нужно запомнить и применять в любом проекте, а также что может изменяться в зависимости от конкретного приложения.

---

## 📋 Содержание

- [🎯 Что такое архитектурный паттерн](#-что-такое-архитектурный-паттерн)
- [📁 Структура папок: что ВСЕГДА одинаково](#-структура-папок-что-всегда-одинаково)
- [🏛️ Domain слой: неизменные принципы](#️-domain-слой-неизменные-принципы)
- [💼 Use Cases слой: универсальные паттерны](#-use-cases-слой-универсальные-паттерны)
- [🔌 Interface Adapters: повторяющиеся структуры](#-interface-adapters-повторяющиеся-структуры)
- [🌐 Infrastructure: что меняется, что нет](#-infrastructure-что-меняется-что-нет)
- [🧩 Dependency Injection: универсальный подход](#-dependency-injection-универсальный-подход)
- [📝 Чек-лист для нового проекта](#-чек-лист-для-нового-проекта)

---

## 🎯 Что такое архитектурный паттерн

**Архитектурный паттерн** - это **повторяемая структура**, которая решает типичные проблемы в разработке ПО. Clean Architecture предоставляет:

### ✅ **Что ВСЕГДА одинаково (паттерн):**
- 🏗️ **4 слоя** и их ответственности
- 🔄 **Направление зависимостей** (всегда внутрь)
- 📁 **Структура папок** и размещение файлов
- 🧩 **Принципы DI** и создания объектов
- 📏 **Правила взаимодействия** между слоями

### 🔀 **Что МОЖЕТ изменяться (реализация):**
- 📊 **Доменные модели** (User, Product, Order, etc.)
- 🎯 **Use Cases** (RegisterUser, ProcessPayment, etc.)
- 🗄️ **Базы данных** (PostgreSQL, MongoDB, etc.)
- 🌐 **Фреймворки** (HTTP, gRPC, GraphQL, etc.)
- 🔧 **Внешние сервисы** (JWT, OAuth, Payment APIs, etc.)

---

## 📁 Структура папок: что ВСЕГДА одинаково

### 🔒 **НЕИЗМЕННАЯ структура (запомните навсегда!):**

```
your-project/
├── cmd/                          # ✅ ВСЕГДА: точки входа
│   └── main.go                   # ✅ ВСЕГДА: основная точка входа
├── configs/                      # ✅ ВСЕГДА: конфигурация
│   └── config.go                 # ✅ ВСЕГДА: загрузка настроек
├── internal/                     # ✅ ВСЕГДА: внутренний код
│   ├── domain/                   # ✅ ВСЕГДА: доменный слой
│   │   ├── {entity1}/            # 🔀 МЕНЯЕТСЯ: названия сущностей
│   │   │   ├── entity.go         # ✅ ВСЕГДА: доменная модель
│   │   │   └── repository.go     # ✅ ВСЕГДА: интерфейс репозитория
│   │   └── {entity2}/
│   │       ├── entity.go
│   │       └── repository.go
│   ├── usecase/                  # ✅ ВСЕГДА: слой Use Cases
│   │   ├── {feature1}/           # 🔀 МЕНЯЕТСЯ: названия функций
│   │   │   ├── interfaces.go     # ✅ ВСЕГДА: интерфейсы зависимостей
│   │   │   ├── {action1}.go      # 🔀 МЕНЯЕТСЯ: конкретные действия
│   │   │   └── {action2}.go
│   │   └── {feature2}/
│   │       ├── interfaces.go
│   │       └── {action}.go
│   ├── infrastructure/           # ✅ ВСЕГДА: инфраструктурный слой
│   │   ├── database/             # 🔀 МЕНЯЕТСЯ: тип БД, но папка есть всегда
│   │   │   ├── connection.go     # ✅ ВСЕГДА: подключение к БД
│   │   │   ├── {entity1}_repository.go  # ✅ ВСЕГДА: реализации репозиториев
│   │   │   └── {entity2}_repository.go
│   │   ├── web/                  # 🔀 МЕНЯЕТСЯ: может быть grpc, graphql
│   │   │   ├── server.go         # ✅ ВСЕГДА: сервер (любой тип)
│   │   │   └── middleware.go     # ✅ ВСЕГДА: middleware для протокола
│   │   └── external/             # ✅ ВСЕГДА: внешние сервисы
│   │       ├── {service1}.go     # 🔀 МЕНЯЕТСЯ: конкретные сервисы
│   │       └── {service2}.go
│   └── interfaces/               # ✅ ВСЕГДА: слой адаптеров
│       ├── controllers/          # 🔀 МЕНЯЕТСЯ: может быть handlers, resolvers
│       │   ├── {feature1}_controller.go  # ✅ ВСЕГДА: контроллеры по фичам
│       │   └── {feature2}_controller.go
│       └── presenters/           # ✅ ВСЕГДА: если нужны презентеры
└── pkg/                          # ✅ ВСЕГДА: публичные пакеты
    └── di/                       # ✅ ВСЕГДА: Dependency Injection
        └── container.go          # ✅ ВСЕГДА: DI контейнер
```

### 📝 **Правила именования (ВСЕГДА соблюдайте):**

```go
// ✅ ВСЕГДА: domain/{entity}/entity.go
internal/domain/user/entity.go
internal/domain/product/entity.go
internal/domain/order/entity.go

// ✅ ВСЕГДА: domain/{entity}/repository.go
internal/domain/user/repository.go
internal/domain/product/repository.go

// ✅ ВСЕГДА: usecase/{feature}/{action}.go
internal/usecase/auth/register.go
internal/usecase/auth/login.go
internal/usecase/payment/process.go
internal/usecase/order/create.go

// ✅ ВСЕГДА: infrastructure/database/{entity}_repository.go
internal/infrastructure/database/user_repository.go
internal/infrastructure/database/product_repository.go

// ✅ ВСЕГДА: interfaces/controllers/{feature}_controller.go
internal/interfaces/controllers/auth_controller.go
internal/interfaces/controllers/product_controller.go
```

---

## 🏛️ Domain слой: неизменные принципы

### ✅ **Что ВСЕГДА одинаково:**

#### 1. Структура Entity файла

```go
// ✅ ВСЕГДА: такая структура entity.go
package {entity}  // 🔀 МЕНЯЕТСЯ: имя сущности

import (
    "errors"      // ✅ ВСЕГДА: для доменных ошибок
    "time"        // ✅ ВСЕГДА: для временных меток
    // НЕТ ВНЕШНИХ ЗАВИСИМОСТЕЙ! ✅ ВСЕГДА
)

// ✅ ВСЕГДА: доменная модель
type {Entity} struct {  // 🔀 МЕНЯЕТСЯ: имя и поля
    ID        uint      // ✅ ВСЕГДА: ID есть почти всегда
    CreatedAt time.Time // ✅ ВСЕГДА: временные метки
    UpdatedAt time.Time
    // 🔀 МЕНЯЕТСЯ: специфичные поля
}

// ✅ ВСЕГДА: доменные ошибки
var (
    Err{Something}       = errors.New("...")  // 🔀 МЕНЯЕТСЯ: конкретные ошибки
    ErrInvalid{Field}    = errors.New("...")
)

// ✅ ВСЕГДА: конструктор с валидацией
func New{Entity}(/* параметры */) (*{Entity}, error) {  // 🔀 МЕНЯЕТСЯ: параметры
    // ✅ ВСЕГДА: валидация доменных правил
    if /* условие */ {
        return nil, Err{Something}
    }
    
    // ✅ ВСЕГДА: установка временных меток
    now := time.Now()
    
    return &{Entity}{
        // поля...
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

// ✅ ВСЕГДА: методы для изменения состояния
func (e *{Entity}) Update{Field}(newValue Type) error {  // 🔀 МЕНЯЕТСЯ: конкретные методы
    // ✅ ВСЕГДА: валидация
    // ✅ ВСЕГДА: обновление UpdatedAt
    e.UpdatedAt = time.Now()
    return nil
}
```

#### 2. Структура Repository интерфейса

```go
// ✅ ВСЕГДА: такая структура repository.go
package {entity}

import "context"  // ✅ ВСЕГДА: контекст для всех методов

// ✅ ВСЕГДА: интерфейс Repository
type Repository interface {
    // ✅ ВСЕГДА: базовые CRUD операции
    Save(ctx context.Context, entity *{Entity}) error
    FindByID(ctx context.Context, id uint) (*{Entity}, error)
    Delete(ctx context.Context, id uint) error
    
    // 🔀 МЕНЯЕТСЯ: специфичные методы поиска
    FindBy{Field}(ctx context.Context, value Type) (*{Entity}, error)
    Exists{Condition}(ctx context.Context, params Type) (bool, error)
}
```

### 🔀 **Что МЕНЯЕТСЯ в зависимости от проекта:**

- **Названия сущностей:** User, Product, Order, Article, etc.
- **Поля моделей:** Email, Name, Price, Title, etc.
- **Доменные правила:** валидация email vs валидация цены
- **Методы репозитория:** FindByEmail vs FindByCategory
- **Доменные ошибки:** конкретные сообщения

---

## 💼 Use Cases слой: универсальные паттерны

### ✅ **Что ВСЕГДА одинаково:**

#### 1. Структура файла interfaces.go

```go
// ✅ ВСЕГДА: интерфейсы зависимостей в usecase слое
package {feature}

// ✅ ВСЕГДА: интерфейсы для внешних сервисов
type {Service}Interface interface {  // 🔀 МЕНЯЕТСЯ: конкретные интерфейсы
    {Method}(params Type) (Type, error)
}

// Примеры ВСЕГДА повторяющихся паттернов:
type EmailSender interface {
    Send(to, subject, body string) error
}

type FileStorage interface {
    Upload(filename string, data []byte) (string, error)
    Delete(filename string) error
}

type PaymentProcessor interface {
    ProcessPayment(amount float64, cardToken string) (*PaymentResult, error)
}
```

#### 2. Структура Use Case файла

```go
// ✅ ВСЕГДА: такая структура Use Case
package {feature}

import (
    "context"                                    // ✅ ВСЕГДА
    "clean-project/internal/domain/{entity}"    // ✅ ВСЕГДА: импорт domain
)

// ✅ ВСЕГДА: Request структура
type {Action}Request struct {  // 🔀 МЕНЯЕТСЯ: конкретные поля
    Field1 Type1 `json:"field1" validate:"required"`
    Field2 Type2 `json:"field2" validate:"omitempty"`
}

// ✅ ВСЕГДА: Response структура
type {Action}Response struct {  // 🔀 МЕНЯЕТСЯ: конкретные поля
    ID      uint   `json:"id"`
    Message string `json:"message"`
}

// ✅ ВСЕГДА: Use Case структура
type {Action}UseCase struct {
    // ✅ ВСЕГДА: зависимости через интерфейсы
    {entity}Repo {entity}.Repository  // 🔀 МЕНЯЕТСЯ: конкретные репозитории
    {service}    {Service}Interface   // 🔀 МЕНЯЕТСЯ: конкретные сервисы
}

// ✅ ВСЕГДА: конструктор с DI
func New{Action}UseCase(
    {entity}Repo {entity}.Repository,  // 🔀 МЕНЯЕТСЯ: параметры
    {service} {Service}Interface,
) *{Action}UseCase {
    return &{Action}UseCase{
        {entity}Repo: {entity}Repo,
        {service}:    {service},
    }
}

// ✅ ВСЕГДА: метод Execute с контекстом
func (uc *{Action}UseCase) Execute(ctx context.Context, req {Action}Request) (*{Action}Response, error) {
    // ✅ ВСЕГДА: валидация входных данных
    // ✅ ВСЕГДА: проверка бизнес-правил
    // ✅ ВСЕГДА: координация между компонентами
    // ✅ ВСЕГДА: возврат структурированного ответа
    
    return &{Action}Response{
        // поля ответа...
    }, nil
}
```

### 🔀 **Что МЕНЯЕТСЯ:**

- **Названия фич:** auth, payment, inventory, notification
- **Названия действий:** register, login, process, send, create
- **Поля Request/Response:** специфичные для каждого Use Case
- **Бизнес-логика:** уникальная для каждого сценария

---

## 🔌 Interface Adapters: повторяющиеся структуры

### ✅ **Что ВСЕГДА одинаково:**

#### 1. Структура Controller файла

```go
// ✅ ВСЕГДА: такая структура контроллера
package controllers

import (
    "encoding/json"                              // ✅ ВСЕГДА: для JSON
    "net/http"                                   // ✅ ВСЕГДА: для HTTP
    "clean-project/internal/usecase/{feature}"  // ✅ ВСЕГДА: импорт Use Cases
    "github.com/go-playground/validator/v10"    // ✅ ВСЕГДА: валидация
)

// ✅ ВСЕГДА: структура контроллера
type {Feature}Controller struct {
    // ✅ ВСЕГДА: зависимости от Use Cases
    {action1}UC *{feature}.{Action1}UseCase  // 🔀 МЕНЯЕТСЯ: конкретные Use Cases
    {action2}UC *{feature}.{Action2}UseCase
    validator   *validator.Validate          // ✅ ВСЕГДА: валидатор
}

// ✅ ВСЕГДА: конструктор контроллера
func New{Feature}Controller(
    {action1}UC *{feature}.{Action1}UseCase,  // 🔀 МЕНЯЕТСЯ: параметры
    {action2}UC *{feature}.{Action2}UseCase,
) *{Feature}Controller {
    return &{Feature}Controller{
        {action1}UC: {action1}UC,
        {action2}UC: {action2}UC,
        validator:   validator.New(),
    }
}

// ✅ ВСЕГДА: паттерн HTTP handler
func (c *{Feature}Controller) {Action}(w http.ResponseWriter, r *http.Request) {
    // ✅ ВСЕГДА: декодирование JSON
    var req {feature}.{Action}Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        c.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // ✅ ВСЕГДА: валидация
    if err := c.validator.Struct(req); err != nil {
        c.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // ✅ ВСЕГДА: вызов Use Case
    response, err := c.{action}UC.Execute(r.Context(), req)
    if err != nil {
        c.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // ✅ ВСЕГДА: возврат JSON ответа
    c.writeJSONResponse(w, response, http.StatusOK)
}

// ✅ ВСЕГДА: метод регистрации маршрутов
func (c *{Feature}Controller) RegisterRoutes(mux *http.ServeMux) {
    // 🔀 МЕНЯЕТСЯ: конкретные маршруты
    mux.HandleFunc("POST /{feature}/{action}", c.{Action})
}

// ✅ ВСЕГДА: вспомогательные методы
func (c *{Feature}Controller) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
    // стандартная реализация...
}

func (c *{Feature}Controller) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
    // стандартная реализация...
}
```

### 🔀 **Что МЕНЯЕТСЯ:**

- **Протокол:** HTTP, gRPC, GraphQL
- **Маршруты:** конкретные URL paths
- **Методы обработки:** названия и логика извлечения данных
- **Формат ответов:** JSON, Protobuf, XML

---

## 🌐 Infrastructure: что меняется, что нет

### ✅ **Что ВСЕГДА одинаково:**

#### 1. Паттерн реализации Repository

```go
// ✅ ВСЕГДА: такая структура репозитория
package database

import (
    "context"                                    // ✅ ВСЕГДА
    "clean-project/internal/domain/{entity}"    // ✅ ВСЕГДА: импорт domain
)

// ✅ ВСЕГДА: структура репозитория
type {Entity}Repository struct {
    db *DB  // 🔀 МЕНЯЕТСЯ: тип подключения (SQL, Mongo, etc.)
}

// ✅ ВСЕГДА: конструктор возвращает интерфейс
func New{Entity}Repository(db *DB) {entity}.Repository {
    return &{Entity}Repository{
        db: db,
    }
}

// ✅ ВСЕГДА: реализация методов интерфейса
func (r *{Entity}Repository) Save(ctx context.Context, e *{entity}.{Entity}) error {
    if e.ID == 0 {
        return r.create(ctx, e)   // ✅ ВСЕГДА: разделение create/update
    }
    return r.update(ctx, e)
}

// ✅ ВСЕГДА: приватные методы для CRUD
func (r *{Entity}Repository) create(ctx context.Context, e *{entity}.{Entity}) error {
    // 🔀 МЕНЯЕТСЯ: SQL, MongoDB, etc.
}

func (r *{Entity}Repository) update(ctx context.Context, e *{entity}.{Entity}) error {
    // 🔀 МЕНЯЕТСЯ: конкретная реализация
}
```

#### 2. Паттерн External Services

```go
// ✅ ВСЕГДА: такая структура внешнего сервиса
package external

import (
    "clean-project/internal/usecase/{feature}"  // ✅ ВСЕГДА: импорт интерфейса
)

// ✅ ВСЕГДА: структура сервиса
type {ConcreteService} struct {
    // 🔀 МЕНЯЕТСЯ: конфигурация конкретного сервиса
    apiKey string
    client HTTPClient
}

// ✅ ВСЕГДА: конструктор возвращает интерфейс
func New{ConcreteService}(config Config) {feature}.{Service}Interface {
    return &{ConcreteService}{
        apiKey: config.APIKey,
        client: newHTTPClient(),
    }
}

// ✅ ВСЕГДА: реализация методов интерфейса
func (s *{ConcreteService}) {Method}(params Type) (Type, error) {
    // 🔀 МЕНЯЕТСЯ: конкретная реализация (HTTP, gRPC, etc.)
}
```

### 🔀 **Что МЕНЯЕТСЯ:**

- **Тип БД:** PostgreSQL, MySQL, MongoDB, Redis
- **Протокол API:** REST, gRPC, GraphQL
- **Внешние сервисы:** JWT, OAuth, Payment gateways, Email services
- **Форматы данных:** JSON, Protobuf, XML

---

## 🧩 Dependency Injection: универсальный подход

### ✅ **Что ВСЕГДА одинаково:**

#### 1. Структура DI Container

```go
// ✅ ВСЕГДА: такая структура container.go
package di

// ✅ ВСЕГДА: импорты всех слоев
import (
    "clean-project/configs"
    "clean-project/internal/domain/{entity1}"      // ✅ ВСЕГДА: domain
    "clean-project/internal/usecase/{feature1}"    // ✅ ВСЕГДА: usecases
    "clean-project/internal/infrastructure/..."    // ✅ ВСЕГДА: infrastructure
    "clean-project/internal/interfaces/..."        // ✅ ВСЕГДА: interfaces
)

// ✅ ВСЕГДА: структура контейнера
type Container struct {
    // ✅ ВСЕГДА: конфигурация
    Config *configs.Config
    
    // ✅ ВСЕГДА: инфраструктура
    DB     *database.DB
    Server *Server
    
    // ✅ ВСЕГДА: репозитории как интерфейсы
    {Entity1}Repo {entity1}.Repository  // 🔀 МЕНЯЕТСЯ: конкретные репозитории
    {Entity2}Repo {entity2}.Repository
    
    // ✅ ВСЕГДА: внешние сервисы как интерфейсы
    {Service1} {feature}.{Service1}Interface  // 🔀 МЕНЯЕТСЯ: конкретные сервисы
    {Service2} {feature}.{Service2}Interface
    
    // ✅ ВСЕГДА: Use Cases
    {Action1}UC *{feature1}.{Action1}UseCase  // 🔀 МЕНЯЕТСЯ: конкретные Use Cases
    {Action2}UC *{feature1}.{Action2}UseCase
    
    // ✅ ВСЕГДА: контроллеры
    {Feature1}Controller *controllers.{Feature1}Controller  // 🔀 МЕНЯЕТСЯ: контроллеры
    {Feature2}Controller *controllers.{Feature2}Controller
}

// ✅ ВСЕГДА: порядок инициализации
func NewContainer(config *configs.Config) (*Container, error) {
    container := &Container{Config: config}
    
    // ✅ ВСЕГДА: такой порядок!
    if err := container.initInfrastructure(); err != nil {
        return nil, err
    }
    
    if err := container.initRepositories(); err != nil {
        return nil, err
    }
    
    if err := container.initExternalServices(); err != nil {
        return nil, err
    }
    
    if err := container.initUseCases(); err != nil {
        return nil, err
    }
    
    if err := container.initControllers(); err != nil {
        return nil, err
    }
    
    if err := container.initServer(); err != nil {
        return nil, err
    }
    
    return container, nil
}
```

#### 2. Паттерн инициализации методов

```go
// ✅ ВСЕГДА: такие методы инициализации
func (c *Container) initRepositories() error {
    // ✅ ВСЕГДА: создание конкретных реализаций
    // ✅ ВСЕГДА: присвоение к интерфейсам
    c.{Entity1}Repo = database.New{Entity1}Repository(c.DB)  // 🔀 МЕНЯЕТСЯ: имена
    c.{Entity2}Repo = database.New{Entity2}Repository(c.DB)
    return nil
}

func (c *Container) initUseCases() error {
    // ✅ ВСЕГДА: передача зависимостей через конструктор
    c.{Action1}UC = {feature}.New{Action1}UseCase(
        c.{Entity}Repo,    // 🔀 МЕНЯЕТСЯ: конкретные зависимости
        c.{Service},
    )
    return nil
}

func (c *Container) initControllers() error {
    // ✅ ВСЕГДА: передача Use Cases в контроллеры
    c.{Feature}Controller = controllers.New{Feature}Controller(
        c.{Action1}UC,  // 🔀 МЕНЯЕТСЯ: конкретные Use Cases
        c.{Action2}UC,
    )
    return nil
}
```

### 🔀 **Что МЕНЯЕТСЯ:**

- **Количество сущностей:** User, Product, Order, etc.
- **Количество фич:** auth, payment, inventory, etc.
- **Типы внешних сервисов:** email, payment, storage, etc.
- **Конфигурационные параметры:** специфичные для проекта

---

## 📝 Чек-лист для нового проекта

### ✅ **ОБЯЗАТЕЛЬНЫЕ шаги (ВСЕГДА делайте):**

#### 1. **Создание структуры папок**
```bash
# ✅ ВСЕГДА создавайте эти папки:
mkdir -p cmd
mkdir -p configs
mkdir -p internal/domain
mkdir -p internal/usecase
mkdir -p internal/infrastructure/{database,external,web}
mkdir -p internal/interfaces/controllers
mkdir -p pkg/di
```

#### 2. **Определение доменных сущностей**
```go
// ✅ ВСЕГДА начинайте с domain слоя:

// 1. Определите основные сущности вашего проекта
internal/domain/user/entity.go
internal/domain/product/entity.go
internal/domain/order/entity.go

// 2. Создайте интерфейсы репозиториев
internal/domain/user/repository.go
internal/domain/product/repository.go
internal/domain/order/repository.go
```

#### 3. **Создание Use Cases**
```go
// ✅ ВСЕГДА определите основные сценарии:

// По фичам/функциональности:
internal/usecase/auth/register.go
internal/usecase/auth/login.go
internal/usecase/catalog/get_products.go
internal/usecase/order/create_order.go
```

#### 4. **Реализация Infrastructure**
```go
// ✅ ВСЕГДА реализуйте интерфейсы:

// Репозитории:
internal/infrastructure/database/user_repository.go
internal/infrastructure/database/product_repository.go

// Внешние сервисы:
internal/infrastructure/external/email_service.go
internal/infrastructure/external/payment_service.go
```

#### 5. **Создание Adapters**
```go
// ✅ ВСЕГДА создавайте контроллеры:

internal/interfaces/controllers/auth_controller.go
internal/interfaces/controllers/catalog_controller.go
internal/interfaces/controllers/order_controller.go
```

#### 6. **Настройка DI**
```go
// ✅ ВСЕГДА создавайте DI контейнер:

pkg/di/container.go
// с методами init* в правильном порядке
```

### 🔀 **Вариативные шаги (зависят от проекта):**

#### Выберите технологии:
- **База данных:** PostgreSQL, MySQL, MongoDB
- **Протокол API:** HTTP REST, gRPC, GraphQL  
- **Аутентификация:** JWT, OAuth, Session
- **Внешние сервисы:** Email, Payment, Storage

#### Адаптируйте структуры:
- **Доменные модели:** под вашу предметную область
- **Use Cases:** под ваши бизнес-сценарии
- **API endpoints:** под ваши требования

---

## 🎯 Резюме: что запомнить навсегда

### 🔒 **НЕИЗМЕННЫЕ принципы Clean Architecture:**

1. **📁 Структура слоев:** domain → usecase → interfaces → infrastructure
2. **🔄 Направление зависимостей:** ВСЕГДА внутрь (к domain)
3. **🧩 Dependency Injection:** зависимости через интерфейсы
4. **📏 Правила именования:** entity.go, repository.go, {action}.go
5. **🏗️ Порядок создания:** infrastructure → repositories → services → usecases → controllers

### 🔀 **ИЗМЕНЯЕМЫЕ элементы:**

1. **📊 Доменные модели:** зависят от предметной области
2. **🎯 Use Cases:** зависят от бизнес-требований  
3. **🗄️ Технологии:** БД, протоколы, внешние сервисы
4. **🌐 API дизайн:** REST, GraphQL, gRPC

### 💡 **Мантра Clean Architecture:**

> **"Структура и принципы - постоянны, реализация - изменяема"**

Запомните паттерн один раз - применяйте везде! 🎉