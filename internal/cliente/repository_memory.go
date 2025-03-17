package cliente

import (
	"errors"
	"strings"
	"sync"
	"time"
)

// MemoryRepository implementa Repository em memória
type MemoryRepository struct {
	mu       sync.RWMutex
	clientes map[int]*Cliente
	nextID   int
}

// NewMemoryRepository cria um novo repositório em memória
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		clientes: make(map[int]*Cliente),
		nextID:   1,
	}
}

// Create adiciona um novo cliente
func (r *MemoryRepository) Create(c *Cliente) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validação básica
	if c.Nome == "" || c.Email == "" {
		return errors.New("nome e email são obrigatórios")
	}

	if c.TipoPessoa == "fisica" && c.CPF == "" {
		return errors.New("CPF é obrigatório para pessoa física")
	}

	if c.TipoPessoa == "juridica" && c.CNPJ == "" {
		return errors.New("CNPJ é obrigatório para pessoa jurídica")
	}

	if c.EmpresaID <= 0 {
		return errors.New("empresa inválida")
	}

	// Configurar campos
	c.ID = r.nextID
	r.nextID++
	c.DataCriacao = time.Now()
	if c.Status == false {
		c.Status = true // padrão ativo
	}

	// Adicionar ao mapa
	r.clientes[c.ID] = c
	return nil
}

// GetByID busca um cliente por ID e empresa
func (r *MemoryRepository) GetByID(id int, empresaID int) (*Cliente, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, exists := r.clientes[id]
	if !exists || c.EmpresaID != empresaID {
		return nil, errors.New("cliente não encontrado")
	}
	return c, nil
}

// Update atualiza um cliente existente
func (r *MemoryRepository) Update(c *Cliente) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.clientes[c.ID]
	if !exists {
		return errors.New("cliente não encontrado")
	}

	// Verificar se pertence à mesma empresa
	if existing.EmpresaID != c.EmpresaID {
		return errors.New("operação não permitida: cliente pertence a outra empresa")
	}

	// Validação básica
	if c.Nome == "" || c.Email == "" {
		return errors.New("nome e email são obrigatórios")
	}

	// Preservar campos que não devem ser alterados
	c.DataCriacao = existing.DataCriacao

	// Atualizar
	r.clientes[c.ID] = c
	return nil
}

// Delete remove um cliente
func (r *MemoryRepository) Delete(id int, empresaID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	c, exists := r.clientes[id]
	if !exists {
		return errors.New("cliente não encontrado")
	}

	// Verificar se pertence à mesma empresa
	if c.EmpresaID != empresaID {
		return errors.New("operação não permitida: cliente pertence a outra empresa")
	}

	delete(r.clientes, id)
	return nil
}

// List retorna uma lista de clientes de uma empresa
func (r *MemoryRepository) List(empresaID int, limit, offset int) ([]*Cliente, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit <= 0 {
		limit = len(r.clientes)
	}

	result := make([]*Cliente, 0, limit)
	count := 0
	skipCount := 0

	for _, c := range r.clientes {
		if c.EmpresaID != empresaID {
			continue
		}

		if skipCount < offset {
			skipCount++
			continue
		}

		if count < limit {
			result = append(result, c)
			count++
		} else {
			break
		}
	}

	return result, nil
}

// Search procura clientes por termos
func (r *MemoryRepository) Search(empresaID int, query string, limit, offset int) ([]*Cliente, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit <= 0 {
		limit = len(r.clientes)
	}

	result := make([]*Cliente, 0, limit)
	count := 0
	skipCount := 0
	query = strings.ToLower(query)

	for _, c := range r.clientes {
		if c.EmpresaID != empresaID {
			continue
		}

		// Buscar por nome, email, CPF ou CNPJ
		match := strings.Contains(strings.ToLower(c.Nome), query) ||
			strings.Contains(strings.ToLower(c.Email), query) ||
			strings.Contains(c.CPF, query) ||
			strings.Contains(c.CNPJ, query)

		if !match {
			continue
		}

		if skipCount < offset {
			skipCount++
			continue
		}

		if count < limit {
			result = append(result, c)
			count++
		} else {
			break
		}
	}

	return result, nil
}
