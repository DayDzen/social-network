package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

func ConnectToDB() error {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		return fmt.Errorf("connectToDB err: %w", err)
	}

	dbConn = db

	return nil
}
