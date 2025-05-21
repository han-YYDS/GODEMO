package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	data := []byte(`{"user_id": 123, "event_type": "login", "time": "2023-10-01"}`)
	mapp := make(map[string]interface{}, 0)
	json.Unmarshal(data, &mapp)

	fmt.Println(mapp)
}
