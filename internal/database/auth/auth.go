package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("missing API key")
	}
	parts := strings.Split(apiKey, " ")
	if len(parts) != 2 {
		return "", errors.New("invalid API key format")
	}
	if parts[0] != "ApiKey" {
		return "", errors.New("invalid API key format")
	}
	return parts[1], nil
}
