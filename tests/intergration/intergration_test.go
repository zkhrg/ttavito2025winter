package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ttavito/config"
	"ttavito/database"
	myHttp "ttavito/delivery/http"
	"ttavito/repository"
	"ttavito/usecase"

	"github.com/stretchr/testify/assert"
)

func TestBuyItem_Success(t *testing.T) {
	cfg := config.LoadConfig()
	pool, err := database.NewPostgresDB(cfg)
	if err != nil {
		t.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	repo := repository.NewEntityRepo(pool)
	api := usecase.NewUsecase(repo)

	mux := http.NewServeMux()
	myHttp.SetupRoutes(api, mux)

	authRequest := map[string]string{
		"username": "test_user",
		"password": "test_pass",
	}

	authReqBody, _ := json.Marshal(authRequest)
	authReq, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authReqBody))
	authReq.Header.Set("Content-Type", "application/json")

	authRec := httptest.NewRecorder()
	mux.ServeHTTP(authRec, authReq)

	assert.Equal(t, http.StatusOK, authRec.Code)

	var authResponse map[string]interface{}
	err = json.NewDecoder(authRec.Body).Decode(&authResponse)
	if err != nil {
		t.Fatalf("Failed to parse auth response body: %v", err)
	}
	token, ok := authResponse["token"].(string)
	if !ok || token == "" {
		t.Fatalf("Failed to get token from auth response")
	}

	req, _ := http.NewRequest("GET", "/api/buy/cup", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&response)
	assert.Nil(t, response)
}

func TestAuth_Success(t *testing.T) {
	cfg := config.LoadConfig()
	pool, err := database.NewPostgresDB(cfg)
	if err != nil {
		t.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	repo := repository.NewEntityRepo(pool)
	api := usecase.NewUsecase(repo)

	mux := http.NewServeMux()
	myHttp.SetupRoutes(api, mux)

	authRequest := map[string]string{
		"username": "test_user",
		"password": "test_pass",
	}

	authReqBody, _ := json.Marshal(authRequest)
	authReq, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authReqBody))
	authReq.Header.Set("Content-Type", "application/json")

	authRec := httptest.NewRecorder()
	mux.ServeHTTP(authRec, authReq)

	assert.Equal(t, http.StatusOK, authRec.Code)

	var authResponse map[string]interface{}
	err = json.NewDecoder(authRec.Body).Decode(&authResponse)
	if err != nil {
		t.Fatalf("Failed to parse auth response body: %v", err)
	}
	token, ok := authResponse["token"].(string)
	if !ok || token == "" {
		t.Fatalf("Failed to get token from auth response")
	}

	t.Logf("Received token: %s", token)

	assert.NotEmpty(t, token, "Token should not be empty")
}

func TestInfo_Success(t *testing.T) {
	cfg := config.LoadConfig()
	pool, err := database.NewPostgresDB(cfg)
	if err != nil {
		t.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	repo := repository.NewEntityRepo(pool)
	api := usecase.NewUsecase(repo)

	mux := http.NewServeMux()
	myHttp.SetupRoutes(api, mux)

	authRequest := map[string]string{
		"username": "test_user",
		"password": "test_pass",
	}

	authReqBody, _ := json.Marshal(authRequest)
	authReq, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authReqBody))
	authReq.Header.Set("Content-Type", "application/json")

	authRec := httptest.NewRecorder()
	mux.ServeHTTP(authRec, authReq)

	assert.Equal(t, http.StatusOK, authRec.Code)

	var authResponse map[string]interface{}
	err = json.NewDecoder(authRec.Body).Decode(&authResponse)
	if err != nil {
		t.Fatalf("Failed to parse auth response body: %v", err)
	}
	token, ok := authResponse["token"].(string)
	if !ok || token == "" {
		t.Fatalf("Failed to get token from auth response")
	}

	req, _ := http.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.NotNil(t, response)
	assert.Contains(t, response, "coins")
	assert.Contains(t, response, "inventory")
}

func TestSendCoin_Success(t *testing.T) {
	cfg := config.LoadConfig()
	pool, err := database.NewPostgresDB(cfg)
	if err != nil {
		t.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	repo := repository.NewEntityRepo(pool)
	api := usecase.NewUsecase(repo)

	mux := http.NewServeMux()
	myHttp.SetupRoutes(api, mux)

	authRequest := map[string]string{
		"username": "test_user",
		"password": "test_pass",
	}

	authReqBody, _ := json.Marshal(authRequest)
	authReq, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authReqBody))
	authReq.Header.Set("Content-Type", "application/json")

	authRec := httptest.NewRecorder()
	mux.ServeHTTP(authRec, authReq)

	assert.Equal(t, http.StatusOK, authRec.Code)

	var authResponse map[string]interface{}
	err = json.NewDecoder(authRec.Body).Decode(&authResponse)
	if err != nil {
		t.Fatalf("Failed to parse auth response body: %v", err)
	}
	token, ok := authResponse["token"].(string)
	if !ok || token == "" {
		t.Fatalf("Failed to get token from auth response")
	}

	sendCoinRequest := map[string]interface{}{
		"toUser": "another_user",
		"amount": 50,
	}

	sendCoinReqBody, _ := json.Marshal(sendCoinRequest)
	sendCoinReq, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendCoinReqBody))
	sendCoinReq.Header.Set("Authorization", "Bearer "+token)
	sendCoinReq.Header.Set("Content-Type", "application/json")

	sendCoinRec := httptest.NewRecorder()
	mux.ServeHTTP(sendCoinRec, sendCoinReq)

	assert.Equal(t, http.StatusInternalServerError, sendCoinRec.Code)

	var sendCoinResponse map[string]interface{}
	_ = json.NewDecoder(sendCoinRec.Body).Decode(&sendCoinResponse)
	assert.Nil(t, sendCoinResponse)
}
