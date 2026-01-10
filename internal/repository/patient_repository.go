package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PatientRepository interface {
	FindAll() ([]entity.Patient, error)
	FindByHospitalID(hospitalID uuid.UUID) ([]entity.Patient, error)
	Create(patient *entity.Patient) (*entity.Patient, error)
	CreateWithGeneratedHN(patient *entity.Patient) (*entity.Patient, error)
	FindByID(id uuid.UUID) (*entity.Patient, error)
	FindByHospitalAndIdentifier(hospitalID uuid.UUID, identifier string) (*entity.Patient, error)
	FindByHospitalAndID(hospitalID uuid.UUID, id uuid.UUID) (*entity.Patient, error)
	ExistsByNationalIDInHospital(hospitalID uuid.UUID, nationalID string) (bool, error)
	ExistsByPassportIDInHospital(hospitalID uuid.UUID, passportID string) (bool, error)
	ExistsByNationalIDInHospitalExcept(hospitalID uuid.UUID, nationalID string, patientID uuid.UUID) (bool, error)
	ExistsByPassportIDInHospitalExcept(hospitalID uuid.UUID, passportID string, patientID uuid.UUID) (bool, error)
	FindByHospitalWithFilters(hospitalID uuid.UUID, filters PatientSearchFilters, page, limit int) ([]entity.Patient, int64, error)
	Update(patient *entity.Patient) (*entity.Patient, error)
	Delete(hospitalID uuid.UUID, id uuid.UUID) error
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

func (r *patientRepository) FindByHospitalID(hospitalID uuid.UUID) ([]entity.Patient, error) {
	var patients []entity.Patient
	err := r.db.Where("hospital_id = ?", hospitalID).Find(&patients).Error
	return patients, err
}

func (r *patientRepository) Create(patient *entity.Patient) (*entity.Patient, error) {
	err := r.db.Create(patient).Error
	return patient, err
}

func (r *patientRepository) CreateWithGeneratedHN(patient *entity.Patient) (*entity.Patient, error) {
	return patient, r.db.Transaction(func(tx *gorm.DB) error {
		var seq entity.PatientSequence
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&seq, "hospital_id = ?", patient.HospitalID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				seq = entity.PatientSequence{HospitalID: patient.HospitalID, NextNumber: 1}
				if err := tx.Create(&seq).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		patient.PatientHN = fmt.Sprintf("HN-%06d", seq.NextNumber)
		seq.NextNumber++

		if err := tx.Create(patient).Error; err != nil {
			return err
		}

		if err := tx.Model(&entity.PatientSequence{}).Where("hospital_id = ?", patient.HospitalID).Update("next_number", seq.NextNumber).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *patientRepository) FindByID(id uuid.UUID) (*entity.Patient, error) {
	var patient entity.Patient
	if err := r.db.First(&patient, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

// FindByHospitalAndIdentifier searches by national_id or passport_id within a hospital.
func (r *patientRepository) FindByHospitalAndIdentifier(hospitalID uuid.UUID, identifier string) (*entity.Patient, error) {
	var patient entity.Patient
	err := r.db.Where("hospital_id = ? AND (national_id = ? OR passport_id = ?)", hospitalID, identifier, identifier).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) FindByHospitalAndID(hospitalID uuid.UUID, id uuid.UUID) (*entity.Patient, error) {
	var patient entity.Patient
	err := r.db.Where("hospital_id = ? AND id = ?", hospitalID, id).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) ExistsByNationalIDInHospital(hospitalID uuid.UUID, nationalID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Patient{}).Where("hospital_id = ? AND national_id = ?", hospitalID, nationalID).Count(&count).Error
	return count > 0, err
}

func (r *patientRepository) ExistsByNationalIDInHospitalExcept(hospitalID uuid.UUID, nationalID string, patientID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Patient{}).
		Where("hospital_id = ? AND national_id = ? AND id <> ?", hospitalID, nationalID, patientID).
		Count(&count).Error
	return count > 0, err
}

func (r *patientRepository) ExistsByPassportIDInHospital(hospitalID uuid.UUID, passportID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Patient{}).Where("hospital_id = ? AND passport_id = ?", hospitalID, passportID).Count(&count).Error
	return count > 0, err
}

func (r *patientRepository) ExistsByPassportIDInHospitalExcept(hospitalID uuid.UUID, passportID string, patientID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Patient{}).
		Where("hospital_id = ? AND passport_id = ? AND id <> ?", hospitalID, passportID, patientID).
		Count(&count).Error
	return count > 0, err
}

// PatientSearchFilters holds optional filters for patient search.
type PatientSearchFilters struct {
	PatientHN    *string
	FirstNameTH  *string
	MiddleNameTH *string
	LastNameTH   *string

	FirstNameEN  *string
	MiddleNameEN *string
	LastNameEN   *string

	DateOfBirth *time.Time

	NationalID  *string
	PassportID  *string
	PhoneNumber *string
	Email       *string
	Gender      *string
}

func (r *patientRepository) FindByHospitalWithFilters(hospitalID uuid.UUID, f PatientSearchFilters, page, limit int) ([]entity.Patient, int64, error) {
	base := r.db.Model(&entity.Patient{}).Where("hospital_id = ?", hospitalID)

	if f.PatientHN != nil {
		base = base.Where("patient_hn = ?", *f.PatientHN)
	}
	if f.FirstNameTH != nil {
		base = base.Where("first_name_th = ?", *f.FirstNameTH)
	}
	if f.MiddleNameTH != nil {
		base = base.Where("middle_name_th = ?", *f.MiddleNameTH)
	}
	if f.LastNameTH != nil {
		base = base.Where("last_name_th = ?", *f.LastNameTH)
	}

	if f.FirstNameEN != nil {
		base = base.Where("first_name_en = ?", *f.FirstNameEN)
	}
	if f.MiddleNameEN != nil {
		base = base.Where("middle_name_en = ?", *f.MiddleNameEN)
	}
	if f.LastNameEN != nil {
		base = base.Where("last_name_en = ?", *f.LastNameEN)
	}

	if f.DateOfBirth != nil {
		base = base.Where("date_of_birth = ?", *f.DateOfBirth)
	}

	if f.NationalID != nil {
		base = base.Where("national_id = ?", *f.NationalID)
	}
	if f.PassportID != nil {
		base = base.Where("passport_id = ?", *f.PassportID)
	}
	if f.PhoneNumber != nil {
		base = base.Where("phone_number = ?", *f.PhoneNumber)
	}
	if f.Email != nil {
		base = base.Where("email = ?", *f.Email)
	}
	if f.Gender != nil {
		base = base.Where("gender = ?", *f.Gender)
	}

	// Count total with filters
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}

	var patients []entity.Patient
	err := base.Offset(offset).Limit(limit).Find(&patients).Error
	return patients, total, err
}

func (r *patientRepository) Update(patient *entity.Patient) (*entity.Patient, error) {
	err := r.db.Save(patient).Error
	return patient, err
}

func (r *patientRepository) Delete(hospitalID uuid.UUID, id uuid.UUID) error {
	res := r.db.Where("hospital_id = ? AND id = ?", hospitalID, id).Delete(&entity.Patient{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
