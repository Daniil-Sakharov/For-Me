package di

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	
	"clean-url-shortener/configs"
	"clean-url-shortener/internal/domain/user"
	"clean-url-shortener/internal/domain/link"
	"clean-url-shortener/internal/usecase/auth"
	linkUC "clean-url-shortener/internal/usecase/link"
	"clean-url-shortener/internal/infrastructure/database"
	"clean-url-shortener/internal/infrastructure/external"
	"clean-url-shortener/internal/infrastructure/web"
	"clean-url-shortener/internal/interfaces/controllers"
)

// Container содержит все зависимости приложения
// ЭТО ГЛАВНЫЙ КОНТЕЙНЕР ДЛЯ DEPENDENCY INJECTION
type Container struct {
	// Конфигурация
	Config *configs.Config
	
	// Инфраструктура
	DB     *database.DB
	Server *web.Server
	
	// Репозитории (реализации интерфейсов из domain слоя)
	UserRepo user.Repository
	LinkRepo link.Repository
	
	// Внешние сервисы (реализации интерфейсов из usecase слоя)
	PasswordHasher auth.PasswordHasher
	TokenGenerator auth.TokenGenerator
	URLValidator   linkUC.URLValidator
	ShortCodeGen   linkUC.ShortCodeGenerator
	
	// Use Cases (бизнес-логика приложения)
	RegisterUC   *auth.RegisterUseCase
	LoginUC      *auth.LoginUseCase
	CreateLinkUC *linkUC.CreateLinkUseCase
	RedirectUC   *linkUC.RedirectUseCase
	
	// Middleware
	AuthMiddleware *web.AuthMiddleware
	
	// Контроллеры (HTTP handlers)
	AuthController *controllers.AuthController
	LinkController *controllers.LinkController
}

// NewContainer создает и настраивает контейнер зависимостей
// ЭТО МЕСТО ГДЕ ПРОИСХОДИТ ВИРИРОВАНИЕ (WIRING) ВСЕХ КОМПОНЕНТОВ
func NewContainer(config *configs.Config) (*Container, error) {
	container := &Container{
		Config: config,
	}
	
	// Инициализируем компоненты в правильном порядке
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
	
	if err := container.initMiddleware(); err != nil {
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

// initInfrastructure инициализирует инфраструктурные компоненты
func (c *Container) initInfrastructure() error {
	// Подключение к базе данных
	dbConfig := database.Config{
		Host:     c.Config.Database.Host,
		Port:     c.Config.Database.Port,
		User:     c.Config.Database.User,
		Password: c.Config.Database.Password,
		DBName:   c.Config.Database.DBName,
		SSLMode:  c.Config.Database.SSLMode,
	}
	
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		return err
	}
	
	// Выполняем миграции
	if err := db.RunMigrations(); err != nil {
		return err
	}
	
	c.DB = db
	return nil
}

// initRepositories инициализирует репозитории
func (c *Container) initRepositories() error {
	// Создаем репозитории, которые реализуют интерфейсы из domain слоя
	c.UserRepo = database.NewUserRepository(c.DB)
	c.LinkRepo = database.NewLinkRepository(c.DB)
	
	return nil
}

// initExternalServices инициализирует внешние сервисы
func (c *Container) initExternalServices() error {
	// Password hasher (реализация интерфейса из usecase/auth)
	c.PasswordHasher = external.NewBcryptPasswordHasher(c.Config.Auth.BcryptCost)
	
	// JWT token generator (реализация интерфейса из usecase/auth)
	c.TokenGenerator = external.NewJWTTokenGenerator(
		c.Config.Auth.JWTSecret,
		c.Config.Auth.JWTIssuer,
		c.Config.Auth.JWTExpiry,
	)
	
	// URL validator (простая реализация для примера)
	c.URLValidator = &SimpleURLValidator{}
	
	// Short code generator (простая реализация для примера)
	c.ShortCodeGen = &SimpleShortCodeGenerator{}
	
	return nil
}

// initUseCases инициализирует Use Cases
func (c *Container) initUseCases() error {
	// Auth Use Cases
	c.RegisterUC = auth.NewRegisterUseCase(
		c.UserRepo,
		c.PasswordHasher,
		c.TokenGenerator,
	)
	
	c.LoginUC = auth.NewLoginUseCase(
		c.UserRepo,
		c.PasswordHasher,
		c.TokenGenerator,
	)
	
	// Link Use Cases
	c.CreateLinkUC = linkUC.NewCreateLinkUseCase(
		c.LinkRepo,
		c.URLValidator,
		c.ShortCodeGen,
		c.Config.App.BaseURL,
	)
	
	// Для RedirectUC нужен EventPublisher, создаем заглушку
	eventPublisher := &NoOpEventPublisher{}
	
	c.RedirectUC = linkUC.NewRedirectUseCase(
		c.LinkRepo,
		eventPublisher,
	)
	
	return nil
}

// initMiddleware инициализирует middleware
func (c *Container) initMiddleware() error {
	c.AuthMiddleware = web.NewAuthMiddleware(c.TokenGenerator)
	return nil
}

// initControllers инициализирует контроллеры
func (c *Container) initControllers() error {
	c.AuthController = controllers.NewAuthController(
		c.RegisterUC,
		c.LoginUC,
	)
	
	c.LinkController = controllers.NewLinkController(
		c.CreateLinkUC,
		c.RedirectUC,
		c.AuthMiddleware,
	)
	
	return nil
}

// initServer инициализирует HTTP сервер
func (c *Container) initServer() error {
	// Создаем роутер и регистрируем маршруты
	mux := c.setupRoutes()
	
	// Создаем HTTP сервер
	serverConfig := web.Config{
		Host:         c.Config.Server.Host,
		Port:         c.Config.Server.Port,
		ReadTimeout:  c.Config.Server.ReadTimeout,
		WriteTimeout: c.Config.Server.WriteTimeout,
		IdleTimeout:  c.Config.Server.IdleTimeout,
	}
	
	c.Server = web.NewServer(serverConfig, mux)
	return nil
}

// setupRoutes настраивает маршруты приложения
func (c *Container) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	
	// Регистрируем маршруты контроллеров
	c.AuthController.RegisterRoutes(mux)
	c.LinkController.RegisterRoutes(mux)
	
	// Добавляем health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"url-shortener"}`))
	})
	
	return mux
}

