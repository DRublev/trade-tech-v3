package indicators__tests

import (
	"main/bot/indicators"
	"testing"
)

// Тестируем рассчет EMA сразу со всеми данными
func TestEma(t *testing.T) {
	// VKCO 1M 07 Jun 2024 19:12 - 19:34
	testPrices := []float64{568, 567.6, 567.2, 567.6, 567.4, 566.6, 567.2, 566.4, 566.8, 566.6, 566.2, 566.2, 564.6, 565, 565.4, 565, 565.2, 564.8, 565.6, 566.2, 566.6}
	expectedEma := 565.82
	ind := indicators.NewEma(9)

	for i := 0; i < len(testPrices); i++ {
		ind.Update(testPrices[i])
	}

	value, err := ind.Latest()
	if err != nil {
		t.Fatalf("Faced error %v", err)
	}
	if value != expectedEma {
		t.Fatalf("Wrong value %v", value)
	}
}
