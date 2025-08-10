package main

import (
	"fmt"
	"net/http"
)

//	func main() {
//		// 1. Зарегистрируй обработчик для пути "/hello".
//		//    Функция-обработчик должна принимать http.ResponseWriter и *http.Request.
//		//    Внутри обработчика используй fmt.Fprintln(w, "Hello, World!") для отправки ответа.
//		// http.HandleFunc("/hello", func(...) { ... })
//
//		// 2. Зарегистрируй обработчик для всех остальных путей ("/").
//		//    Этот обработчик должен возвращать статус 404 Not Found.
//		//    Используй w.WriteHeader(http.StatusNotFound) перед отправкой тела ответа.
//		// http.HandleFunc("/", func(...) { ... })
//
//		// 3. Запусти сервер на порту :8080.
//		//    Используй http.ListenAndServe. Не забудь обработать возможную ошибку.
//		fmt.Println("Запуск сервера на http://localhost:8080")
//		// http.ListenAndServe(...)
//	}
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		for i := 0; i < 1000; i++ {
			fmt.Fprintln(writer, "Darya, I love you!!!!")
		}
	})
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(writer, "Not Found")
	})
	return mux
}

func main() {
	fmt.Println("Запуск сервера на http://localhost:8080")
	http.ListenAndServe(":8080", NewMux())
}
