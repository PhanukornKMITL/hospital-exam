package config

import (
	"log"

	"gorm.io/gorm"
)

// EnsurePatientIndexes switches patient ID uniqueness from global to per-hospital.
func EnsurePatientIndexes(db *gorm.DB) error {
	stmts := []string{
		// Drop legacy global unique constraints if they exist
		"ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_national_id_key",
		"ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_passport_id_key",

		// Drop possible legacy unique indexes with common names
		"DROP INDEX IF EXISTS idx_patients_national_id",
		"DROP INDEX IF EXISTS idx_patients_passport_id",
		"DROP INDEX IF EXISTS uix_patients_national_id",
		"DROP INDEX IF EXISTS uix_patients_passport_id",
		"DROP INDEX IF EXISTS unique_patients_national_id",
		"DROP INDEX IF EXISTS unique_patients_passport_id",

		// Create composite unique indexes per hospital
		"CREATE UNIQUE INDEX IF NOT EXISTS uni_patients_hospital_nid ON patients (hospital_id, national_id)",
		"CREATE UNIQUE INDEX IF NOT EXISTS uni_patients_hospital_pid ON patients (hospital_id, passport_id)",
	}

	for _, s := range stmts {
		if err := db.Exec(s).Error; err != nil {
			log.Printf("EnsurePatientIndexes step failed: %s: %v", s, err)
			return err
		}
	}
	return nil
}
