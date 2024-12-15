package database

import (
	"os"
	"path"
	"sync"

	"github.com/meanii/teagit/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Db   *gorm.DB
	Path string
}

var name = "teagit.db"
var DatabaseInstance *Database
var once sync.Once

func NewDatabase() *Database {
	once.Do(func() {

		DatabaseInstance = &Database{}

		// Set the path
		DatabaseInstance.setPath(name)

		// Connect to the database
		DatabaseInstance.connect()

		// Migrate the schema
		DatabaseInstance.autoMigrate()
	})

	return DatabaseInstance
}

func (d *Database) Close() {
	sqlDB, err := d.Db.DB()
	if err != nil {
		panic("failed to close database")
	}
	sqlDB.Close()
}

// setPath sets the path of the database
func (d *Database) setPath(name string) {
	// get the path from the environment variable
	teagitwd := os.Getenv("TEAGITHOME")
	if teagitwd == "" {
		wd, _ := os.UserHomeDir()
		teagitwd = path.Join(wd, ".config", "teagit")
	}

	// make sure the directory exists
	os.MkdirAll(teagitwd, os.ModePerm)

	path := path.Join(teagitwd, name)
	d.Path = path
}

// connect connects to the database
func (d *Database) connect() {
	db, err := gorm.Open(sqlite.Open(d.Path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	d.Db = db
}

// autoMigrate migrates the schema
func (d *Database) autoMigrate() {
	d.Db.AutoMigrate(&models.Configs{})
	d.Db.AutoMigrate(&models.Users{})
	d.Db.AutoMigrate(&models.Keys{})
}
