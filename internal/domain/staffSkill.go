package domain

import "github.com/google/uuid"

type StaffSkill struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	User        *User
	TreatmentID uuid.UUID
	Treatment   *Treatment
}

type StaffSkillRepository interface {
	Assign(userID, treatmentID uuid.UUID) (*StaffSkill, error)
	Unassign(id uuid.UUID) error
	GetByID(id uuid.UUID) (*StaffSkill, error)
	ListByStaff(userID uuid.UUID) ([]StaffSkill, error)
	ListByTreatment(treatmentID uuid.UUID) ([]StaffSkill, error)
	ListAll(offset, limit int, search string) ([]StaffSkill, int64, error)
}

type StaffSkillService interface {
	Assign(userID, treatmentID uuid.UUID) (*StaffSkill, error)
	Unassign(id uuid.UUID) error
	GetByID(id uuid.UUID) (*StaffSkill, error)
	ListByStaff(userID uuid.UUID) ([]StaffSkill, error)
	ListByTreatment(treatmentID uuid.UUID) ([]StaffSkill, error)
	ListAll(offset, limit int, search string) ([]StaffSkill, int64, error)
}
