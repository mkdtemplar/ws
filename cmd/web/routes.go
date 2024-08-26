package main

import (
	"ws/internal/handlers"

	"github.com/gin-gonic/gin"
)

func routes() *gin.Engine {
	mux := gin.Default()

	mux.GET("/", handlers.Home)
	mux.GET("/ws", handlers.WsEndpoint)

	return mux
}
