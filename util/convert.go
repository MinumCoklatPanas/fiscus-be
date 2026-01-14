package util

import (
	"encoding/json"
	"time"
)

func StructToMap(s interface{}) map[string]interface{} {
	jsonData, _ := json.Marshal(s)

	var result map[string]interface{}
	_ = json.Unmarshal(jsonData, &result)

	return result
}

func StringToTime(s string) (time.Time, error) {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Now(), err
	}

	return d, err
}
