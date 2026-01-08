package mocks

import (
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

func (m *MockPatientRepository) ExistsByNationalIDInHospital(hospitalID uuid.UUID, nationalID string) (bool, error) {
	for _, p := range m.patients {
		if p.HospitalID == hospitalID && p.NationalID != nil && *p.NationalID == nationalID {
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

func (m *MockPatientRepository) FindByHospitalWithFilters(hospitalID uuid.UUID, filters repository.PatientSearchFilters, page, limit int) ([]entity.Patient, int64, error) {
	return nil, 0, nil
}
