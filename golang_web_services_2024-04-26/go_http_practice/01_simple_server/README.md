# Создай http-сервер на Go с одним маршрутом /hello, который возвращает "Hello, World!".

# Используй net/http, не нужны никакие сторонние пакеты.

package main

import (
"fmt"
"net/http"
)

func main() {
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
fmt.Fprintln(w, "Hello, World!")
})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not Found")
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
