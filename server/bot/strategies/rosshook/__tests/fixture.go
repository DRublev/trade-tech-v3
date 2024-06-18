package rosshook__tests

import (
	"main/types"

	"encoding/json"
	"fmt"
	"io/ioutil"
)

func GetMock() []types.OHLC {
	var result []types.OHLC
	file, err := ioutil.ReadFile("./fixture.json")
	if err != nil {
		return result
	}

	err = json.Unmarshal(file, &result)
	if err != nil {
		fmt.Println("error:", err)
	}

	return result
}
