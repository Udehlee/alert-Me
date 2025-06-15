package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func UnmarshalJSON(data []byte, destination interface{}) error {
	if err := json.Unmarshal(data, destination); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

func ToFloatPrice(price string) (float64, error) {
	p := strings.TrimSpace(price)
	return strconv.ParseFloat(p, 64)
}
