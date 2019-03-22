package kita

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/gormigrate.v1"
)

func Migrate() {
	// TODO: put this under configuration
	db, err := gorm.Open("mysql", "kita:kita@tcp(127.0.0.1:3306)/kita?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.LogMode(true)
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
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
				return tx.CreateTable(&District{}, &Kita{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("districts", "kitas").Error
			},
		},
	})
	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration sucess")
}
