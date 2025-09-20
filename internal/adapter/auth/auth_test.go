package auth_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"cleo.com/internal/adapter/auth"
	"cleo.com/testsupport"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateToken(t *testing.T, secret, role string) string {
	t.Helper()

	claims := auth.CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	var token *jwt.Token
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

func TestService_HasRole(t *testing.T) {
	secret := "a-test-string-secret-at-least-256-bits-long"
	svc := auth.NewService(testsupport.Logger(), auth.Config{APISecret: secret})

	tests := []struct {
		desc               string
		bearer             string
		expectedToHaveRole bool
		expectErr          bool
	}{
		{
			desc:               "valid token plus correct role",
			bearer:             "Bearer " + generateToken(t, secret, "CLINICAL-EDITOR"),
			expectedToHaveRole: true,
		},
		{
			desc:               "valid token but without correct role",
			bearer:             "Bearer " + generateToken(t, secret, "NON_CLINICAL_EDITOR"),
			expectedToHaveRole: false,
		},
		{
			desc:               "missing header",
			bearer:             "",
			expectedToHaveRole: false,
			expectErr:          true,
		},
		{
			desc:               "invalid bearer format",
			bearer:             "Token abc123",
			expectedToHaveRole: false,
			expectErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			c, _ := testsupport.NewTestContext(nil)
			if tt.bearer != "" {
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("Authorization", tt.bearer)
			} else {
				c.Request = httptest.NewRequest("GET", "/", nil)
			}

			ok, err := svc.HasRole(c, "CLINICAL-EDITOR")

			assert.Equal(t, tt.expectedToHaveRole, ok, "expected role should match expectations")
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
