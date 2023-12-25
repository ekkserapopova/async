package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const password = "documents_for_change" //какой-то пароль, вероятно каневский токен

func StartServer() {
	log.Println("Server start up")

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/name", func(c *gin.Context) { //переделать на мфц
		var data NameData

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		AppID := data.AppID //эпп айди

		// Запуск горутины для отправки статуса
		go sendName(AppID, password, fmt.Sprintf("http://127.0.0.1:8000/applications/%d/mfc/", AppID))

		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated"})
	})
	router.Run(":8080")

	log.Println("Server down")
}

func genRandomName(password string) Result {
	time.Sleep(10 * time.Second)
	mfc := [2]string{
		"Отправлена",
		"Не отправлена",
	}

	rand.Seed(time.Now().UnixNano())   // Устанавливаем seed для генератора случайных чисел
	randomIndex := rand.Intn(len(mfc)) // Генерируем случайный индекс в пределах длины массива
	randomFIO := mfc[randomIndex]      // Выбираем случайную запись из массива

	fmt.Println(randomFIO)

	return Result{randomFIO, password}
}

// Функция для отправки статуса в отдельной горутине
func sendName(AppID int, password string, url string) {
	// Выполнение расчётов с randomStatus
	result := genRandomName(password)

	// Отправка PUT-запроса к основному серверу
	_, err := performPUTRequest(url, result)
	if err != nil {
		fmt.Println("Error sending Name:", err)
		return
	}

	fmt.Println("Name sent successfully for AppID:", AppID)
}

type Result struct {
	MFC      string `json:"mfc_status"`
	Password string `json:"password"`
}

type NameData struct {
	AppID int `json:"application_id"`
}

func performPUTRequest(url string, data Result) (*http.Response, error) {
	// Сериализация структуры в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Создание PUT-запроса
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}

func main() {
	log.Println("App start")
	StartServer()
	log.Println("App stop")
}
