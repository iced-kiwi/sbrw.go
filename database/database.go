package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitializeDatabase(host string, port int, user string, password string, dbname string) error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		DB.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("PostgreSQL Connected successfully!")

	if err := CreateServerInformationTable(); err != nil {
		log.Printf("Failed to create/verify server_information table: %v", err)
	}

	if err := CreateUserTable(); err != nil {
		log.Printf("Failed to create/verify users table: %v", err)
	}

	if err := CreateTokenTable(); err != nil {
		log.Printf("Failed to create/verify token table: %v", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return DB
}
func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Error closing PostgreSQL connection: %v", err)
			return
		}
		log.Println("PostgreSQL connection closed.")
	}
}

func CheckError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
