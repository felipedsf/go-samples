package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	DB_FILE      = "exchange.db"
	CREATE_TABLE = "CREATE TABLE IF NOT EXISTS exchange " +
		"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
		"code TEXT, " +
		"codein TEXT, " +
		"name TEXT," +
		"high TEXT," +
		"low TEXT," +
		"var_bid TEXT," +
		"pct_change TEXT," +
		"bid TEXT," +
		"ask TEXT," +
		"timestamp TEXT, " +
		"create_date TEXT)"
)

func GetDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		log.Fatalf("error openning db: %v", err)
		return nil
	}

	stmt, err := db.Prepare(CREATE_TABLE)
	if err != nil {
		log.Fatalf("error preparing statement db: %v", err)
		return nil
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("error preparing statement db: %v", err)
		return nil
	}
	return db
}
