package auth

import (
    "errors"
    "time"
)

// Errors
var (
    ErrInvalidCredentials = errors.New("credenciais inválidas")
    ErrUserNotFound = errors.New("usuário não encontrado")
    ErrUnauthorized = errors.New("não autorizado")
)

// User representa um usuário autenticado
type User struct {
    ID       int
    Email    string
    Role     string
    Empresa  int
}

// Authenticator define a interface para o sistema de autenticação
type Authenticator interface {
    Login(email, password string) (User, error)
    Verify(token string) (User, error)
    Refresh(token string) (string, error)
}

// Token representa um token JWT
type Token struct {
    Value     string
    ExpiresAt time.Time
    UserID    int
}

// Implementação será expandida conforme avançarmos no projeto
