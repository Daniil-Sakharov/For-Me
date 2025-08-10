package main

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"slices"
	"strings"
)

type Client struct {
	ID            int    `xml:"id"`
	Guid          string `xml:"guid"`
	IsActive      bool   `xml:"isActive"`
	Balance       string `xml:"balance"`
	Picture       string `xml:"picture"`
	Age           int    `xml:"age"`
	EyeColor      string `xml:"eyeColor"`
	FirstName     string `xml:"first_name"`
	LastName      string `xml:"last_name"`
	Gender        string `xml:"gender"`
	Company       string `xml:"company"`
	Email         string `xml:"email"`
	Phone         string `xml:"phone"`
	Address       string `xml:"address"`
	About         string `xml:"about"`
	Registered    string `xml:"registered"`
	FavoriteFruit string `xml:"favoriteFruit"`
}

type Root struct {
	Users []Client `xml:"row"`
}

type SearchParams struct {
	Query      string
	Limit      int
	Offset     int
	OrderField string
	OrderBy    int
}

func XMLParser(filename string) ([]Client, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("wrong path to file")
	}

	var root Root

	err = xml.Unmarshal(data, &root)
	if err != nil {
		return nil, errors.New("parsing error")
	}

	return root.Users, nil
}

func FindClients(clients []Client, params SearchParams) ([]User, error) {
	// 1. Инициализируем как пустой, а не nil срез
	result := make([]User, 0)

	query := strings.ToLower(params.Query)

	for _, user := range clients {
		// Сохраняем оригинальное имя
		fullName := user.FirstName + " " + user.LastName
		// Ищем по нижнему регистру
		nameForSearch := strings.ToLower(fullName)
		about := strings.ToLower(user.About)

		if query == "" ||
			strings.Contains(nameForSearch, query) ||
			strings.Contains(about, query) {
			// В результат кладём оригинальное имя
			result = append(result, User{
				Id:     user.ID,
				Name:   fullName, // Используем имя с заглавными буквами
				Age:    user.Age,
				About:  user.About,
				Gender: user.Gender,
			})
		}
	}

	// ... остальной код сортировки и пагинации ...
	orderField := params.OrderField
	if orderField == "" {
		orderField = "Name"
	}

	orderBy := params.OrderBy

	slices.SortFunc(result, func(a, b User) int {
		var cmp int
		switch orderField {
		case "Id":
			cmp = a.Id - b.Id
		case "Age":
			cmp = a.Age - b.Age
		case "Name":
			fallthrough
		default:
			// Сортируем по полю Name, которое уже есть в структуре User
			cmp = strings.Compare(a.Name, b.Name)
		}
		if orderBy == -1 {
			cmp = -cmp
		}
		return cmp
	})
	start := params.Offset
	if start > len(result) {
		start = len(result)
	}
	end := start + params.Limit
	if params.Limit == 0 || end > len(result) {
		end = len(result)
	}
	result = result[start:end]
	return result, nil
}
