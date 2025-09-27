package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const pathLogFile = "./log.txt"

func main() {
	fmt.Println("Hello Storage!")
	router := gin.Default()

	router.Use(gin.Logger())
	router.POST("/log", postLog)
	router.GET("/log", getLog)
	router.GET("/health", getHealth)
	router.Run()
}

// postLog handles the POST /log request and stores the log message in a text file
func postLog(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Could not extract the log message from the body: %v\n", err)
		c.String(http.StatusInternalServerError, "Could not extract the log message from the body")
		return
	}
	if len(body) == 0 {
		c.String(http.StatusBadRequest, "The message to be logged can not be empty")
		return
	}

	f, err := os.OpenFile(pathLogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Could not open the log file: %v\n", err)
		c.String(http.StatusInternalServerError, "Could not open the log file")
		return
	}
	defer f.Close()
	_, err = f.WriteString(string(body) + "\n")
	if err != nil {
		log.Printf("Could not append log message to the log file: %v\n", err)
		c.String(http.StatusInternalServerError, "Could not append log message to the log file")
		return
	}
}

// getLog handles the GET /log request and sends back all the log messages as text/plain
func getLog(c *gin.Context) {
	f, err := os.OpenFile(pathLogFile, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("Could not open or create the log file: %v\n", err)
		c.String(http.StatusInternalServerError, "Issue accessing the log file")
		return
	}
	defer f.Close()
	msg, err := io.ReadAll(f)
	if err != nil {
		log.Printf("Could not read the content of the log file: %v\n", err)
		c.String(http.StatusInternalServerError, "Issue accessing the log file")
		return
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", msg)
}

// getHealth is a endpoint used to check if the storage service is already up and running
func getHealth(c *gin.Context) {
	c.String(http.StatusOK, "up")
}
