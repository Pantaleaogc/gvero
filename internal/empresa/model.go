package empresa

import (
	"time"
)

// Empresa representa um perfil empresarial no sistema
type Empresa struct {
	ID           int       `json:"id"`
	Nome         string    `json:"nome"`
	CNPJ         string    `json:"cnpj"`
	Endereco     string    `json:"endereco"`
	Telefone     string    `json:"telefone"`
	Email        string    `json:"email"`
	Status       bool      `json:"status"`
	DataCriacao  time.Time `json:"data_criacao"`
	Plano        string    `json:"plano"` // "básico", "premium", etc.
	DataExpira   time.Time `json:"data_expira"`
	LogoURL      string    `json:"logo_url,omitempty"`
	Modulos      []string  `json:"modulos"`
	MaxUsuarios  int       `json:"max_usuarios"`
}

// Repository define a interface para acesso aos dados de empresas
type Repository interface {
	Create(e *Empresa) error
	GetByID(id int) (*Empresa, error)
	Update(e *Empresa) error
	Delete(id int) error
	List(limit, offset int) ([]*Empresa, error)
}
