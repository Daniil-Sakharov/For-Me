package main

import (
	"fmt"
	"sync"
	"time"
)

// UserID - тип для идентификатора пользователя
type UserID int

// UserData - структура с данными пользователя
type UserData struct {
	ID   UserID
	Name string
	Age  int
}

// fetchUserData - имитирует долгий запрос к внешнему API для получения данных пользователя.
// В реальном приложении это мог бы быть HTTP-запрос.
func fetchUserData(id UserID) UserData {
	// Имитация сетевой задержки
	time.Sleep(100 * time.Millisecond)
	return UserData{
		ID:   id,
		Name: fmt.Sprintf("User-%d", id),
		Age:  int(id) + 20, // Просто для примера
	}
}

// --- ВАШЕ ЗАДАНИЕ ---

/*
	Необходимо реализовать функцию `ProcessUserBatch`, которая параллельно
	обрабатывает слайс идентификаторов пользователей (`userIDs`).

	Требования:
	1. Используйте пул воркеров (`worker pool`) для параллельной обработки.
	   Количество воркеров должно быть настраиваемым (`numWorkers`).

	2. Реализуйте ограничение частоты запросов (`rate limiting`).
	   Функция `fetchUserData` должна вызываться не чаще, чем `rateLimit` раз в секунду.
	   Для этого используйте `time.Ticker`.

	3. Функция должна возвращать мапу `map[UserID]UserData`, содержащую данные
	   всех успешно обработанных пользователей.

	4. Используйте `sync.WaitGroup` для ожидания завершения всех воркеров.
*/

// ProcessUserBatch - основная функция, которую вам нужно реализовать.
func ProcessUserBatch(userIDs []UserID, numWorkers int, rateLimit int) map[UserID]UserData {
	users := make(chan UserID)
	results := make(chan UserData)
	ticker := time.NewTicker(time.Second / time.Duration(rateLimit))
	go func() {
		defer close(users)
		for _, user := range userIDs {
			users <- user
		}
	}()
	wg := &sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go workker(users, results, ticker, wg)
	}
	resultMap := make(map[UserID]UserData)
	collectorWg := &sync.WaitGroup{}
	collectorWg.Add(1)
	go func() {
		defer collectorWg.Done()
		for resultData := range results {
			resultMap[resultData.ID] = resultData
		}
	}()

	wg.Wait()
	close(results)
	collectorWg.Wait()

	return resultMap
}

func workker(users <-chan UserID, results chan<- UserData, ticker *time.Ticker, wg *sync.WaitGroup) {
	defer wg.Done()
	for usersID := range users {
		<-ticker.C
		results <- fetchUserData(usersID)

	}
}

func main() {
	// Пример использования
	userIDs := make([]UserID, 0, 100)
	for i := 0; i < 100; i++ {
		userIDs = append(userIDs, UserID(i))
	}

	numWorkers := 5
	rateLimit := 100 // не более 10 запросов в секунду

	fmt.Println("Начинаем обработку пользователей...")
	startTime := time.Now()

	results := ProcessUserBatch(userIDs, numWorkers, rateLimit)

	elapsedTime := time.Since(startTime)
	fmt.Printf("Обработка завершена за %v\n", elapsedTime)
	if results != nil {
		fmt.Printf("Получено %d результатов\n", len(results))
		// Проверка нескольких результатов для примера
		if len(results) > 5 {
			fmt.Printf("Пример результата для UserID 5: %+v\n", results[5])
		}
	}
}
