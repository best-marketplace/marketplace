package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidToken = errors.New("недействительный токен авторизации")
)

type TokenInfo struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func extractTokenInfo(token string) (*TokenInfo, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования токена: %w", err)
	}

	var claims struct {
		Sub  string `json:"sub"`
		Role string `json:"role"`
		Exp  int64  `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("ошибка парсинга данных токена: %w", err)
	}

	return &TokenInfo{
		UserID:   claims.Sub,
		Username: "user",
		Role:     claims.Role,
	}, nil
}

func VerifyToken(authServiceURL string, token string) (*TokenInfo, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	token = strings.TrimPrefix(token, "Bearer ")

	tokenInfo, err := extractTokenInfo(token)
	if err == nil && tokenInfo != nil {
		return tokenInfo, nil
	}

	return nil, fmt.Errorf("не удалось извлечь информацию из токена: %w", err)
}
