package postgres

import "gorm.io/gorm"

type user struct {
	gorm.Model
	FirstName string
	LastName  string
	Address   address `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type address struct {
	gorm.Model
	Street     string
	PostalCode string
	UserId     uint
}
