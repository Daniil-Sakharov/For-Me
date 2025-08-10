package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Структура заметки
type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// Здесь будет твоя "БД"
var notes = make(map[int]Note)
var currentID = 1

func main() {
	fmt.Println("Запуск сервера на http://localhost:8080")
	http.ListenAndServe(":8080", NewMux())
}

// handlers: для POST, GET /notes, GET /notes/{id}, DELETE /notes/{id}
// не забудь Content-Type: application/json и всегда отдавай json-ответ!
func notesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetNoteHandler(w, r)
	case http.MethodPost:
		PostNoteHandler(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
		return
	}

}

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	var allNote []Note
	for _, n := range notes {
		allNote = append(allNote, n)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allNote)
}

func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	var newNote Note
	err := json.NewDecoder(r.Body).Decode(&newNote)
	if err != nil || newNote.Text == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
		return
	}

	newNote.ID = currentID
	notes[newNote.ID] = newNote
	currentID++
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newNote)
}

func noteByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
	}
	switch r.Method {
	case http.MethodGet:
		note, ok := notes[id]
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"note found"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(note)
	case http.MethodDelete:
		_, ok := notes[id]
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"note found"}`))
			return
		}
		delete(notes, id)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "deleted"}`))
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
	}
}

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/notes", func(writer http.ResponseWriter, request *http.Request) {
		notesHandler(writer, request)
	})
	mux.HandleFunc("/notes/", func(writer http.ResponseWriter, request *http.Request) {
		noteByIdHandler(writer, request)
	})
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(writer, "Not Found")
	})
	return mux
}
