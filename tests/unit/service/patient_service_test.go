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

	if result.Gender != "M" {
		t.Errorf("Gender = %v, want M", result.Gender)
	}
}

// Test 4: Create patient with gender 'F'
func TestCreatePatientWithFemaleGender(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockPatientRepository()
	svc := service.NewPatientService(mockRepo)

	hospitalID := uuid.New()
	dob := time.Date(1995, 5, 20, 0, 0, 0, 0, time.UTC)
	passportID := "AB1234567"
	input := service.PatientCreateInput{
		HospitalID:  hospitalID,
		FirstNameTH: "สมหญิง",
		LastNameTH:  "สวยงาม",
		FirstNameEN: "Somying",
		LastNameEN:  "Suayngam",
		DateOfBirth: &dob,
		PassportID:  passportID,
		PhoneNumber: "0898765432",
		Email:       "somying@example.com",
		Gender:      "F",
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

	if result.Gender != "F" {
		t.Errorf("Gender = %v, want F", result.Gender)
	}

	if result.FirstNameTH != input.FirstNameTH {
		t.Errorf("FirstNameTH = %v, want %v", result.FirstNameTH, input.FirstNameTH)
	}
}

// Test 5: Create patient with invalid gender 
func TestCreatePatientWithInvalidGender(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockPatientRepository()
	svc := service.NewPatientService(mockRepo)

	hospitalID := uuid.New()
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	nationalID := "9876543210987"
	
	// Test with invalid gender value
	input := service.PatientCreateInput{
		HospitalID:  hospitalID,
		FirstNameTH: "ทดสอบ",
		LastNameTH:  "ระบบ",
		FirstNameEN: "Test",
		LastNameEN:  "System",
		DateOfBirth: &dob,
		NationalID:  nationalID,
		PhoneNumber: "0811111111",
		Gender:      "X", // Invalid gender
	}

	// Act
	result, err := svc.CreatePatient(input)

	// Assert - should return error
	if err == nil {
		t.Error("CreatePatient() with invalid gender should return error, but got nil")
	} else {
		t.Logf("Got expected error: %v", err)
	}

	if err != nil && err.Error() != "gender must be either 'M' or 'F'" {
		t.Errorf("Error message = %v, want 'gender must be either 'M' or 'F''", err.Error())
	}

	if result != nil {
		t.Error("CreatePatient() should return nil patient when error occurs")
	}
}

// Test 6: Create patient without gender should return error
func TestCreatePatientWithoutGender(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockPatientRepository()
	svc := service.NewPatientService(mockRepo)

	hospitalID := uuid.New()
	dob := time.Date(1985, 12, 1, 0, 0, 0, 0, time.UTC)
	nationalID := "5555555555555"

	input := service.PatientCreateInput{
		HospitalID:  hospitalID,
		FirstNameTH: "ไม่ระบุ",
		LastNameTH:  "เพศ",
		FirstNameEN: "No",
		LastNameEN:  "Gender",
		DateOfBirth: &dob,
		NationalID:  nationalID,
		PhoneNumber: "0822222222",
		Gender:      "", // Empty gender
	}

	// Act
	result, err := svc.CreatePatient(input)

	// Assert - should return error
	if err == nil {
		t.Error("CreatePatient() without gender should return error, but got nil")
	}

	if err.Error() != "gender must be either 'M' or 'F'" {
		t.Errorf("Error message = %v, want 'gender must be either 'M' or 'F''", err.Error())
	}

	if result != nil {
		t.Error("CreatePatient() should return nil patient when error occurs")
	}
}

// Test 7: Search patient by identifier within same hospital
func TestSearchPatientByIdentifierSameHospital(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockPatientRepository()
	svc := service.NewPatientService(mockRepo)

	hospitalID := uuid.New()
	nationalID := "1234567890123"

	patient := &entity.Patient{
		ID:         uuid.New(),
		HospitalID: hospitalID,
		FirstNameTH: "สมชาย",
		LastNameTH:  "ใจดี",
		NationalID:  &nationalID,
	}

	mockRepo.Create(patient)

	// Act
	result, err := svc.SearchPatientByIdentifier(hospitalID, nationalID)

	// Assert
	if err != nil {
		t.Errorf("SearchPatientByIdentifier() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("SearchPatientByIdentifier() returned nil, want patient")
	}

	if result.ID != patient.ID {
		t.Errorf("Patient ID = %v, want %v", result.ID, patient.ID)
	}
}

// Test 8: Search patient by identifier across different hospitals should not return data
func TestSearchPatientByIdentifierDifferentHospital(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockPatientRepository()
	svc := service.NewPatientService(mockRepo)

	hospitalA := uuid.New()
	hospitalB := uuid.New()
	nationalID := "9876543210123"

	patient := &entity.Patient{
		ID:         uuid.New(),
		HospitalID: hospitalA,
		FirstNameTH: "สมหญิง",
		LastNameTH:  "สวยงาม",
		NationalID:  &nationalID,
	}

	mockRepo.Create(patient)

	// Act
	result, err := svc.SearchPatientByIdentifier(hospitalB, nationalID)

	// Assert
	if err != nil {
		t.Errorf("SearchPatientByIdentifier() error = %v, want nil", err)
	}

	if result != nil {
		t.Errorf("SearchPatientByIdentifier() should return nil when hospital mismatched, got %v", result.ID)
	}
}
