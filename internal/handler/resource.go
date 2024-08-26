package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/malamsyah/go-skele/internal/dto"
	"github.com/malamsyah/go-skele/internal/middleware"
	"github.com/malamsyah/go-skele/internal/model"
	"github.com/malamsyah/go-skele/internal/service"
	"gorm.io/gorm"
)

type ResourceHandler struct {
	service service.ResourceService
}

func NewResourceHandler(service service.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		service: service,
	}
}

func (h *ResourceHandler) RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.ValidateID())
	r.GET("/resources", h.GetResources)
	r.POST("/resources", h.CreateResource)
	r.GET("/resources/:id", h.GetResource)
	r.PUT("/resources/:id", h.UpdateResource)
	r.DELETE("/resources/:id", h.DeleteResource)
}

func (h *ResourceHandler) GetResources(c *gin.Context) {
	var resources []model.Resource
	var err error

	resources, err = h.service.GetResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	resp := dto.SuccessResponse{Data: resources}
	c.JSON(http.StatusOK, resp)
}

func (h *ResourceHandler) CreateResource(c *gin.Context) {
	var request dto.Resource
	var err error

	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	resource, err := h.service.CreateResource(request.ResourceType, request.Payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	resp := dto.SuccessResponse{Data: resource}
	c.JSON(http.StatusCreated, resp)

}

func (h *ResourceHandler) GetResource(c *gin.Context) {
	parsedID, exists := c.Get("parsedID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID not found in context"})
		return
	}

	id := uint(parsedID.(int))

	resource, err := h.service.GetResource(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	resp := dto.SuccessResponse{Data: resource}
	c.JSON(http.StatusOK, resp)
}

func (h *ResourceHandler) UpdateResource(c *gin.Context) {
	parsedID, exists := c.Get("parsedID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID not found in context"})
		return
	}

	id := uint(parsedID.(int))

	var request dto.Resource
	var err error

	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	if request.ID != id {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "ID in path and body should be the same"})
		return
	}

	res, err := h.service.UpdateResource(id, model.Resource{
		ResourceType: request.ResourceType,
		Payload:      request.Payload,
		Model: gorm.Model{
			ID: request.ID,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ResourceHandler) DeleteResource(c *gin.Context) {
	parsedID, exists := c.Get("parsedID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID not found in context"})
		return
	}

	id := uint(parsedID.(int))

	err := h.service.DeleteResource(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
