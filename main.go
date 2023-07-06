package main

import (
	logs_server "LogsServer/logs-server"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

func main() {
	http.HandleFunc("/logs", logs_server.GetLogs)
	http.HandleFunc("/", handleWebSocket)

	port := 8081
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server listening on port %d", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

//func handleLogs(w http.ResponseWriter, r *http.Request) {
//	appName := r.URL.Query().Get("appName")
//	logControllerGetLogs(w, r, appName)
//}
//
//func logControllerGetLogs(w http.ResponseWriter, r *http.Request, appName string) {
//	// Your implementation of logController.getLogs(appName) goes here
//	// Replace it with the logic to fetch the logs based on the appName.
//	// You can use the response writer to send the logs as the response.
//	// Example: w.Write([]byte("Logs for " + appName))
//}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from WebSocket:", err)
			break
		}

		if messageType == websocket.TextMessage {
			msg := string(message)
			splitMsg := strings.Split(msg, "*****")
			if len(splitMsg) != 2 {
				log.Println("Invalid message format")
				break
			}

			devstackLabel := splitMsg[0]
			service := splitMsg[1]

			kubectl := exec.Command("kubectl", "logs", "-f", "-l", fmt.Sprintf("devstack_label=%s", devstackLabel), "-n", service)

			stdout, err := kubectl.StdoutPipe()
			if err != nil {
				log.Println("Error creating stdout pipe:", err)
				break
			}

			if err := kubectl.Start(); err != nil {
				log.Println("Error starting kubectl command:", err)
				break
			}

			go func() {
				defer kubectl.Wait()

				buf := make([]byte, 1024)
				for {
					n, err := stdout.Read(buf)
					if n > 0 {
						conn.WriteMessage(websocket.TextMessage, buf[:n])
					}
					if err != nil {
						break
					}
				}
			}()
		}
	}
}
