package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

var (
	numCPU = runtime.GOMAXPROCS(0)
)

func main() {
	wp := New()
	defer wp.Stop()

	wp.AddGroupWorker(numCPU)

	for {
		req, err := http.NewRequest("GET", "http://localhost:8080/containers", nil)
		if err != nil {
			panic(err)
		}

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		respBody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}

		var body []Container
		json.Unmarshal(respBody, &body)

		var divided [][]Container

		chunkSize := (len(body) + numCPU - 1) / numCPU

		for i := 0; i < len(body); i += chunkSize {
			end := i + chunkSize

			if end > len(body) {
				end = len(body)
			}

			divided = append(divided, body[i:end])
		}

		for _, c := range divided {
			wp.SendMsg(c)
		}

		time.Sleep(10 * time.Minute)
	}

}
