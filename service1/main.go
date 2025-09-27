package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

var startTime time.Time

func main() {
	fmt.Println("Hello Service1!")
	router := gin.Default()

	router.Use(gin.Logger())
	router.GET("/status", getStatus)
	router.GET("/log", getLog)

	startTime = time.Now()
	router.Run() // default port is 8080
}

// getStatus handles the GET /status request
func getStatus(c *gin.Context) {
	/*cType := c.Request.Header.Get("Content-Type")
	if cType != "application/plain" {

	}*/
	freeDiskSpace := getFreeDiskSpaceOfRootInMb()
	statusMsg := fmt.Sprintf("%s: uptime %s hours, free disk in root: %s MBytes", time.Now().UTC().Format(time.RFC3339), uptime(), freeDiskSpace)

	// post status msg to storage service
	resStorage, err := http.Post("http://storage:8080/log", "text/plain", bytes.NewBufferString(statusMsg))
	if err != nil {
		log.Printf("Could not log status to storage service: %v\n", err)
		c.String(http.StatusInternalServerError, "%s", statusMsg)
		return
	}
	defer resStorage.Body.Close()
	if resStorage.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resStorage.Body)
		log.Printf("Logging status to storage failed: %v\n", body)
		c.String(http.StatusInternalServerError, "%s", statusMsg)
		return
	}

	// write status to vstorage log file
	f, err := os.OpenFile("./externalData/vstorage", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Could open/create vstorage: %v\n", err)
		c.String(http.StatusInternalServerError, "Could not open/create vstorage")
		return
	}
	defer f.Close()
	_, err = f.WriteString(statusMsg + "\n")
	if err != nil {
		log.Printf("Could not log to vstorage: %v\n", err)
		c.String(http.StatusInternalServerError, "Could not log to vstorage")
		return
	}

	// get status of service2
	res, err := http.Get("http://service2:8080/status")
	if err != nil {
		log.Printf("Could not reache service2 service: %v\n", err)
		c.String(http.StatusOK, "%s", statusMsg)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Could not read response body of service2: %v\n", err)
		c.String(http.StatusOK, "%s", statusMsg)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Service2 had an issue processing, status code: %v\n", res.StatusCode)
		return
	}
	data := []byte(fmt.Sprintf("%s\n%s", statusMsg, string(body)))
	c.Data(http.StatusOK, "text/plain; charset=utf-8", data)
}

// getLog handels the GET /log request
func getLog(c *gin.Context) {
	// get log of storage service
	res, err := http.Get("http://storage:8080/log")
	if err != nil {
		log.Printf("Could not reache storage service: %v\n", err)
		c.String(http.StatusServiceUnavailable, "Could not reache storage service")
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Could not read response body: %v\n", err)
		c.String(http.StatusInternalServerError, "Could not read response body of storage service")
		return
	}
	if res.StatusCode != http.StatusOK {
		c.String(http.StatusInternalServerError, "Storage service had problems %s", res.StatusCode)
		return
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", body)
}

func uptime() string {
	return strconv.FormatFloat(time.Since(startTime).Hours(), 'f', 8, 64)
}

// getFreeDiskSpaceOfRootInMb calles the df linux command to get the available space in the root file system in MB.
func getFreeDiskSpaceOfRootInMb() string {
	cmd := exec.Command("df", "-m", "/")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Could not call df to get free diskspace: %v\n", err)
		return "-1"
	}

	// A fixed strucutred layout can be assumend
	// ----
	// ➜  test df -m /
	// Filesystem     1M-blocks  Used Available Use% Mounted on
	// /dev/sdd         1031019  6288    972287   1% /
	// ----
	// therefore the correct element to extrext is index 10
	elems := strings.FieldsFunc(string(out), unicode.IsSpace)
	if len(elems) >= 10 {
		freeMb := elems[10]
		//fmt.Printf("\n%s\n", freeMb)
		return freeMb
	}
	return "-1"
}
