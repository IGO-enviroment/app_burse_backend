package service

import (
	"errors"
	"strings"
)

func GetTokenFromHeader(header string) (string, error) {
	authHeader := strings.Split(header, " ")
	if len(authHeader) != 2 || authHeader[0] != "Bearer" {
		return "", errors.New("Invalid authorization header")
	}

	return authHeader[1], nil
}
