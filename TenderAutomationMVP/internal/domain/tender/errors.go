// =====================================================================
// 🚨 ДОМЕННЫЕ ОШИБКИ ДЛЯ TENDER - Специфичные ошибки предметной области
// =====================================================================
//
// Этот файл содержит все ошибки, связанные с доменной логикой тендеров.
// Следует принципам Clean Architecture:
// 1. Ошибки определяются в доменном слое
// 2. Внешние слои могут проверять тип ошибки и обрабатывать соответственно
// 3. Ошибки содержат бизнес-контекст, а не технические детали
//
// TODO: При расширении функциональности добавить:
// - Ошибки для новых бизнес-правил
// - Коды ошибок для интеграции с внешними системами
// - Локализация сообщений об ошибках
// - Structured errors с дополнительными полями

package tender

import (
	"errors"
	"fmt"
)

// =====================================================================
// 🚨 ОСНОВНЫЕ ДОМЕННЫЕ ОШИБКИ
// =====================================================================

var (
	// ✋ Ошибки валидации входных данных
	ErrEmptyExternalID      = errors.New("external ID cannot be empty")
	ErrEmptyTitle          = errors.New("tender title cannot be empty")
	ErrEmptyURL            = errors.New("tender URL cannot be empty")
	ErrInvalidURL          = errors.New("invalid tender URL format")
	ErrInvalidPlatform     = errors.New("unsupported tender platform")

	// 🤖 Ошибки AI анализа
	ErrInvalidAIScore         = errors.New("AI score must be between 0.0 and 1.0")
	ErrInvalidAIRecommendation = errors.New("invalid AI recommendation")
	ErrAlreadyAnalyzed        = errors.New("tender has already been analyzed")
	ErrCannotAnalyze          = errors.New("tender cannot be analyzed in current state")

	// 🔄 Ошибки изменения состояния
	ErrInvalidStatusTransition = errors.New("invalid tender status transition")
	ErrTenderNotActive        = errors.New("tender is not active")
	ErrTenderExpired          = errors.New("tender submission deadline has expired")

	// 🔍 Ошибки поиска и доступа
	ErrTenderNotFound = errors.New("tender not found")
	ErrDuplicateTender = errors.New("tender with this external ID already exists")

	// 💰 Ошибки финансовых данных
	ErrNegativePrice    = errors.New("tender price cannot be negative")
	ErrInvalidCurrency  = errors.New("unsupported currency")

	// 📅 Ошибки дат и времени
	ErrInvalidDeadline    = errors.New("deadline cannot be before published date")
	ErrDeadlineInPast     = errors.New("deadline cannot be in the past")
)

// =====================================================================
// 🏗️ КОНСТРУКТОРЫ ОШИБОК С КОНТЕКСТОМ
// =====================================================================

// NewValidationError создает ошибку валидации с дополнительным контекстом
//
// TODO: Добавить structured error с полями для детализации
// TODO: Добавить поддержку множественных ошибок валидации
func NewValidationError(field, message string) error {
	return fmt.Errorf("validation error for field '%s': %s", field, message)
}

// NewBusinessRuleError создает ошибку нарушения бизнес-правила
//
// TODO: Добавить коды ошибок для категоризации
// TODO: Добавить контекст для логирования и отладки
func NewBusinessRuleError(rule, context string) error {
	return fmt.Errorf("business rule violation '%s': %s", rule, context)
}

// NewNotFoundError создает ошибку "не найдено" с контекстом
//
// TODO: Добавить типизированные параметры поиска
func NewNotFoundError(entityType, identifier string) error {
	return fmt.Errorf("%s not found: %s", entityType, identifier)
}

// =====================================================================
// 🔍 ПРОВЕРКИ ТИПОВ ОШИБОК
// =====================================================================

// IsValidationError проверяет, является ли ошибка ошибкой валидации
//
// TODO: Реализовать через typed errors или error wrapping
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}

	// Проверяем известные ошибки валидации
	validationErrors := []error{
		ErrEmptyExternalID,
		ErrEmptyTitle,
		ErrEmptyURL,
		ErrInvalidURL,
		ErrInvalidPlatform,
		ErrInvalidAIScore,
		ErrInvalidAIRecommendation,
		ErrNegativePrice,
		ErrInvalidCurrency,
		ErrInvalidDeadline,
		ErrDeadlineInPast,
	}

	for _, validationErr := range validationErrors {
		if errors.Is(err, validationErr) {
			return true
		}
	}

	return false
}

// IsBusinessRuleError проверяет, является ли ошибка ошибкой бизнес-логики
//
// TODO: Добавить более точную категоризацию бизнес-ошибок
func IsBusinessRuleError(err error) bool {
	if err == nil {
		return false
	}

	businessErrors := []error{
		ErrAlreadyAnalyzed,
		ErrCannotAnalyze,
		ErrInvalidStatusTransition,
		ErrTenderNotActive,
		ErrTenderExpired,
	}

	for _, businessErr := range businessErrors {
		if errors.Is(err, businessErr) {
			return true
		}
	}

	return false
}

