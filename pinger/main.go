package main

import (
	"encoding/json"
	"io"
	"net/http"
	"runtime"
	"time"
)

var (
	numCPU   = runtime.GOMAXPROCS(0)
	url      = "http://api:8080/containers"
	interval = 10 * time.Minute
)

func GetContainers() []Container {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var containers []Container
	json.Unmarshal(respBody, &containers)

	return containers
}

func main() {
	time.Sleep(25 * time.Second) // waiting postgresql container and rest api server container

	wp := New()
	defer wp.Stop()

	wp.AddGroupWorker(numCPU)

	for {
		containers := GetContainers()

		var piecesContainers [][]Container

		chunkSize := (len(containers) + numCPU - 1) / numCPU

		for i := 0; i < len(containers); i += chunkSize {
			end := i + chunkSize

			if end > len(containers) {
				end = len(containers)
			}

			piecesContainers = append(piecesContainers, containers[i:end])
		}

		for _, c := range piecesContainers {
			wp.SendMsg(c)
		}

		time.Sleep(interval)
	}

}
