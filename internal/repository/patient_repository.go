package repository

import (
	"errors"
	"fmt"

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
	ExistsByNationalIDInHospital(hospitalID uuid.UUID, nationalID string) (bool, error)
	ExistsByPassportIDInHospital(hospitalID uuid.UUID, passportID string) (bool, error)
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

func (r *patientRepository) ExistsByNationalIDInHospital(hospitalID uuid.UUID, nationalID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Patient{}).Where("hospital_id = ? AND national_id = ?", hospitalID, nationalID).Count(&count).Error
	return count > 0, err
}

func (r *patientRepository) ExistsByPassportIDInHospital(hospitalID uuid.UUID, passportID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Patient{}).Where("hospital_id = ? AND passport_id = ?", hospitalID, passportID).Count(&count).Error
	return count > 0, err
}
