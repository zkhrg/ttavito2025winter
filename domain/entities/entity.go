package entities

type InfoResponse struct {
	Coins       int                 `json:"coins"`
	Inventory   []ItemResponse      `json:"inventory"`
	CoinHistory CoinHistoryResponse `json:"coinHistory"`
}

type ItemResponse struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistoryResponse struct {
	Received []ReceivedResponse `json:"received"`
	Sent     []SentResponse     `json:"sent"`
}

type ReceivedResponse struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentResponse struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type SendCoinRequest SentResponse

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}
