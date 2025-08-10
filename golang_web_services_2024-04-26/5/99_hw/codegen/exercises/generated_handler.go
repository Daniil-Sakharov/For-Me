package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)


func (h *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	// Проверка метода
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Проверка авторизации
	
	
	// Создание структуры параметров
	params := MyApiParams{}
	
	// Заполнение параметров из запроса
	
	
	Login := r.URL.Query().Get("login")
	
	if Login == "" {
		http.Error(w, "Login must be not empty", http.StatusBadRequest)
		return
	}
	
	
	
	
	params.Login = Login
	
	
	
	
	
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
