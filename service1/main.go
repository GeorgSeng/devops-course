package main

import (
	"bytes"
	"fmt"
	"io"
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

	// post status msg to storage service
	resStorage, err := http.Post("http://storage:8080/log", "text/plain", bytes.NewBufferString(statusMsg))
	if err != nil {
		fmt.Println("Could not log status to storage service: %s", err)
	}
	defer resStorage.Body.Close()
	if resStorage.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resStorage.Body)
		fmt.Println("Logging status to storage failed: %s", body)
	}

	// write status to vstorage log file
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
	res, err := http.Get("http://service2:8080/status")
	if err != nil {
		println("Could not reache service2 service")
		println(err)
		c.String(http.StatusOK, "%s", statusMsg)
		return
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			println("Could not read response body")
			println(err)
		} else {
			c.String(http.StatusOK, "%s%s", statusMsg, string(body))
		}
	}
}

func getLog(c *gin.Context) {
	// get log of storage service
	res, err := http.Get("http://storage:8080/log")
	if err != nil {
		println("Could not reache storage service: %s", err)
		c.String(http.StatusServiceUnavailable, "Could not reache storage service")
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		println("Could not read response body: %s", err)
		c.String(http.StatusInternalServerError, "Could not read response body of storage service")
		return
	}
	if res.StatusCode != http.StatusOK {
		c.String(http.StatusInternalServerError, "Storage service had problems %s", res.StatusCode)
		return
	}
	c.String(http.StatusOK, string(body))
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
