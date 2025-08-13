// 🔗 УНИВЕРСАЛЬНЫЙ ШАБЛОН USE CASE INTERFACES для Clean Architecture
//
// ✅ ВСЕГДА одинаково:
// - Интерфейсы для внешних зависимостей в usecase слое
// - НЕТ импортов конкретных реализаций
// - Каждый интерфейс решает одну задачу (SRP)
//
// 🔀 МЕНЯЕТСЯ в зависимости от проекта:
// - Названия интерфейсов (EmailSender, PaymentProcessor, etc.)
// - Методы интерфейсов
// - Структуры данных

package example_feature // 🔀 ИЗМЕНИТЕ: название фичи (auth, payment, notification, etc.)

// ✅ ВСЕГДА: только стандартная библиотека, никаких внешних зависимостей!
// ❌ НЕ ИМПОРТИРУЙТЕ: smtp, stripe, aws-sdk, redis и другие конкретные библиотеки

// 🔧 ИНТЕРФЕЙСЫ ДЛЯ ВНЕШНИХ СЕРВИСОВ
// ✅ ВСЕГДА: каждый интерфейс представляет одну возможность (capability)

// 🔀 ИЗМЕНИТЕ: интерфейсы под вашу фичу
// Примеры для ExampleFeature (замените на реальные):

type ExampleService interface {
	ProcessSomething(data string) (string, error) // 🔀 ИЗМЕНИТЕ: метод под вашу логику
	ValidateSomething(input string) (bool, error) // 🔀 ИЗМЕНИТЕ: метод под вашу логику
}

type ExampleNotifier interface {
	SendNotification(message string) error // 🔀 ИЗМЕНИТЕ: метод под вашу логику
}

// 📝 ПРИМЕРЫ ИНТЕРФЕЙСОВ для разных фич:

// ДЛЯ AUTH (аутентификация):
// type PasswordHasher interface {
//     Hash(password string) (string, error)
//     Compare(hashedPassword, password string) error
// }
//
// type TokenGenerator interface {
//     Generate(userID uint, email string) (string, error)
//     Validate(token string) (*TokenData, error)
//     Refresh(refreshToken string) (string, error)
// }
//
// type EmailVerifier interface {
//     SendVerificationEmail(email, token string) error
//     VerifyEmail(token string) (string, error) // возвращает email
// }

// ДЛЯ PAYMENT (платежи):
// type PaymentProcessor interface {
//     ProcessPayment(amount float64, cardToken string) (*PaymentResult, error)
//     RefundPayment(transactionID string, amount float64) (*RefundResult, error)
//     GetPaymentStatus(transactionID string) (*PaymentStatus, error)
// }
//
// type PaymentValidator interface {
//     ValidateCard(cardNumber, cvv, expiry string) error
//     ValidateAmount(amount float64, currency string) error
// }

// ДЛЯ NOTIFICATION (уведомления):
// type EmailSender interface {
//     SendEmail(to, subject, body string) error
//     SendEmailWithTemplate(to string, templateID string, data map[string]interface{}) error
// }
//
// type SMSSender interface {
//     SendSMS(phoneNumber, message string) error
//     SendSMSWithTemplate(phoneNumber, templateID string, data map[string]interface{}) error
// }
//
// type PushNotificationSender interface {
//     SendPushNotification(deviceToken, title, body string, data map[string]interface{}) error
// }

// ДЛЯ FILE_UPLOAD (загрузка файлов):
// type FileStorage interface {
//     Upload(filename string, data []byte) (string, error) // возвращает URL
//     Download(filename string) ([]byte, error)
//     Delete(filename string) error
//     GetURL(filename string) (string, error)
// }
//
// type ImageProcessor interface {
//     Resize(image []byte, width, height int) ([]byte, error)
//     Compress(image []byte, quality int) ([]byte, error)
//     GenerateThumbnail(image []byte) ([]byte, error)
// }

// ДЛЯ ANALYTICS (аналитика):
// type EventTracker interface {
//     TrackEvent(userID uint, eventName string, properties map[string]interface{}) error
//     TrackPageView(userID uint, page string) error
//     TrackConversion(userID uint, conversionType string, value float64) error
// }
//
// type MetricsCollector interface {
//     IncrementCounter(name string, tags map[string]string) error
//     RecordGauge(name string, value float64, tags map[string]string) error
//     RecordHistogram(name string, value float64, tags map[string]string) error
// }

// ДЛЯ CACHING (кеширование):
// type CacheManager interface {
//     Set(key string, value interface{}, expiration time.Duration) error
//     Get(key string, dest interface{}) error
//     Delete(key string) error
//     Exists(key string) (bool, error)
//     Clear() error
// }

// ДЛЯ SEARCH (поиск):
// type SearchEngine interface {
//     Index(documentID string, document interface{}) error
//     Search(query string, filters map[string]interface{}) (*SearchResults, error)
//     Delete(documentID string) error
//     BulkIndex(documents map[string]interface{}) error
// }

