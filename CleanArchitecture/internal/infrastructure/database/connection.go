package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// Config содержит конфигурацию для подключения к базе данных
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DB обертка над sql.DB с дополнительными методами
type DB struct {
	*sql.DB
	config Config
}

// NewConnection создает новое подключение к PostgreSQL базе данных
func NewConnection(config Config) (*DB, error) {
	// Формируем строку подключения для PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	// Подключаемся к базе данных
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настраиваем пул соединений
	sqlDB.SetMaxOpenConns(25)                 // Максимум открытых соединений
	sqlDB.SetMaxIdleConns(25)                 // Максимум idle соединений
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Время жизни соединения

	// Проверяем подключение
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		DB:     sqlDB,
		config: config,
	}, nil
}

// Close закрывает подключение к базе данных
func (db *DB) Close() error {
	return db.DB.Close()
}

// GetStats возвращает статистику пула соединений
func (db *DB) GetStats() sql.DBStats {
	return db.DB.Stats()
}

// RunMigrations выполняет миграции базы данных
func (db *DB) RunMigrations() error {
	// Создаем таблицы если они не существуют
	migrations := []string{
		createUsersTable,
		createLinksTable,
		createStatsTable,
		createIndexes,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// SQL МИГРАЦИИ
// В реальном проекте эти миграции лучше вынести в отдельные файлы

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);`

const createLinksTable = `
CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    clicks_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);`

const createStatsTable = `
CREATE TABLE IF NOT EXISTS stats (
    id SERIAL PRIMARY KEY,
    link_id INTEGER NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    user_agent TEXT,
    ip_address INET,
    referer TEXT,
    country VARCHAR(100),
    city VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);`

const createIndexes = `
-- Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_links_short_code ON links(short_code);
CREATE INDEX IF NOT EXISTS idx_links_user_id ON links(user_id);
CREATE INDEX IF NOT EXISTS idx_stats_link_id ON stats(link_id);
CREATE INDEX IF NOT EXISTS idx_stats_created_at ON stats(created_at);
CREATE INDEX IF NOT EXISTS idx_stats_ip_address ON stats(ip_address);
`

// ПРИНЦИПЫ INFRASTRUCTURE СЛОЯ:
// 1. Содержит технические детали (SQL, драйверы БД)
// 2. Реализует интерфейсы из domain слоя
// 3. Не содержит бизнес-логики
// 4. Может быть заменен без изменения других слоев