package mocks

import (
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/google/uuid"
)

// MockStaffRepository - จำลอง staff repository สำหรับ unit test
type MockStaffRepository struct {
	staffs map[string]*entity.Staff // private Map<String, Staff> staffs;
}

// NewMockStaffRepository - สร้าง mock repository ใหม่
func NewMockStaffRepository() *MockStaffRepository {
	return &MockStaffRepository{
		staffs: make(map[string]*entity.Staff),
	}
}

// ===== Implement StaffRepository interface =====

func (m *MockStaffRepository) FindAll() ([]entity.Staff, error) {
	var staffs []entity.Staff
	for _, s := range m.staffs {
		staffs = append(staffs, *s)
	}
	return staffs, nil
}

func (m *MockStaffRepository) FindByUsername(username string) (*entity.Staff, error) {
	for _, s := range m.staffs {
		if s.Username == username {
			return s, nil
		}
	}
	return nil, nil
}

func (m *MockStaffRepository) FindByHospitalAndUsername(hospitalID uuid.UUID, username string) (*entity.Staff, error) {
	for _, s := range m.staffs {
		if s.HospitalID == hospitalID && s.Username == username {
			return s, nil
		}
	}
	return nil, nil
}

func (m *MockStaffRepository) Create(staff *entity.Staff) (*entity.Staff, error) {
	if m.staffs == nil {
		m.staffs = make(map[string]*entity.Staff)
	}
	m.staffs[staff.ID.String()] = staff
	return staff, nil
}

func (m *MockStaffRepository) DeleteByID(id uuid.UUID) error {
	delete(m.staffs, id.String())
	return nil
}

func (m *MockStaffRepository) ExistsByUsernameInHospital(hospitalID uuid.UUID, username string) (bool, error) {
	for _, s := range m.staffs {
		if s.HospitalID == hospitalID && s.Username == username {
			return true, nil
		}
	}
	return false, nil
}
