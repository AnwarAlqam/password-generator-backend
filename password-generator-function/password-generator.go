package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

func getPassword(c *gin.Context) {
	log.Println("Invoke ROOT")

	var queryParams QueryParams

	if err := c.ShouldBindQuery(&queryParams); err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if (queryParams.PasswordLength < 1)  {
		queryParams.PasswordLength = 1
	}

	numNonAlphabetic := min(20, queryParams.PasswordLength / 2)

	res, err := password.Generate(queryParams.PasswordLength, numNonAlphabetic / 2, numNonAlphabetic / 2, false, false)

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

	port_info := APIPort()
	router.Run(port_info)
	log.Println("API is up & running - " + port_info)
}
