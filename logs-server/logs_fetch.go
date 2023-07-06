package logs_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//func HelloController(c *gin.Context) {
//	c.JSON(200, gin.H{"message": "Hello, World!"})
//}

func GetLogs(w http.ResponseWriter, r *http.Request) {
	//r := c.Request
	//w := c.Writer
	appName := r.URL.Query().Get("appName")
	devstackLabel := r.URL.Query().Get("devstackLabel")
	logs, err := getLogs(appName, devstackLabel)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		return
	}
	response := struct {
		Logs []interface{} `json:"logs"`
	}{
		Logs: logs,
	}
	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
