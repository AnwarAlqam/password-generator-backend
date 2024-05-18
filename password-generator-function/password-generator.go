package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-password/password"
)

type QueryParams struct {
	PasswordLength int `form:"passwordLength" binding:"required"`
}

func APIPort() string {
	port := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		port = ":" + val
	}
	return port
}

func getCount() int {
	// Attempt to connect to database
	database, err := sql.Open("sqlite3", "./password-generator.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		return -1
	}

	rows, err := database.Query("SELECT count FROM password_count")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var numPasswordGenerated int

	for rows.Next() {
		if err := rows.Scan(&numPasswordGenerated); err != nil {
			log.Fatal(err)
		}
	}

	return numPasswordGenerated
}

func updateCount() {
	// Attempt to connect to database
	database, err := sql.Open("sqlite3", "./password-generator.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		return
	}

	numPasswordGenerated := getCount()
	numPasswordGenerated += 1

	statement, err := database.Prepare("INSERT INTO password_count (count) VALUES (?)")
	if err != nil {
		log.Fatalf("Failed to prepare INSERT statement: %v", err)
		return
	}
	statement.Exec(numPasswordGenerated)
}

func getPasswordGeneratedCount(c *gin.Context) {
	numPasswordGenerated := getCount()

	c.JSON(http.StatusOK, gin.H{
		"message": numPasswordGenerated,
	})
}

func getPassword(c *gin.Context) {
	log.Println("Invoke ROOT")

	var queryParams QueryParams

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if queryParams.PasswordLength < 1 {
		queryParams.PasswordLength = 1
	}

	numNonAlphabetic := min(20, queryParams.PasswordLength/2)

	res, err := password.Generate(queryParams.PasswordLength, numNonAlphabetic/2, numNonAlphabetic/2, false, false)

	if err != nil {
		log.Fatal(err)
	}

	updateCount()

	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
	// return
}

func setUpDatabase() {
	// Attempt to connect to database
	database, err := sql.Open("sqlite3", "./password-generator.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		return
	}

	fmt.Println("Connected to database")

	// Create table
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS password_count (count INTEGER)")
	if err != nil {
		log.Fatalf("Failed to prepare CREATE TABLE statement: %v", err)
		return
	}
	statement.Exec()
	fmt.Println("Created table")

	// Check if a row exists
	rows, err := database.Query("SELECT COUNT(*) FROM password_count")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var numRows int

	for rows.Next() {
		if err := rows.Scan(&numRows); err != nil {
			log.Fatal(err)
		}
	}

	if numRows == 0 {
		// Add initial row
		statement, err = database.Prepare("INSERT INTO password_count (count) VALUES (0)")
		if err != nil {
			log.Fatalf("Failed to prepare INSERT statement: %v", err)
			return
		}
		statement.Exec()
		fmt.Println("Added initial row")
	}
}

func main() {
	setUpDatabase()

	router := gin.Default()
	router.GET("api/getPassword", getPassword)
	router.GET("api/getPasswordGeneratedCount", getPasswordGeneratedCount)

	port_info := APIPort()
	router.Run(port_info)
	log.Println("API is up & running - " + port_info)
}
