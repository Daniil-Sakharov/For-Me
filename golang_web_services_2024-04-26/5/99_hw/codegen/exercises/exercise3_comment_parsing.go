package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

// Упражнение 3: Парсинг комментариев и меток
// Задача: Написать программу, которая находит методы с меткой apigen:api
// и извлекает из неё JSON конфигурацию

type ApiConfig struct {
	URL    string `json:"url"`
	Auth   bool   `json:"auth"`
	Method string `json:"method"`
}

func parseApiComment(comment string) (*ApiConfig, error) {
	// Ищем метку apigen:api
	if !strings.Contains(comment, "apigen:api") {
		return nil, nil
	}
	
	// Извлекаем JSON часть
	start := strings.Index(comment, "{")
	end := strings.LastIndex(comment, "}")
	
	if start == -1 || end == -1 {
		return nil, fmt.Errorf("invalid JSON in comment")
	}
	
	jsonStr := comment[start:end+1]
	
	var config ApiConfig
	err := json.Unmarshal([]byte(jsonStr), &config)
	if err != nil {
		return nil, err
	}
	
	return &config, nil
}

func findApiMethods(filename string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("=== ПОИСК API МЕТОДОВ В %s ===\n", filename)
	
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// Проверяем, есть ли у функции комментарии
			if funcDecl.Doc != nil {
				for _, comment := range funcDecl.Doc.List {
					config, err := parseApiComment(comment.Text)
					if err != nil {
						fmt.Printf("Ошибка парсинга комментария для %s: %v\n", funcDecl.Name.Name, err)
						continue
					}
					
					if config != nil {
						fmt.Printf("Найден API метод: %s\n", funcDecl.Name.Name)
						fmt.Printf("  URL: %s\n", config.URL)
						fmt.Printf("  Auth: %v\n", config.Auth)
						fmt.Printf("  Method: %s\n", config.Method)
						fmt.Println()
					}
				}
			}
		}
	}
}

// Пример структуры с API методами
type ExampleApi struct{}

// apigen:api {"url": "/user/profile", "auth": false}
func (api *ExampleApi) Profile() {
	// Реализация метода
}

// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
func (api *ExampleApi) Create() {
	// Реализация метода
}

func main() {
	// Анализируем этот же файл
	findApiMethods("exercise3_comment_parsing.go")
	
	fmt.Println("=== ДОПОЛНИТЕЛЬНОЕ ЗАДАНИЕ ===")
	fmt.Println("Напиши функцию, которая:")
	fmt.Println("1. Находит все методы с меткой apigen:api")
	fmt.Println("2. Извлекает информацию о структуре параметров")
	fmt.Println("3. Генерирует HTTP-обработчики")
}