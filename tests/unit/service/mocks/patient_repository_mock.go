package mocks

import (
	"time"

	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/repository"
	"github.com/google/uuid"
)

// MockPatientRepository - จำลอง patient repository สำหรับ unit test
// โดยเก็บ data ไว้ใน memory แทนที่จะเชื่อมต่อกับ database จริง
type MockPatientRepository struct {
	patients map[string]*entity.Patient
}

// NewMockPatientRepository - สร้าง mock repository ใหม่
func NewMockPatientRepository() *MockPatientRepository {
	return &MockPatientRepository{
		patients: make(map[string]*entity.Patient),
	}
}

// ===== Implement PatientRepository interface =====

func (m *MockPatientRepository) FindAll() ([]entity.Patient, error) {
	var patients []entity.Patient
	for _, p := range m.patients {
		patients = append(patients, *p)
	}
	return patients, nil
}

func (m *MockPatientRepository) FindByHospitalID(hospitalID uuid.UUID) ([]entity.Patient, error) {
	var patients []entity.Patient
	for _, p := range m.patients {
		if p.HospitalID == hospitalID {
			patients = append(patients, *p)
		}
	}
	return patients, nil
}

func (m *MockPatientRepository) Create(patient *entity.Patient) (*entity.Patient, error) {
	if m.patients == nil {
		m.patients = make(map[string]*entity.Patient)
	}
	m.patients[patient.ID.String()] = patient
	return patient, nil
}

func (m *MockPatientRepository) CreateWithGeneratedHN(patient *entity.Patient) (*entity.Patient, error) {
	if m.patients == nil {
		m.patients = make(map[string]*entity.Patient)
	}
	m.patients[patient.ID.String()] = patient
	return patient, nil
}

func (m *MockPatientRepository) FindByID(id uuid.UUID) (*entity.Patient, error) {
	if p, ok := m.patients[id.String()]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *MockPatientRepository) FindByHospitalAndIdentifier(hospitalID uuid.UUID, identifier string) (*entity.Patient, error) {
	for _, p := range m.patients {
		if p.HospitalID == hospitalID {
			if p.NationalID != nil && *p.NationalID == identifier {
				return p, nil
			}
			if p.PassportID != nil && *p.PassportID == identifier {
				return p, nil
			}
		}
	}
	return nil, nil
}

func (m *MockPatientRepository) FindByHospitalAndID(hospitalID uuid.UUID, id uuid.UUID) (*entity.Patient, error) {
	for _, p := range m.patients {
		if p.HospitalID == hospitalID && p.ID == id {
			return p, nil
		}
	}
	return nil, nil
}

func (m *MockPatientRepository) ExistsByNationalIDInHospital(hospitalID uuid.UUID, nationalID string) (bool, error) {
	for _, p := range m.patients {
		if p.HospitalID == hospitalID && p.NationalID != nil && *p.NationalID == nationalID {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockPatientRepository) ExistsByNationalIDInHospitalExcept(hospitalID uuid.UUID, nationalID string, patientID uuid.UUID) (bool, error) {
	for _, p := range m.patients {
		if p.HospitalID == hospitalID && p.ID != patientID && p.NationalID != nil && *p.NationalID == nationalID {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockPatientRepository) ExistsByPassportIDInHospital(hospitalID uuid.UUID, passportID string) (bool, error) {
	for _, p := range m.patients {
		if p.HospitalID == hospitalID && p.PassportID != nil && *p.PassportID == passportID {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockPatientRepository) ExistsByPassportIDInHospitalExcept(hospitalID uuid.UUID, passportID string, patientID uuid.UUID) (bool, error) {
	for _, p := range m.patients {
		if p.HospitalID == hospitalID && p.ID != patientID && p.PassportID != nil && *p.PassportID == passportID {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockPatientRepository) FindByHospitalWithFilters(hospitalID uuid.UUID, filters repository.PatientSearchFilters, page, limit int) ([]entity.Patient, int64, error) {
	var results []entity.Patient

	// Filter by hospital and fields
	for _, p := range m.patients {
		if p.HospitalID != hospitalID {
			continue
		}

		if filters.PatientHN != nil && p.PatientHN != *filters.PatientHN {
			continue
		}
		if filters.FirstNameTH != nil && p.FirstNameTH != *filters.FirstNameTH {
			continue
		}
		if filters.MiddleNameTH != nil && p.MiddleNameTH != *filters.MiddleNameTH {
			continue
		}
		if filters.LastNameTH != nil && p.LastNameTH != *filters.LastNameTH {
			continue
		}
		if filters.FirstNameEN != nil && p.FirstNameEN != *filters.FirstNameEN {
			continue
		}
		if filters.MiddleNameEN != nil && p.MiddleNameEN != *filters.MiddleNameEN {
			continue
		}
		if filters.LastNameEN != nil && p.LastNameEN != *filters.LastNameEN {
			continue
		}
		if filters.DateOfBirth != nil && !sameDate(p.DateOfBirth, filters.DateOfBirth) {
			continue
		}
		if filters.NationalID != nil {
			if p.NationalID == nil || *p.NationalID != *filters.NationalID {
				continue
			}
		}
		if filters.PassportID != nil {
			if p.PassportID == nil || *p.PassportID != *filters.PassportID {
				continue
			}
		}
		if filters.PhoneNumber != nil && p.PhoneNumber != *filters.PhoneNumber {
			continue
		}
		if filters.Email != nil && p.Email != *filters.Email {
			continue
		}
		if filters.Gender != nil && p.Gender != *filters.Gender {
			continue
		}

		results = append(results, *p)
	}

	total := int64(len(results))

	// Pagination
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	start := (page - 1) * limit
	if start >= len(results) {
		return []entity.Patient{}, total, nil
	}
	end := start + limit
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], total, nil
}

func (m *MockPatientRepository) Update(patient *entity.Patient) (*entity.Patient, error) {
	m.patients[patient.ID.String()] = patient
	return patient, nil
}

func (m *MockPatientRepository) Delete(hospitalID uuid.UUID, id uuid.UUID) error {
	for key, p := range m.patients {
		if p.HospitalID == hospitalID && p.ID == id {
			delete(m.patients, key)
			return nil
		}
	}
	return nil
}

// sameDate compares date by YYYY-MM-DD ignoring time of day and timezone.
func sameDate(a *time.Time, b *time.Time) bool {
	if a == nil || b == nil {
		return a == b
	}
	aa := a.UTC()
	bb := b.UTC()
	return aa.Year() == bb.Year() && aa.Month() == bb.Month() && aa.Day() == bb.Day()
}