// IsNotFoundError проверяет, является ли ошибка ошибкой "не найдено"
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrTenderNotFound)
}

// IsConflictError проверяет, является ли ошибка ошибкой конфликта
func IsConflictError(err error) bool {
	return errors.Is(err, ErrDuplicateTender)
}

// =====================================================================
// 📊 ОШИБКИ С ДОПОЛНИТЕЛЬНЫМИ ДАННЫМИ (для расширения)
// =====================================================================

// ValidationError представляет ошибку валидации с дополнительными данными
//
// TODO: Реализовать полностью, если потребуется детальная валидация
type ValidationError struct {
	Field   string                 // Поле, которое не прошло валидацию
	Value   interface{}            // Значение, которое не прошло валидацию
	Rule    string                 // Правило валидации, которое было нарушено
	Message string                 // Сообщение об ошибке
	Context map[string]interface{} // Дополнительный контекст
}

func (e ValidationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("validation failed for field '%s' with rule '%s'", e.Field, e.Rule)
}

// BusinessRuleError представляет ошибку нарушения бизнес-правила
//
// TODO: Добавить категории бизнес-правил
// TODO: Добавить рекомендации по исправлению
type BusinessRuleError struct {
	Rule        string                 // Название бизнес-правила
	Entity      string                 // Сущность, для которой нарушено правило
	Context     map[string]interface{} // Контекст нарушения
	Message     string                 // Сообщение об ошибке
	Suggestion  string                 // Рекомендация по исправлению
}

func (e BusinessRuleError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("business rule '%s' violated for %s", e.Rule, e.Entity)
}

// =====================================================================
// 🔧 ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ДЛЯ РАБОТЫ С ОШИБКАМИ
// =====================================================================

// WrapError оборачивает ошибку с дополнительным контекстом
//
// TODO: Использовать стандартный fmt.Errorf с %w для wrapping
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// CombineErrors объединяет несколько ошибок в одну
//
// TODO: Реализовать proper error aggregation
func CombineErrors(errors ...error) error {
	var validErrors []error
	for _, err := range errors {
		if err != nil {
			validErrors = append(validErrors, err)
		}
	}

	if len(validErrors) == 0 {
		return nil
	}

	if len(validErrors) == 1 {
		return validErrors[0]
	}

	// Простое объединение - TODO: улучшить
	var messages []string
	for _, err := range validErrors {
		messages = append(messages, err.Error())
	}

	return fmt.Errorf("multiple errors: %v", messages)
}

// =====================================================================
// 📝 ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ ОШИБОК
// =====================================================================

// TODO: Примеры использования в коде:
//
// 1. Проверка типа ошибки:
//    if IsValidationError(err) {
//        // Обработка ошибки валидации
//        return 400, err
//    }
//
// 2. Создание ошибки с контекстом:
//    return NewValidationError("title", "must not be empty")
//
// 3. Оборачивание ошибки:
//    return WrapError(err, "failed to create tender")
//
// 4. Проверка конкретной ошибки:
//    if errors.Is(err, ErrTenderNotFound) {
//        return 404, "Tender not found"
//    }
//
// 5. Создание бизнес-ошибки:
//    return NewBusinessRuleError("cannot_analyze_inactive", "tender must be active")

// =====================================================================
// 🔄 КОНСТАНТЫ ДЛЯ КАТЕГОРИЗАЦИИ ОШИБОК (для будущего расширения)
// =====================================================================

// ErrorCategory представляет категорию ошибки
type ErrorCategory string

const (
	CategoryValidation   ErrorCategory = "validation"    // Ошибки валидации
	CategoryBusinessRule ErrorCategory = "business_rule" // Нарушения бизнес-правил  
	CategoryNotFound     ErrorCategory = "not_found"     // Сущность не найдена
	CategoryConflict     ErrorCategory = "conflict"      // Конфликт данных
	CategoryPermission   ErrorCategory = "permission"    // Нет прав доступа
	CategoryInternal     ErrorCategory = "internal"      // Внутренние ошибки
)

// GetErrorCategory возвращает категорию ошибки
//
// TODO: Реализовать полную категоризацию
func GetErrorCategory(err error) ErrorCategory {
	if IsValidationError(err) {
		return CategoryValidation
	}
	if IsBusinessRuleError(err) {
		return CategoryBusinessRule
	}
	if IsNotFoundError(err) {
		return CategoryNotFound
	}
	if IsConflictError(err) {
		return CategoryConflict
	}
	return CategoryInternal
}

// =====================================================================
// ✅ ЗАКЛЮЧЕНИЕ
// =====================================================================

// Этот файл предоставляет полный набор ошибок для доменной логики тендеров.
// При расширении функциональности:
// 1. Добавляйте новые ошибки в соответствующие секции
// 2. Создавайте конструкторы для ошибок с контекстом
// 3. Добавляйте проверки типов для новых категорий ошибок
// 4. Документируйте примеры использования
//
// Важно: Доменные ошибки должны быть независимы от технических деталей
// и содержать только бизнес-контекст.