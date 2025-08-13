package stat

import (
	"net"
	"time"
)

// Stat представляет доменную модель статистики переходов по ссылке
// Каждый клик по ссылке создает новую запись статистики
type Stat struct {
	// ID - уникальный идентификатор записи статистики
	ID uint

	// LinkID - ID ссылки, по которой произошел переход
	LinkID uint

	// UserAgent - информация о браузере/устройстве пользователя
	UserAgent string

	// IPAddress - IP адрес пользователя, совершившего переход
	IPAddress string

	// Referer - страница, с которой пришел пользователь
	Referer string

	// Country - страна пользователя (определяется по IP)
	Country string

	// City - город пользователя (определяется по IP)
	City string

	// CreatedAt - время перехода
	CreatedAt time.Time
}

// NewStat создает новую запись статистики
func NewStat(linkID uint, userAgent, ipAddress, referer string) *Stat {
	return &Stat{
		LinkID:    linkID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		Referer:   referer,
		CreatedAt: time.Now(),
	}
}

// SetLocation устанавливает географическую информацию
// Эта информация может быть получена через внешний сервис геолокации
func (s *Stat) SetLocation(country, city string) {
	s.Country = country
	s.City = city
}

// IsValidIP проверяет валидность IP адреса
func (s *Stat) IsValidIP() bool {
	return net.ParseIP(s.IPAddress) != nil
}

// GetAnonymizedIP возвращает анонимизированный IP адрес
// Убирает последний октет для IPv4 или последние 80 бит для IPv6
// Это важно для соблюдения GDPR и других требований приватности
func (s *Stat) GetAnonymizedIP() string {
	ip := net.ParseIP(s.IPAddress)
	if ip == nil {
		return ""
	}

	if ip.To4() != nil {
		// IPv4 - заменяем последний октет на 0
		ip[3] = 0
	} else {
		// IPv6 - заменяем последние 80 бит на 0
		for i := 6; i < 16; i++ {
			ip[i] = 0
		}
	}

	return ip.String()
}

// LinkStatistics представляет агрегированную статистику для ссылки
type LinkStatistics struct {
	// LinkID - ID ссылки
	LinkID uint

	// TotalClicks - общее количество переходов
	TotalClicks int

	// UniqueClicks - количество уникальных переходов (по IP)
	UniqueClicks int

	// ClicksByCountry - статистика по странам
	ClicksByCountry map[string]int

	// ClicksByHour - статистика по часам (для последних 24 часов)
	ClicksByHour []int

	// TopReferers - топ источников трафика
	TopReferers []RefererStat

	// Period - период, за который собрана статистика
	Period Period
}

// RefererStat представляет статистику по источнику трафика
type RefererStat struct {
	Referer string
	Count   int
}

// Period представляет период времени для статистики
type Period struct {
	StartDate time.Time
	EndDate   time.Time
}

// NewPeriod создает новый период
func NewPeriod(startDate, endDate time.Time) Period {
	return Period{
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// Last24Hours создает период для последних 24 часов
func Last24Hours() Period {
	now := time.Now()
	return Period{
		StartDate: now.Add(-24 * time.Hour),
		EndDate:   now,
	}
}

// Last7Days создает период для последних 7 дней
func Last7Days() Period {
	now := time.Now()
	return Period{
		StartDate: now.Add(-7 * 24 * time.Hour),
		EndDate:   now,
	}
}

// Last30Days создает период для последних 30 дней
func Last30Days() Period {
	now := time.Now()
	return Period{
		StartDate: now.Add(-30 * 24 * time.Hour),
		EndDate:   now,
	}
}