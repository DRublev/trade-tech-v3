package indicators__tests

import (
	"main/bot/indicators"
	"testing"
)

// Тестируем рассчет EMA сразу со всеми данными и имитируем поступление данных
func TestEma(t *testing.T) {
	testPrices := []float64{
		100, 110, 105, 115, 120, 130, 140, 150, 145, 155,
	}
	expectedEma := float64(128.7459360425953)
	emaOnce := indicators.NewEma(12)
	emaTwice := indicators.NewEma(12)

	emaOnce.Update(testPrices)

	emaTwice.Update(testPrices[:7])
	emaTwice.Update(testPrices[7:])

	valueOnce, err := emaOnce.Latest()
	valueTwice, err := emaOnce.Latest()

	if err != nil {
		t.Fatalf("Faced error %v", err)
	}

	if valueOnce != expectedEma {
		t.Fatalf("Wrong value %v", valueOnce)
	}

	if valueOnce != valueTwice {
		t.Fatalf("Vflues not equal %v != %v", valueOnce, valueTwice)
	}
}
