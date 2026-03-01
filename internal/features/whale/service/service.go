package service

import (
	"log"
	"strconv"
)

func IsWhale(price string, quantity string) bool {
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Println("fail to convert price:", err)
		return false
	}
	quantityFloat, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		log.Println("fail to convert quantity:", err)
		return false
	}
	if priceFloat*quantityFloat >= 50000 {
		return true
	}
	return false
}
