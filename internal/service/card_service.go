package service

import "github.com/gomarchy/estimate/internal/domain"

func FibonacciCards() []domain.Card {
	return []domain.Card{
		{DisplayValue: "1", Value: 1},
		{DisplayValue: "2", Value: 2},
		{DisplayValue: "3", Value: 3},
		{DisplayValue: "5", Value: 5},
		{DisplayValue: "8", Value: 8},
		{DisplayValue: "13", Value: 13},
		{DisplayValue: "?", Value: 0},
	}
}
