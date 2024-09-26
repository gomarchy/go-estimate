package api

import (
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gomarchy/estimate/internal/api/controller"
	"github.com/gomarchy/estimate/internal/infrastructure/websocket"
	"github.com/gomarchy/estimate/internal/service"
	"github.com/olahol/melody"
)

var router = gin.Default()

func Start() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	// Load templates
	router.LoadHTMLGlob("templates/**/*")
	// Load static files
	router.Static("/static", "public/")
	router.StaticFS("/.well-known/acme-challenge", http.Dir("/var/www/html/.well-known/acme-challenge"))
	// Setup API routes
	setupHomePageRoutes()
	setupBreakoutRoutes()
	setupWebSocket()
	// Run the API
	router.Run()
}

func setupHomePageRoutes() {
	controller := controller.NewHomePageController()

	router.
		GET("", controller.HomePage)
}

func setupBreakoutRoutes() {
	controller := controller.NewBreakoutController()

	router.
		Group("breakout").
		GET("", controller.Index).
		PUT("", controller.Create)
}

func setupWebSocket() {
	controller := controller.NewWebSocketController()
	breakoutService := service.NewBreakoutService()

	router.
		Group("ws").
		GET(":channel", controller.Handle)

	websocket.Manager.HandleConnect(func(s *melody.Session) {
		channel, _ := s.Get("channel")
		breakoutService.Add(channel.(string))
	})

	websocket.Manager.HandleDisconnect(func(s *melody.Session) {
		channel, _ := s.Get("channel")
		breakoutService.Remove(channel.(string))
	})
}
