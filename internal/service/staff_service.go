package service

import (
	"errors"
	"os"
	"time"

	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type StaffService interface {
	GetStaffs() ([]entity.Staff, error)
	CreateStaff(input StaffCreateInput) (*entity.Staff, error)
	DeleteStaff(id uuid.UUID) error
	Login(input StaffLoginInput) (string, error)
}

type staffService struct {
	repo repository.StaffRepository
}

type StaffCreateInput struct {
	Username   string
	Password   string
	HospitalID uuid.UUID
}

type StaffLoginInput struct {
	Username string
	Password string
}

func NewStaffService(repo repository.StaffRepository) StaffService {
	return &staffService{repo: repo}
}

func (s *staffService) GetStaffs() ([]entity.Staff, error) {
	return s.repo.FindAll()
}

func (s *staffService) CreateStaff(input StaffCreateInput) (*entity.Staff, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	staff := &entity.Staff{
		Username:   input.Username,
		Password:   string(hashed),
		HospitalID: input.HospitalID,
	}
	return s.repo.Create(staff)
}

func (s *staffService) DeleteStaff(id uuid.UUID) error {
	return s.repo.DeleteByID(id)
}

func (s *staffService) Login(input StaffLoginInput) (string, error) {
	staff, err := s.repo.FindByUsername(input.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(input.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	secret := os.Getenv("APP_JWT_SECRET")
	if secret == "" {
		return "", errors.New("jwt secret not configured")
	}

	claims := jwt.MapClaims{
		"sub":        staff.ID.String(),
		"username":   staff.Username,
		"hospitalId": staff.HospitalID.String(),
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}
