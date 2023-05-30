package utils

import (
	"encoding/json"
	"fmt"
)

// Output interface in a pretty json in terminal
func ToJson(i interface{}) {
	bytes, _ := json.MarshalIndent(i, "", "    ")
	fmt.Println(string(bytes))
}
