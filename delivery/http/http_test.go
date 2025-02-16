package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ttavito/domain/entities"
	"ttavito/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUsecase struct {
	mock.Mock
}

func (m *MockUsecase) GetInfo(ctx context.Context, username string) (*entities.InfoResponse, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entities.InfoResponse), args.Error(1)
}

func (m *MockUsecase) SendCoin(ctx context.Context, username, toUser string, amount int) error {
	args := m.Called(ctx, username, toUser, amount)
	return args.Error(0)
}

func (m *MockUsecase) BuyItem(ctx context.Context, username, item string) error {
	args := m.Called(ctx, username, item)
	return args.Error(0)
}

func (m *MockUsecase) Auth(ctx context.Context, username, password string) error {
	args := m.Called(ctx, username, password)
	return args.Error(0)
}

func TestGetInfoHandler_Success(t *testing.T) {
	mockUsecase := new(MockUsecase)
	mockUsecase.On("GetInfo", mock.Anything, "test_user").Return(&entities.InfoResponse{
		Coins: 100,
		Inventory: []entities.ItemResponse{
			{Type: "item1", Quantity: 10},
			{Type: "item2", Quantity: 5},
		},
		CoinHistory: entities.CoinHistoryResponse{
			Received: []entities.ReceivedResponse{
				{FromUser: "user1", Amount: 20},
			},
			Sent: []entities.SentResponse{
				{ToUser: "user2", Amount: 15},
			},
		},
	}, nil)

	req := httptest.NewRequest("GET", "/api/info", nil)
	req = req.WithContext(context.WithValue(req.Context(), internal.UsernameContextKey, "test_user"))

	rr := httptest.NewRecorder()
	handler := GetInfoHandler(mockUsecase)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response entities.InfoResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 100, response.Coins)
	assert.Len(t, response.Inventory, 2)
	assert.Len(t, response.CoinHistory.Received, 1)
	assert.Len(t, response.CoinHistory.Sent, 1)

	mockUsecase.AssertExpectations(t)
}

func TestSendCoinHandler_Success(t *testing.T) {
	mockUsecase := new(MockUsecase)
	mockUsecase.On("SendCoin", mock.Anything, "test_user", "recipient_user", 50).Return(nil)

	sendCoinRequest := entities.SendCoinRequest{
		ToUser: "recipient_user",
		Amount: 50,
	}
	sendCoinRequestBody, _ := json.Marshal(sendCoinRequest)

	req := httptest.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendCoinRequestBody))
	req = req.WithContext(context.WithValue(req.Context(), internal.UsernameContextKey, "test_user"))
	req = req.WithContext(context.WithValue(req.Context(), internal.ValidSendCoinKey, sendCoinRequest))

	rr := httptest.NewRecorder()
	handler := SendCoinHandler(mockUsecase)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response string
	_ = json.NewDecoder(rr.Body).Decode(&response)
	// assert.NoError(t, err)
	assert.Equal(t, "", response)

	mockUsecase.AssertExpectations(t)
}

func TestBuyItemHandler_Success(t *testing.T) {
	mockUsecase := new(MockUsecase)
	mockUsecase.On("BuyItem", mock.Anything, "test_user", "item_cup").Return(nil)

	req := httptest.NewRequest("POST", "/api/buy/cup", nil)
	req = req.WithContext(context.WithValue(req.Context(), internal.UsernameContextKey, "test_user"))
	req = req.WithContext(context.WithValue(req.Context(), internal.ValidBuyItemKey, "item_cup"))

	rr := httptest.NewRecorder()
	handler := BuyItemHandler(mockUsecase)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_Success(t *testing.T) {
	mockUsecase := new(MockUsecase)
	mockUsecase.On("Auth", mock.Anything, "test_user", "test_pass").Return(nil)

	authRequest := entities.AuthRequest{
		Username: "test_user",
		Password: "test_pass",
	}
	authRequestBody, _ := json.Marshal(authRequest)

	req := httptest.NewRequest("POST", "/api/auth", bytes.NewBuffer(authRequestBody))
	req = req.WithContext(context.WithValue(req.Context(), internal.ValidAuthReqKey, authRequest))

	rr := httptest.NewRecorder()
	handler := AuthHandler(mockUsecase)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response entities.AuthResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	// Проверка, что токен не пустой
	assert.NotEmpty(t, response.Token)

	mockUsecase.AssertExpectations(t)
}
