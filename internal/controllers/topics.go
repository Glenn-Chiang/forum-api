package controllers

import (
	"cvwo-backend/internal/services"
	"net/http"

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
