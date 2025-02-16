package usecase

import (
	"context"
	"fmt"
	"testing"
	"ttavito/domain/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockShopRepository struct {
	mock.Mock
}

func (m *MockShopRepository) GetInfo(ctx context.Context, username string) (*entities.InfoResponse, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entities.InfoResponse), args.Error(1)
}

func (m *MockShopRepository) BuyItem(ctx context.Context, username, item string) error {
	args := m.Called(ctx, username, item)
	return args.Error(0)
}

func (m *MockShopRepository) SendCoin(ctx context.Context, senderUsername, recipientUsername string, amount int) error {
	args := m.Called(ctx, senderUsername, recipientUsername, amount)
	return args.Error(0)
}

func (m *MockShopRepository) Auth(ctx context.Context, username, password string) (bool, error) {
	args := m.Called(ctx, username, password)
	return args.Bool(0), args.Error(1)
}

func TestGetInfo(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("GetInfo", mock.Anything, "testUser").Return(&entities.InfoResponse{
		Coins: 1000,
		Inventory: []entities.ItemResponse{
			{Type: "t-shirt", Quantity: 2},
			{Type: "cup", Quantity: 1},
		},
		CoinHistory: entities.CoinHistoryResponse{
			Received: []entities.ReceivedResponse{
				{FromUser: "userA", Amount: 100},
			},
			Sent: []entities.SentResponse{
				{ToUser: "userB", Amount: 50},
			},
		},
	}, nil)

	info, err := uc.GetInfo(context.Background(), "testUser")

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, 1000, info.Coins)
	assert.Len(t, info.Inventory, 2)
	assert.Equal(t, "t-shirt", info.Inventory[0].Type)
	assert.Equal(t, 2, info.Inventory[0].Quantity)
	assert.Equal(t, "userA", info.CoinHistory.Received[0].FromUser)
	assert.Equal(t, 100, info.CoinHistory.Received[0].Amount)

	mockRepo.AssertExpectations(t)
}

func TestBuyItem(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("BuyItem", mock.Anything, "testUser", "t-shirt").Return(nil)

	err := uc.BuyItem(context.Background(), "testUser", "t-shirt")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSendCoin(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("SendCoin", mock.Anything, "user1", "user2", 100).Return(nil)

	err := uc.SendCoin(context.Background(), "user1", "user2", 100)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAuth(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("Auth", mock.Anything, "testUser", "password123").Return(true, nil)

	err := uc.Auth(context.Background(), "testUser", "password123")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAuthFailure(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("Auth", mock.Anything, "testUser", "wrongPassword").Return(false, fmt.Errorf("invalid credentials"))

	err := uc.Auth(context.Background(), "testUser", "wrongPassword")

	assert.Error(t, err)
	assert.Equal(t, "passwords does not match", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestBuyItemInsufficientCoins(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("BuyItem", mock.Anything, "testUser", "powerbank").Return(fmt.Errorf("not enough coins"))

	err := uc.BuyItem(context.Background(), "testUser", "powerbank")

	assert.Error(t, err)
	assert.Equal(t, "not enough coins", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestSendCoinInsufficientBalance(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("SendCoin", mock.Anything, "user1", "user2", 100).Return(fmt.Errorf("not enough coins"))

	err := uc.SendCoin(context.Background(), "user1", "user2", 100)

	assert.Error(t, err)
	assert.Equal(t, "not enough coins", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetInfoUserNotFound(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("GetInfo", mock.Anything, "nonExistentUser").Return((*entities.InfoResponse)(nil), fmt.Errorf("user not found"))

	info, err := uc.GetInfo(context.Background(), "nonExistentUser")

	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAuthWrongUsername(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("Auth", mock.Anything, "wrongUser", "password123").Return(false, fmt.Errorf("invalid credentials"))

	err := uc.Auth(context.Background(), "wrongUser", "password123")

	assert.Error(t, err)
	assert.Equal(t, "passwords does not match", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAuthWrongPassword(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("Auth", mock.Anything, "testUser", "wrongPassword").Return(false, fmt.Errorf("invalid credentials"))

	err := uc.Auth(context.Background(), "testUser", "wrongPassword")

	assert.Error(t, err)
	assert.Equal(t, "passwords does not match", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestBuyItemSuccess(t *testing.T) {
	mockRepo := new(MockShopRepository)
	uc := NewUsecase(mockRepo)

	mockRepo.On("BuyItem", mock.Anything, "testUser", "hoody").Return(nil)

	err := uc.BuyItem(context.Background(), "testUser", "hoody")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
