package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
)

const Path = "BRACKETS_DB_PATH"
const DatabaseFile = "brackets-sqlite-database.db"

var DB *gorm.DB

func getPath() string {
	base, ok := os.LookupEnv(Path)
	if !ok {
		log.Fatal("BRACKETS_DB_PATH not set, aborting")
	}
	return path.Join(base, DatabaseFile)

}

func ensureDBExists() {
	dbPath := getPath()

	file, err := os.OpenFile(dbPath, os.O_RDWR, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(dbPath)
			if err != nil {
				log.Fatal("Could not create database file, aborting.", err.Error())
			}
			log.Println("Database file created.", dbPath)
		} else {
			log.Fatal("Could not read database file, aborting.", err.Error())
		}
	}
	defer file.Close()
}

func Init() {
	var err error
	ensureDBExists()
	DB, err = gorm.Open(sqlite.Open(getPath()), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not load sqlite3 database.", err.Error())
	}
}
