-- =====================================================================
-- 🗃️ СОЗДАНИЕ ТАБЛИЦЫ TENDERS - Основная таблица для хранения тендеров
-- =====================================================================
-- 
-- Эта миграция создает основную таблицу для хранения информации о тендерах
-- Следует принципам:
-- 1. Эффективная индексация для быстрого поиска
-- 2. Правильные типы данных для оптимизации хранения
-- 3. Ограничения целостности для защиты от некорректных данных
-- 4. Значения по умолчанию для упрощения работы
--
-- TODO: При необходимости расширения добавить:
-- - Таблицы для документов тендеров
-- - Таблицы для товаров в тендерах  
-- - Таблицы для поставщиков
-- - Таблицы для email кампаний
-- - Таблицы для аналитики

-- =====================================================================
-- 📊 ОСНОВНАЯ ТАБЛИЦА TENDERS
-- =====================================================================

CREATE TABLE tenders (
    -- 🔑 Первичный ключ
    id SERIAL PRIMARY KEY,
    
    -- 🆔 Внешний идентификатор тендера (с площадки)
    external_id VARCHAR(100) NOT NULL UNIQUE,
    
    -- 📝 Основная информация о тендере
    title TEXT NOT NULL,                    -- Название тендера
    description TEXT,                       -- Описание тендера
    
    -- 🌐 Информация о платформе
    platform VARCHAR(50) NOT NULL,         -- Платформа (zakupki, szvo, spb и т.д.)
    url TEXT,                               -- Ссылка на тендер
    
    -- 👤 Информация о заказчике
    customer TEXT,                          -- Наименование заказчика
    customer_inn VARCHAR(20),               -- ИНН заказчика
    
    -- 💰 Финансовая информация
    start_price DECIMAL(15,2),              -- Начальная (максимальная) цена
    currency VARCHAR(3) DEFAULT 'RUB',      -- Валюта
    
    -- 📅 Временные рамки
    published_at TIMESTAMP,                 -- Дата публикации тендера
    deadline_at TIMESTAMP,                  -- Дедлайн подачи заявок
    
    -- 🏷️ Статус и категоризация
    status VARCHAR(20) NOT NULL DEFAULT 'active',  -- Статус тендера
    category VARCHAR(100),                  -- Категория товаров/услуг
    
    -- 🤖 AI анализ результаты
    ai_score DECIMAL(3,2),                  -- Оценка релевантности от AI (0.00-1.00)
    ai_recommendation VARCHAR(20),          -- Рекомендация AI (participate/skip/analyze)
    ai_analysis_reason TEXT,                -- Обоснование анализа AI
    ai_analyzed_at TIMESTAMP,               -- Дата анализа AI
    
    -- 📊 Служебные поля
    created_at TIMESTAMP DEFAULT NOW(),     -- Дата создания записи
    updated_at TIMESTAMP DEFAULT NOW(),     -- Дата последнего обновления
    
    -- ✅ Ограничения целостности данных
    CONSTRAINT valid_status 
        CHECK (status IN ('active', 'completed', 'cancelled', 'expired', 'draft')),
    
    CONSTRAINT valid_ai_score 
        CHECK (ai_score IS NULL OR (ai_score >= 0 AND ai_score <= 1)),
    
    CONSTRAINT valid_ai_recommendation 
        CHECK (ai_recommendation IS NULL OR 
               ai_recommendation IN ('participate', 'skip', 'analyze')),
    
    CONSTRAINT valid_currency 
        CHECK (currency IN ('RUB', 'USD', 'EUR')),
    
    CONSTRAINT valid_dates 
        CHECK (published_at IS NULL OR deadline_at IS NULL OR published_at <= deadline_at),
    
    CONSTRAINT positive_price 
        CHECK (start_price IS NULL OR start_price >= 0)
);

-- =====================================================================
-- 🚀 ИНДЕКСЫ ДЛЯ ПРОИЗВОДИТЕЛЬНОСТИ
-- =====================================================================

-- 🔍 Уникальный индекс по external_id (основной поиск)
CREATE UNIQUE INDEX idx_tenders_external_id ON tenders(external_id);

-- 🌐 Индекс по платформе (фильтрация по источнику)
CREATE INDEX idx_tenders_platform ON tenders(platform);

-- 🏷️ Индекс по статусу (основные запросы по активным тендерам)
CREATE INDEX idx_tenders_status ON tenders(status);

-- 📅 Индекс по дедлайну (поиск истекающих тендеров)
CREATE INDEX idx_tenders_deadline ON tenders(deadline_at) 
WHERE deadline_at IS NOT NULL;

-- 🤖 Индекс по AI оценке (поиск релевантных тендеров)
CREATE INDEX idx_tenders_ai_score ON tenders(ai_score) 
WHERE ai_score IS NOT NULL;

