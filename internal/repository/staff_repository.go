package repository

import (
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StaffRepository interface {
	FindAll() ([]entity.Staff, error)
	FindByUsername(username string) (*entity.Staff, error)
	Create(staff *entity.Staff) (*entity.Staff, error)
	DeleteByID(id uuid.UUID) error
}

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) FindAll() ([]entity.Staff, error) {
	var staffs []entity.Staff
	err := r.db.Find(&staffs).Error
	return staffs, err
}

func (r *staffRepository) FindByUsername(username string) (*entity.Staff, error) {
	var staff entity.Staff
	if err := r.db.Where("username = ?", username).First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepository) Create(staff *entity.Staff) (*entity.Staff, error) {
	err := r.db.Create(staff).Error
	return staff, err
}

func (r *staffRepository) DeleteByID(id uuid.UUID) error {
	return r.db.Delete(&entity.Staff{}, "id = ?", id).Error
}
