package controllers

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TopicController struct {
	service services.TopicService
}

func NewTopicController(service services.TopicService) *TopicController {
	return &TopicController{service}
}

// GET /topics
func (controller *TopicController) GetAll(ctx *gin.Context) {
	topics, err := controller.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, topics)
}

// POST /topics
func (controller *TopicController) Create(ctx *gin.Context) {
	var topic models.Topic

	if err := ctx.ShouldBindJSON(&topic); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTopic, err := controller.service.Create(&topic)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newTopic)
}

// DELETE /topics/:id
func (controller *TopicController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic ID"})		
		return
	}

	if err:= controller.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
