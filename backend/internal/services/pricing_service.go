package services

import "math"

const (
	VolatilityFactor = 5.0
	BaselinePrice    = 10.0
	DecayRate        = 0.005
	MinPrice         = 0.01
	MaxPrice         = 10000.0
)

// CalculateNewPrice returns the new market price after a trade.
// Used by market maker for price nudges (no execution price needed).
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

// CalculateTradeExecution returns the new market price AND the execution price.
// The execution price is the average of pre-impact and post-impact price.
// This prevents buy-sell arbitrage: a buy→sell cycle always results in a net loss
// because you buy at the midpoint going up and sell at the midpoint going down.
//
// Math proof: buy at avg(P, P+d) then sell at avg(P+d, P+d-d') → net = -P*d²/2 < 0
func CalculateTradeExecution(currentPrice, netShares, totalShares float64) (newMarketPrice, executionPrice float64) {
	if totalShares == 0 {
		return currentPrice, currentPrice
	}

	priceChange := (netShares / totalShares) * VolatilityFactor
	newMarketPrice = currentPrice * (1 + priceChange)
	newMarketPrice = math.Max(MinPrice, math.Min(MaxPrice, newMarketPrice))
	newMarketPrice = math.Round(newMarketPrice*100) / 100

	// Execution price = average of pre-impact and post-impact price
	// This simulates slippage like a real AMM / order book
	executionPrice = (currentPrice + newMarketPrice) / 2
	executionPrice = math.Round(executionPrice*100) / 100

	return newMarketPrice, executionPrice
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
