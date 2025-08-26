package auth

import "github.com/golang-jwt/jwt/v5"

type Authorizer interface {
	AuthenticateUser(name, password string) (string, error)
	RegisterUser(name, password string) (string, error)
	AuthorizeUser(tokenString string) (*jwt.Token, error)
}
