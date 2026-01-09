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

// Test 6: Update hospital successfully
func TestUpdateHospitalSuccess(t *testing.T) {
	mockRepo := mocks.NewMockHospitalRepository()
	svc := service.NewHospitalService(mockRepo)

	createInput := service.HospitalCreateInput{
		Name:    "Old Hospital",
		Address: "Old Address",
	}

	hospital, err := svc.CreateHospital(createInput)
	if err != nil {
		t.Fatalf("CreateHospital() failed: %v", err)
	}

	updateInput := service.HospitalUpdateInput{
		Name:    "Updated Hospital",
		Address: "New Address",
	}

	result, err := svc.UpdateHospital(hospital.ID, updateInput)

	if err != nil {
		t.Errorf("UpdateHospital() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("UpdateHospital() returned nil hospital")
	}

	if result.Name != updateInput.Name {
		t.Errorf("Name = %v, want %v", result.Name, updateInput.Name)
	}

	if result.Address != updateInput.Address {
		t.Errorf("Address = %v, want %v", result.Address, updateInput.Address)
	}
}

// Test 7: Update hospital with empty name should return error
func TestUpdateHospitalWithEmptyName(t *testing.T) {
	mockRepo := mocks.NewMockHospitalRepository()
	svc := service.NewHospitalService(mockRepo)

	createInput := service.HospitalCreateInput{
		Name:    "Test Hospital",
		Address: "Test Address",
	}

	hospital, _ := svc.CreateHospital(createInput)

	updateInput := service.HospitalUpdateInput{
		Name:    "",
		Address: "New Address",
	}

	result, err := svc.UpdateHospital(hospital.ID, updateInput)

	if err == nil {
		t.Error("UpdateHospital() with empty name should return error, but got nil")
	}

	if err != nil && err.Error() != "hospital name is required" {
		t.Errorf("Error message = %v, want 'hospital name is required'", err.Error())
	}

	if result != nil {
		t.Error("UpdateHospital() should return nil when error occurs")
	}
}

// Test 8: Update non-existent hospital should return error
func TestUpdateHospitalNotFound(t *testing.T) {
	mockRepo := mocks.NewMockHospitalRepository()
	svc := service.NewHospitalService(mockRepo)

	nonExistentID := uuid.New()
	updateInput := service.HospitalUpdateInput{
		Name:    "Updated Hospital",
		Address: "New Address",
	}

	result, err := svc.UpdateHospital(nonExistentID, updateInput)

	if err == nil {
		t.Error("UpdateHospital() with non-existent ID should return error, but got nil")
	}

	if err != nil && err.Error() != "hospital not found" {
		t.Errorf("Error message = %v, want 'hospital not found'", err.Error())
	}

	if result != nil {
		t.Error("UpdateHospital() should return nil when hospital not found")
	}
}

// Test 9: Delete hospital successfully
func TestDeleteHospitalSuccess(t *testing.T) {
	mockRepo := mocks.NewMockHospitalRepository()
	svc := service.NewHospitalService(mockRepo)

	createInput := service.HospitalCreateInput{
		Name:    "Hospital to Delete",
		Address: "Delete Address",
	}

	hospital, _ := svc.CreateHospital(createInput)

	err := svc.DeleteHospital(hospital.ID)

	if err != nil {
		t.Errorf("DeleteHospital() error = %v, want nil", err)
	}

	hospitals, _ := svc.GetHospitals()
	if len(hospitals) != 0 {
		t.Errorf("After delete, GetHospitals() returned %d hospitals, want 0", len(hospitals))
	}
}

// Test 10: Delete non-existent hospital should return error
func TestDeleteHospitalNotFound(t *testing.T) {
	mockRepo := mocks.NewMockHospitalRepository()
	svc := service.NewHospitalService(mockRepo)

	nonExistentID := uuid.New()

	err := svc.DeleteHospital(nonExistentID)

	if err == nil {
		t.Error("DeleteHospital() with non-existent ID should return error, but got nil")
	}

	if err != nil && err.Error() != "hospital not found" {
		t.Errorf("Error message = %v, want 'hospital not found'", err.Error())
	}
}
