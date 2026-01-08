package service

import (
	"testing"
	"time"

	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/service"
	"github.com/PhanukornKMITL/hospital-exam/tests/unit/service/mocks"
	"github.com/google/uuid"
)

// ============= TEST CASES =============

// Test 1: Get all patients successfully
func TestGetPatientsSuccess(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockPatientRepository()

	patient1 := &entity.Patient{
		ID:          uuid.New(),
		FirstNameTH: "สมชาย",
		LastNameTH:  "ใจดี",
	}
	patient2 := &entity.Patient{
		ID:          uuid.New(),
		FirstNameTH: "สมหญิง",
		LastNameTH:  "สวยงาม",
	}
	mockRepo.Create(patient1)
	mockRepo.Create(patient2)

	svc := service.NewPatientService(mockRepo)

	// Act
	patients, err := svc.GetPatients()

	// Assert
	if err != nil {
		t.Errorf("GetPatients() error = %v, want nil", err)
	}

	if len(patients) != 2 {
		t.Errorf("GetPatients() returned %d patients, want 2", len(patients))
	}
}

// Test 2: Get patients by hospital successfully
func TestGetPatientsByHospitalSuccess(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockPatientRepository()

	hospitalID1 := uuid.New()
	hospitalID2 := uuid.New()

	patient1 := &entity.Patient{
		ID:          uuid.New(),
		HospitalID:  hospitalID1,
		FirstNameTH: "สมชาย",
	}
	patient2 := &entity.Patient{
		ID:          uuid.New(),
		HospitalID:  hospitalID2,
		FirstNameTH: "สมหญิง",
	}
	mockRepo.Create(patient1)
	mockRepo.Create(patient2)

	svc := service.NewPatientService(mockRepo)

	// Act
	patients, err := svc.GetPatientsByHospital(hospitalID1)

	// Assert
	if err != nil {
		t.Errorf("GetPatientsByHospital() error = %v, want nil", err)
	}

	if len(patients) != 1 {
		t.Errorf("GetPatientsByHospital() returned %d patients, want 1", len(patients))
	}

	if patients[0].HospitalID != hospitalID1 {
		t.Errorf("Patient HospitalID = %v, want %v", patients[0].HospitalID, hospitalID1)
	}
}

// Test 3: Create patient successfully
func TestCreatePatientSuccess(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockPatientRepository()
	svc := service.NewPatientService(mockRepo)

	hospitalID := uuid.New()
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	nationalID := "1234567890123"
	input := service.PatientCreateInput{
		HospitalID:  hospitalID,
		FirstNameTH: "สมชาย",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Somchai",
		LastNameEN:  "Jaidit",
		DateOfBirth: &dob,
		NationalID:  nationalID,
		PhoneNumber: "0812345678",
		Email:       "somchai@example.com",
		Gender:      "M",
	}

	// Act
	result, err := svc.CreatePatient(input)

	// Assert
	if err != nil {
		t.Errorf("CreatePatient() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("CreatePatient() returned nil patient")
	}

	if result.FirstNameTH != input.FirstNameTH {
		t.Errorf("FirstNameTH = %v, want %v", result.FirstNameTH, input.FirstNameTH)
	}

	if result.NationalID != nil && *result.NationalID != nationalID {
		t.Errorf("NationalID = %v, want %v", *result.NationalID, nationalID)
	}
}
