package empresa

import (
	"errors"
	"sync"
	"time"
)

// MemoryRepository implementa Repository em memória
type MemoryRepository struct {
	mu       sync.RWMutex
	empresas map[int]*Empresa
	nextID   int
}

// NewMemoryRepository cria um novo repositório em memória
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		empresas: make(map[int]*Empresa),
		nextID:   1,
	}
}

// Create adiciona uma nova empresa
func (r *MemoryRepository) Create(e *Empresa) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validação básica
	if e.Nome == "" || e.CNPJ == "" || e.Email == "" {
		return errors.New("nome, CNPJ e email são obrigatórios")
	}

	// Verificar se CNPJ já existe
	for _, existing := range r.empresas {
		if existing.CNPJ == e.CNPJ {
			return errors.New("CNPJ já cadastrado")
		}
	}

	// Configurar campos
	e.ID = r.nextID
	r.nextID++
	e.DataCriacao = time.Now()
	if e.Status == false {
		e.Status = true // padrão ativo
	}

	// Configurações padrão
	if e.Plano == "" {
		e.Plano = "básico"
	}
	if e.MaxUsuarios <= 0 {
		e.MaxUsuarios = 5 // padrão para plano básico
	}
	if len(e.Modulos) == 0 {
		e.Modulos = []string{"clientes", "financeiro"} // módulos padrão
	}
	
	// Data de expiração padrão: 1 ano
	if e.DataExpira.IsZero() {
		e.DataExpira = time.Now().AddDate(1, 0, 0)
	}

	// Adicionar ao mapa
	r.empresas[e.ID] = e
	return nil
}

// GetByID busca uma empresa por ID
func (r *MemoryRepository) GetByID(id int) (*Empresa, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	e, exists := r.empresas[id]
	if !exists {
		return nil, errors.New("empresa não encontrada")
	}
	return e, nil
}

// Update atualiza uma empresa existente
func (r *MemoryRepository) Update(e *Empresa) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.empresas[e.ID]
	if !exists {
		return errors.New("empresa não encontrada")
	}

	// Validação básica
	if e.Nome == "" || e.CNPJ == "" || e.Email == "" {
		return errors.New("nome, CNPJ e email são obrigatórios")
	}

	// Verificar se CNPJ já existe em outra empresa
	for id, existingEmpresa := range r.empresas {
		if existingEmpresa.CNPJ == e.CNPJ && id != e.ID {
			return errors.New("CNPJ já cadastrado em outra empresa")
		}
	}

	// Preservar campos que não devem ser alterados
	e.DataCriacao = existing.DataCriacao

	// Atualizar
	r.empresas[e.ID] = e
	return nil
}

// Delete remove uma empresa
func (r *MemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.empresas[id]; !exists {
		return errors.New("empresa não encontrada")
	}

	delete(r.empresas, id)
	return nil
}

// List retorna uma lista de empresas
func (r *MemoryRepository) List(limit, offset int) ([]*Empresa, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit <= 0 {
		limit = len(r.empresas)
	}

	result := make([]*Empresa, 0, limit)
	count := 0
	skipCount := 0

	for _, e := range r.empresas {
		if skipCount < offset {
			skipCount++
			continue
		}

		if count < limit {
			result = append(result, e)
			count++
		} else {
			break
		}
	}

	return result, nil
}
