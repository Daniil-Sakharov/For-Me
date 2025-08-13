package link

// ИНТЕРФЕЙСЫ ДЛЯ ЗАВИСИМОСТЕЙ LINK USE CASES

// EventPublisher определяет интерфейс для публикации событий
// Используется для уведомления о переходах по ссылкам
type EventPublisher interface {
	// PublishLinkClicked публикует событие о клике по ссылке
	PublishLinkClicked(linkID uint, userAgent, ipAddress, referer string) error
}

// ShortCodeGenerator определяет интерфейс для генерации коротких кодов
// Может использовать различные алгоритмы (случайный, инкрементальный, hash-based)
type ShortCodeGenerator interface {
	// Generate генерирует новый короткий код
	Generate() (string, error)
	
	// GenerateCustom проверяет и адаптирует пользовательский код
	GenerateCustom(customCode string) (string, error)
}

// URLValidator определяет интерфейс для валидации URL
// Может включать проверки на безопасность, доступность и т.д.
type URLValidator interface {
	// Validate проверяет валидность и безопасность URL
	Validate(url string) error
	
	// IsSafe проверяет, что URL безопасен (не фишинг, не malware)
	IsSafe(url string) (bool, error)
}

// ПРИНЦИПЫ:
// 1. Каждый интерфейс имеет единственную ответственность
// 2. Интерфейсы определяют ЧТО нужно делать, а не КАК
// 3. Реализации будут в infrastructure слое
// 4. Это обеспечивает тестируемость и гибкость