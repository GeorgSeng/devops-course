package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	NginxConfigPath  = "/etc/nginx/nginx.conf"
	BackupConfigPath = "/etc/nginx/nginx.conf.bak"
	NginxHealthURL   = "http://localhost:8081/nginx-health"
)

func main() {
	fmt.Println("Starting API Service for NGINX configuration update...")
	router := gin.Default()

	router.GET("/health", getHealth)
	router.POST("/config/nginx", updateNginxConfig)

	router.Run(":8080")
}

// region Handlers

// getHealth is an endpoint used to check if the gateway and NGINX are up and running
func getHealth(c *gin.Context) {
	// get an http client
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// check if NGINX is up
	resp, err := client.Get(NginxHealthURL)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"gateway_status": "up",
			"nginx_status":   "down",
			"error":          err.Error(),
		})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body) // close the client

	// check if NGINX is healthy
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"gateway_status":    "up",
			"nginx_status":      "unhealthy",
			"nginx_http_status": resp.StatusCode,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"gateway_status": "up",
		"nginx_status":   "up",
	})
}

// updateNginxConfig is an endpoint used to update the NGINX configuration file and reload NGINX
func updateNginxConfig(c *gin.Context) {
	// read new config file from request
	newConfig, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to read request body: %v", err)
		return
	}

	// get the old config file
	oldConfig, err := os.ReadFile(NginxConfigPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read old config file: %v", err)
		return
	}

	// save the old config file
	if err := os.WriteFile(BackupConfigPath, oldConfig, 0644); err != nil {
		c.String(http.StatusInternalServerError, "Failed to create backup config file: %v", err)
		return
	}

	// write the new config file
	if err := os.WriteFile(NginxConfigPath, newConfig, 0644); err != nil {
		c.String(http.StatusInternalServerError, "Failed to write new config file: %v", err)
		return
	}

	// reload the NGINX configuration
	if err := reloadNginx(oldConfig); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "NGINX configuration updated and successfully loaded")
}

// endregion Handlers

// region Utils

// reloadNginx tries to reload the NGINX configuration in case there is an issue the old config is restored
func reloadNginx(oldConfig []byte) error {
	cmd := exec.Command("nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// restore old file
		if restoreErr := os.WriteFile(NginxConfigPath, oldConfig, 0644); restoreErr != nil {
			return fmt.Errorf("NGINX reload failed. Writing old config file fails as well. NGINX realod error: %v. Write file error: %s", string(output), restoreErr)
		}
		// try old config
		retryCmd := exec.Command("nginx", "-s", "reload")
		output, err = retryCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to reload NGINX with old config. NGINX might be bricked. NGINX Reload error: %s", string(output))
		}
		return fmt.Errorf("failed to reload the new config file to NGINX. The old conig was sucesfully restored. NGINX Reload error: %s", string(output))
	}
	return nil
}

// endregion Utils
