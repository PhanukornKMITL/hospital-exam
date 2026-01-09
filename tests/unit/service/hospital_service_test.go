package service

import (
	"testing"

	"github.com/PhanukornKMITL/hospital-exam/internal/service"
	"github.com/PhanukornKMITL/hospital-exam/tests/unit/service/mocks"
	"github.com/google/uuid"
)

// Test 1: Get all hospitals successfully (empty list)
func TestGetHospitalsSuccess_Empty(t *testing.T) {
	mockRepo := mocks.NewMockHospitalRepository()
	svc := service.NewHospitalService(mockRepo)

	result, err := svc.GetHospitals()

	if err != nil {
		t.Errorf("GetHospitals() error = %v, want nil", err)
	}

	// result can be nil or empty slice
	if result != nil && len(result) != 0 {
		t.Errorf("GetHospitals() length = %d, want 0", len(result))
	}
}

// Test 2: Get all hospitals successfully (with data)
func TestGetHospitalsSuccess_WithData(t *testing.T) {
	mockRepo := mocks.NewMockHospitalRepository()
	svc := service.NewHospitalService(mockRepo)

	// Create 2 hospitals
	input1 := service.HospitalCreateInput{
		Name:    "Hospital A",
		Address: "Bangkok",
	}
	input2 := service.HospitalCreateInput{
		Name:    "Hospital B",
		Address: "Chiang Mai",
	}

	_, err1 := svc.CreateHospital(input1)
	_, err2 := svc.CreateHospital(input2)

	if err1 != nil || err2 != nil {
		t.Fatalf("CreateHospital() failed: err1=%v, err2=%v", err1, err2)
	}

	result, err := svc.GetHospitals()

	if err != nil {
		t.Errorf("GetHospitals() error = %v, want nil", err)
	}

	if len(result) != 2 {
		t.Errorf("GetHospitals() length = %d, want 2", len(result))
	}

	// Check names
	names := make(map[string]bool)
	for _, h := range result {
		names[h.Name] = true
	}

	if !names["Hospital A"] || !names["Hospital B"] {
		t.Error("GetHospitals() missing expected hospital names")
	}
}

// Test 3: Create hospital successfully
func TestCreateHospitalSuccess(t *testing.T) {
	mockRepo := mocks.NewMockHospitalRepository()
	svc := service.NewHospitalService(mockRepo)

	input := service.HospitalCreateInput{
		Name:    "Central Hospital",
		Address: "123 Main St, Bangkok",
	}

	result, err := svc.CreateHospital(input)

	if err != nil {
		t.Errorf("CreateHospital() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("CreateHospital() returned nil hospital")
	}

	if result.Name != input.Name {
		t.Errorf("Name = %v, want %v", result.Name, input.Name)
	}

	if result.Address != input.Address {
		t.Errorf("Address = %v, want %v", result.Address, input.Address)
	}

	// ID ต้องถูกสร้าง
	if result.ID == uuid.Nil {
		t.Error("Hospital ID should not be empty")
	}
}

// Test 4: Create hospital with empty name should return error
func TestCreateHospitalWithEmptyName(t *testing.T) {
	mockRepo := mocks.NewMockHospitalRepository()
	svc := service.NewHospitalService(mockRepo)

	input := service.HospitalCreateInput{
		Name:    "",
		Address: "123 Main St, Bangkok",
	}

	result, err := svc.CreateHospital(input)

	if err == nil {
		t.Error("CreateHospital() with empty name should return error, but got nil")
	}

	if err != nil && err.Error() != "hospital name is required" {
		t.Errorf("Error message = %v, want 'hospital name is required'", err.Error())
	}

	if result != nil {
		t.Error("CreateHospital() should return nil hospital when error occurs")
	}
}

// Test 5: Create multiple hospitals
func TestCreateMultipleHospitals(t *testing.T) {
	mockRepo := mocks.NewMockHospitalRepository()
	svc := service.NewHospitalService(mockRepo)

	hospitals := []service.HospitalCreateInput{
		{Name: "Hospital 1", Address: "Address 1"},
		{Name: "Hospital 2", Address: "Address 2"},
		{Name: "Hospital 3", Address: "Address 3"},
	}

	for _, input := range hospitals {
		_, err := svc.CreateHospital(input)
		if err != nil {
			t.Errorf("CreateHospital() for %s failed: %v", input.Name, err)
		}
	}

	result, err := svc.GetHospitals()

	if err != nil {
		t.Errorf("GetHospitals() error = %v, want nil", err)
	}

	if len(result) != 3 {
		t.Errorf("GetHospitals() length = %d, want 3", len(result))
	}
}
