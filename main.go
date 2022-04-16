package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func init() {
	if err := godotenv.Load("prod_config.yaml"); err != nil {
		if err = godotenv.Load("local_config.yaml"); err != nil {
			log.Print("No config file found")
		}
	}
}

func main() {
	fmt.Print("Hello")
}
