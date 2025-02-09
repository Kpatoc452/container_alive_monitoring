package main

import (
	"github.com/gin-contrib/cors"
	handlers "github.com/Kpatoc452/container_manager/controllers"
	"github.com/Kpatoc452/container_manager/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := storage.MustNew()
	handler := handlers.New(db)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	r.Use(cors.New(config))

	r.GET("/containers", handler.GetAllContainers)
	r.GET("/container/:id", handler.GetContainerByID)
	r.POST("/container", handler.CreateContainer)
	r.PUT("/container", handler.UpdateContainerByID)
	r.DELETE("/container/:id", handler.DeleteContainerByID)

	r.PUT("/pinger", handler.UpdateTimeContainers) // Endpoint for update time pings for pinger

	r.Run()
}
