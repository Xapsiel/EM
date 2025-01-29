package handler

import (
	_ "github.com/Xapsiel/EffectiveMobile/docs"
	"github.com/Xapsiel/EffectiveMobile/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/info", h.GetSongs)
	router.POST("/songs", h.AddSong)
	router.GET("/info/verse", h.GetSongVerse)
	router.DELETE("/songs", h.DeleteSong)
	router.PUT("/songs", h.UpdateSong)
	return router
}
