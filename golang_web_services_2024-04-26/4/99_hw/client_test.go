package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type TestCase struct {
	Name     string
	Request  SearchRequest
	Expected []User
	IsError  bool
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	order := r.URL.Query().Get("order_field")
	orderByStr := r.URL.Query().Get("order_by")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	orderBy, _ := strconv.Atoi(orderByStr)

	params := SearchParams{
		Query:      query,
		Limit:      limit,
		Offset:     offset,
		OrderField: order,
		OrderBy:    orderBy,
	}

	users, err := XMLParser("dataset.xml")
	if err != nil {
		http.Error(w, "Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := FindClients(users, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "JSON marshal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func TestFindUsersBasic(t *testing.T) {
	testCases := []TestCase{
		{
			Name: "Обычный поиск по имени",
			Request: SearchRequest{
				Query:      "Boyd",
				Limit:      1,
				Offset:     0,
				OrderField: "Name",
				OrderBy:    1,
			},
			Expected: []User{
				{Id: 0, Name: "Boyd Wolf", Age: 22, About: "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.", Gender: "male"},
			},
			IsError: false,
		},
		{
			Name: "Пустой результат (нет совпадений)",
			Request: SearchRequest{
				Query: "НетТакогоИмени",
				Limit: 1,
			},
			Expected: []User{},
			IsError:  false,
		},
		{
			Name: "Ошибка limit - отрицательное значение",
			Request: SearchRequest{
				Query: "Boyd",
				Limit: -1,
			},
			Expected: nil,
			IsError:  true,
		},
		// Добавляй кейсы: limit>25, плохой order_field, fail по offset, сервер 500 и т.д.
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	client := SearchClient{URL: ts.URL}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := client.FindUsers(tc.Request)
			if (err != nil) != tc.IsError {
				t.Fatalf("error: got %v, want error? %v", err, tc.IsError)
			}
			if !tc.IsError && len(result.Users) > 0 {
				got := result.Users[0]
				want := tc.Expected[0]
				if got.Id != want.Id || got.Name != want.Name || got.Age != want.Age || got.Gender != want.Gender {
					t.Errorf("got = %+v, want = %+v", got, want)
				}
			}
		})
	}
}
