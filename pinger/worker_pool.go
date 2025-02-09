package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"sync"
	"time"
)

type workerPool struct {
	msgChan          (chan []Container)
	destroyChan      chan struct{}
	waitGroupWorkers sync.WaitGroup

	mutex        sync.Mutex
	opts         OptionWP
	currentId    int
	countWorkers int
}

func New(opts ...Option) *workerPool {
	return &workerPool{
		opts:        NewOptionWP(opts...),
		msgChan:     make(chan []Container),
		destroyChan: make(chan struct{}),
	}
}

func (wp *workerPool) AddWorker() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	if wp.countWorkers < wp.opts.Max {
		wp.waitGroupWorkers.Add(1)

		wp.countWorkers++
		wp.currentId++

		go wp.process(wp.currentId)
		wp.opts.Logger.Logf("[ADD] Worker %d created", wp.currentId)
	}
}

func (wp *workerPool) AddGroupWorker(countWorkers int) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	for range countWorkers {
		if wp.countWorkers < wp.opts.Max {
			wp.waitGroupWorkers.Add(1)

			wp.countWorkers++
			wp.currentId++

			go wp.process(wp.currentId)
			wp.opts.Logger.Logf("[ADD] Worker %d created", wp.currentId)
		}
	}
}

func (wp *workerPool) DestroyWorker() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	if wp.countWorkers > 0 {
		wp.destroyChan <- struct{}{}

		wp.countWorkers--
	}
}

func (wp *workerPool) process(id int) {
	defer wp.waitGroupWorkers.Done()
	for {
		select {
		case msg := <-wp.msgChan:
			wp.opts.Logger.Logf("Worker %d is running\n", id)
			if len(msg) == 0 {
				wp.opts.Logger.Logf("Len msg is 0\n")
				break
			}

			for i := 0; i < len(msg); i++ {
				container := msg[i]

				// TODO: make faster parsing
				address, err := netip.ParseAddrPort(container.Address)
				if err != nil {
					wp.opts.Logger.Logf("%s Couldn't parse netip\n", container)
					break
				}

				host := address.Addr().String()
				port := strconv.Itoa(int(address.Port()))

				if wp.ping(host, port) {
					msg[i].LastSuccessPing = time.Now()
				}

				msg[i].LastPing = time.Now()
			}

            err := wp.sendContainers(msg)
			if err != nil {
				wp.opts.Logger.Logf("error send containers %v\n", err)
				break
			}

			wp.opts.Logger.Logf("Worker %d sended data", id)
		case <-wp.destroyChan:
			fmt.Printf("[DELETE] Worker %d destroyed\n", id)
			wp.opts.Logger.Logf("[DELETE] Worker %d destroyed\n", id)
			return
		}
	}
}

func (wp *workerPool) Stop() {
	close(wp.destroyChan)
	wp.waitGroupWorkers.Wait()
	close(wp.msgChan)
	wp.countWorkers = 0
}

func (wp *workerPool) SendMsg(msg []Container) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	if wp.countWorkers > 0 {
		wp.msgChan <- msg
	}
}

func (wp *workerPool) GetCountWorkers() int {
	return wp.countWorkers
}

func (wp *workerPool) ping(ip string, port string) bool {
	timeout := time.Second

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		wp.opts.Logger.Logf("Connection error %v\n", err)
		return false
	}

	if conn != nil {
		defer conn.Close()
		wp.opts.Logger.Logf("Connection successful!")
		return true
	}

	return false
}

func (wp *workerPool) sendContainers(data []Container) error {
	url := "http://localhost:8080/pinger"

	marshaled, err := json.Marshal(data)
	if err != nil {
		wp.opts.Logger.Logf("Error Marshalling to json\n")
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(marshaled))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	_, err = client.Do(req)

	return err
}
