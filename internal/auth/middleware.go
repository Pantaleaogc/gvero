package auth

import (
	"context"
	"net/http"
	"strings"
	
	"github.com/Pantaleaogc/pantaleaocrmerp/pkg/logger"
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

		// TODO: Implementar verificação real do token JWT
		// Por enquanto, apenas simulamos um usuário autenticado para desenvolvimento
		user := User{
			ID:      1,
			Email:   "admin@example.com",
			Role:    "admin",
			Empresa: 1,
		}

		logger.DebugLogger.Printf("Usuário autenticado: %s (ID: %d, Empresa: %d)", user.Email, user.ID, user.Empresa)

		// Injetar usuário no contexto
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
