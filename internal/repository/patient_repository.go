package repository

import (
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PatientRepository interface {
	FindAll() ([]entity.Patient, error)
	Create(patient *entity.Patient) (*entity.Patient, error)
	FindByID(id uuid.UUID) (*entity.Patient, error)
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) FindAll() ([]entity.Patient, error) {
	var patients []entity.Patient
	err := r.db.Find(&patients).Error
	return patients, err
}

func (r *patientRepository) Create(patient *entity.Patient) (*entity.Patient, error) {
	err := r.db.Create(patient).Error
	return patient, err
}

func (r *patientRepository) FindByID(id uuid.UUID) (*entity.Patient, error) {
	var patient entity.Patient
	if err := r.db.First(&patient, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}
