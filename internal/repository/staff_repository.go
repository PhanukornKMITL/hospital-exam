package repository

import (
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StaffRepository interface {
	FindAll() ([]entity.Staff, error)
	FindByUsername(username string) (*entity.Staff, error)
	FindByHospitalAndUsername(hospitalID uuid.UUID, username string) (*entity.Staff, error)
	Create(staff *entity.Staff) (*entity.Staff, error)
	DeleteByID(id uuid.UUID) error
	ExistsByUsernameInHospital(hospitalID uuid.UUID, username string) (bool, error)
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

func (r *staffRepository) FindByHospitalAndUsername(hospitalID uuid.UUID, username string) (*entity.Staff, error) {
	var staff entity.Staff
	if err := r.db.Where("hospital_id = ? AND username = ?", hospitalID, username).First(&staff).Error; err != nil {
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

func (r *staffRepository) ExistsByUsernameInHospital(hospitalID uuid.UUID, username string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Staff{}).Where("hospital_id = ? AND username = ?", hospitalID, username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
