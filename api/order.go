package api

type OrderParameters struct {
	OrderID    string      `json:"orderId"`
	CustomerID string      `json:"customerId"`
	OrderItems []OrderItem `json:"orderItems"`
}

type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"orderId"`
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderCreateResponse struct {
	ID string `json:"id"`
}

type OrderGetResponse struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customerId"`
	OrderItems []OrderItem `json:"orderItems"`
}
