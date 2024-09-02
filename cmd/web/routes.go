package main

import (
	"net/http"
	"ws/internal/handlers"

	"github.com/gin-gonic/gin"
)

func routes() *gin.Engine {
	mux := gin.Default()

	mux.GET("/", handlers.Home)

	mux.GET("/ws", handlers.WsEndpoint)

	mux.StaticFS("/static/", http.Dir("./static/"))

	return mux
}
