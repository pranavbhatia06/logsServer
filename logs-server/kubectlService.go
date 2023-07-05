package logs_server

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const devstackLabel = "hdfc-test"

func getPods(appName string) ([]string, error) {
	podsNames := make([]string, 0)

	command := fmt.Sprintf("kubectl get pods -n %s | grep %s", appName, devstackLabel)
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
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
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return podsNames, nil
}

func getLogs(appName string) ([]string, error) {
	pods, err := getPods(appName)
	if err != nil {
		return nil, err
	}

	logs := make([]string, 0)

	for _, pod := range pods {
		if !strings.Contains(pod, appName) {
			continue
		}

		command := fmt.Sprintf("kubectl logs %s -n %s", pod, appName)
		cmd := exec.Command("bash", "-c", command)
		stdout, err := cmd.Output()
		if err != nil {
			log.Println(err)
			continue
		}

		podLogs := strings.Split(string(stdout), "\n")
		for _, podLog := range podLogs {
			logs = append(logs, podLog)
		}
	}

	return logs, nil
}
