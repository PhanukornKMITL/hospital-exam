package repository

import (
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"

	"gorm.io/gorm"
)

type HospitalRepository interface {
	FindAll() ([]entity.Hospital, error)
	Create(hospital *entity.Hospital) (*entity.Hospital, error)
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

func (r *hospitalRepository) Create(hospital *entity.Hospital) (*entity.Hospital, error) {
	err := r.db.Create(hospital).Error
	return hospital, err
}
