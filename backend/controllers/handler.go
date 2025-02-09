package handlers

import (
	"log"
	"net/http"
	"net/netip"
	"strconv"

	"github.com/Kpatoc452/container_manager/models"
	"github.com/gin-gonic/gin"
)

type database interface {
	GetContainerByID(id int) (models.Container, error)
	GetAllContainers() ([]models.Container, error)
	CreateContainer(address string) error
	UpdateContainerByID(id int, newAddress string) error
	DeleteContainerByID(id int) error
	UpdateTimeContainers(containers []models.Container) error
}

type Handler struct {
	db database
}

func New(db database) *Handler {
	return &Handler{db}
}

func (h *Handler) GetContainerByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	container, err := h.db.GetContainerByID(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, container)
}

func (h *Handler) GetAllContainers(c *gin.Context) {
	containers, err := h.db.GetAllContainers()
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	c.JSON(http.StatusOK, containers)
}

func (h *Handler) CreateContainer(c *gin.Context) {
	type createJson struct {
		Address string `json:"address"`
	}
	var msg createJson

	err := c.ShouldBindJSON(&msg)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	_, err = netip.ParseAddrPort(msg.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid IP address")
		return
	}

	err = h.db.CreateContainer(msg.Address)
	if err != nil {
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}

	log.Println("Handler Created container")
	c.JSON(http.StatusOK, "ok")
}

func (h *Handler) UpdateContainerByID(c *gin.Context) {
	type updateJson struct {
		Id      int    `json:"id"`
		Address string `json:"address"`
	}

	var msg updateJson
	err := c.ShouldBindJSON(&msg)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = netip.ParseAddrPort(msg.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid IP address")
		return
	}

	

	err = h.db.UpdateContainerByID(msg.Id, msg.Address)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	c.String(http.StatusOK, "ok")
}

func (h *Handler) UpdateTimeContainers(c *gin.Context) {
	containers := make([]models.Container, 0)

	err := c.ShouldBindJSON(&containers)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't unmarshall json")
		return
	}

	err = h.db.UpdateTimeContainers(containers)
	if err != nil {
		c.String(http.StatusBadRequest, "Error update time in psql")
		return
	}
	c.String(http.StatusOK, "Time updated")
}

func (h *Handler) DeleteContainerByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = h.db.DeleteContainerByID(id)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	c.JSON(http.StatusOK, "ok")
}
