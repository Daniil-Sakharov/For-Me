# Резюме изученного материала

## Что ты изучил

### 1. AST (Abstract Syntax Tree) в Go
- **Пакеты**: `go/ast`, `go/parser`, `go/token`
- **Ключевые структуры**:
  - `ast.File` - файл Go
  - `ast.FuncDecl` - объявление функции/метода
  - `ast.GenDecl` - объявление типа, переменной, константы
  - `ast.StructType` - структура
  - `ast.Field` - поле структуры
- **Парсинг**: `parser.ParseFile(fset, filename, nil, parser.ParseComments)`

### 2. Структурные теги
- **Пакет**: `reflect`
- **Извлечение**: `field.Tag.Get("apivalidator")`
- **Парсинг**: разбор строки по запятой и обработка каждого правила
- **Поддерживаемые правила**:
  - `required` - обязательное поле
  - `min=X` - минимальное значение/длина
  - `max=X` - максимальное значение
  - `enum=val1|val2|val3` - допустимые значения
  - `default=value` - значение по умолчанию
  - `paramname=name` - имя параметра в запросе

### 3. Парсинг комментариев
- **Доступ**: `funcDecl.Doc.List` - список комментариев
- **Поиск меток**: `strings.Contains(comment.Text, "apigen:api")`
- **Извлечение JSON**: поиск `{` и `}` в комментарии
- **Парсинг JSON**: `json.Unmarshal([]byte(jsonStr), &config)`

### 4. Генерация кода
- **Пакет**: `text/template`
- **Создание шаблона**: `template.New("name").Parse(templateString)`
- **Выполнение**: `tmpl.Execute(outputFile, data)`
- **Условия в шаблонах**: `{{if .Condition}}...{{end}}`
- **Циклы**: `{{range .Items}}...{{end}}`

## Алгоритм работы кодогенератора

### 1. Парсинг входного файла
```go
fset := token.NewFileSet()
node, err := parser.ParseFile(fset, inputFile, nil, parser.ParseComments)
```

### 2. Поиск методов с меткой apigen:api
```go
for _, decl := range node.Decls {
    if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Doc != nil {
        for _, comment := range funcDecl.Doc.List {
            config, err := parseApiComment(comment.Text)
            if config != nil {
                // Обработка метода
            }
        }
    }
}
```

### 3. Анализ структуры параметров
```go
// Находим второй параметр метода (структура параметров)
paramType := funcDecl.Type.Params.List[1].Type
if starExpr, ok := paramType.(*ast.StarExpr); ok {
    if ident, ok := starExpr.X.(*ast.Ident); ok {
        structType := findStructByName(node, ident.Name)
        // Анализируем поля структуры
    }
}
```

### 4. Извлечение правил валидации
```go
for _, field := range structType.Fields.List {
    if field.Tag != nil {
        tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
        validationTag := tag.Get("apivalidator")
        if validationTag != "" {
            rule := parseValidationTag(validationTag)
            // Сохраняем информацию о параметре
        }
    }
}
```

### 5. Генерация кода
```go
// Записываем заголовок файла
fmt.Fprintln(out, "package main")
fmt.Fprintln(out, `import (...)`)

// Генерируем ServeHTTP метод
fmt.Fprintf(out, "func (h *%s) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n", structName)

// Генерируем обработчики для каждого метода
for _, method := range methods {
    generateHandler(method, out)
}
```

## Ключевые шаблоны для генерации

### ServeHTTP метод
```go
func (h *{{.StructName}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    {{range .Methods}}
    case "{{.URL}}":
        h.handler{{.MethodName}}(w, r)
    {{end}}
    default:
        http.Error(w, "Not Found", http.StatusNotFound)
    }
}
```

### Обработчик метода
```go
func (h *{{.StructName}}) handler{{.MethodName}}(w http.ResponseWriter, r *http.Request) {
    {{if ne .Method ""}}
    // Проверка метода
    if r.Method != "{{.Method}}" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    {{end}}
    
    {{if .Auth}}
    // Проверка авторизации
    if r.Header.Get("X-Auth") != "100500" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    {{end}}
    
    // Валидация параметров
    {{range .Params}}
    // Код валидации для каждого параметра
    {{end}}
    
    // Вызов метода и обработка результата
    result, err := h.{{.MethodName}}(r.Context(), params)
    if err != nil {
        if apiErr, ok := err.(ApiError); ok {
            http.Error(w, apiErr.Error(), apiErr.HTTPStatus)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error": "",
        "response": result,
    })
}
```

## Важные моменты

### 1. Порядок проверок
1. Метод HTTP (если указан)
2. Авторизация (если требуется)
3. Валидация параметров (в порядке полей структуры)

### 2. Формат ошибок
- Валидация: `http.Error(w, "field must be not empty", http.StatusBadRequest)`
- ApiError: `http.Error(w, apiErr.Error(), apiErr.HTTPStatus)`
- Общие ошибки: `http.Error(w, err.Error(), http.StatusInternalServerError)`

### 3. Формат ответа
```json
{
    "error": "",
    "response": { ... }
}
```

### 4. Авторизация
Проверка заголовка: `X-Auth: 100500`

## Готов к выполнению задания!

Теперь у тебя есть все необходимые знания и инструменты для выполнения основного задания. Ты понимаешь:

- ✅ Как парсить Go код с помощью AST
- ✅ Как извлекать информацию из структурных тегов
- ✅ Как находить и парсить комментарии с метками
- ✅ Как генерировать код с помощью шаблонов
- ✅ Как структурировать кодогенератор

Удачи! 🚀