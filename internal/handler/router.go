package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/malamsyah/go-skele/internal/db"
	"github.com/malamsyah/go-skele/internal/model"
	"github.com/malamsyah/go-skele/internal/repository"
	"github.com/malamsyah/go-skele/internal/service"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", Health)

	dbConn, err := db.ConnectPostgres()
	if err != nil {
		panic(err)
	}

	resourceRepository := repository.NewRepository[model.Resource](dbConn)

	// resource CRUD routes
	resourceHandler := NewResourceHandler(service.NewResourceService(resourceRepository))
	resourceHandler.RegisterRoutes(r)

	return r
}
