package model

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	Name       string
	PatientID  string
	Status     string
	Department string
}
