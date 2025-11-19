package token

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type UserAuthToken struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// ValidateRefreshToken parses and validates a JWT refresh token string,
// returning the user ID if valid, or an error if invalid.
func ValidateRefreshToken(token string) (int, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtConfig.publicKey, nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return 0, errors.New("invalid token claims or token not valid")
	}

	dataRaw, ok := claims["data"]
	if !ok {
		return 0, errors.New(`missing "data" field in token claims`)
	}

	data, ok := dataRaw.(map[string]any)
	if !ok {
		return 0, errors.New(`invalid "data" field format in token claims`)
	}

	idVal, ok := data["id"]
	if !ok {
		return 0, errors.New(`missing "id" field in token data`)
	}

	idFloat, ok := idVal.(float64)
	if !ok {
		return 0, fmt.Errorf(`"id" field is not a number: %v`, idVal)
	}
	if idFloat < float64(int(^uint(0)>>1)*-1-1) || idFloat > float64(int(^uint(0)>>1)) {
		return 0, fmt.Errorf(`"id" field value out of int range: %v`, idFloat)
	}
	if idFloat != float64(int(idFloat)) {
		return 0, fmt.Errorf(`"id" field is not an integer: %v`, idFloat)
	}

	return int(idFloat), nil
}

func ValidateAccessToken(token string) (*UserAuthToken, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtConfig.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token claims or token not valid")
	}

	data, ok := claims["data"]
	if !ok {
		return nil, errors.New(`missing "data" field in token claims`)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal 'data' field: %w", err)
	}

	var user UserAuthToken
	if err := json.Unmarshal(jsonData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data into UserAuthToken: %w", err)
	}

	return &user, nil
}
