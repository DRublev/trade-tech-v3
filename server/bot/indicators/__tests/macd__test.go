package indicators__tests

import (
	"main/bot/indicators"
	"testing"
)

func TestMacd(t *testing.T) {
	testPrices := []float64{100, 110, 105, 115, 120, 130, 140, 150, 145, 155}
	expectedValue := float64(1.4345525513726614)
	expectedSignal := float64(0.5196116513769141)

	ind := indicators.NewMacd(21, 16, 9)

	ind.Update(testPrices[1:5])
	ind.Update(testPrices[5:])

	macd, err := ind.Latest()

	if err != nil {
		t.Fatalf("Faced error %v", err)
	}

	if macd.Signal != expectedSignal || macd.Value != expectedValue {
		t.Fatalf("Wrong value %v", macd)
	}
}
