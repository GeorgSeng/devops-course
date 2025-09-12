package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"unicode"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var startTime time.Time

func main() {
	fmt.Println("Hello Gabr: relay")
	router := gin.Default()

	// Add CORS middleware with default settings
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Update this to match your Nuxt.js frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(gin.Logger())
	router.GET("/status", getStatus)
	router.GET("/log", getLog)

	startTime = time.Now()
	router.Run() // default port is 8080
}

func getStatus(c *gin.Context) {
	/*cType := c.Request.Header.Get("Content-Type")
	if cType != "application/plain" {

	}*/
	freeDiskSpace := getFreeDiskSpaceOfRootInMb()
	statusMsg := fmt.Sprintf("%s: uptime %s hours, free disk in root: %s MBytes", time.Now().Local(), uptime(), freeDiskSpace)
	// post msg to Status service
	// write status to log file
	// get status of service2
	c.String(http.StatusOK, "%s", statusMsg)
}

func getLog(c *gin.Context) {
	// get log of storage service
	c.String(http.StatusOK, "<get log of storage service>")
}

func uptime() string {
	return time.Since(startTime).String()
}

func getFreeDiskSpaceOfRootInMb() string {
	cmd := exec.Command("df", "-m", "/")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("problem in output 2")
		log.Fatal(err)
	}

	// A fixed strucutred layout can be assumend
	// ----
	// ➜  test df -m /
	// Filesystem     1M-blocks  Used Available Use% Mounted on
	// /dev/sdd         1031019  6288    972287   1% /
	// ----
	// therefore the correct element to extrext is index 10
	freeMb := strings.FieldsFunc(string(out), unicode.IsSpace)[10]
	fmt.Printf("\n%s\n", freeMb)
	return freeMb
}
