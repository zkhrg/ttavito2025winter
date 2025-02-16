package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Stats struct {
	successfulRequests int
	failedRequests     int
	totalRequests      int
	totalResponseTime  time.Duration
}

type AuthResponse struct {
	Token string `json:"token"`
}

func getAuthToken(username, password string) (string, error) {
	// Формируем тело запроса для получения токена
	authData := map[string]string{
		"username": username,
		"password": password,
	}
	authJSON, err := json.Marshal(authData)
	if err != nil {
		return "", err
	}

	// Отправляем POST-запрос на /api/auth
	resp, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(authJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Если ответ успешен, парсим токен
	if resp.StatusCode == http.StatusOK {
		var authResp AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			return "", err
		}
		return authResp.Token, nil
	}

	return "", fmt.Errorf("failed to get auth token, status code: %d", resp.StatusCode)
}

func makeRequest(url, method, token string, stats *Stats, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()

	// Создаем новый запрос с авторизационным токеном
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		stats.failedRequests++
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Выполняем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		stats.failedRequests++
	} else {
		stats.successfulRequests++
		stats.totalResponseTime += time.Since(startTime)
		resp.Body.Close() // Закрытие тела ответа
	}
	stats.totalRequests++
}

func makeNonStatRequest(url, token string) {
	// Создаем новый запрос с авторизационным токеном, но не учитываем статистику
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Выполняем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error performing non-stat request:", err)
	} else {
		resp.Body.Close() // Закрытие тела ответа
	}
}

func main() {
	// Получаем токен
	username := "your-username"
	password := "your-password"
	token, err := getAuthToken(username, password)
	if err != nil {
		fmt.Println("Error getting auth token:", err)
		return
	}

	url := "http://localhost:8080/api/info" // Укажите свою ручку
	requestsPerSecond := 1000               // Запросов в секунду (1k RPS)
	duration := 10 * time.Second            // Длительность теста

	url1 := "http://localhost:8080/api/buy/cup"
	makeNonStatRequest(url1, token)
	makeNonStatRequest(url1, token)
	makeNonStatRequest(url1, token)
	makeNonStatRequest(url1, token)

	var wg sync.WaitGroup
	stats := &Stats{}
	ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))

	// Выполняем нагрузочное тестирование заданное количество времени
	endTime := time.Now().Add(duration)
	for time.Now().Before(endTime) {
		wg.Add(1)
		<-ticker.C
		go makeRequest(url, http.MethodGet, token, stats, &wg)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	// Рассчитываем статистику
	rps := float64(stats.totalRequests) / duration.Seconds()
	avgResponseTime := float64(stats.totalResponseTime.Milliseconds()) / float64(stats.successfulRequests)
	successRate := float64(stats.successfulRequests) / float64(stats.totalRequests) * 100

	// Выводим результаты
	fmt.Printf("RPS: %.2f\n", rps)
	fmt.Printf("SLI (время ответа) (ms): %.2f\n", avgResponseTime)
	fmt.Printf("SLI (успешность): %.4f%%\n", successRate)
}
