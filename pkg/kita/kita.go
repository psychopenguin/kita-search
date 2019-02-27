package kita

import "github.com/jinzhu/gorm"

type Kita struct {
	gorm.Model
	Name string
}
