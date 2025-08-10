package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type EchoMessage struct {
	Message string `json:"message"`
}

func main() {
	fmt.Println("Запуск сервера на http://localhost:8080")
	http.ListenAndServe(":8080", NewMux())
}

// Обработчик /echojson
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/echojson", func(writer http.ResponseWriter, request *http.Request) {
		echoHandler(writer, request)
	})
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(writer, "Not Found")
	})
	return mux
}

//1. Должен принимать POST только с application/json
// 2. Разбирать json в EchoMessage
// 3. Переводить message в upper-case и отправлять обратно
// 4. При ошибке — отдавать 400 и {"error": "bad request"}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
		return
	}

	var msg EchoMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil || msg.Message == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
		return
	}

	resp := EchoMessage{Message: strings.ToUpper(msg.Message)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
