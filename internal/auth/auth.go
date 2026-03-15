package auth

type Authorizer interface {
	AuthenticateUser(name, password string) (string, error)
	RegisterUser(name, password string) (string, error)
	AuthorizeUser(tokenString string) error
}
