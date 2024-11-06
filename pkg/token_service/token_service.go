package tokenservice

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secretKey  string
	expSeconds int
}

func NewTokenService(secretKey string, expSeconds int) *TokenService {
	return &TokenService{secretKey: secretKey, expSeconds: expSeconds}
}

func (s *TokenService) GenerateToken(userID int) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().UTC().Add(time.Duration(s.expSeconds) * time.Second).Unix(),
	})

	tokenString, err := jwtToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *TokenService) VerifyToken(tokenString string) (*jwt.Token, bool) {
	token, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, false
	}

	if !token.Valid {
		return nil, false
	}

	return token, true
}

func (s *TokenService) ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.secretKey), nil
	})

	return token, err
}

func (s *TokenService) TokenFromRequest(r *http.Request, headerField string, cookieField string) (string, error) {
	var field string

	if headerField != "" {
		field = r.Header.Get(headerField)
	} else {
		cookie, err := r.Cookie(cookieField)
		if err != nil {
			return "", err
		}

		field = cookie.Value
	}

	fmt.Println("field", field)
	if field == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.SplitN(field, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format is invalid")
	}

	return parts[1], nil
}
