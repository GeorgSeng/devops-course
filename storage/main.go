package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(gin.Logger())
	router.POST("/log", postLog)
	router.GET("/log", getLog)
	router.Run()
}

func postLog(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(string(body))
	if err != nil {
		log.Fatal(err)
	}
}

func getLog(c *gin.Context) {
	f, err := os.OpenFile("./log.txt", os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		c.String(http.StatusInternalServerError, "Issue accessing the log file")
		fmt.Println(err)
		return
	}
	defer f.Close()
	msg, err := io.ReadAll(f)
	if err != nil {
		c.String(http.StatusInternalServerError, "Issue accessing the log file")
		fmt.Println(err)
		return
	}
	c.String(http.StatusOK, string(msg))
}
