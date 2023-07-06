package logs_server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func getPods(appName string, devstackLabel string) ([]string, error) {
	podsNames := make([]string, 0)

	command := fmt.Sprintf("kubectl get pods -n %s | grep %s", appName, devstackLabel)
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("123")
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("1234")
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		pod := scanner.Text()
		if !strings.Contains(pod, "mysql") {
			podsNames = append(podsNames, strings.Fields(pod)[0])
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("12345")
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return podsNames, nil
}

func getLogs(appName string, devStackLabel string) ([]interface{}, error) {
	pods, err := getPods(appName, devStackLabel)
	if err != nil {
		return nil, err
	}
	if len(pods) == 0 {
		return nil, fmt.Errorf("unable to find pods for service %s", appName)
	}

	command := fmt.Sprintf("kubectl logs -l devstack_label=%s -n %s", devStackLabel, appName)
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var logs []interface{}

	podLogs := strings.Split(string(stdout), "\n")
	for _, podLog := range podLogs {
		var data interface{}
		fmt.Println("LOGS", podLog)
		err := json.Unmarshal([]byte(podLog), &data)
		if err != nil {
			logs = append(logs, podLog)
			continue
		}
		logs = append(logs, data)
	}

	return logs, nil
}
