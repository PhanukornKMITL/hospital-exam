package mocks

import (
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/google/uuid"
)

// MockHospitalRepository - จำลอง hospital repository สำหรับ unit test
type MockHospitalRepository struct {
	hospitals map[string]*entity.Hospital
}

// NewMockHospitalRepository - สร้าง mock repository ใหม่
func NewMockHospitalRepository() *MockHospitalRepository {
	return &MockHospitalRepository{
		hospitals: make(map[string]*entity.Hospital),
	}
}

// ===== Implement HospitalRepository interface =====

func (m *MockHospitalRepository) FindAll() ([]entity.Hospital, error) {
	var hospitals []entity.Hospital
	if m.hospitals == nil {
		return hospitals, nil
	}
	for _, h := range m.hospitals {
		hospitals = append(hospitals, *h)
	}
	return hospitals, nil
}

func (m *MockHospitalRepository) FindByID(id uuid.UUID) (*entity.Hospital, error) {
	if h, ok := m.hospitals[id.String()]; ok {
		return h, nil
	}
	return nil, nil
}

func (m *MockHospitalRepository) FindByName(name string) (*entity.Hospital, error) {
	for _, h := range m.hospitals {
		if h.Name == name {
			return h, nil
		}
	}
	return nil, nil
}

func (m *MockHospitalRepository) Create(hospital *entity.Hospital) (*entity.Hospital, error) {
	if m.hospitals == nil {
		m.hospitals = make(map[string]*entity.Hospital)
	}

	// Generate UUID if not provided
	if hospital.ID == uuid.Nil {
		hospital.ID = uuid.New()
	}

	m.hospitals[hospital.ID.String()] = hospital
	return hospital, nil
}

func (m *MockHospitalRepository) Update(hospital *entity.Hospital) (*entity.Hospital, error) {
	if _, ok := m.hospitals[hospital.ID.String()]; ok {
		m.hospitals[hospital.ID.String()] = hospital
		return hospital, nil
	}
	return nil, nil
}

func (m *MockHospitalRepository) Delete(id uuid.UUID) error {
	delete(m.hospitals, id.String())
	return nil
}
