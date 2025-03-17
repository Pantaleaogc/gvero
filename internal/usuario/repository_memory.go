package usuario

import (
	"errors"
	"sync"
	"time"
)

// MemoryRepository implementa Repository para testes
type MemoryRepository struct {
	mu       sync.RWMutex
	usuarios map[int]*Usuario
	nextID   int
}

// NewMemoryRepository cria um novo repositório em memória
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		usuarios: make(map[int]*Usuario),
		nextID:   1,
	}
}

// Create adiciona um novo usuário
func (r *MemoryRepository) Create(u *Usuario) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validação básica
	if u.Nome == "" || u.Email == "" {
		return errors.New("nome e email são obrigatórios")
	}

	// Verificar se email já existe
	for _, existingUser := range r.usuarios {
		if existingUser.Email == u.Email {
			return errors.New("email já está em uso")
		}
	}

	// Configurar campos
	u.ID = r.nextID
	r.nextID++
	u.DataCriacao = time.Now()
	if u.Status == false {
		u.Status = true // padrão ativo
	}

	// Simular senha (na implementação real deve ser hasheada)
	if u.Senha == "" {
		u.Senha = "senha_padrao" // Apenas para teste
	}

	// Adicionar no mapa
	r.usuarios[u.ID] = u
	return nil
}

// GetByID busca um usuário por ID
func (r *MemoryRepository) GetByID(id int) (*Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.usuarios[id]
	if !exists {
		return nil, errors.New("usuário não encontrado")
	}
	return u, nil
}

// GetByEmail busca um usuário por email
func (r *MemoryRepository) GetByEmail(email string) (*Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.usuarios {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("usuário não encontrado")
}

// Update atualiza um usuário existente
func (r *MemoryRepository) Update(u *Usuario) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.usuarios[u.ID]; !exists {
		return errors.New("usuário não encontrado")
	}

	// Validação básica
	if u.Nome == "" || u.Email == "" {
		return errors.New("nome e email são obrigatórios")
	}

	// Verificar email duplicado
	for id, existingUser := range r.usuarios {
		if existingUser.Email == u.Email && id != u.ID {
			return errors.New("email já está em uso")
		}
	}

	// Preservar campos que não devem ser alterados
	u.DataCriacao = r.usuarios[u.ID].DataCriacao

	// Atualizar
	r.usuarios[u.ID] = u
	return nil
}

// Delete remove um usuário
func (r *MemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.usuarios[id]; !exists {
		return errors.New("usuário não encontrado")
	}

	delete(r.usuarios, id)
	return nil
}

// List retorna uma lista de usuários
func (r *MemoryRepository) List(limit, offset int) ([]*Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit <= 0 {
		limit = len(r.usuarios)
	}

	result := make([]*Usuario, 0, limit)
	count := 0
	skipCount := 0

	for _, u := range r.usuarios {
		if skipCount < offset {
			skipCount++
			continue
		}

		if count < limit {
			result = append(result, u)
			count++
		} else {
			break
		}
	}

	return result, nil
}
