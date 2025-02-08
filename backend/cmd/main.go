package main

import (
	handlers "github.com/Kpatoc452/container_manager/controllers"
	"github.com/Kpatoc452/container_manager/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := storage.MustNew()
	handler := handlers.New(db)

	r.GET("/containers", handler.GetAll)
	r.GET("/container/:id", handler.Get)
	r.POST("/container", handler.Create)
	r.PUT("/container", handler.Update)
	r.DELETE("/container/:id", handler.Delete)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