// ДЛЯ MESSAGE_QUEUE (очереди сообщений):
// type MessagePublisher interface {
//     Publish(topic string, message interface{}) error
//     PublishWithDelay(topic string, message interface{}, delay time.Duration) error
// }
//
// type MessageSubscriber interface {
//     Subscribe(topic string, handler func(message interface{}) error) error
//     Unsubscribe(topic string) error
// }

// 📊 СТРУКТУРЫ ДАННЫХ ДЛЯ ИНТЕРФЕЙСОВ
// ✅ ВСЕГДА: простые структуры для передачи данных

// 🔀 ИЗМЕНИТЕ: структуры под ваши интерфейсы
type ExampleData struct {
	Field1 string `json:"field1"` // 🔀 ИЗМЕНИТЕ: поля под ваши данные
	Field2 int    `json:"field2"`
}

// 📝 ПРИМЕРЫ структур для разных фич:

// ДЛЯ AUTH:
// type TokenData struct {
//     UserID uint   `json:"user_id"`
//     Email  string `json:"email"`
//     Role   string `json:"role"`
// }

// ДЛЯ PAYMENT:
// type PaymentResult struct {
//     TransactionID string    `json:"transaction_id"`
//     Status        string    `json:"status"`
//     Amount        float64   `json:"amount"`
//     Currency      string    `json:"currency"`
//     ProcessedAt   time.Time `json:"processed_at"`
// }
//
// type PaymentStatus struct {
//     TransactionID string    `json:"transaction_id"`
//     Status        string    `json:"status"` // pending, completed, failed, refunded
//     Amount        float64   `json:"amount"`
//     CreatedAt     time.Time `json:"created_at"`
//     UpdatedAt     time.Time `json:"updated_at"`
// }

// ДЛЯ FILE_UPLOAD:
// type UploadResult struct {
//     URL      string `json:"url"`
//     Filename string `json:"filename"`
//     Size     int64  `json:"size"`
//     MimeType string `json:"mime_type"`
// }

// ДЛЯ SEARCH:
// type SearchResults struct {
//     Results []interface{} `json:"results"`
//     Total   int64         `json:"total"`
//     Page    int           `json:"page"`
//     Size    int           `json:"size"`
// }

// 📚 ИНСТРУКЦИЯ ПО АДАПТАЦИИ:
//
// 1. 🔀 ПЕРЕИМЕНУЙТЕ файл и пакет:
//    example_feature → auth (или payment, notification, etc.)
//
// 2. 🔀 ЗАМЕНИТЕ примеры интерфейсов:
//    Выберите подходящие для вашей фичи из примеров выше
//
// 3. 🔀 СОЗДАЙТЕ специфичные интерфейсы:
//    Добавьте интерфейсы уникальные для вашей фичи
//
// 4. 🔀 ОПРЕДЕЛИТЕ структуры данных:
//    Создайте типы для передачи данных между интерфейсами
//
// 5. ✅ СЛЕДУЙТЕ принципам:
//    - Один интерфейс = одна ответственность
//    - Методы интерфейса решают конкретную задачу
//    - НЕТ зависимостей от конкретных реализаций
//
// 6. 📝 ДОКУМЕНТИРУЙТЕ интерфейсы:
//    Добавьте комментарии к сложным методам
//
// 📋 ПРИНЦИПЫ INTERFACE DESIGN:
//
// ✅ ДЕЛАЙТЕ:
// - Маленькие интерфейсы (1-3 метода)
// - Понятные названия методов
// - Возвращайте ошибки там, где они возможны
// - Используйте context.Context для длительных операций
// - Группируйте связанные методы в один интерфейс
//
// ❌ НЕ ДЕЛАЙТЕ:
// - Большие интерфейсы с множеством методов
// - Смешивание разных ответственностей в одном интерфейсе
// - Импорт конкретных реализаций
// - Слишком специфичные интерфейсы под одну реализацию
//
// 🔄 СВЯЗЬ С ДРУГИМИ СЛОЯМИ:
//
// Use Cases → Используют интерфейсы (этот файл)
//    ↑
// Infrastructure → Реализует интерфейсы
//    ↑
// DI Container → Связывает интерфейсы с реализациями
//
// 📝 ПРИМЕРЫ РЕАЛИЗАЦИИ (в infrastructure/external):
//
// type SMTPEmailSender struct {
//     host     string
//     port     int
//     username string
//     password string
// }
//
// func NewSMTPEmailSender(config SMTPConfig) EmailSender {
//     return &SMTPEmailSender{...}
// }
//
// func (s *SMTPEmailSender) SendEmail(to, subject, body string) error {
//     // SMTP реализация
// }
//
// 📝 ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ (в use cases):
//
// type RegisterUseCase struct {
//     userRepo      user.Repository
//     emailSender   EmailSender      // ← интерфейс из этого файла
//     passwordHasher PasswordHasher  // ← интерфейс из этого файла
// }
//
// func (uc *RegisterUseCase) Execute(req RegisterRequest) error {
//     // Используем интерфейсы без знания о реализации
//     hashedPassword, _ := uc.passwordHasher.Hash(req.Password)
//     _ = uc.emailSender.SendEmail(req.Email, "Welcome", "...")
// }