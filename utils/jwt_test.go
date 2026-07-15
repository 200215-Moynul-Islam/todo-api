package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("expected token")
	}
}

func TestValidateToken(t *testing.T) {
	validToken, err := GenerateToken(123)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	expiredClaims := Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).
		SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	wrongSecret := []byte("wrong-secret")

	wrongSignatureToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString(wrongSecret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	tests := []struct {
		name       string
		token      string
		wantUserID int
		wantErr    error
	}{
		{
			name:       "valid token",
			token:      validToken,
			wantUserID: 123,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "invalid token",
			token:   "invalid-token",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "expired token",
			token:   expiredToken,
			wantErr: ErrInvalidToken,
		},
		{
			name:    "wrong signature",
			token:   wrongSignatureToken,
			wantErr: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := ValidateToken(tt.token)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if userID != tt.wantUserID {
				t.Fatalf("expected userID %d, got %d", tt.wantUserID, userID)
			}
		})
	}
}
