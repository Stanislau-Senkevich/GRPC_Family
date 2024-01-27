package jwt

import (
	"context"
	"errors"
	"fmt"
	"strings"

	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

type Manager struct {
	signingKey []byte
}

// New creates and returns a new instance of the Manager with the provided
// signing key and tokenTTL.
func New(signingKey []byte) *Manager {
	return &Manager{
		signingKey: signingKey,
	}
}

// ParseToken parses the provided JWT token string and validates its signature
// using the configured signing key. It returns the claims embedded in the token
// if the signature is valid.
func (m *Manager) ParseToken(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tkn.Header["alg"]) //nolint
		}
		return m.signingKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, grpcerror.ErrNoToken
	}

	if err = verify(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func verify(claims jwt.MapClaims) error {
	if _, ok := claims["user_id"]; !ok {
		return fmt.Errorf("user_id was not found") //nolint
	}

	if _, ok := claims["role"]; !ok {
		return errors.New("role was not found") //nolint
	}

	if _, ok := claims["email"]; !ok {
		return errors.New("email was not found") //nolint
	}

	if _, ok := claims["user_id"].(float64); !ok {
		return errors.New("invalid user_id") //nolint
	}

	return nil
}

// GetClaims extracts and returns the JWT claims from the authorization token
// in the provided context. It relies on the ParseToken method to parse and
// validate the token's signature.
func (m *Manager) GetClaims(ctx context.Context) (jwt.MapClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpcerror.ErrTokenClaims
	}
	values := md["authorization"]
	if len(values) == 0 {
		return nil, grpcerror.ErrNoToken
	}

	accessToken := strings.Fields(values[0])[1]

	claims, err := m.ParseToken(accessToken)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (m *Manager) GetUserIDFromContext(ctx context.Context) int64 {
	claims, _ := m.GetClaims(ctx)

	id, _ := claims["user_id"].(float64)

	return int64(id)
}

func (m *Manager) IsAdmin(ctx context.Context) bool {
	claims, _ := m.GetClaims(ctx)

	return claims["role"] == "admin"
}
