package logs_server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

//func HelloController(c *gin.Context) {
//	c.JSON(200, gin.H{"message": "Hello, World!"})
//}

func GetLogs(c *gin.Context) {
	r := c.Request
	w := c.Writer
	appName := r.URL.Query().Get("appName")
	devstackLabel := r.URL.Query().Get("devstackLabel")
	logs, err := getLogs(appName, devstackLabel)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		return
	}
	response := struct {
		Logs []string `json:"logs"`
	}{
		Logs: logs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
