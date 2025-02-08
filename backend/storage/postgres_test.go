package storage

import (
	"log"
	"testing"
)

func TestMustNew(t *testing.T) {
	db := MustNew() 
	log.Println("Connected")
	err := db.Create("192.168.0.1:12345")
	if err != nil {
		log.Println(err)
	}
}
