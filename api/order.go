package api

// TODO: fill in order parameters

type OrderParameters struct {
}

type OrderCreateResponse struct {
	ID string `json:"id"`
}

type OrderGetResponse struct {
	ID string `json:"id"`
	//Name                   string                       `json:"name"`
}

type OrdersGetResponse struct {
	Orders []OrderGetResponse `json:"orders"`
}
