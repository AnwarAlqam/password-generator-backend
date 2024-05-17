package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-password/password"
)

func APIPort() string {
	port := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		port = ":" + val
	}
	return port
}

func getPassword(c *gin.Context) {
	log.Println("Invoke ROOT")

	res, err := password.Generate(32, 10, 10, false, false)

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
	// return
}

func main() {
	router := gin.Default()
	router.GET("api/getPassword", getPassword)
	fmt.Println("HERE")

	port_info := APIPort()
	router.Run(port_info)
	log.Println("API is up & running - " + port_info)
}
