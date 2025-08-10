package main

import (
	"fmt"
	"reflect"
	"strings"
)

// Упражнение 2: Работа со структурными тегами
// Задача: Написать программу, которая анализирует структурные теги
// и извлекает из них информацию для валидации

type ValidationRule struct {
	Required bool
	Min      int
	Max      int
	Enum     []string
	Default  string
	ParamName string
}

func parseValidationTag(tag string) ValidationRule {
	rule := ValidationRule{}
	
	// Убираем кавычки и разделяем по запятой
	tag = strings.Trim(tag, `"`)
	parts := strings.Split(tag, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		
		switch {
		case part == "required":
			rule.Required = true
		case strings.HasPrefix(part, "min="):
			fmt.Sscanf(part, "min=%d", &rule.Min)
		case strings.HasPrefix(part, "max="):
			fmt.Sscanf(part, "max=%d", &rule.Max)
		case strings.HasPrefix(part, "enum="):
			enumPart := strings.TrimPrefix(part, "enum=")
			rule.Enum = strings.Split(enumPart, "|")
		case strings.HasPrefix(part, "default="):
			rule.Default = strings.TrimPrefix(part, "default=")
		case strings.HasPrefix(part, "paramname="):
			rule.ParamName = strings.TrimPrefix(part, "paramname=")
		}
	}
	
	return rule
}

func analyzeStruct(s interface{}) {
	t := reflect.TypeOf(s)
	
	fmt.Printf("=== АНАЛИЗ СТРУКТУРЫ: %s ===\n", t.Name())
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("Поле: %s (тип: %s)\n", field.Name, field.Type)
		
		// Получаем тег apivalidator
		tag := field.Tag.Get("apivalidator")
		if tag != "" {
			fmt.Printf("  Тег apivalidator: %s\n", tag)
			
			// Парсим правила валидации
			rule := parseValidationTag(tag)
			fmt.Printf("  Правила валидации:\n")
			if rule.Required {
				fmt.Printf("    - Обязательное поле\n")
			}
			if rule.Min > 0 {
				fmt.Printf("    - Минимум: %d\n", rule.Min)
			}
			if rule.Max > 0 {
				fmt.Printf("    - Максимум: %d\n", rule.Max)
			}
			if len(rule.Enum) > 0 {
				fmt.Printf("    - Допустимые значения: %v\n", rule.Enum)
			}
			if rule.Default != "" {
				fmt.Printf("    - Значение по умолчанию: %s\n", rule.Default)
			}
			if rule.ParamName != "" {
				fmt.Printf("    - Имя параметра: %s\n", rule.ParamName)
			}
		}
		fmt.Println()
	}
}

// Пример структуры с тегами валидации
type UserParams struct {
	Login  string `apivalidator:"required,min=3"`
	Name   string `apivalidator:"paramname=full_name"`
	Status string `apivalidator:"enum=user|moderator|admin,default=user"`
	Age    int    `apivalidator:"min=0,max=120"`
}

func main() {
	// Анализируем структуру
	userParams := UserParams{}
	analyzeStruct(userParams)
	
	// Дополнительное упражнение: попробуй написать функцию валидации
	fmt.Println("=== ДОПОЛНИТЕЛЬНОЕ ЗАДАНИЕ ===")
	fmt.Println("Напиши функцию validateStruct, которая:")
	fmt.Println("1. Принимает структуру")
	fmt.Println("2. Проверяет все правила валидации")
	fmt.Println("3. Возвращает список ошибок")
}