package helpers

import (
	"fmt"
	"strings"
)

func ExtractJSON(input string) (string, error) {
	startMarker := "```json"
	endMarker := "```"

	startIdx := strings.Index(input, startMarker)
	if startIdx == -1 {
		return input, nil
	}

	endIdx := strings.Index(input[startIdx+len(startMarker):], endMarker)
	if endIdx == -1 {
		return "", fmt.Errorf("не найден закрывающий маркер ```")
	}

	jsonStr := input[startIdx+len(startMarker) : startIdx+len(startMarker)+endIdx]
	return strings.TrimSpace(jsonStr), nil
}
