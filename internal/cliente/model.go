package cliente

import (
	"time"
)

// Cliente representa um cliente no sistema
type Cliente struct {
	ID           int       `json:"id"`
	Nome         string    `json:"nome"`
	CNPJ         string    `json:"cnpj,omitempty"`
	CPF          string    `json:"cpf,omitempty"`
	Email        string    `json:"email"`
	Telefone     string    `json:"telefone"`
	Endereco     string    `json:"endereco"`
	Status       bool      `json:"status"`
	DataCriacao  time.Time `json:"data_criacao"`
	EmpresaID    int       `json:"empresa_id"`
	Observacoes  string    `json:"observacoes,omitempty"`
	UltimaCompra time.Time `json:"ultima_compra,omitempty"`
	TipoPessoa   string    `json:"tipo_pessoa"` // "fisica" ou "juridica"
}

// Repository define a interface para acesso aos dados de clientes
type Repository interface {
	Create(c *Cliente) error
	GetByID(id int, empresaID int) (*Cliente, error)
	Update(c *Cliente) error
	Delete(id int, empresaID int) error
	List(empresaID int, limit, offset int) ([]*Cliente, error)
	Search(empresaID int, query string, limit, offset int) ([]*Cliente, error)
}
