package internalAuth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zenrot/CryptoService/internal/storage"
	"time"
)

type internalCustomClaims struct {
	Username string `json:"name"`
	jwt.RegisteredClaims
}
type internalAuthorizer struct {
	Store  storage.Auth
	JwtKey string
}

func New(store storage.Auth, jwtKey string) *internalAuthorizer {
	return &internalAuthorizer{
		Store:  store,
		JwtKey: jwtKey,
	}
}

func (au *internalAuthorizer) AuthenticateUser(name, password string) (string, error) {
	_, err := au.Store.LoginUser(name, password)
	if err != nil {
		return "", err
	}

	claims := internalCustomClaims{
		Username: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-auth",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "auth_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(au.JwtKey))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (au *internalAuthorizer) RegisterUser(name, password string) (string, error) {
	if err := au.Store.RegisterUser(name, password); err != nil {
		return "", err
	}
	return au.AuthenticateUser(name, password)
}

func (au *internalAuthorizer) AuthorizeUser(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(au.JwtKey), nil
	})

	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return fmt.Errorf("invalid token claims")
}
