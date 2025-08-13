package stat

// ИНТЕРФЕЙСЫ ДЛЯ ЗАВИСИМОСТЕЙ STAT USE CASES

// GeoLocationService определяет интерфейс для определения геолокации по IP
type GeoLocationService interface {
	// GetLocation определяет страну и город по IP адресу
	GetLocation(ipAddress string) (country, city string, err error)
	
	// IsValidIP проверяет валидность IP адреса
	IsValidIP(ipAddress string) bool
}

// EventSubscriber определяет интерфейс для подписки на события
type EventSubscriber interface {
	// SubscribeLinkClicked подписывается на события кликов по ссылкам
	SubscribeLinkClicked() (<-chan LinkClickedEvent, error)
	
	// Close закрывает подписку
	Close() error
}

// LinkClickedEvent представляет событие клика по ссылке
type LinkClickedEvent struct {
	LinkID    uint   `json:"link_id"`
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
	Referer   string `json:"referer"`
	Timestamp string `json:"timestamp"`
}

// ПРИНЦИПЫ:
// 1. Слой Use Cases определяет ЧТО ему нужно от внешних сервисов
// 2. Infrastructure слой предоставляет реализации этих интерфейсов
// 3. Это обеспечивает тестируемость и гибкость архитектуры