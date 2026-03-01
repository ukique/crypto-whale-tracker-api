package models

type Whale struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
	Symbol   string `json:"symbol"`
	Time     int64  `json:"time"`
}
