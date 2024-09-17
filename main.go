package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Структура для хранения информации о статусе пользователя
type UserOnlineStatus struct {
	UUID             string `json:"uuid"`
	LastEntrance     int64  `json:"lastEntrance"`
	Status           int    `json:"status"`
	ShowOnlineStatus bool   `json:"showOnlineStatus"`
}

// Контекст для операций с Redis
var ctx = context.Background()

func main() {
	// Настройка клиента Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",
		DB:       0,
	})

	// Получение JSON-данных пользователя из Redis по ключу
	userKey := "user:123"
	userJSON, err := rdb.Get(ctx, userKey).Result()
	if err == redis.Nil {
		log.Println("No data found for key:", userKey)
		return
	} else if err != nil {
		log.Fatalf("Could not get user data: %v", err)
	} else {
		fmt.Println("User data BEFORE update:", userJSON)
	}

	// Проверка на случай, если данные пустые
	if len(userJSON) == 0 {
		log.Println("Empty data received for key:", userKey)
		return
	}

	// Преобразование JSON-строки в структуру Go (десериализация)
	var userStatus UserOnlineStatus
	err = json.Unmarshal([]byte(userJSON), &userStatus)
	if err != nil {
		log.Fatalf("Could not unmarshal JSON: %v", err)
	}

	// Обновление статуса пользователя
	userStatus.Status = 1

	// Преобразование обновленной структуры обратно в JSON (сериализация)
	updatedUserJSON, err := json.Marshal(userStatus)
	if err != nil {
		log.Fatalf("Could not marshal updated struct to JSON: %v", err)
	}

	// Сохранение обновленных данных пользователя обратно в Redis
	err = rdb.Set(ctx, userKey, updatedUserJSON, 0).Err()
	if err != nil {
		log.Fatalf("Could not save updated JSON to Redis: %v", err)
	}

	// Сообщение об успешном обновлении статуса пользователя
	fmt.Println("Updated user status successfully!")

	// Вывод данных после обновления
	fmt.Println("User data AFTER update:", string(updatedUserJSON))

}
