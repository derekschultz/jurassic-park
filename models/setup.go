package models

import (
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Create logger for verbose mode
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info,
			Colorful: true,
		},
	)

	db, err := gorm.Open(sqlite.Open("park.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect to database")
	}

	// seed the database with default data
	// FIXME: make idempotent across app restarts
	var species = []Species{
		{Name: "Tyrannosaurus", Diet: "Carnivore", CageID: 1},
		{Name: "Velociraptor", Diet: "Carnivore", CageID: 2},
		{Name: "Spinosaurus", Diet: "Carnivore", CageID: 3},
		{Name: "Megalosaurus", Diet: "Carnivore", CageID: 4},
		{Name: "Brachiosaurus", Diet: "Herbivore", CageID: 5},
		{Name: "Stegosaurus", Diet: "Herbivore", CageID: 5},
		{Name: "Ankylosaurus", Diet: "Herbivore", CageID: 5},
		{Name: "Triceratops", Diet: "Herbivore", CageID: 5},
	}
	// Create cages for each carnivore species and a cage for all herbivores
	var cages = []Cage{
		{Status: "ACTIVE", Capacity: 1, MaxCapacity: 4},
		{Status: "ACTIVE", Capacity: 1, MaxCapacity: 4},
		{Status: "ACTIVE", Capacity: 1, MaxCapacity: 4},
		{Status: "ACTIVE", Capacity: 1, MaxCapacity: 4},
		{Status: "ACTIVE", Capacity: 4, MaxCapacity: 16},
	}

	db.AutoMigrate(&Species{}, &Cage{})
	db.Create(&species)
	db.Create(&cages)

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
