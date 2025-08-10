package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupServer() *http.ServeMux {
	mux := http.NewServeMux()
	// Ты должен реализовать соответствующие handlers в main.go
	mux.HandleFunc("/notes", notesHandler)     // Для GET и POST списком
	mux.HandleFunc("/notes/", noteByIdHandler) // Для GET/DELETE по id
	return mux
}

func TestCRUDAPI(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	// 1. Создание заметки (POST)
	resp, err := http.Post(ts.URL+"/notes", "application/json", bytes.NewBufferString(`{"text": "abc"}`))
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST expected 200, got %d", resp.StatusCode)
	}
	var note Note
	json.NewDecoder(resp.Body).Decode(&note)
	resp.Body.Close()
	if note.Text != "abc" || note.ID <= 0 {
		t.Errorf("POST got %+v", note)
	}

	// 2. Получение одной заметки (GET)
	resp, err = http.Get(fmt.Sprintf("%s/notes/%d", ts.URL, note.ID))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("GET note failed")
	}
	var note2 Note
	json.NewDecoder(resp.Body).Decode(&note2)
	resp.Body.Close()
	if note2 != note {
		t.Errorf("GET got %+v, want %+v", note2, note)
	}

	// 3. Получение всех заметок (GET список)
	resp, err = http.Get(ts.URL + "/notes")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("GET all notes failed")
	}
	var notes []Note
	json.NewDecoder(resp.Body).Decode(&notes)
	resp.Body.Close()
	if len(notes) == 0 || notes[0] != note {
		t.Errorf("GET all notes, want first %v, got %v", note, notes)
	}

	// 4. Удаление заметки
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/notes/%d", ts.URL, note.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed")
	}
	var delResp map[string]string
	json.NewDecoder(resp.Body).Decode(&delResp)
	resp.Body.Close()
	if delResp["result"] != "deleted" {
		t.Errorf("DELETE wrong response: %v", delResp)
	}

	// 5. Проверка отсутствия заметки после удаления
	resp, _ = http.Get(fmt.Sprintf("%s/notes/%d", ts.URL, note.ID))
	if resp.StatusCode == http.StatusOK {
		t.Errorf("deleted note still accessible")
	}
}
