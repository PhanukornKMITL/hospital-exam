package repository

import (
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type HospitalRepository interface {
	FindAll() ([]entity.Hospital, error)
	FindByID(id uuid.UUID) (*entity.Hospital, error)
	Create(hospital *entity.Hospital) (*entity.Hospital, error)
	Update(hospital *entity.Hospital) (*entity.Hospital, error)
	Delete(id uuid.UUID) error
}

type hospitalRepository struct {
	db *gorm.DB
}

func NewHospitalRepository(db *gorm.DB) HospitalRepository {
	return &hospitalRepository{db: db}
}

func (r *hospitalRepository) FindAll() ([]entity.Hospital, error) {
	var hospitals []entity.Hospital
	err := r.db.Find(&hospitals).Error
	return hospitals, err
}

func (r *hospitalRepository) FindByID(id uuid.UUID) (*entity.Hospital, error) {
	var hospital entity.Hospital
	err := r.db.First(&hospital, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &hospital, nil
}

func (r *hospitalRepository) Create(hospital *entity.Hospital) (*entity.Hospital, error) {
	err := r.db.Create(hospital).Error
	return hospital, err
}

func (r *hospitalRepository) Update(hospital *entity.Hospital) (*entity.Hospital, error) {
	err := r.db.Save(hospital).Error
	return hospital, err
}

func (r *hospitalRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Hospital{}, "id = ?", id).Error
}
