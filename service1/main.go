package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

var startTime time.Time

func main() {
	fmt.Println("Hello Gabr: relay")
	router := gin.Default()

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
	statusMsg := fmt.Sprintf("%s: uptime %s hours, free disk in root: %s MBytes\n", time.Now().Local().UTC().Format(time.RFC3339), uptime(), freeDiskSpace)
	// post msg to Status service
	// write status to log file
	f, err := os.OpenFile("./externalData/vstorage", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(statusMsg)
	if err != nil {
		log.Fatal(err)
	}
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
	//fmt.Printf("\n%s\n", freeMb)
	return freeMb
}
