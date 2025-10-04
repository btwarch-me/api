package utils

import (
	"encoding/json"
	"strings"
)

type CloudflareError struct {
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func ExtractErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	// Try to find JSON part inside the error string
	if idx := strings.Index(err.Error(), "{"); idx != -1 {
		jsonPart := err.Error()[idx:]
		var cfErr CloudflareError
		if json.Unmarshal([]byte(jsonPart), &cfErr) == nil {
			if len(cfErr.Errors) > 0 {
				return cfErr.Errors[0].Message
			}
		}
	}

	return err.Error()
}
