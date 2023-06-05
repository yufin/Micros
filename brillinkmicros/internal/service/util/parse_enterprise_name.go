package util

import (
	"encoding/json"
	"fmt"
)

func ParseEnterpriseName(content string) (string, error) {
	var contentMap map[string]any
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		return "", err
	}
	s := contentMap["impExpEntReport"].(map[string]any)["businessInfo"].(map[string]any)["QYMC"]
	// determine enterprise name is nil or string
	if _, ok := s.(string); ok {
		return s.(string), nil
	}
	return "", fmt.Errorf("enterprise name is not string")
}
