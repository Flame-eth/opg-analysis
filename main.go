package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
)

type Stock struct {
	Ticker string
	Gap float64
	OpeningPrice float64
}

type Position struct {
	EntryPrice float64
	Shares int
	TakeProfitPrice float64
	StopLossPrice float64
	Profit float64
}

var accountBalance = 10000.0

var lossTolerance = 0.02

var maxLossPerTrade = accountBalance * lossTolerance

var profitPercent = 0.8

func CalculatePosition(gapPercent, openingPrice float64) Position {
	closingPrice := openingPrice / (1 + gapPercent)
	gapValue := closingPrice - openingPrice
	profitFromGap := gapValue * profitPercent

	stopLoss := openingPrice - profitFromGap
	takeProfit := openingPrice + profitFromGap

	shares := int(maxLossPerTrade / math.Abs(stopLoss - openingPrice))

	profit := math.Abs(openingPrice - takeProfit) * float64(shares)
	profit = math.Round(profit * 100) / 100

	return Position{
		EntryPrice: math.Round(openingPrice * 100) / 100,
		Shares: shares,
		TakeProfitPrice: math.Round(takeProfit * 100) / 100,
		StopLossPrice: math.Round(stopLoss * 100) / 100,
		Profit: math.Round(profit * 100) / 100,
	}
}


type Selection struct {
	Ticker string
	Position
}

func Load (path string) ([]Stock, error){
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()

	if err != nil {
		return nil, err
	}

	records = slices.Delete(records, 0, 1)

	stocks := make([]Stock, len(records))

	for _, record := range records {
		ticker := record[0]
		gap, err := strconv.ParseFloat(record[1], 64)

		if err != nil {
			return nil, err
		}


		openingPrice, err := strconv.ParseFloat(record[2], 64)

		if err != nil {
			return nil, err
		}

		

		stocks = append(stocks, Stock{Ticker: ticker, Gap: gap, OpeningPrice: openingPrice})
	}

	return stocks, nil
} 
func main() {
	

	stocks, err := Load("./opg.csv")

	if err != nil {
		fmt.Println("Error", err)
	}

	slices.DeleteFunc(stocks, func(stock Stock) bool {
		return math.Abs(stock.Gap) > 0.1
	})


	var selections []Selection

	for _, stock := range stocks {
		position := CalculatePosition(stock.Gap, stock.OpeningPrice)
		selection := Selection{stock.Ticker, position}
		selections = append(selections, selection)
	}


}