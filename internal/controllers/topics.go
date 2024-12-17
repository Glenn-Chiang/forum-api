package controllers

import (
	errs "cvwo-backend/internal/errors"
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
		errs.HTTPErrorResponse(ctx, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, topics)
}
