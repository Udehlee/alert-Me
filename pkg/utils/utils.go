package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func UnmarshalJSON(data []byte, destination interface{}) error {
	if err := json.Unmarshal(data, destination); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

// ToFloat converts a price string to float64
func ToFloat(price string) (float64, error) {
	p := strings.TrimSpace(price)
	re := regexp.MustCompile(`[^\d.]`)
	r := re.ReplaceAllString(p, "")
	return strconv.ParseFloat(r, 64)
}
