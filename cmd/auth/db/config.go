package db

import (
	"database/sql"
	// "fmt"

	// "fmt"
	"log"
	"net/url"

	// "os"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func InitializeDatabase() (*sql.DB, error) {
	AIVEN_URL := "postgres://avnadmin:AVNS_vQEeq6UsUmcVW69WV5T@pg-1d133d20-lokeshkatari921-5634.a.aivencloud.com:21652/defaultdb?sslmode=require"

	conn, _ := url.Parse(AIVEN_URL)
	conn.RawQuery = "sslmode=verify-ca;sslrootcert=ca.pem"

	db, err := sql.Open("postgres", conn.String())
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return nil, err
	}

	log.Println("Connected to PostgreSQL database")

	return db, nil
}
