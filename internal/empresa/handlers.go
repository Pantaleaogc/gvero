package empresa

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Pantaleaogc/gvero/internal/auth"
	"github.com/go-chi/chi/v5"
)

// Handlers contém os manipuladores HTTP para empresas
type Handlers struct {
	repo Repository
}

// NewHandlers cria uma nova instância de Handlers
func NewHandlers(repo Repository) *Handlers {
	return &Handlers{
		repo: repo,
	}
}

// Routes retorna as rotas para empresas
func Routes() http.Handler {
	repo := NewMemoryRepository()
	h := NewHandlers(repo)

	r := chi.NewRouter()
	// Middleware de autenticação
	r.Use(auth.Middleware)
	
	// Apenas administradores podem acessar
	r.Use(auth.RequireRole("admin"))

	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)

	return r
}

// List lista todas as empresas
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 100 // valor padrão
	}

	empresas, err := h.repo.List(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(empresas)
}

// GetByID retorna uma empresa por ID
func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	empresa, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(empresa)
}

// Create cria uma nova empresa
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var e Empresa
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

// Update atualiza uma empresa existente
func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var e Empresa
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e.ID = id
	if err := h.repo.Update(&e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

// Delete remove uma empresa
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
