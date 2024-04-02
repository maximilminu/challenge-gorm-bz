package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetConnection() {
	// var err error
	// db, err := sql.Open("sqlite3", "challenge.db")
	// if err != nil {
	// 	panic(err)
	// }
	// return db
	var err error
	DB, err = gorm.Open(sqlite.Open("challenge.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

}

func MakeMigrations(db *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS endpoint1 (
	);
	
	CREATE TABLE IF NOT EXISTS endpoint2 (
	);`

	_, err := db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}
