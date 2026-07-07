package flags

import "errors"

var ErrFlagNotFound = errors.New("flag not found")
var ErrFlagAlreadyExists = errors.New("flag already exists")

type Repository interface {
	Create(flag Flag) error
	List() ([]Flag, error)
	GetByKey(key string) (Flag, error)
	Update(flag Flag) error
}

type MemoryRepository struct {
	flags map[string]Flag
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		flags: map[string]Flag{},
	}
}

func (r *MemoryRepository) Create(flag Flag) error {
	if _, exists := r.flags[flag.Key]; exists {
		return ErrFlagAlreadyExists
	}

	r.flags[flag.Key] = flag
	return nil
}

func (r *MemoryRepository) List() ([]Flag, error) {
	result := make([]Flag, 0, len(r.flags))

	for _, flag := range r.flags {
		result = append(result, flag)
	}

	return result, nil
}

func (r *MemoryRepository) GetByKey(key string) (Flag, error) {
	flag, exists := r.flags[key]
	if !exists {
		return Flag{}, ErrFlagNotFound
	}

	return flag, nil
}

func (r *MemoryRepository) Update(flag Flag) error {
	if _, exists := r.flags[flag.Key]; !exists {
		return ErrFlagNotFound
	}

	r.flags[flag.Key] = flag
	return nil
}
