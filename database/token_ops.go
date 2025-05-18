package database

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"gosbrw/database/structs"
	"log"
	"time"
)

func GetUserByToken(token string) (structs.Token, error) {
	db := GetDB()
	var userToken structs.Token
	
	query := "SELECT * FROM tokens WHERE token=$1"
	err := db.QueryRow(query, token).Scan(
		&userToken.UserID,
		&userToken.Token,
		&userToken.ExpiresAt,
	)
	
	if err != nil {
		return structs.Token{}, err
	}
	
	return userToken, nil
}

func GenerateUserToken(userID int, expiresAt string) string {
	token := GenerateRandomToken()
	db := GetDB()
	
	query := "INSERT INTO tokens (user_id, token, expires_at) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, userID, token, expiresAt)
	
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return ""
	}
	
	return token
}

func GenerateRandomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error generating random token: %v", err)
		return ""
	}
	return fmt.Sprintf("%x", b)
}

func DeleteToken(token string) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	
	query := "DELETE FROM tokens WHERE token=$1"
	_, err := db.Exec(query, token)
	return err
}


func CreateTokenTable() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS token (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL UNIQUE,
		expires_at TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Printf("Error creating token table: %v", err)
		return fmt.Errorf("error creating token table: %w", err)
	}
	log.Println("Token table checked/created successfully.")
	return nil
}

func GenerateToken(UserID int, expiresAt string) string {
	token := generateRandomToken()
	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return ""
	}

	_, err = tx.Exec("DELETE FROM token WHERE user_id = $1", UserID)
	if err != nil {
		log.Printf("Error deleting existing tokens: %v", err)
		tx.Rollback()
		return ""
	}

	_, err = tx.Exec("INSERT INTO token (user_id, token, expires_at) VALUES ($1, $2, $3)",
		UserID, token, expiresAt)
	if err != nil {
		log.Printf("Error inserting new token: %v", err)
		tx.Rollback()
		return ""
	}

	_, err = tx.Exec("UPDATE users SET security_token = $1 WHERE id = $2", token, UserID)
	if err != nil {
		log.Printf("Error updating user security token: %v", err)
		tx.Rollback()
		return ""
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return ""
	}

	return token
}

func generateRandomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error generating random token: %v", err)
		return ""
	}

	return fmt.Sprintf("%x", b)
}

func VerifyToken(tokenString string) (int, bool, error) {
	db := GetDB()
	if db == nil {
		return 0, false, fmt.Errorf("database connection is not initialized")
	}

	query := `
		SELECT t.user_id, t.expires_at 
		FROM token t 
		WHERE t.token = $1
	`

	var userID int
	var expiresAt string

	err := db.QueryRow(query, tokenString).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("error querying token: %w", err)
	}

	expTime, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return 0, false, fmt.Errorf("error parsing expiration time: %w", err)
	}

	if time.Now().After(expTime) {
		go cleanupExpiredToken(userID)
		return 0, false, nil
	}

	return userID, true, nil
}

func cleanupExpiredToken(userID int) {
	db := GetDB()
	if db == nil {
		log.Printf("Cannot cleanup expired token: database connection is not initialized")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction for token cleanup: %v", err)
		return
	}

	_, err = tx.Exec("DELETE FROM token WHERE user_id = $1", userID)
	if err != nil {
		log.Printf("Error deleting expired token: %v", err)
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE users SET security_token = NULL WHERE id = $1", userID)
	if err != nil {
		log.Printf("Error clearing user security token: %v", err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction for token cleanup: %v", err)
		return
	}

	log.Printf("Successfully cleaned up expired token for user ID: %d", userID)
}
