package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/profile":
		h.handlerProfile(w, r)
	case "/user/create":
		h.handlerCreate(w, r)
	case "/user/create":
		h.handlerCreate(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}


func (h *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	
	
	
	
	// Создание структуры параметров
	params := MyApiParams{}
	
	// Заполнение параметров из запроса
	
	
	// Вызов метода
	result, err := h.Profile(r.Context(), params)
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

func (h *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	
	// Проверка метода
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	
	
	// Проверка авторизации
	if r.Header.Get("X-Auth") != "100500" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	
	// Создание структуры параметров
	params := MyApiParams{}
	
	// Заполнение параметров из запроса
	
	
	// Вызов метода
	result, err := h.Create(r.Context(), params)
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

func (h *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	
	// Проверка метода
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	
	
	// Проверка авторизации
	if r.Header.Get("X-Auth") != "100500" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	
	// Создание структуры параметров
	params := OtherApiParams{}
	
	// Заполнение параметров из запроса
	
	
	// Вызов метода
	result, err := h.Create(r.Context(), params)
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
