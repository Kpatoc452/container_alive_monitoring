package storage

import (
	"log"
	"testing"
	"time"

	"github.com/Kpatoc452/container_manager/models"
)

func TestPsqlContainers(t *testing.T) {
	db := MustNew()
	log.Println("Connected")
	err := db.CreateContainer("192.168.0.1:12345")

	if err != nil {
		panic(err)
	}
	log.Println("Created")

	test_cases := []models.Container{
		{Id: 60, Address: "127.1.0.1:8080",LastPing:  time.Now(),LastSuccessPing:  time.Now()},
	}
	
	err = db.UpdateTimeContainers(test_cases)
	log.Println(err)
}
