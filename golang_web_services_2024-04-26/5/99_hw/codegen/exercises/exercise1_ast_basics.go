package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

// Упражнение 1: Основы работы с AST
// Задача: Написать программу, которая парсит Go файл и выводит:
// 1. Все имена функций
// 2. Все имена структур
// 3. Все поля структур с их типами

func main() {
	// Создаем файловый набор токенов
	fset := token.NewFileSet()
	
	// Парсим файл (используем этот же файл как пример)
	node, err := parser.ParseFile(fset, "exercise1_ast_basics.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== АНАЛИЗ AST ===")
	
	// Проходим по всем объявлениям в файле
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			// Это функция или метод
			if d.Recv != nil {
				fmt.Printf("Метод: %s\n", d.Name.Name)
			} else {
				fmt.Printf("Функция: %s\n", d.Name.Name)
			}
			
		case *ast.GenDecl:
			// Это может быть объявление типа, переменной, константы
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						fmt.Printf("Тип: %s\n", typeSpec.Name.Name)
						
						// Если это структура, выводим её поля
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							fmt.Printf("  Поля структуры %s:\n", typeSpec.Name.Name)
							for _, field := range structType.Fields.List {
								fieldName := field.Names[0].Name
								fieldType := fmt.Sprintf("%v", field.Type)
								fmt.Printf("    %s: %s\n", fieldName, fieldType)
							}
						}
					}
				}
			}
		}
	}
}

// Пример структуры для анализа
type ExampleStruct struct {
	Name string
	Age  int
}

// Пример функции для анализа
func exampleFunction() {
	fmt.Println("Hello, AST!")
}

// Пример метода для анализа
func (e *ExampleStruct) exampleMethod() {
	fmt.Println("Method called")
}