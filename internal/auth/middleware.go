// Substitua o conteúdo do arquivo internal/auth/middleware.go

package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	
	"github.com/Pantaleaogc/gvero/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

// Chave para contexto do usuário
type contextKey string

const UserContextKey = contextKey("user")

// Middleware verifica o token JWT e injeta o usuário no contexto
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obter token do cabeçalho Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.DebugLogger.Printf("Requisição sem token de autorização: %s", r.URL.Path)
			http.Error(w, "Autorização requerida", http.StatusUnauthorized)
			return
		}

		// Verificar formato do token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.DebugLogger.Printf("Formato de token inválido: %s", authHeader)
			http.Error(w, "Formato de autorização inválido", http.StatusUnauthorized)
			return
		}

		// Verificar e validar o token JWT
		tokenStr := parts[1]
		user, err := validateJWT(tokenStr)
		if err != nil {
			logger.DebugLogger.Printf("Token inválido: %v", err)
			http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
			return
		}

		logger.DebugLogger.Printf("Usuário autenticado: %s (ID: %d)", user.Email, user.ID)

		// Injetar usuário no contexto
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateJWT valida o token JWT e retorna os dados do usuário
func validateJWT(tokenStr string) (User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Verificar o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return User{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extrair dados do token
		var user User
		
		if id, ok := claims["id"].(float64); ok {
			user.ID = int(id)
		} else {
			return User{}, errors.New("ID inválido no token")
		}
		
		if email, ok := claims["email"].(string); ok {
			user.Email = email
		} else {
			return User{}, errors.New("email inválido no token")
		}
		
		if role, ok := claims["tipo"].(string); ok {
			user.Role = role
		} else {
			return User{}, errors.New("tipo inválido no token")
		}
		
		// Empresa é opcional
		if empresa, ok := claims["empresa"].(float64); ok {
			user.Empresa = int(empresa)
		} else {
			user.Empresa = 1 // Valor padrão
		}
		
		return user, nil
	}

	return User{}, errors.New("token inválido")
}

// FromContext extrai o usuário do contexto
func FromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(UserContextKey).(User)
	return user, ok
}

// RequireRole é um middleware para verificar a função do usuário
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := FromContext(r.Context())
			if !ok {
				logger.DebugLogger.Printf("Acesso negado: usuário não autenticado em rota restrita: %s", r.URL.Path)
				http.Error(w, "Não autorizado", http.StatusUnauthorized)
				return
			}

			if user.Role != role && role != "any" {
				logger.DebugLogger.Printf("Acesso negado: usuário %s com role %s tentou acessar recurso que requer role %s", 
					user.Email, user.Role, role)
				http.Error(w, "Acesso negado", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}