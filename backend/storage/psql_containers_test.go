package storage

import (
	"fmt"
	"log"
	"testing"
)

func TestPsqlContainers(t *testing.T) {
	db := MustNew()
	log.Println("Connected")
	err := db.Create("192.168.0.1:12345")

	if err != nil {
		panic(err)
	}
	log.Println("Created")

	containers, err := db.GetAll()
	if err != nil {
		panic(err)
	}

	fmt.Println(containers)

	for _, c := range containers {
		err = db.Delete(c.Id)
		if err != nil { 
			panic(err)
		}
	}

	containers, err = db.GetAll()

	fmt.Println(containers)
}
