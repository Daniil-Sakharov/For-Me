package main

import (
	"fmt"
	"os"
	"text/template"
)

// Упражнение 4: Генерация кода с использованием шаблонов
// Задача: Написать программу, которая генерирует код на основе шаблонов

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

// Шаблон для HTTP-обработчика
const handlerTemplate = `
func (h *{{.StructName}}) handler{{.MethodName}}(w http.ResponseWriter, r *http.Request) {
	// Проверка метода
	if r.Method != "{{.Method}}" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Проверка авторизации
	{{if .Auth}}
	if r.Header.Get("X-Auth") != "100500" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	{{end}}
	
	// Создание структуры параметров
	params := {{.StructName}}Params{}
	
	// Заполнение параметров из запроса
	{{range .Params}}
	{{if eq .Type "string"}}
	{{.Name}} := r.URL.Query().Get("{{.ParamName}}")
	{{if .Required}}
	if {{.Name}} == "" {
		http.Error(w, "{{.Name}} must be not empty", http.StatusBadRequest)
		return
	}
	{{end}}
	{{if gt .Min 0}}
	if len({{.Name}}) < {{.Min}} {
		http.Error(w, "{{.Name}} too short", http.StatusBadRequest)
		return
	}
	{{end}}
	{{if .Default}}
	if {{.Name}} == "" {
		{{.Name}} = "{{.Default}}"
	}
	{{end}}
	{{if .Enum}}
	valid{{.Name}} := false
	{{range .Enum}}
	if {{$.Name}} == "{{.}}" {
		valid{{$.Name}} = true
	}
	{{end}}
	if !valid{{.Name}} {
		http.Error(w, "{{.Name}} must be one of [{{range .Enum}}{{.}} {{end}}]", http.StatusBadRequest)
		return
	}
	{{end}}
	params.{{.Name}} = {{.Name}}
	{{end}}
	
	{{if eq .Type "int"}}
	{{.Name}}Str := r.URL.Query().Get("{{.ParamName}}")
	{{if .Required}}
	if {{.Name}}Str == "" {
		http.Error(w, "{{.Name}} must be not empty", http.StatusBadRequest)
		return
	}
	{{end}}
	{{.Name}}, err := strconv.Atoi({{.Name}}Str)
	if err != nil {
		http.Error(w, "{{.Name}} must be int", http.StatusBadRequest)
		return
	}
	{{if gt .Min 0}}
	if {{.Name}} < {{.Min}} {
		http.Error(w, "{{.Name}} too small", http.StatusBadRequest)
		return
	}
	{{end}}
	{{if gt .Max 0}}
	if {{.Name}} > {{.Max}} {
		http.Error(w, "{{.Name}} too big", http.StatusBadRequest)
		return
	}
	{{end}}
	{{if .Default}}
	if {{.Name}}Str == "" {
		{{.Name}} = {{.Default}}
	}
	{{end}}
	params.{{.Name}} = {{.Name}}
	{{end}}
	{{end}}
	
	// Вызов метода
	result, err := h.{{.MethodName}}(r.Context(), params)
	if err != nil {
		// Обработка ошибок
		if apiErr, ok := err.(ApiError); ok {
			http.Error(w, apiErr.Error(), apiErr.HTTPStatus)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Возврат результата
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": "",
		"response": result,
	})
}
`

func generateHandler(methodInfo MethodInfo, outputFile *os.File) error {
	// Создаем шаблон
	tmpl, err := template.New("handler").Parse(handlerTemplate)
	if err != nil {
		return err
	}
	
	// Выполняем шаблон
	return tmpl.Execute(outputFile, methodInfo)
}

func main() {
	// Пример данных для генерации
	methodInfo := MethodInfo{
		StructName: "MyApi",
		MethodName: "Profile",
		URL:        "/user/profile",
		Auth:       false,
		Method:     "GET",
		Params: []ParamInfo{
			{
				Name:      "Login",
				Type:      "string",
				Required:  true,
				ParamName: "login",
			},
		},
	}
	
	// Создаем файл для вывода
	outputFile, err := os.Create("generated_handler.go")
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return
	}
	defer outputFile.Close()
	
	// Записываем заголовок файла
	fmt.Fprintln(outputFile, "package main")
	fmt.Fprintln(outputFile)
	fmt.Fprintln(outputFile, `import (
	"encoding/json"
	"net/http"
	"strconv"
)`)
	fmt.Fprintln(outputFile)
	
	// Генерируем обработчик
	err = generateHandler(methodInfo, outputFile)
	if err != nil {
		fmt.Printf("Ошибка генерации: %v\n", err)
		return
	}
	
	fmt.Println("Код успешно сгенерирован в файл generated_handler.go")
	
	fmt.Println("\n=== ДОПОЛНИТЕЛЬНОЕ ЗАДАНИЕ ===")
	fmt.Println("1. Добавь генерацию ServeHTTP метода")
	fmt.Println("2. Добавь поддержку POST запросов с JSON")
	fmt.Println("3. Добавь более детальную обработку ошибок")
}