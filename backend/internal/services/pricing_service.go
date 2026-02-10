package services

import "math"

const (
	VolatilityFactor = 0.1
	BaselinePrice    = 10.0
	DecayRate        = 0.005
	MinPrice         = 1.0
	MaxPrice         = 1000.0
)

func CalculateNewPrice(currentPrice, netShares, totalShares float64) float64 {
	if totalShares == 0 {
		return currentPrice
	}

	priceChange := (netShares / totalShares) * VolatilityFactor
	newPrice := currentPrice * (1 + priceChange)

	newPrice = math.Max(MinPrice, math.Min(MaxPrice, newPrice))
	newPrice = math.Round(newPrice*100) / 100

	return newPrice
}

func ApplyDecay(currentPrice float64) float64 {
	if math.Abs(currentPrice-BaselinePrice) < 0.01 {
		return currentPrice
	}

	var newPrice float64
	if currentPrice > BaselinePrice {
		newPrice = currentPrice * (1 - DecayRate)
		if newPrice < BaselinePrice {
			newPrice = BaselinePrice
		}
	} else {
		newPrice = currentPrice * (1 + DecayRate)
		if newPrice > BaselinePrice {
			newPrice = BaselinePrice
		}
	}

	newPrice = math.Round(newPrice*100) / 100
	return math.Max(MinPrice, math.Min(MaxPrice, newPrice))
}
