package rosshook__tests

import (
	"main/types"

	"encoding/json"
	"fmt"
	"io/ioutil"
)

func readJSON(fixturePath string) []types.OHLC {
	var result []types.OHLC
	file, err := ioutil.ReadFile(fixturePath)
	if err != nil {
		return result
	}

	err = json.Unmarshal(file, &result)
	if err != nil {
		fmt.Println("error:", err)
	}
	return result
}

func getShouldBuyMock() []types.OHLC {
	fixturePath := "./vkco_2024-06-11_1min_13-04_13-34.json"

	return readJSON(fixturePath)
}

func getShouldNotBuyMock() []types.OHLC {
	fixturePath := "./vkco_2024-06-24_1min_12-05_12-25.json"
	return readJSON(fixturePath)
}

func getShouldCloseBuyWhenNotExecutedMock() []types.OHLC {
	fixturePath := "./vkco_2024-06-24_1min_10-38_11-20.json"
	return readJSON(fixturePath)
}
func getBuyAndStopLossMock() []types.OHLC {
	fixturePath := "./vkco_2024-07-02_1min_11-14_12-10.json"
	return readJSON(fixturePath)

}
