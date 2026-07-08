package domain

import "github.com/google/uuid"

type Treatment struct {
	ID          uuid.UUID
	Name        string
	Description string
	Duration    int
	Price       float64
}

type TreatmentRepository interface {
	Create(treatment *Treatment) (*Treatment, error)
	GetByID(id uuid.UUID) (*Treatment, error)
	List(offset, limit int, search string) ([]Treatment, int64, error)
	Update(treatment *Treatment) error
	Delete(id uuid.UUID) error
}

type TreatmentService interface {
	Create(treatment *Treatment) (*Treatment, error)
	GetByID(id uuid.UUID) (*Treatment, error)
	List(offset, limit int, search string) ([]Treatment, int64, error)
	Update(treatment *Treatment) error
	Delete(id uuid.UUID) error
}
