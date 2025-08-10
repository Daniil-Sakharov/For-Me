package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// setupServer - это вспомогательная функция, которая настраивает ваш сервер для тестов.
// В реальном приложении ваши хендлеры будут в функции main(), но для тестов их лучше вынести.
func setupServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/echojson", echoHandler) // Предполагаем, что вы назовете ваш хендлер echoHandler
	return mux
}

func TestEchoJSON(t *testing.T) {
	testCases := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Success",
			body:         `{"message": "hello"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"Message":"HELLO"}`,
		},
		{
			name:         "BadRequest - Empty Body",
			body:         ``,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"bad request"}`,
		},
		{
			name:         "BadRequest - Invalid JSON",
			body:         `{"message":}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"bad request"}`,
		},
	}

	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := bytes.NewBufferString(tc.body)
			resp, err := http.Post(ts.URL+"/echojson", "application/json", reqBody)
			if err != nil {
				t.Fatalf("http.Post failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("expected status %d, got %d", tc.expectedCode, resp.StatusCode)
			}

			body, _ := ioutil.ReadAll(resp.Body)

			// Для успешного кейса сравниваем JSON-объекты, чтобы порядок полей не влиял
			if tc.expectedCode == http.StatusOK {
				var expected, actual EchoMessage
				json.Unmarshal([]byte(tc.expectedBody), &expected)
				json.Unmarshal(body, &actual)
				if expected != actual {
					t.Errorf("expected body %v, got %v", expected, actual)
				}
			} else {
				// Для ошибок можно сравнивать как строки
				if strings.TrimSpace(string(body)) != tc.expectedBody {
					t.Errorf("expected body %q, got %q", tc.expectedBody, body)
				}
			}
		})
	}
}
