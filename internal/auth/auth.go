package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Pantaleaogc/gvero/internal/usuario"
	"github.com/Pantaleaogc/gvero/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtKey = []byte("sua_chave_secreta_deve_ser_substituida_em_producao") // Substitua em produção!
)

// LoginRequest representa os dados enviados no login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
	Token string         `json:"token"`
	User  UserResponse   `json:"user"`
}

// UserResponse representa os dados do usuário na resposta
type UserResponse struct {
	ID        int    `json:"id"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Tipo      string `json:"tipo"`
	EmpresaID int    `json:"empresa_id,omitempty"`
}

// Handlers contém os handlers HTTP para autenticação
type Handlers struct {
	userRepo usuario.Repository
}

// NewHandlers cria uma nova instância de Handlers
func NewHandlers(userRepo usuario.Repository) *Handlers {
	return &Handlers{
		userRepo: userRepo,
	}
}

// Routes retorna as rotas para autenticação
func Routes() http.Handler {
	userRepo := usuario.NewMemoryRepository()
	h := NewHandlers(userRepo)

	r := chi.NewRouter()
	r.Post("/login", h.Login)
	r.With(Middleware).Get("/verify", h.Verify)
	r.With(Middleware).Post("/logout", h.Logout)

	return r
}

// Login processa a requisição de login
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Printf("Erro ao decodificar requisição: %v", err)
		http.Error(w, "Formato inválido", http.StatusBadRequest)
		return
	}

	// Buscar usuário por email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		logger.InfoLogger.Printf("Tentativa de login com email não encontrado: %s", req.Email)
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	// Verificar senha (em produção, você deve usar bcrypt ou similar)
	// Nota: Este é um exemplo simples, em produção use hashing adequado
	if user.Senha != req.Password {
		logger.InfoLogger.Printf("Tentativa de login com senha incorreta para: %s", req.Email)
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	// Gerar token JWT
	token, err := generateJWT(user)
	if err != nil {
		logger.ErrorLogger.Printf("Erro ao gerar token: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// Atualizar último acesso
	user.UltimoAcesso = time.Now()
	h.userRepo.Update(user)

	// Responder com token e dados do usuário
	resp := LoginResponse{
		Token: token,
		User: UserResponse{
			ID:        user.ID,
			Nome:      user.Nome,
			Email:     user.Email,
			Tipo:      user.Tipo,
			EmpresaID: 1, // Você deve obter isso do usuário em uma implementação real
		},
	}

	logger.InfoLogger.Printf("Login bem-sucedido: %s (ID: %d)", user.Email, user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Verify verifica se o token do usuário é válido
func (h *Handlers) Verify(w http.ResponseWriter, r *http.Request) {
	// O middleware já verifica o token, então se chegamos aqui, o token é válido
	user, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Não autorizado", http.StatusUnauthorized)
		return
	}

	// Buscar dados atualizados do usuário
	userObj, err := h.userRepo.GetByID(user.ID)
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}

	resp := UserResponse{
		ID:        userObj.ID,
		Nome:      userObj.Nome,
		Email:     userObj.Email,
		Tipo:      userObj.Tipo,
		EmpresaID: 1, // Você deve obter isso do usuário em uma implementação real
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Logout faz o logout do usuário (opcional no backend)
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	// Em uma implementação real, você pode querer invalidar o token
	// Para JWT, isso geralmente é feito com uma lista negra ou usando tokens de vida curta
	// Para simplificar, apenas retornamos sucesso

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// generateJWT gera um token JWT para o usuário
func generateJWT(user *usuario.Usuario) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"tipo":  user.Tipo,
		"exp":   expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}