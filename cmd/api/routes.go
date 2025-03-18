package main

import (
	"net/http"

	_"github.com/Pantaleaogc/gvero/internal/auth"
	"github.com/Pantaleaogc/gvero/internal/cliente"
	"github.com/Pantaleaogc/gvero/internal/empresa"
	"github.com/Pantaleaogc/gvero/internal/usuario"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes configura todas as rotas da aplicação
func setupRoutes() http.Handler {
	r := chi.NewRouter()

	// Middleware básicos
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Rota raiz para verificação de saúde
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sistema CRM/ERP em Go - API funcionando! Versão 0.1.0"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Versão da API
		r.Route("/v1", func(r chi.Router) {
			// Rotas de usuários
			r.Mount("/usuarios", usuario.Routes())
			
			// Rotas de clientes
			r.Mount("/clientes", cliente.Routes())
			
			// Rotas de empresas
			r.Mount("/empresas", empresa.Routes())
		})
	})

	return r
}
