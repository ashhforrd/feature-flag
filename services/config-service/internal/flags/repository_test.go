package flags

import (
	"errors"
	"testing"
)

func TestRepositoryCreateAndGetByKey(t *testing.T) {
	repo := NewRepository()

	flag := Flag{
		Key:     "new-checkout",
		Name:    "New Checkout",
		Enabled: true,
	}

	if err := repo.Create(flag); err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	stored, err := repo.GetByKey("new-checkout")
	if err != nil {
		t.Fatalf("expected get by key to succeed, got %v", err)
	}

	if stored.Key != flag.Key {
		t.Fatalf("expected key %s, got %s", flag.Key, stored.Key)
	}
}

func TestRepositoryCreateDuplicateReturnsError(t *testing.T) {
	repo := NewRepository()

	flag := Flag{
		Key:  "new-checkout",
		Name: "New Checkout",
	}

	if err := repo.Create(flag); err != nil {
		t.Fatalf("expected first create to succeed, got %v", err)
	}

	err := repo.Create(flag)
	if !errors.Is(err, ErrFlagAlreadyExists) {
		t.Fatalf("expected ErrFlagAlreadyExists, got %v", err)
	}
}

func TestRepositoryGetMissingFlagReturnsError(t *testing.T) {
	repo := NewRepository()

	_, err := repo.GetByKey("missing-flag")
	if !errors.Is(err, ErrFlagNotFound) {
		t.Fatalf("expected ErrFlagNotFound, got %v", err)
	}
}