// Cleanup освобождает ресурсы
func (c *Container) Cleanup() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// ПРОСТЫЕ РЕАЛИЗАЦИИ ДЛЯ ПРИМЕРА
// В реальном проекте эти компоненты были бы более сложными

// SimpleURLValidator простая реализация валидатора URL
type SimpleURLValidator struct{}

func (v *SimpleURLValidator) Validate(rawURL string) error {
	_, err := url.Parse(rawURL)
	return err
}

func (v *SimpleURLValidator) IsSafe(rawURL string) (bool, error) {
	// Простая проверка - в реальном проекте здесь была бы проверка на malware/phishing
	return true, nil
}

// SimpleShortCodeGenerator простая реализация генератора коротких кодов
type SimpleShortCodeGenerator struct{}

func (g *SimpleShortCodeGenerator) Generate() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	
	return string(bytes), nil
}

func (g *SimpleShortCodeGenerator) GenerateCustom(customCode string) (string, error) {
	// Простая валидация - в реальном проекте здесь была бы более сложная логика
	if len(customCode) < 3 || len(customCode) > 10 {
		return "", fmt.Errorf("custom code must be between 3 and 10 characters")
	}
	return customCode, nil
}

// NoOpEventPublisher заглушка для event publisher
type NoOpEventPublisher struct{}

func (p *NoOpEventPublisher) PublishLinkClicked(linkID uint, userAgent, ipAddress, referer string) error {
	// В реальном проекте здесь была бы публикация событий
	return nil
}

// ПРИНЦИПЫ DEPENDENCY INJECTION:
// 1. Все зависимости создаются в одном месте
// 2. Соблюдается правильный порядок инициализации
// 3. Интерфейсы используются для развязки компонентов
// 4. Конфигурация централизована
// 5. Легко тестировать (можно заменить любую зависимость)