-- 📊 Составной индекс для аналитики (статус + дата создания)
CREATE INDEX idx_tenders_status_created ON tenders(status, created_at);

-- 🔍 Составной индекс для поиска по платформе и статусу
CREATE INDEX idx_tenders_platform_status ON tenders(platform, status);

-- 📅 Индекс по дате создания (сортировка по времени)
CREATE INDEX idx_tenders_created_at ON tenders(created_at);

-- 📝 Полнотекстовый поиск по названию и описанию (для PostgreSQL)
-- TODO: Раскомментировать если нужен полнотекстовый поиск
-- CREATE INDEX idx_tenders_fulltext ON tenders 
-- USING gin(to_tsvector('russian', title || ' ' || COALESCE(description, '')));

-- =====================================================================
-- 🕒 ТРИГГЕР ДЛЯ АВТОМАТИЧЕСКОГО ОБНОВЛЕНИЯ updated_at
-- =====================================================================

-- Функция для обновления поля updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггер на таблицу tenders
CREATE TRIGGER trigger_tenders_updated_at
    BEFORE UPDATE ON tenders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================================
-- 📊 КОММЕНТАРИИ К ТАБЛИЦЕ И ПОЛЯМ (для документации)
-- =====================================================================

COMMENT ON TABLE tenders IS 'Основная таблица для хранения информации о тендерах';

COMMENT ON COLUMN tenders.id IS 'Внутренний ID тендера';
COMMENT ON COLUMN tenders.external_id IS 'ID тендера на внешней платформе (уникальный)';
COMMENT ON COLUMN tenders.title IS 'Название тендера';
COMMENT ON COLUMN tenders.description IS 'Подробное описание тендера';
COMMENT ON COLUMN tenders.platform IS 'Платформа-источник (zakupki.gov.ru и т.д.)';
COMMENT ON COLUMN tenders.url IS 'Прямая ссылка на тендер';
COMMENT ON COLUMN tenders.customer IS 'Наименование организации-заказчика';
COMMENT ON COLUMN tenders.customer_inn IS 'ИНН заказчика';
COMMENT ON COLUMN tenders.start_price IS 'Начальная максимальная цена';
COMMENT ON COLUMN tenders.currency IS 'Валюта тендера';
COMMENT ON COLUMN tenders.published_at IS 'Дата публикации тендера';
COMMENT ON COLUMN tenders.deadline_at IS 'Крайний срок подачи заявок';
COMMENT ON COLUMN tenders.status IS 'Текущий статус тендера';
COMMENT ON COLUMN tenders.category IS 'Категория товаров/услуг';
COMMENT ON COLUMN tenders.ai_score IS 'Оценка релевантности от AI (0-1)';
COMMENT ON COLUMN tenders.ai_recommendation IS 'Рекомендация AI по участию';
COMMENT ON COLUMN tenders.ai_analysis_reason IS 'Объяснение решения AI';
COMMENT ON COLUMN tenders.ai_analyzed_at IS 'Время анализа AI';
COMMENT ON COLUMN tenders.created_at IS 'Время создания записи';
COMMENT ON COLUMN tenders.updated_at IS 'Время последнего обновления';

-- =====================================================================
-- 📋 ПРИМЕРЫ ЗАПРОСОВ ДЛЯ ТЕСТИРОВАНИЯ
-- =====================================================================

-- TODO: Использовать эти запросы для тестирования после создания таблицы:

-- 1. Поиск активных тендеров:
-- SELECT * FROM tenders WHERE status = 'active' ORDER BY created_at DESC LIMIT 10;

-- 2. Поиск тендеров с высокой AI оценкой:
-- SELECT * FROM tenders WHERE ai_score > 0.8 ORDER BY ai_score DESC;

-- 3. Статистика по платформам:
-- SELECT platform, COUNT(*), AVG(ai_score) FROM tenders GROUP BY platform;

-- 4. Поиск истекающих тендеров:
-- SELECT * FROM tenders 
-- WHERE deadline_at > NOW() 
--   AND deadline_at <= NOW() + INTERVAL '7 days' 
--   AND status = 'active';

-- =====================================================================
-- ✅ МИГРАЦИЯ ЗАВЕРШЕНА
-- =====================================================================

-- Эта миграция создает полнофункциональную таблицу для MVP
-- Таблица оптимизирована для:
-- 1. Быстрого поиска тендеров по различным критериям
-- 2. Эффективного хранения AI анализа
-- 3. Аналитических запросов
-- 4. Будущего расширения функциональности
--
-- Следующие шаги:
-- 1. Создать соответствующую down миграцию
-- 2. Реализовать доменную модель Tender в Go
-- 3. Создать repository для работы с таблицей
-- 4. Добавить тесты для проверки схемы БД