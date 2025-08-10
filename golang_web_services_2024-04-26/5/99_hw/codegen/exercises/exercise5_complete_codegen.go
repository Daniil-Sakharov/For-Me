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

// Упражнение 5: Полный кодогенератор
// Задача: Создать простой кодогенератор, который объединяет все изученные концепции

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

func parseValidationTag(tag string) ValidationRule {
	rule := ValidationRule{}
	
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

func parseApiComment(comment string) (*ApiConfig, error) {
	if !strings.Contains(comment, "apigen:api") {
		return nil, nil
	}
	
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

func findStructByName(node *ast.File, structName string) *ast.StructType {
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.Name == structName {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						return structType
					}
				}
			}
		}
	}
	return nil
}

func extractParamInfo(field *ast.Field, rule ValidationRule) ParamInfo {
	paramName := rule.ParamName
	if paramName == "" {
		paramName = strings.ToLower(field.Names[0].Name)
	}
	
	fieldType := "string"
	if ident, ok := field.Type.(*ast.Ident); ok {
		fieldType = ident.Name
	}
	
	return ParamInfo{
		Name:       field.Names[0].Name,
		Type:       fieldType,
		Required:   rule.Required,
		Min:        rule.Min,
		Max:        rule.Max,
		Enum:       rule.Enum,
		Default:    rule.Default,
		ParamName:  paramName,
	}
}

func generateCode(inputFile, outputFile string) error {
	// Парсим входной файл
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, inputFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	
	// Собираем информацию о методах
	var methods []MethodInfo
	
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Doc != nil {
			// Проверяем комментарии на наличие apigen:api
			for _, comment := range funcDecl.Doc.List {
				config, err := parseApiComment(comment.Text)
				if err != nil {
					continue
				}
				
				if config != nil {
					// Находим структуру параметров
					var params []ParamInfo
					if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 1 {
						// Второй параметр - это структура параметров
						paramType := funcDecl.Type.Params.List[1].Type
						if starExpr, ok := paramType.(*ast.StarExpr); ok {
							if ident, ok := starExpr.X.(*ast.Ident); ok {
								structType := findStructByName(node, ident.Name)
								if structType != nil {
									for _, field := range structType.Fields.List {
										if field.Tag != nil {
											tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
											validationTag := tag.Get("apivalidator")
											if validationTag != "" {
												rule := parseValidationTag(validationTag)
												paramInfo := extractParamInfo(field, rule)
												params = append(params, paramInfo)
											}
										}
									}
								}
							}
						}
					}
					
					// Определяем имя структуры
					structName := "MyApi"
					if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
						if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
							if ident, ok := starExpr.X.(*ast.Ident); ok {
								structName = ident.Name
							}
						}
					}
					
					method := MethodInfo{
						StructName: structName,
						MethodName: funcDecl.Name.Name,
						URL:        config.URL,
						Auth:       config.Auth,
						Method:     config.Method,
						Params:     params,
					}
					
					methods = append(methods, method)
				}
			}
		}
	}
	
	// Генерируем код
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()
	
	// Записываем заголовок
	fmt.Fprintln(out, "package main")
	fmt.Fprintln(out)
	fmt.Fprintln(out, `import (
	"encoding/json"
	"net/http"
	"strconv"
)`)
	fmt.Fprintln(out)
	
	// Генерируем ServeHTTP метод
	if len(methods) == 0 {
		fmt.Fprintln(out, "// Нет методов для генерации")
		return nil
	}
	fmt.Fprintf(out, "func (h *%s) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n", methods[0].StructName)
	fmt.Fprintln(out, "\tswitch r.URL.Path {")
	for _, method := range methods {
		fmt.Fprintf(out, "\tcase \"%s\":\n", method.URL)
		fmt.Fprintf(out, "\t\th.handler%s(w, r)\n", method.MethodName)
	}
	fmt.Fprintln(out, "\tdefault:")
	fmt.Fprintln(out, "\t\thttp.Error(w, \"Not Found\", http.StatusNotFound)")
	fmt.Fprintln(out, "\t}")
	fmt.Fprintln(out, "}")
	fmt.Fprintln(out)
	
	// Генерируем обработчики для каждого метода
	for _, method := range methods {
		generateHandler(method, out)
	}
	
	return nil
}

func generateHandler(method MethodInfo, out *os.File) {
	// Шаблон для обработчика
	templateStr := `
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
	
	tmpl, err := template.New("handler").Parse(templateStr)
	if err != nil {
		log.Printf("Ошибка создания шаблона: %v", err)
		return
	}
	
	err = tmpl.Execute(out, method)
	if err != nil {
		log.Printf("Ошибка выполнения шаблона: %v", err)
	}
}

func main() {
	fmt.Printf("Аргументы: %v\n", os.Args)
	if len(os.Args) != 3 {
		fmt.Println("Использование: go run exercise5_complete_codegen.go <входной_файл> <выходной_файл>")
		fmt.Println("Пример: go run exercise5_complete_codegen.go api.go api_handlers.go")
		return
	}
	
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	
	fmt.Printf("Входной файл: %s\n", inputFile)
	fmt.Printf("Выходной файл: %s\n", outputFile)
	
	// Проверяем существование входного файла
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		log.Fatalf("Входной файл %s не найден", inputFile)
	}
	
	err := generateCode(inputFile, outputFile)
	if err != nil {
		log.Fatalf("Ошибка генерации кода: %v", err)
	}
	
	fmt.Printf("Код успешно сгенерирован в файл %s\n", outputFile)
}