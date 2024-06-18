package rosshook__tests

import (
	"main/types"

	"encoding/json"
	"fmt"
	"io/ioutil"
)

func GetMock() []types.OHLC {
	fixturePath := "./fixture.json" // VKCO 11.06.2024 1min 15:51-16:35

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
