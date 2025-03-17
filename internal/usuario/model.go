package usuario

import (
    "time"
)

// Usuario representa um usuário do sistema
type Usuario struct {
    ID           int       `json:"id"`
    Nome         string    `json:"nome"`
    Email        string    `json:"email"`
    Senha        string    `json:"-"` // Não serializar senha
    UltimoAcesso time.Time `json:"ultimo_acesso"`
    Status       bool      `json:"status"`
    DataCriacao  time.Time `json:"data_criacao"`
    Tipo         string    `json:"tipo"` // "admin", "cliente", etc.
}

// Repository define a interface para acesso aos dados de usuários
type Repository interface {
    Create(u *Usuario) error
    GetByID(id int) (*Usuario, error)
    GetByEmail(email string) (*Usuario, error)
    Update(u *Usuario) error
    Delete(id int) error
    List(limit, offset int) ([]*Usuario, error)
}

// Implementação do repositório será expandida mais tarde
