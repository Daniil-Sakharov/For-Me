# Финальное руководство по выполнению задания

Теперь, когда ты изучил все упражнения, ты готов к выполнению основного задания! Вот пошаговый план:

## Что ты уже знаешь

После изучения упражнений ты понимаешь:

1. **AST в Go** - как парсить код и извлекать информацию
2. **Структурные теги** - как извлекать правила валидации
3. **Парсинг комментариев** - как находить метки `apigen:api`
4. **Генерация кода** - как создавать код с помощью шаблонов

## Пошаговый план выполнения основного задания

### Шаг 1: Анализ требований

Твое задание требует создать кодогенератор, который:

1. **Находит методы** с меткой `apigen:api`
2. **Извлекает конфигурацию** из JSON в комментарии
3. **Анализирует структуры параметров** и их теги валидации
4. **Генерирует HTTP-обработчики** с валидацией
5. **Создает ServeHTTP метод** для маршрутизации

### Шаг 2: Структура кодогенератора

Создай файл `handlers_gen/codegen.go` со следующей структурой:

```go
package main

import (
    "encoding/json"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "log"
    "os"
    "reflect"
    "strings"
    "text/template"
)

// Структуры для хранения информации
type ApiConfig struct {
    URL    string `json:"url"`
    Auth   bool   `json:"auth"`
    Method string `json:"method"`
}

type ValidationRule struct {
    Required  bool
    Min       int
    Max       int
    Enum      []string
    Default   string
    ParamName string
}

type MethodInfo struct {
    StructName string
    MethodName string
    URL        string
    Auth       bool
    Method     string
    Params     []ParamInfo
}

type ParamInfo struct {
    Name       string
    Type       string
    Required   bool
    Min        int
    Max        int
    Enum       []string
    Default    string
    ParamName  string
}

func main() {
    if len(os.Args) != 3 {
        log.Fatal("Usage: codegen <input_file> <output_file>")
    }
    
    inputFile := os.Args[1]
    outputFile := os.Args[2]
    
    err := generateCode(inputFile, outputFile)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Шаг 3: Функции парсинга

Добавь функции из упражнений:

```go
func parseValidationTag(tag string) ValidationRule {
    // Код из упражнения 2
}

func parseApiComment(comment string) (*ApiConfig, error) {
    // Код из упражнения 3
}

func findStructByName(node *ast.File, structName string) *ast.StructType {
    // Код из упражнения 5
}
```

### Шаг 4: Основная функция генерации

```go
func generateCode(inputFile, outputFile string) error {
    // 1. Парсим входной файл
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, inputFile, nil, parser.ParseComments)
    if err != nil {
        return err
    }
    
    // 2. Собираем информацию о методах
    var methods []MethodInfo
    // ... код из упражнения 5
    
    // 3. Генерируем код
    out, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer out.Close()
    
    // 4. Записываем заголовок и импорты
    // 5. Генерируем ServeHTTP метод
    // 6. Генерируем обработчики для каждого метода
    
    return nil
}
```

### Шаг 5: Шаблоны для генерации

Используй шаблоны из упражнения 4, но адаптируй их под требования задания:

1. **ServeHTTP метод** - маршрутизация по URL
2. **Обработчики методов** - валидация, авторизация, вызов методов
3. **Обработка ошибок** - поддержка ApiError и общих ошибок

### Шаг 6: Особенности задания

Обрати внимание на специфические требования:

1. **Авторизация** - проверка заголовка `X-Auth: 100500`
2. **Валидация** - все правила из тегов `apivalidator`
3. **Формат ответа** - JSON с полями `error` и `response`
4. **Обработка ошибок** - правильные HTTP статусы
5. **Порядок проверок** - метод → авторизация → параметры

### Шаг 7: Тестирование

1. **Собери кодогенератор**:
   ```bash
   go build handlers_gen/* && ./codegen api.go api_handlers.go
   ```

2. **Запусти тесты**:
   ```bash
   go test -v
   ```

3. **Исправляй ошибки** по мере их появления

## Ключевые моменты для успеха

### 1. Правильный порядок проверок
```go
// 1. Проверка метода (если указан)
if r.Method != "POST" {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
}

// 2. Проверка авторизации (если требуется)
if r.Header.Get("X-Auth") != "100500" {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
}

// 3. Валидация параметров (в порядке полей структуры)
```

### 2. Правильный формат ошибок
```go
// Для валидации
http.Error(w, "login must me not empty", http.StatusBadRequest)

// Для ApiError
if apiErr, ok := err.(ApiError); ok {
    http.Error(w, apiErr.Error(), apiErr.HTTPStatus)
    return
}

// Для общих ошибок
http.Error(w, err.Error(), http.StatusInternalServerError)
```

### 3. Правильный формат ответа
```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{
    "error": "",
    "response": result,
})
```

## Отладка

Если что-то не работает:

1. **Добавляй отладочные выводы**:
   ```go
   fmt.Printf("Найден метод: %s\n", methodName)
   fmt.Printf("Параметры: %+v\n", params)
   ```

2. **Проверяй типы**:
   ```go
   fmt.Printf("type: %T, value: %+v\n", value, value)
   ```

3. **Тестируй по частям** - сначала парсинг, потом генерация

## Готовые решения

Если застрянешь, можешь посмотреть на:
- `exercise5_complete_codegen.go` - полный пример
- `example/gen/codegen.go` - пример из лекции
- Сгенерированный код в `test_handlers.go`

## Финальная проверка

Перед сдачей убедись, что:

1. ✅ Кодогенератор собирается и запускается
2. ✅ Генерирует правильный код
3. ✅ Все тесты проходят
4. ✅ Код работает с разными структурами (MyApi, OtherApi)
5. ✅ Все правила валидации поддерживаются
6. ✅ Обработка ошибок корректна

Удачи! Ты готов к выполнению задания! 🚀