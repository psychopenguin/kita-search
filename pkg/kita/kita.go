package kita

import "github.com/jinzhu/gorm"

type District struct {
	gorm.Model
	Name  string `gorm:"UNIQUE"`
	Kitas []Kita
}

type Kita struct {
	gorm.Model
	Name       string
	Email      string
	Permalink  string `gorm:"UNIQUE"`
	District   District
	DistrictID int
}
