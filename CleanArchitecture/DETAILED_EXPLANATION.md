# 🔍 Подробное объяснение Clean Architecture на примерах

Этот файл содержит **максимально подробное** объяснение того, как работает наше приложение с Clean Architecture. Мы разберем каждый шаг, каждую зависимость и каждый принцип на конкретных примерах.

---

## 📚 Содержание

- [🚀 Запуск приложения: что происходит?](#-запуск-приложения-что-происходит)
- [👤 Регистрация пользователя: пошаговый разбор](#-регистрация-пользователя-пошаговый-разбор)
- [🔐 Вход пользователя: как работает аутентификация](#-вход-пользователя-как-работает-аутентификация)
- [🔗 Создание ссылки: полный цикл](#-создание-ссылки-полный-цикл)
- [🌐 Переход по короткой ссылке: редирект](#-переход-по-короткой-ссылке-редирект)
- [🧩 Dependency Injection: зачем и как](#-dependency-injection-зачем-и-как)
- [🏗️ Слои архитектуры: что где находится](#️-слои-архитектуры-что-где-находится)
- [🔄 Инверсия зависимостей: почему это важно](#-инверсия-зависимостей-почему-это-важно)

---

## 🚀 Запуск приложения: что происходит?

### Шаг 1: Пользователь запускает `go run cmd/main.go`

Когда вы выполняете эту команду, начинается выполнение функции `main()` в файле `cmd/main.go`:

```go
func main() {
    fmt.Println("🚀 Starting URL Shortener Service...")
    
    // ШАГ 1: Загружаем конфигурацию
    config, err := configs.LoadConfig()
```

**Что происходит:** Программа читает переменные окружения и создает структуру `Config` со всеми настройками.

### Шаг 2: Создание DI контейнера

```go
// ШАГ 2: Создаем контейнер зависимостей
container, err := di.NewContainer(config)
```

**Это самая важная часть!** Здесь создаются **ВСЕ** компоненты нашего приложения. Давайте разберем подробно:

#### 2.1. Подключение к базе данных

```go
// В pkg/di/container.go, метод initInfrastructure()
dbConfig := database.Config{
    Host:     "localhost",
    Port:     5432,
    User:     "postgres",
    Password: "postgres",
    DBName:   "url_shortener",
    SSLMode:  "disable",
}

db, err := database.NewConnection(dbConfig)
```

**Что происходит:** Создается подключение к PostgreSQL и выполняются SQL миграции для создания таблиц.

#### 2.2. Создание репозиториев

```go
// В pkg/di/container.go, метод initRepositories()
c.UserRepo = database.NewUserRepository(c.DB)
c.LinkRepo = database.NewLinkRepository(c.DB)
```

**Важно:** Мы создаем **конкретные реализации** репозиториев, но сохраняем их как **интерфейсы**:

```go
type Container struct {
    UserRepo user.Repository  // ← Это ИНТЕРФЕЙС из domain слоя!
    LinkRepo link.Repository  // ← Это тоже ИНТЕРФЕЙС!
}
```

#### 2.3. Создание внешних сервисов

```go
// В pkg/di/container.go, метод initExternalServices()
c.PasswordHasher = external.NewBcryptPasswordHasher(12)
c.TokenGenerator = external.NewJWTTokenGenerator("secret", "issuer", 24*time.Hour)
```

**Зачем:** Use Cases будут использовать эти сервисы через интерфейсы, не зная о конкретных реализациях.

#### 2.4. Создание Use Cases

```go
// В pkg/di/container.go, метод initUseCases()
c.RegisterUC = auth.NewRegisterUseCase(
    c.UserRepo,        // ← Передаем интерфейс репозитория
    c.PasswordHasher,  // ← Передаем интерфейс хешера
    c.TokenGenerator,  // ← Передаем интерфейс токен-генератора
)
```

**Принцип:** Use Case получает ВСЕ свои зависимости через конструктор. Он не создает их сам!

#### 2.5. Создание контроллеров

```go
// В pkg/di/container.go, метод initControllers()
c.AuthController = controllers.NewAuthController(
    c.RegisterUC,  // ← Передаем Use Case для регистрации
    c.LoginUC,     // ← Передаем Use Case для входа
)
```

### Шаг 3: Запуск HTTP сервера

```go
// В main.go
go func() {
    if err := container.Server.Start(); err != nil {
        serverErrors <- err
    }
}()
```

**Результат:** HTTP сервер слушает на порту 8080 и готов принимать запросы.

---

## 👤 Регистрация пользователя: пошаговый разбор

Теперь давайте проследим **весь путь** запроса на регистрацию:

### Пользователь отправляет HTTP запрос:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "strongpassword123",
    "name": "John Doe"
  }'
```

### ШАГ 1: HTTP Router направляет запрос в контроллер

В `pkg/di/container.go` мы зарегистрировали маршрут:

```go
mux.HandleFunc("POST /auth/register", web.ChainMiddleware(
    c.AuthController.Register,
    web.CORSMiddleware,
    web.LoggingMiddleware,
    web.RecoveryMiddleware,
))
```

**Что происходит:** HTTP запрос проходит через middleware (CORS, логирование, recovery) и попадает в `AuthController.Register()`.

### ШАГ 2: AuthController обрабатывает HTTP запрос

Файл: `internal/interfaces/controllers/auth_controller.go`

```go
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
    // 2.1: Декодируем JSON из тела запроса
    var req auth.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        c.writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }
    
    // Теперь req содержит:
    // req.Email = "john@example.com"
    // req.Password = "strongpassword123"  
    // req.Name = "John Doe"
    
    // 2.2: Валидируем входные данные
    if err := c.validator.Struct(req); err != nil {
        c.writeErrorResponse(w, "Validation failed", http.StatusBadRequest)
        return
    }
    
    // 2.3: Вызываем Use Case для регистрации
    response, err := c.registerUC.Execute(r.Context(), req)
```

**Важно:** Контроллер НЕ содержит бизнес-логики! Он только:
- Извлекает данные из HTTP запроса
- Валидирует формат
- Вызывает Use Case
- Возвращает HTTP ответ

### ШАГ 3: RegisterUseCase выполняет бизнес-логику

Файл: `internal/usecase/auth/register.go`

```go
func (uc *RegisterUseCase) Execute(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
    // 3.1: Проверяем, что пользователь с таким email не существует
    exists, err := uc.userRepo.Exists(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, ErrUserAlreadyExists  // Бизнес-правило!
    }
```

**Вызов базы данных:** `uc.userRepo.Exists()` - это вызов метода интерфейса:

```go
type Repository interface {
    Exists(ctx context.Context, email string) (bool, error)
}
```

Реальная реализация находится в `internal/infrastructure/database/user_repository.go`:

```go
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    
    var exists bool
    err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
    return exists, err
}
```

**SQL запрос:** `SELECT EXISTS(SELECT 1 FROM users WHERE email = 'john@example.com')`

Продолжаем Use Case:

```go
    // 3.2: Валидируем силу пароля (бизнес-правило)
    if len(req.Password) < 8 {
        return nil, ErrWeakPassword
    }
    
    // 3.3: Создаем доменную модель пользователя
    newUser, err := user.NewUser(req.Email, req.Name)
    if err != nil {
        return nil, err
    }
```

### ШАГ 4: Создание доменной модели

Файл: `internal/domain/user/entity.go`

```go
func NewUser(email, name string) (*User, error) {
    // Валидация email согласно доменным правилам
    if !isValidEmail(email) {
        return nil, ErrInvalidEmail
    }
    
    // Валидация имени
    if name == "" {
        return nil, ErrEmptyName
    }
    
    now := time.Now()
    
    return &User{
        Email:     email,      // "john@example.com"
        Name:      name,       // "John Doe"
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}
```

**Важно:** Доменная модель проверяет БИЗНЕС-ПРАВИЛА (валидный email, непустое имя).

Возвращаемся в Use Case:

```go
    // 3.4: Хешируем пароль через внешний сервис
    hashedPassword, err := uc.passwordHasher.Hash(req.Password)
    if err != nil {
        return nil, err
    }
```

### ШАГ 5: Хеширование пароля

`uc.passwordHasher` - это интерфейс:

```go
type PasswordHasher interface {
    Hash(password string) (string, error)
}
```

Реализация в `internal/infrastructure/external/password_hasher.go`:

```go
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
    return string(hashedBytes), nil
}
```

**Результат:** `"strongpassword123"` превращается в `"$2a$12$..."` (bcrypt hash).

Продолжаем Use Case:

```go
    // 3.5: Устанавливаем захешированный пароль
    err = newUser.SetPassword(hashedPassword)
    if err != nil {
        return nil, err
    }
    
    // 3.6: Сохраняем пользователя в базе данных
    err = uc.userRepo.Save(ctx, newUser)
    if err != nil {
        return nil, err
    }
```

### ШАГ 6: Сохранение в базе данных

`uc.userRepo.Save()` вызывает реализацию в `internal/infrastructure/database/user_repository.go`:

```go
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
    if u.ID == 0 {
        return r.create(ctx, u)  // Создаем нового пользователя
    }
    return r.update(ctx, u)      // Обновляем существующего
}

func (r *UserRepository) create(ctx context.Context, u *user.User) error {
    query := `
        INSERT INTO users (email, hashed_password, name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

    err := r.db.QueryRowContext(
        ctx,
        query,
        u.Email,         // "john@example.com"
        u.HashedPassword, // "$2a$12$..."
        u.Name,          // "John Doe"
        u.CreatedAt,
        u.UpdatedAt,
    ).Scan(&u.ID)
    
    return err
}
```

**SQL запрос:**
```sql
INSERT INTO users (email, hashed_password, name, created_at, updated_at)
VALUES ('john@example.com', '$2a$12$...', 'John Doe', '2024-01-01 12:00:00', '2024-01-01 12:00:00')
RETURNING id
```

**Результат:** Пользователь сохранен в БД, получен ID (например, 1).

Финал Use Case:

```go
    // 3.7: Генерируем JWT токен для автоматического входа
    token, err := uc.tokenGenerator.Generate(newUser.ID, newUser.Email)
    if err != nil {
        return nil, err
    }
    
    // 3.8: Возвращаем ответ
    return &RegisterResponse{
        Token:  token,
        UserID: newUser.ID,
        Email:  newUser.Email,
        Name:   newUser.Name,
    }, nil
}
```

### ШАГ 7: Генерация JWT токена

`uc.tokenGenerator.Generate()` вызывает реализацию в `internal/infrastructure/external/jwt_generator.go`:

```go
func (g *JWTTokenGenerator) Generate(userID uint, email string) (string, error) {
    claims := Claims{
        UserID: userID,    // 1
        Email:  email,     // "john@example.com"
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    g.issuer,
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.expiry)),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(g.secretKey)
}
```

**Результат:** JWT токен типа `"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

### ШАГ 8: Возврат HTTP ответа

Контроллер получает ответ от Use Case и отправляет HTTP ответ:

```go
// В AuthController.Register()
response, err := c.registerUC.Execute(r.Context(), req)
if err != nil {
    c.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
    return
}

c.writeJSONResponse(w, response, http.StatusCreated)
```

**HTTP ответ:**
```json
HTTP/1.1 201 Created
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": 1,
  "email": "john@example.com",
  "name": "John Doe"
}
```

---

## 🔐 Вход пользователя: как работает аутентификация

Теперь пользователь хочет войти в систему:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "strongpassword123"
  }'
```

### ШАГ 1-2: HTTP Router → AuthController.Login()

Аналогично регистрации, запрос попадает в `AuthController.Login()`.

### ШАГ 3: LoginUseCase выполняет аутентификацию

Файл: `internal/usecase/auth/login.go`

```go
func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
    // 3.1: Находим пользователя по email
    existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser == nil {
        return nil, ErrUserNotFound
    }
```

**SQL запрос в UserRepository:**
```sql
SELECT id, email, hashed_password, name, created_at, updated_at
FROM users
WHERE email = 'john@example.com'
```

**Результат:** Найден пользователь с ID=1, именем "John Doe" и хешем пароля.

```go
    // 3.2: Проверяем пароль
    err = uc.passwordHasher.Compare(existingUser.HashedPassword, req.Password)
    if err != nil {
        return nil, ErrInvalidCredentials
    }
```

### ШАГ 4: Проверка пароля

`uc.passwordHasher.Compare()` в `internal/infrastructure/external/password_hasher.go`:

```go
func (h *BcryptPasswordHasher) Compare(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
```

**Что происходит:** 
- `hashedPassword` = `"$2a$12$..."` (из БД)
- `password` = `"strongpassword123"` (от пользователя)
- bcrypt проверяет соответствие

**Результат:** Если пароли совпадают, ошибки нет. Если не совпадают - ошибка.

```go
    // 3.3: Генерируем JWT токен
    token, err := uc.tokenGenerator.Generate(existingUser.ID, existingUser.Email)
    
    return &LoginResponse{
        Token:  token,
        UserID: existingUser.ID,
        Email:  existingUser.Email,
        Name:   existingUser.Name,
    }, nil
}
```

**HTTP ответ:** Аналогичен регистрации, но статус 200 OK.

---

## 🔗 Создание ссылки: полный цикл

Теперь авторизованный пользователь хочет создать короткую ссылку:

```bash
curl -X POST http://localhost:8080/links \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "original_url": "https://very-long-url.com/some/very/long/path",
    "custom_code": "my-link"
  }'
```

### ШАГ 1: HTTP Router → Middleware → LinkController

```go
// В pkg/di/container.go
mux.HandleFunc("POST /links", web.ChainMiddleware(
    c.LinkController.CreateLink,
    c.authMiddleware.RequireAuth, // ← ВАЖНО: Сначала аутентификация!
    web.CORSMiddleware,
    web.LoggingMiddleware,
    web.RecoveryMiddleware,
))
```

### ШАГ 2: AuthMiddleware проверяет JWT токен

Файл: `internal/infrastructure/web/middleware.go`

```go
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 2.1: Извлекаем токен из заголовка Authorization
        authHeader := r.Header.Get("Authorization")
        // authHeader = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        
        // 2.2: Проверяем формат "Bearer <token>"
        parts := strings.SplitN(authHeader, " ", 2)
        token := parts[1] // "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        
        // 2.3: Валидируем токен
        tokenData, err := m.tokenGenerator.Validate(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
```

**Валидация JWT токена** в `internal/infrastructure/external/jwt_generator.go`:

```go
func (g *JWTTokenGenerator) Validate(tokenString string) (*auth.TokenData, error) {
    // Парсим и проверяем подпись токена
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return g.secretKey, nil
    })
    
    claims, ok := token.Claims.(*Claims)
    
    return &auth.TokenData{
        UserID: claims.UserID,  // 1
        Email:  claims.Email,   // "john@example.com"
    }, nil
}
```

**Результат:** Токен валиден, получены UserID=1 и Email="john@example.com".

```go
        // 2.4: Добавляем данные пользователя в контекст
        ctx := context.WithValue(r.Context(), UserIDKey, tokenData.UserID)
        ctx = context.WithValue(ctx, UserEmailKey, tokenData.Email)
        r = r.WithContext(ctx)
        
        // 2.5: Передаем управление следующему handler
        next(w, r)
    }
}
```

**Важно:** Данные пользователя сохраняются в контексте запроса и будут доступны в контроллере!

### ШАГ 3: LinkController.CreateLink()

Файл: `internal/interfaces/controllers/link_controller.go`

```go
func (c *LinkController) CreateLink(w http.ResponseWriter, r *http.Request) {
    // 3.1: Извлекаем ID пользователя из контекста
    userID, ok := web.GetUserIDFromContext(r.Context())
    if !ok {
        c.writeErrorResponse(w, "User not authenticated", http.StatusUnauthorized)
        return
    }
    // userID = 1 (из JWT токена)
    
    // 3.2: Декодируем JSON из тела запроса
    var req link.CreateLinkRequest
    json.NewDecoder(r.Body).Decode(&req)
    // req.OriginalURL = "https://very-long-url.com/some/very/long/path"
    // req.CustomCode = "my-link"
    
    // 3.3: Устанавливаем ID пользователя в запрос
    req.UserID = userID  // 1
    
    // 3.4: Вызываем Use Case
    response, err := c.createLinkUC.Execute(r.Context(), req)
```

### ШАГ 4: CreateLinkUseCase выполняет бизнес-логику

Файл: `internal/usecase/link/create.go`

```go
func (uc *CreateLinkUseCase) Execute(ctx context.Context, req CreateLinkRequest) (*CreateLinkResponse, error) {
    // 4.1: Валидируем URL на корректность
    err := uc.urlValidator.Validate(req.OriginalURL)
    if err != nil {
        return nil, err
    }
```

**Валидация URL** в `pkg/di/container.go` (SimpleURLValidator):

```go
func (v *SimpleURLValidator) Validate(rawURL string) error {
    _, err := url.Parse(rawURL)  // Проверяем, что URL корректный
    return err
}
```

```go
    // 4.2: Проверяем безопасность URL
    isSafe, err := uc.urlValidator.IsSafe(req.OriginalURL)
    if !isSafe {
        return nil, ErrUnsafeURL
    }
    
    // 4.3: Создаем ссылку с пользовательским кодом
    if req.CustomCode != "" {
        newLink, err = uc.createLinkWithCustomCode(ctx, req)
    }
```

### ШАГ 5: Создание ссылки с пользовательским кодом

```go
func (uc *CreateLinkUseCase) createLinkWithCustomCode(ctx context.Context, req CreateLinkRequest) (*link.Link, error) {
    // 5.1: Проверяем уникальность пользовательского кода
    exists, err := uc.linkRepo.ExistsByShortCode(ctx, req.CustomCode)
    if exists {
        return nil, ErrShortCodeAlreadyExists
    }
```

**SQL запрос в LinkRepository:**
```sql
SELECT EXISTS(SELECT 1 FROM links WHERE short_code = 'my-link')
```

**Результат:** Предположим, код "my-link" еще не занят.

```go
    // 5.2: Создаем доменную модель ссылки
    return link.NewLinkWithCustomCode(req.OriginalURL, req.CustomCode, req.UserID)
}
```

### ШАГ 6: Создание доменной модели Link

Файл: `internal/domain/link/entity.go`

```go
func NewLinkWithCustomCode(originalURL, shortCode string, userID uint) (*Link, error) {
    // 6.1: Валидация URL согласно доменным правилам
    if err := validateURL(originalURL); err != nil {
        return nil, err
    }
    
    // 6.2: Валидация пользовательского короткого кода
    if err := validateShortCode(shortCode); err != nil {
        return nil, err
    }
    
    now := time.Now()
    
    return &Link{
        OriginalURL: originalURL,  // "https://very-long-url.com/some/very/long/path"
        ShortCode:   shortCode,    // "my-link"
        UserID:      userID,       // 1
        CreatedAt:   now,
        UpdatedAt:   now,
        ClicksCount: 0,
    }, nil
}
```

**Доменная валидация:**
- URL должен иметь схему (http/https) и хост
- Короткий код должен содержать только буквы и цифры, длина 3-10 символов

Возвращаемся в Use Case:

```go
    // 4.4: Сохраняем ссылку в базе данных
    err = uc.linkRepo.Save(ctx, newLink)
    if err != nil {
        return nil, err
    }
```

### ШАГ 7: Сохранение ссылки в БД

`uc.linkRepo.Save()` в `internal/infrastructure/database/link_repository.go`:

```go
func (r *LinkRepository) Save(ctx context.Context, l *link.Link) error {
    if l.ID == 0 {
        return r.create(ctx, l)
    }
    return r.update(ctx, l)
}

func (r *LinkRepository) create(ctx context.Context, l *link.Link) error {
    query := `
        INSERT INTO links (original_url, short_code, user_id, clicks_count, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

    err := r.db.QueryRowContext(
        ctx,
        query,
        l.OriginalURL,  // "https://very-long-url.com/some/very/long/path"
        l.ShortCode,    // "my-link"
        l.UserID,       // 1
        l.ClicksCount,  // 0
        l.CreatedAt,
        l.UpdatedAt,
    ).Scan(&l.ID)
    
    return err
}
```

**SQL запрос:**
```sql
INSERT INTO links (original_url, short_code, user_id, clicks_count, created_at, updated_at)
VALUES ('https://very-long-url.com/some/very/long/path', 'my-link', 1, 0, '2024-01-01 12:00:00', '2024-01-01 12:00:00')
RETURNING id
```

**Результат:** Ссылка сохранена с ID=1.

Финал Use Case:

```go
    // 4.5: Формируем ответ
    return &CreateLinkResponse{
        ID:          newLink.ID,          // 1
        OriginalURL: newLink.OriginalURL, // "https://very-long-url.com/some/very/long/path"
        ShortCode:   newLink.ShortCode,   // "my-link"
        ShortURL:    uc.baseURL + "/" + newLink.ShortCode,  // "http://localhost:8080/my-link"
        CreatedAt:   newLink.CreatedAt.Format("2006-01-02T15:04:05Z"),
    }, nil
}
```

**HTTP ответ:**
```json
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": 1,
  "original_url": "https://very-long-url.com/some/very/long/path",
  "short_code": "my-link",
  "short_url": "http://localhost:8080/my-link",
  "created_at": "2024-01-01T12:00:00Z"
}
```

---

## 🌐 Переход по короткой ссылке: редирект

Теперь любой пользователь (даже неавторизованный) переходит по короткой ссылке:

```bash
curl -L http://localhost:8080/my-link
```

### ШАГ 1: HTTP Router → LinkController.Redirect()

```go
// В pkg/di/container.go
mux.HandleFunc("GET /", web.ChainMiddleware(
    c.LinkController.Redirect,
    // НЕТ authMiddleware! Редирект доступен всем
    web.CORSMiddleware,
    web.LoggingMiddleware,
    web.RecoveryMiddleware,
))
```

**Важно:** Маршрут `GET /` обрабатывает любой путь, включая `/my-link`.

### ШАГ 2: LinkController.Redirect()

Файл: `internal/interfaces/controllers/link_controller.go`

```go
func (c *LinkController) Redirect(w http.ResponseWriter, r *http.Request) {
    // 2.1: Извлекаем короткий код из URL пути
    path := strings.TrimPrefix(r.URL.Path, "/")
    // r.URL.Path = "/my-link"
    // path = "my-link"
    
    if path == "" {
        c.writeErrorResponse(w, "Short code is required", http.StatusBadRequest)
        return
    }
    
    // 2.2: Создаем запрос для Use Case
    req := link.RedirectRequest{
        ShortCode: path,                           // "my-link"
        UserAgent: r.Header.Get("User-Agent"),     // "Mozilla/5.0..."
        IPAddress: c.getClientIP(r),               // "192.168.1.100"
        Referer:   r.Header.Get("Referer"),        // "https://google.com"
    }
    
    // 2.3: Вызываем Use Case
    response, err := c.redirectUC.Execute(r.Context(), req)
```

### ШАГ 3: RedirectUseCase находит ссылку

Файл: `internal/usecase/link/redirect.go`

```go
func (uc *RedirectUseCase) Execute(ctx context.Context, req RedirectRequest) (*RedirectResponse, error) {
    // 3.1: Находим ссылку по короткому коду
    foundLink, err := uc.linkRepo.FindByShortCode(ctx, req.ShortCode)
    if foundLink == nil {
        return nil, ErrLinkNotFound
    }
```

**SQL запрос в LinkRepository:**
```sql
SELECT id, original_url, short_code, user_id, clicks_count, created_at, updated_at
FROM links
WHERE short_code = 'my-link'
```

**Результат:** Найдена ссылка с ID=1, OriginalURL="https://very-long-url.com/some/very/long/path".

```go
    // 3.2: Увеличиваем счетчик кликов в доменной модели
    foundLink.IncrementClicks()
```

### ШАГ 4: Доменная операция увеличения счетчика

Файл: `internal/domain/link/entity.go`

```go
func (l *Link) IncrementClicks() {
    l.ClicksCount++           // Было 0, стало 1
    l.UpdatedAt = time.Now()  // Обновляем время
}
```

**Важно:** Это доменная операция! Бизнес-правило: каждый переход увеличивает счетчик.

Продолжаем Use Case:

```go
    // 3.3: Обновляем счетчик в базе данных
    err = uc.linkRepo.IncrementClicks(ctx, foundLink.ID)
    if err != nil {
        // Логируем ошибку, но НЕ прерываем редирект
        // Пользователь должен быть перенаправлен даже если счетчик не обновился
    }
```

**SQL запрос для обновления счетчика:**
```sql
UPDATE links SET clicks_count = clicks_count + 1, updated_at = NOW() WHERE id = 1
```

```go
    // 3.4: Публикуем событие о клике для сбора статистики
    err = uc.eventPublisher.PublishLinkClicked(
        foundLink.ID,    // 1
        req.UserAgent,   // "Mozilla/5.0..."
        req.IPAddress,   // "192.168.1.100"
        req.Referer,     // "https://google.com"
    )
    // В нашей реализации это заглушка, но в реальном проекте
    // здесь была бы публикация события для асинхронной обработки
    
    // 3.5: Возвращаем URL для редиректа
    return &RedirectResponse{
        OriginalURL: foundLink.OriginalURL,  // "https://very-long-url.com/some/very/long/path"
        LinkID:      foundLink.ID,           // 1
    }, nil
}
```

### ШАГ 5: HTTP редирект

Контроллер получает ответ и выполняет редирект:

```go
// В LinkController.Redirect()
response, err := c.redirectUC.Execute(r.Context(), req)
if err != nil {
    c.writeErrorResponse(w, err.Error(), http.StatusNotFound)
    return
}

// Выполняем HTTP редирект
http.Redirect(w, r, response.OriginalURL, http.StatusFound)
```

**HTTP ответ:**
```
HTTP/1.1 302 Found
Location: https://very-long-url.com/some/very/long/path
```

**Результат:** Браузер автоматически переходит на оригинальный URL!

---

## 🧩 Dependency Injection: зачем и как

### Проблема без DI

Представьте, если бы наш `RegisterUseCase` создавал зависимости сам:

```go
// ❌ ПЛОХОЙ ПРИМЕР (без DI)
func (uc *RegisterUseCase) Execute(req RegisterRequest) {
    // Создаем зависимости прямо в Use Case
    db, _ := sql.Open("postgres", "host=localhost...")
    userRepo := &PostgreSQLUserRepository{db: db}
    
    hasher := &BcryptPasswordHasher{cost: 12}
    tokenGen := &JWTTokenGenerator{secret: "secret"}
    
    // Используем зависимости...
}
```

**Проблемы:**
1. **Невозможно тестировать** - нельзя подставить моки
2. **Жесткая связанность** - Use Case знает о PostgreSQL, bcrypt, JWT
3. **Нарушение Single Responsibility** - Use Case отвечает за создание зависимостей
4. **Дублирование кода** - каждый Use Case создает одни и те же объекты

### Решение с DI

```go
// ✅ ХОРОШИЙ ПРИМЕР (с DI)
type RegisterUseCase struct {
    userRepo       user.Repository    // ← ИНТЕРФЕЙС!
    passwordHasher PasswordHasher     // ← ИНТЕРФЕЙС!
    tokenGenerator TokenGenerator     // ← ИНТЕРФЕЙС!
}

func NewRegisterUseCase(
    userRepo user.Repository,
    passwordHasher PasswordHasher,
    tokenGenerator TokenGenerator,
) *RegisterUseCase {
    return &RegisterUseCase{
        userRepo:       userRepo,
        passwordHasher: passwordHasher,
        tokenGenerator: tokenGenerator,
    }
}
```

**Преимущества:**
1. **Легко тестировать:**
```go
func TestRegisterUseCase(t *testing.T) {
    mockRepo := &MockUserRepository{}
    mockHasher := &MockPasswordHasher{}
    mockTokenGen := &MockTokenGenerator{}
    
    useCase := NewRegisterUseCase(mockRepo, mockHasher, mockTokenGen)
    // Тестируем без реальной БД!
}
```

2. **Легко менять реализации:**
```go
// Можем заменить PostgreSQL на MongoDB
userRepo := &MongoUserRepository{client: mongoClient}

// Можем заменить bcrypt на argon2
passwordHasher := &Argon2PasswordHasher{}

// Можем добавить кеширование
userRepo = &CachedUserRepository{
    inner: userRepo,
    cache: redisClient,
}
```

### Как работает наш DI контейнер

В `pkg/di/container.go`:

```go
func NewContainer(config *configs.Config) (*Container, error) {
    container := &Container{Config: config}
    
    // 1. Создаем инфраструктуру (БД, сервер)
    container.initInfrastructure()
    
    // 2. Создаем репозитории (используют БД)
    container.initRepositories()
    
    // 3. Создаем внешние сервисы (bcrypt, JWT)
    container.initExternalServices()
    
    // 4. Создаем Use Cases (используют репозитории и сервисы)
    container.initUseCases()
    
    // 5. Создаем контроллеры (используют Use Cases)
    container.initControllers()
    
    return container, nil
}
```

**Принцип:** Каждый компонент получает свои зависимости готовыми через конструктор.

**Граф зависимостей:**
```
AuthController
     ↓ (зависит от)
RegisterUseCase
     ↓ (зависит от)
UserRepository, PasswordHasher, TokenGenerator
     ↓ (зависит от)
Database, Config
```

**Порядок создания:** снизу вверх!
1. Database
2. UserRepository (нужна Database)
3. PasswordHasher, TokenGenerator
4. RegisterUseCase (нужны Repository, PasswordHasher, TokenGenerator)
5. AuthController (нужен RegisterUseCase)

---

## 🏗️ Слои архитектуры: что где находится

### 🏛️ Domain Layer (Самый внутренний)

**Расположение:** `internal/domain/`

**Что содержит:**
- **Entities** (доменные модели)
- **Repository interfaces** (контракты для данных)
- **Domain errors** (бизнес-ошибки)

**Пример - User Entity:**
```go
// internal/domain/user/entity.go
type User struct {
    ID             uint
    Email          string
    HashedPassword string
    Name           string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func NewUser(email, name string) (*User, error) {
    if !isValidEmail(email) {
        return nil, ErrInvalidEmail  // Доменная ошибка
    }
    // Бизнес-правила валидации
}
```

**Правила Domain слоя:**
- ✅ НЕТ импортов внешних библиотек (кроме стандартной библиотеки Go)
- ✅ НЕ ЗНАЕТ о базах данных, HTTP, JSON
- ✅ Содержит только чистые бизнес-правила
- ✅ Может использовать другие entities из того же слоя

### 💼 Use Cases Layer (Application Business Rules)

**Расположение:** `internal/usecase/`

**Что содержит:**
- **Use Cases** (сценарии использования)
- **Request/Response structures** (входные и выходные данные)
- **Interfaces для зависимостей** (контракты для внешних сервисов)

**Пример - RegisterUseCase:**
```go
// internal/usecase/auth/register.go
func (uc *RegisterUseCase) Execute(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
    // ОРКЕСТРАЦИЯ: координирует работу между компонентами
    
    // 1. Проверяем через репозиторий
    exists, _ := uc.userRepo.Exists(ctx, req.Email)
    
    // 2. Создаем доменную модель
    newUser, _ := user.NewUser(req.Email, req.Name)
    
    // 3. Используем внешний сервис
    hashedPassword, _ := uc.passwordHasher.Hash(req.Password)
    
    // 4. Сохраняем через репозиторий
    _ = uc.userRepo.Save(ctx, newUser)
    
    return response, nil
}
```

**Правила Use Cases слоя:**
- ✅ МОЖЕТ импортировать domain слой
- ✅ НЕ ЗНАЕТ о HTTP, SQL, конкретных библиотеках
- ✅ Зависит только от интерфейсов
- ✅ Содержит application-специфичную логику

### 🔌 Interface Adapters Layer

**Расположение:** `internal/interfaces/`

**Что содержит:**
- **Controllers** (HTTP handlers)
- **Presenters** (форматирование ответов)
- **Gateways** (адаптеры для внешних API)

**Пример - AuthController:**
```go
// internal/interfaces/controllers/auth_controller.go
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
    // АДАПТАЦИЯ: преобразует HTTP в Use Case и обратно
    
    // 1. HTTP → Go struct
    var req auth.RegisterRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 2. Вызываем Use Case
    response, err := c.registerUC.Execute(r.Context(), req)
    
    // 3. Go struct → HTTP
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

**Правила Interface Adapters слоя:**
- ✅ МОЖЕТ импортировать usecase и domain слои
- ✅ ЗНАЕТ о HTTP протоколе
- ✅ НЕ содержит бизнес-логики
- ✅ Только адаптация форматов данных

### 🌐 Infrastructure Layer (Самый внешний)

**Расположение:** `internal/infrastructure/`

**Что содержит:**
- **Database implementations** (PostgreSQL репозитории)
- **Web server** (HTTP сервер и middleware)
- **External services** (JWT, bcrypt, API клиенты)

**Пример - UserRepository:**
```go
// internal/infrastructure/database/user_repository.go
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
    // РЕАЛИЗАЦИЯ: конкретная работа с PostgreSQL
    
    query := `INSERT INTO users (email, hashed_password, name) VALUES ($1, $2, $3)`
    _, err := r.db.ExecContext(ctx, query, u.Email, u.HashedPassword, u.Name)
    return err
}
```

**Правила Infrastructure слоя:**
- ✅ МОЖЕТ импортировать ВСЕ слои
- ✅ ЗНАЕТ о конкретных технологиях (PostgreSQL, bcrypt, JWT)
- ✅ РЕАЛИЗУЕТ интерфейсы из внутренних слоев
- ✅ НЕ содержит бизнес-логики

### 🔄 Поток зависимостей

```
┌─────────────────────────────────────────────────────────────┐
│                    🌐 INFRASTRUCTURE                        │
│              (PostgreSQL, bcrypt, JWT, HTTP)               │
│                           │                                 │
│                           │ implements                      │
│                           ↓                                 │
├─────────────────────────────────────────────────────────────┤
│                   🔌 INTERFACE ADAPTERS                     │
│                    (HTTP Controllers)                      │
│                           │                                 │
│                           │ calls                           │
│                           ↓                                 │
├─────────────────────────────────────────────────────────────┤
│                     💼 USE CASES                            │
│                (RegisterUseCase, LoginUseCase)             │
│                           │                                 │
│                           │ uses                            │
│                           ↓                                 │
├─────────────────────────────────────────────────────────────┤
│                     🏛️ DOMAIN                              │
│                   (User, Link, interfaces)                 │
└─────────────────────────────────────────────────────────────┘
```

**Правило:** Стрелки показывают направление **импортов**. Внешние слои знают о внутренних, но не наоборот!

---

## 🔄 Инверсия зависимостей: почему это важно

### Без инверсии зависимостей (ПЛОХО)

```go
// ❌ Use Case зависит от конкретной реализации
type RegisterUseCase struct {
    userRepo *PostgreSQLUserRepository  // ← Конкретный тип!
}

func (uc *RegisterUseCase) Execute(req RegisterRequest) {
    // Прямая зависимость от PostgreSQL
    exists := uc.userRepo.ExistsInPostgreSQL(req.Email)
}
```

**Проблемы:**
1. Нельзя заменить PostgreSQL на MongoDB
2. Нельзя протестировать без реальной БД
3. Use Case знает о технических деталях PostgreSQL

### С инверсией зависимостей (ХОРОШО)

```go
// ✅ Use Case зависит от интерфейса
type RegisterUseCase struct {
    userRepo user.Repository  // ← ИНТЕРФЕЙС из domain слоя!
}

// Интерфейс определен в domain слое
type Repository interface {
    Exists(ctx context.Context, email string) (bool, error)
}
```

**Схема инверсии:**

**Без инверсии:**
```
Use Case → PostgreSQLRepository
    ↑              ↑
 (знает о)    (реализует)
```

**С инверсией:**
```
Use Case → Repository Interface ← PostgreSQLRepository
    ↑              ↑                        ↑
 (знает о)    (определяет)              (реализует)
```

### Конкретный пример из нашего кода

**1. Domain слой определяет интерфейс:**
```go
// internal/domain/user/repository.go
package user

type Repository interface {
    Save(ctx context.Context, user *User) error
    FindByEmail(ctx context.Context, email string) (*User, error)
    Exists(ctx context.Context, email string) (bool, error)
}
```

**2. Use Case использует интерфейс:**
```go
// internal/usecase/auth/register.go
package auth

import "clean-url-shortener/internal/domain/user"

type RegisterUseCase struct {
    userRepo user.Repository  // ← Интерфейс из domain!
}
```

**3. Infrastructure реализует интерфейс:**
```go
// internal/infrastructure/database/user_repository.go
package database

import "clean-url-shortener/internal/domain/user"

type UserRepository struct {
    db *DB
}

// Реализуем интерфейс user.Repository
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
    // PostgreSQL-специфичная реализация
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
    // PostgreSQL-специфичная реализация
}

func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
    // PostgreSQL-специфичная реализация
}
```

**4. DI контейнер связывает все вместе:**
```go
// pkg/di/container.go
func (c *Container) initRepositories() error {
    // Создаем конкретную реализацию
    postgresRepo := database.NewUserRepository(c.DB)
    
    // Сохраняем как интерфейс!
    c.UserRepo = postgresRepo  // user.Repository = *database.UserRepository
    
    return nil
}

func (c *Container) initUseCases() error {
    // Use Case получает интерфейс, не зная о конкретной реализации
    c.RegisterUC = auth.NewRegisterUseCase(
        c.UserRepo,  // Передаем интерфейс
        c.PasswordHasher,
        c.TokenGenerator,
    )
    return nil
}
```

### Преимущества инверсии

**1. Легкая замена реализации:**
```go
// Можем заменить PostgreSQL на MongoDB
c.UserRepo = mongodb.NewUserRepository(mongoClient)

// Или добавить кеширование
c.UserRepo = cache.NewCachedUserRepository(c.UserRepo, redisClient)

// Use Case не изменится!
```

**2. Простое тестирование:**
```go
// Тест с мок-объектом
func TestRegisterUseCase(t *testing.T) {
    mockRepo := &MockUserRepository{
        ExistsFunc: func(email string) bool {
            return email == "existing@example.com"
        },
    }
    
    useCase := auth.NewRegisterUseCase(mockRepo, mockHasher, mockTokenGen)
    
    // Тестируем без реальной БД!
    result, err := useCase.Execute(RegisterRequest{
        Email: "new@example.com",
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

**3. Модульность:**
```go
// Можем создать разные конфигурации
if config.Database.Type == "postgres" {
    c.UserRepo = database.NewPostgreSQLUserRepository(db)
} else if config.Database.Type == "mongodb" {
    c.UserRepo = database.NewMongoUserRepository(mongoClient)
}

// Use Cases остаются теми же!
```

---

## 🎯 Заключение

Наша Clean Architecture обеспечивает:

### ✅ Независимость от фреймворков
- Можем заменить `net/http` на Gin без изменения Use Cases
- Можем добавить gRPC API рядом с HTTP

### ✅ Тестируемость
- Каждый слой можно тестировать изолированно
- Легко создавать моки для всех зависимостей

### ✅ Независимость от базы данных
- Можем заменить PostgreSQL на MongoDB
- Можем добавить Redis для кеширования

### ✅ Независимость от внешних сервисов
- Можем заменить bcrypt на argon2
- Можем заменить локальный JWT на внешний сервис аутентификации

### ✅ Бизнес-правила защищены
- Доменная логика не зависит от технических деталей
- Изменения во внешних библиотеках не влияют на бизнес-логику

### 🔄 Направление зависимостей
```
Infrastructure → Interfaces → Use Cases → Domain
     ↑              ↑            ↑         ↑
(PostgreSQL,    (HTTP,       (Business  (Business
 bcrypt,        JSON,        Logic,     Rules,
 JWT)           REST)        Orchestr.) Entities)
```

**Главный принцип:** Внутренние слои НЕ ЗНАЮТ о внешних. Все зависимости направлены внутрь!

Это позволяет нам легко адаптироваться к изменениям, тестировать код и развивать приложение без нарушения существующей функциональности.