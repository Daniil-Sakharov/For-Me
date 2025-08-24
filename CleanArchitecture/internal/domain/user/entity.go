package user

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// User представляет доменную модель пользователя
// ЭТО ЧИСТАЯ ДОМЕННАЯ СУЩНОСТЬ - НЕТ ЗАВИСИМОСТЕЙ ОТ ВНЕШНИХ БИБЛИОТЕК!
// Содержит только бизнес-логику и правила домена
type User struct {
	// ID - уникальный идентификатор пользователя
	ID uint

	// Email - электронная почта (должна быть уникальной)
	Email string

	// HashedPassword - захешированный пароль (не храним пароль в открытом виде)
	HashedPassword string

	// Name - имя пользователя
	Name string

	// CreatedAt - время создания аккаунта
	CreatedAt time.Time

	// UpdatedAt - время последнего обновления
	UpdatedAt time.Time
}

// 📚 ИЗУЧИТЬ: Domain-Driven Design Validation
// https://martinfowler.com/articles/replaceThrowWithNotification.html
// https://enterprisecraftsmanship.com/posts/validation-in-domain-driven-design/
// https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/domain-model-layer-validations
//
// ПРИНЦИПЫ ДОМЕННОЙ ВАЛИДАЦИИ:
// 1. Валидация в домене защищает ИНВАРИАНТЫ
// 2. Валидация должна быть ЯВНОЙ и ПОНЯТНОЙ
// 3. Ошибки валидации - это ДОМЕННЫЕ ОШИБКИ
// 4. Валидация должна быть ДЕТЕРМИНИРОВАННОЙ
// 5. Валидация НЕ ЗАВИСИТ от внешних сервисов

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// ValidationErrors коллекция ошибок валидации
type ValidationErrors []ValidationError

func (errs ValidationErrors) Error() string {
	if len(errs) == 0 {
		return "no validation errors"
	}
	
	if len(errs) == 1 {
		return errs[0].Error()
	}
	
	var messages []string
	for _, err := range errs {
		messages = append(messages, err.Error())
	}
	
	return fmt.Sprintf("multiple validation errors: %s", strings.Join(messages, "; "))
}

// HasField проверяет наличие ошибки для конкретного поля
func (errs ValidationErrors) HasField(field string) bool {
	for _, err := range errs {
		if err.Field == field {
			return true
		}
	}
	return false
}

// GetField возвращает ошибку для конкретного поля
func (errs ValidationErrors) GetField(field string) *ValidationError {
	for _, err := range errs {
		if err.Field == field {
			return &err
		}
	}
	return nil
}

// Доменные ошибки - определяем их в доменном слое
var (
	// Базовые ошибки валидации
	ErrInvalidEmail    = ValidationError{Field: "email", Message: "invalid email format", Code: "INVALID_EMAIL"}
	ErrEmptyName       = ValidationError{Field: "name", Message: "name cannot be empty", Code: "EMPTY_NAME"}
	ErrEmptyPassword   = ValidationError{Field: "password", Message: "password cannot be empty", Code: "EMPTY_PASSWORD"}
	ErrPasswordTooWeak = ValidationError{Field: "password", Message: "password must be at least 8 characters long", Code: "WEAK_PASSWORD"}
	
	// Расширенные ошибки валидации
	ErrNameTooShort     = ValidationError{Field: "name", Message: "name must be at least 2 characters long", Code: "NAME_TOO_SHORT"}
	ErrNameTooLong      = ValidationError{Field: "name", Message: "name cannot exceed 100 characters", Code: "NAME_TOO_LONG"}
	ErrInvalidNameChars = ValidationError{Field: "name", Message: "name contains invalid characters", Code: "INVALID_NAME_CHARS"}
	
	// Ошибки email валидации
	ErrEmailTooLong     = ValidationError{Field: "email", Message: "email cannot exceed 255 characters", Code: "EMAIL_TOO_LONG"}
	ErrEmailInvalidTLD  = ValidationError{Field: "email", Message: "email has invalid top-level domain", Code: "INVALID_TLD"}
	ErrEmailDisposable  = ValidationError{Field: "email", Message: "disposable email addresses are not allowed", Code: "DISPOSABLE_EMAIL"}
	
	// Ошибки пароля
	ErrPasswordTooLong        = ValidationError{Field: "password", Message: "password cannot exceed 128 characters", Code: "PASSWORD_TOO_LONG"}
	ErrPasswordNoUppercase    = ValidationError{Field: "password", Message: "password must contain at least one uppercase letter", Code: "NO_UPPERCASE"}
	ErrPasswordNoLowercase    = ValidationError{Field: "password", Message: "password must contain at least one lowercase letter", Code: "NO_LOWERCASE"}
	ErrPasswordNoDigit        = ValidationError{Field: "password", Message: "password must contain at least one digit", Code: "NO_DIGIT"}
	ErrPasswordNoSpecialChar  = ValidationError{Field: "password", Message: "password must contain at least one special character", Code: "NO_SPECIAL_CHAR"}
	ErrPasswordCommonPassword = ValidationError{Field: "password", Message: "password is too common", Code: "COMMON_PASSWORD"}
	ErrPasswordContainsName   = ValidationError{Field: "password", Message: "password cannot contain user's name", Code: "PASSWORD_CONTAINS_NAME"}
	ErrPasswordContainsEmail  = ValidationError{Field: "password", Message: "password cannot contain user's email", Code: "PASSWORD_CONTAINS_EMAIL"}
)

