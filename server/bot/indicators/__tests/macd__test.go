package indicators__tests

import (
	"main/bot/indicators"
	"testing"
)

func TestMacdSignal(t *testing.T) {
	// VKCO 1M 07 Jun 2024 19:12 - 19:34
	testPrices := []float64{568, 567.6, 567.2, 567.6, 567.4, 566.6, 567.2, 566.4, 566.8, 566.6, 566.2, 566.2, 564.6, 565, 565.4, 565, 565.2, 564.8, 565.6, 566.2, 566.6}
	expectedValue := 0.08
	expectedSignal := -0.07

	ind := indicators.NewMacd(10, 5, 3)

	for i := 0; i < len(testPrices); i++ {
		ind.Update(testPrices[i])
	}

	macd, signal, err := ind.Latest()
	if err != nil {
		t.Fatalf("Faced error %v", err)
	}
	if signal != expectedSignal || macd != expectedValue {
		t.Fatalf("Wrong value %v != %v || %v != %v", macd, expectedValue, signal, expectedSignal)
	}
}
