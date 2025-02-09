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
	Get(id int) (models.Container, error)
	GetAll() ([]models.Container, error)
	Create(address string) error
	Update(id int, newAddress string) error
	Delete(id int) error
	UpdateTime(containers []models.Container) error
}

type Handler struct {
	db database
}

func New(db database) *Handler {
	return &Handler{db}
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	container, err := h.db.Get(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, container)
}

func (h *Handler) GetAll(c *gin.Context) {
	containers, err := h.db.GetAll()
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	c.JSON(http.StatusOK, containers)
}

func (h *Handler) Create(c *gin.Context) {
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

	err = h.db.Create(msg.Address)
	if err != nil {
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}

	log.Println("Handler Created container")
	c.JSON(http.StatusOK, "ok")
}

func (h *Handler) Update(c *gin.Context) {
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

	

	err = h.db.Update(msg.Id, msg.Address)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	c.String(http.StatusOK, "ok")
}

func (h *Handler) UpdateTime(c *gin.Context) {
	containers := make([]models.Container, 0)

	err := c.ShouldBindJSON(&containers)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't unmarshall json")
		return
	}

	err = h.db.UpdateTime(containers)
	if err != nil {
		c.String(http.StatusBadRequest, "Error update time in psql")
		return
	}
	c.String(http.StatusOK, "Time updated")
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = h.db.Delete(id)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	c.JSON(http.StatusOK, "ok")
}
