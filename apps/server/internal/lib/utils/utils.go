package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v any) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Println("JSON:", string(jsonData))
}
