package tokenservice_test

import (
	tokenservice "app_burse_backend/pkg/token_service"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_ValidUserID(t *testing.T) {
	secretKey := "testsecret"
	expSeconds := 3600
	userID := 123

	tokenService := tokenservice.NewTokenService(secretKey, expSeconds)
	token, err := tokenService.GenerateToken(userID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected a token, got an empty string")
	}

	parsedToken, ok := tokenService.VerifyToken(token)
	if !ok {
		t.Fatalf("expected token to be valid, got error %v", ok)
	}

	claims, parsed := parsedToken.Claims.(*jwt.MapClaims)
	if !parsed {
		t.Fatalf("expected claims to be of type jwt.MapClaims, got %T", parsedToken.Claims)
	}

	assert.Equal(t, (*claims)["id"], float64(userID))
}

func TestTokenFromRequest_MissingAuthorizationHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tokenService := tokenservice.NewTokenService("testsecret", 3600)
	token, err := tokenService.TokenFromRequest(req)

	t.Run("вернуть nil token", func(t *testing.T) {
		assert.Equal(t, "", token)
	})

	t.Run("вернуть ошибку при попытке расшифровать token", func(t *testing.T) {
		assert.Error(t, err)
	})
}

func TestTokenFromRequest_InvalidAuthorizationHeaderFormat(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	req.Header.Set("Authorization", "InvalidHeaderFormat")

	tokenService := tokenservice.NewTokenService("testsecret", 3600)
	token, err := tokenService.TokenFromRequest(req)

	t.Run("вернет пустой токен", func(t *testing.T) {
		assert.Equal(t, "", token)
	})

	t.Run("вернуть ошибку при попытке расшифровать token", func(t *testing.T) {
		assert.Error(t, err)
	})
}

func TestTokenFromRequest_ValidAuthorizationHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectToken := "testtoken"
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", expectToken))

	tokenService := tokenservice.NewTokenService("testsecret", 3600)
	token, err := tokenService.TokenFromRequest(req)

	t.Run("вернуть расшифрованный токен", func(t *testing.T) {
		assert.Equal(t, expectToken, token)
	})

	t.Run("вернуть nil ошибку при попытке расшифровать token", func(t *testing.T) {
		assert.Nil(t, err)
	})
}
