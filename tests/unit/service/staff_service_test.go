package service

import (
	"testing"

	"github.com/PhanukornKMITL/hospital-exam/internal/service"
	"github.com/PhanukornKMITL/hospital-exam/tests/unit/service/mocks"
	"github.com/google/uuid"
)

// ============= TEST CASES =============

// Test 1: Create staff successfully
func TestCreateStaffSuccess(t *testing.T) {
	
	mockRepo := mocks.NewMockStaffRepository()
	svc := service.NewStaffService(mockRepo)

	hospitalID := uuid.New()
	input := service.StaffCreateInput{
		Username:   "admin1",
		Password:   "password123",
		HospitalID: hospitalID,
	}

	// Act
	result, err := svc.CreateStaff(input)

	// Assert
	if err != nil {
		t.Errorf("CreateStaff() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("CreateStaff() returned nil staff")
	}

	if result.Username != input.Username {
		t.Errorf("Username = %v, want %v", result.Username, input.Username)
	}

	if result.HospitalID != input.HospitalID {
		t.Errorf("HospitalID = %v, want %v", result.HospitalID, input.HospitalID)
	}

	// Password ต้องถูก hash ไม่ใช่ plain text
	if result.Password == input.Password {
		t.Error("Password should be hashed, but it's plain text")
	}

	// ID ต้องถูกสร้าง
	if result.ID == uuid.Nil {
		t.Error("Staff ID should not be empty")
	}
}

// Test 2: Create staff with duplicate username should fail
func TestCreateStaffDuplicateUsername(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockStaffRepository()
	svc := service.NewStaffService(mockRepo)

	hospitalID := uuid.New()
	input := service.StaffCreateInput{
		Username:   "admin1",
		Password:   "password123",
		HospitalID: hospitalID,
	}

	// สร้าง staff ตัวแรก
	_, err := svc.CreateStaff(input)
	if err != nil {
		t.Fatalf("First CreateStaff() failed: %v", err)
	}

	// Act: ลองสร้าง staff ที่มี username เดียวกัน
	_, err = svc.CreateStaff(input)

	// Assert: ควรจะ error
	if err == nil {
		t.Error("CreateStaff() with duplicate username should return error, but got nil")
	}

	if err.Error() != "username already exists in this hospital" {
		t.Errorf("Error message = %v, want 'username already exists in this hospital'", err.Error())
	}
}

// Test 3: Get all staffs successfully
func TestGetStaffsSuccess(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockStaffRepository()
	svc := service.NewStaffService(mockRepo)

	hospitalID := uuid.New()
	input1 := service.StaffCreateInput{
		Username:   "admin1",
		Password:   "password123",
		HospitalID: hospitalID,
	}
	input2 := service.StaffCreateInput{
		Username:   "admin2",
		Password:   "password456",
		HospitalID: hospitalID,
	}

	svc.CreateStaff(input1)
	svc.CreateStaff(input2)

	// Act
	staffs, err := svc.GetStaffs()

	// Assert
	if err != nil {
		t.Errorf("GetStaffs() error = %v, want nil", err)
	}

	if len(staffs) != 2 {
		t.Errorf("GetStaffs() returned %d staffs, want 2", len(staffs))
	}
}

// Test 4: Create staff with empty username should fail
func TestCreateStaffEmptyUsername(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockStaffRepository()
	svc := service.NewStaffService(mockRepo)

	input := service.StaffCreateInput{
		Username:   "", // ชื่อว่าง
		Password:   "password123",
		HospitalID: uuid.New(),
	}

	// Act
	result, _ := svc.CreateStaff(input)

	// Assert
	// CreateStaff ไม่ validate username ว่างใน layer นี้
	// แต่ควรจะ fail ที่ repo หรือ database layer
	// ขึ้นอยู่กับ implementation
	if result != nil {
		t.Errorf("CreateStaff() with empty username should return nil or error")
	}
}

// Test 5: Delete staff successfully
func TestDeleteStaffSuccess(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockStaffRepository()
	svc := service.NewStaffService(mockRepo)

	hospitalID := uuid.New()
	input := service.StaffCreateInput{
		Username:   "admin1",
		Password:   "password123",
		HospitalID: hospitalID,
	}

	staff, _ := svc.CreateStaff(input)

	// Act
	err := svc.DeleteStaff(staff.ID)

	// Assert
	if err != nil {
		t.Errorf("DeleteStaff() error = %v, want nil", err)
	}

	// ยืนยันว่า staff ถูกลบแล้ว
	staffs, _ := svc.GetStaffs()
	if len(staffs) != 0 {
		t.Errorf("After delete, GetStaffs() returned %d staffs, want 0", len(staffs))
	}
}