// NewUser - конструктор для создания нового пользователя
// Содержит валидацию доменных правил
// 📚 ИЗУЧИТЬ: Notification Pattern для валидации
// https://martinfowler.com/articles/replaceThrowWithNotification.html
func NewUser(email, name string) (*User, error) {
	var validationErrors ValidationErrors
	
	// Валидация email согласно доменным правилам
	if emailErrors := validateEmail(email); len(emailErrors) > 0 {
		validationErrors = append(validationErrors, emailErrors...)
	}

	// Валидация имени согласно доменным правилам
	if nameErrors := validateName(name); len(nameErrors) > 0 {
		validationErrors = append(validationErrors, nameErrors...)
	}
	
	// Если есть ошибки валидации, возвращаем их все
	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	now := time.Now()

	return &User{
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// SetPassword устанавливает захешированный пароль
// Принимает уже захешированный пароль (хеширование - это техническая деталь,
// которая должна быть в infrastructure слое)
func (u *User) SetPassword(hashedPassword string) error {
	if hashedPassword == "" {
		return ErrEmptyPassword
	}

	u.HashedPassword = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateName обновляет имя пользователя
func (u *User) UpdateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}

	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateEmail обновляет email пользователя с валидацией
func (u *User) UpdateEmail(email string) error {
	if emailErrors := validateEmail(email); len(emailErrors) > 0 {
		return emailErrors
	}

	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateName обновляет имя пользователя с валидацией
func (u *User) UpdateName(name string) error {
	if nameErrors := validateName(name); len(nameErrors) > 0 {
		return nameErrors
	}

	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

// ValidatePassword валидирует пароль согласно доменным правилам
// Этот метод НЕ сохраняет пароль, только валидирует
func (u *User) ValidatePassword(password string) error {
	return validatePassword(password, u.Name, u.Email)
}

// Touch обновляет время последнего изменения
func (u *User) Touch() {
	u.UpdatedAt = time.Now()
}

// IsValid проверяет валидность всей сущности пользователя
func (u *User) IsValid() error {
	var validationErrors ValidationErrors
	
	// Проверяем email
	if emailErrors := validateEmail(u.Email); len(emailErrors) > 0 {
		validationErrors = append(validationErrors, emailErrors...)
	}
	
	// Проверяем имя
	if nameErrors := validateName(u.Name); len(nameErrors) > 0 {
		validationErrors = append(validationErrors, nameErrors...)
	}
	
	// Проверяем базовые поля
	if u.ID == 0 && !u.CreatedAt.IsZero() {
		// Если ID = 0, то это новый пользователь, CreatedAt должен быть установлен
	}
	
	if u.HashedPassword == "" {
		validationErrors = append(validationErrors, ErrEmptyPassword)
	}
	
	if len(validationErrors) > 0 {
		return validationErrors
	}
	
	return nil
}

// РАСШИРЕННЫЕ МЕТОДЫ ВАЛИДАЦИИ

// validateEmail выполняет комплексную валидацию email
func validateEmail(email string) ValidationErrors {
	var errors ValidationErrors
	
	// Проверка на пустоту
	if email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "email cannot be empty",
			Code:    "EMPTY_EMAIL",
		})
		return errors
	}
	
	// Проверка длины
	if len(email) > 255 {
		errors = append(errors, ErrEmailTooLong)
	}
	
	// Базовая проверка формата
	if !isValidEmailFormat(email) {
		errors = append(errors, ErrInvalidEmail)
		return errors // Если формат неверный, остальные проверки бессмысленны
	}
	
	// Проверка домена
	if !isValidEmailDomain(email) {
		errors = append(errors, ErrEmailInvalidTLD)
	}
	
	// Проверка на одноразовые email (простая реализация)
	if isDisposableEmail(email) {
		errors = append(errors, ErrEmailDisposable)
	}
	
	return errors
}

// validateName выполняет комплексную валидацию имени
func validateName(name string) ValidationErrors {
	var errors ValidationErrors
	
	// Удаляем лишние пробелы
	name = strings.TrimSpace(name)
	
	// Проверка на пустоту
	if name == "" {
		errors = append(errors, ErrEmptyName)
		return errors
	}
	
	// Проверка минимальной длины
	if len(name) < 2 {
		errors = append(errors, ErrNameTooShort)
	}
	
	// Проверка максимальной длины
	if len(name) > 100 {
		errors = append(errors, ErrNameTooLong)
	}
	
	// Проверка на допустимые символы
	if !isValidNameCharacters(name) {
		errors = append(errors, ErrInvalidNameChars)
	}
	
	return errors
}

// validatePassword выполняет комплексную валидацию пароля
func validatePassword(password, userName, userEmail string) error {
	var errors ValidationErrors
	
	// Проверка на пустоту
	if password == "" {
		return ErrEmptyPassword
	}
	
	// Проверка минимальной длины
	if len(password) < 8 {
		errors = append(errors, ErrPasswordTooWeak)
	}
	
	// Проверка максимальной длины
	if len(password) > 128 {
		errors = append(errors, ErrPasswordTooLong)
	}
	
	// Проверка на наличие заглавных букв
	if !hasUppercase(password) {
		errors = append(errors, ErrPasswordNoUppercase)
	}
	
	// Проверка на наличие строчных букв
	if !hasLowercase(password) {
		errors = append(errors, ErrPasswordNoLowercase)
	}
	
	// Проверка на наличие цифр
	if !hasDigit(password) {
		errors = append(errors, ErrPasswordNoDigit)
	}
	
	// Проверка на наличие специальных символов
	if !hasSpecialChar(password) {
		errors = append(errors, ErrPasswordNoSpecialChar)
	}
	
	// Проверка на общие пароли
	if isCommonPassword(password) {
		errors = append(errors, ErrPasswordCommonPassword)
	}
	
	// Проверка, что пароль не содержит имя пользователя
	if userName != "" && containsIgnoreCase(password, userName) {
		errors = append(errors, ErrPasswordContainsName)
	}
	
	// Проверка, что пароль не содержит email
	if userEmail != "" && containsIgnoreCase(password, strings.Split(userEmail, "@")[0]) {
		errors = append(errors, ErrPasswordContainsEmail)
	}
	
	if len(errors) > 0 {
		return errors
	}
	
	return nil
}

// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ДЛЯ ВАЛИДАЦИИ

// isValidEmailFormat проверяет базовый формат email
func isValidEmailFormat(email string) bool {
	// RFC 5322 compliant regex (упрощенная версия)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	return emailRegex.MatchString(email)
}

// isValidEmailDomain проверяет домен email
func isValidEmailDomain(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	domain := parts[1]
	
	// Проверка на валидные TLD (упрощенная)
	validTLDs := []string{".com", ".org", ".net", ".edu", ".gov", ".mil", ".int", ".ru", ".рф"}
	
	for _, tld := range validTLDs {
		if strings.HasSuffix(strings.ToLower(domain), tld) {
			return true
		}
	}
	
	// Проверка на двухбуквенные TLD (страны)
	if len(domain) >= 3 && domain[len(domain)-3] == '.' {
		return true
	}
	
	return false
}

// isDisposableEmail проверяет, является ли email одноразовым
func isDisposableEmail(email string) bool {
	// Простой список одноразовых доменов
	disposableDomains := []string{
		"10minutemail.com", "guerrillamail.com", "mailinator.com",
		"tempmail.org", "yopmail.com", "throwaway.email",
	}
	
	domain := strings.ToLower(strings.Split(email, "@")[1])
	
	for _, disposable := range disposableDomains {
		if domain == disposable {
			return true
		}
	}
	
	return false
}

// isValidNameCharacters проверяет допустимые символы в имени
func isValidNameCharacters(name string) bool {
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) && r != '-' && r != '\'' && r != '.' {
			return false
		}
	}
	return true
}

// hasUppercase проверяет наличие заглавных букв
func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// hasLowercase проверяет наличие строчных букв
func hasLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// hasDigit проверяет наличие цифр
func hasDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// hasSpecialChar проверяет наличие специальных символов
func hasSpecialChar(s string) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, r := range s {
		if strings.ContainsRune(specialChars, r) {
			return true
		}
	}
	return false
}

// isCommonPassword проверяет, является ли пароль слишком простым
func isCommonPassword(password string) bool {
	// Список самых распространенных паролей
	commonPasswords := []string{
		"password", "123456", "password123", "admin", "qwerty",
		"letmein", "welcome", "monkey", "1234567890", "abc123",
		"Password1", "password1", "123456789", "welcome123",
	}
	
	passwordLower := strings.ToLower(password)
	
	for _, common := range commonPasswords {
		if passwordLower == strings.ToLower(common) {
			return true
		}
	}
	
	return false
}

// containsIgnoreCase проверяет, содержит ли строка подстроку (без учета регистра)
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// isValidEmail проверяет валидность email согласно доменным правилам (legacy)
// Оставлено для совместимости, используйте validateEmail
func isValidEmail(email string) bool {
	return len(validateEmail(email)) == 0
}