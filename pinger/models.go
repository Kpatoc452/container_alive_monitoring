package main

import "time"

type Container struct {
	Id              int       `json:"id"`
	Address         string    `json:"address"`
	LastSuccessPing time.Time `json:"last_success_ping"`
	LastPing        time.Time `json:"last_ping"`
}
