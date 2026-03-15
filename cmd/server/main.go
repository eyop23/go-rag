package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/eyop23/insurance-go/config"
	"github.com/eyop23/insurance-go/handlers"
	"github.com/eyop23/insurance-go/middleware"
)

func main() {
	cfg := config.Load()

	r := gin.Default()
	r.Use(middleware.CORS())
	r.POST("/ask", handlers.AskHandler(cfg))

	log.Printf("Insurance RAG server running on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
