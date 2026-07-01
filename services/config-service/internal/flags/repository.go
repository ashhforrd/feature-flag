package flags

import "errors"

var ErrFlagNotFound = errors.New("flag not found")
var ErrFlagAlreadyExists = errors.New("flag already exists")

type Repository struct {
	flags map[string]Flag
}

func NewRepository() *Repository {
	return &Repository{
		flags: map[string]Flag{},
	}
}

func (r *Repository) Create(flag Flag) error {
	if _, exists := r.flags[flag.Key]; exists {
		return ErrFlagAlreadyExists
	}

	r.flags[flag.Key] = flag
	return nil
}

func (r *Repository) List() []Flag {
	result := make([]Flag, 0, len(r.flags))

	for _, flag := range r.flags {
		result = append(result, flag)
	}

	return result
}

func (r *Repository) GetByKey(key string) (Flag, error) {
	flag, exists := r.flags[key]
	if !exists {
		return Flag{}, ErrFlagNotFound
	}

	return flag, nil
}

func (r *Repository) Update(flag Flag) error {
	if _, exists := r.flags[flag.Key]; !exists {
		return ErrFlagNotFound
	}

	r.flags[flag.Key] = flag
	return nil
}
