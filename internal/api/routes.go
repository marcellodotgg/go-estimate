package api

import (
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gomarchy/estimate/internal/api/controller"
	"github.com/gomarchy/estimate/internal/api/middleware"
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
	router.Use(middleware.Audit())
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
		GET(":id", controller.Index).
		PUT("", controller.Create)
	router.Group("breakout/:id").
		POST("vote", controller.Vote).
		POST("reset", controller.Reset).
		GET("update-display-name", controller.UpdateDisplayNameModal).
		PATCH("update-display-name", controller.UpdateDisplayName)
}

func setupWebSocket() {
	controller := controller.NewWebSocketController()
	breakoutService := service.NewBreakoutService()

	router.
		Group("ws").
		GET(":channel", controller.Handle)

	websocket.Manager.HandleConnect(func(s *melody.Session) {
		channel, _ := s.Get("channel")
		userID := s.Request.URL.Query().Get("user_id")

		if userID == "" {
			return
		}

		breakoutService.AddUser(channel.(string), userID)
	})

	websocket.Manager.HandleDisconnect(func(s *melody.Session) {
		channel, _ := s.Get("channel")
		userID := s.Request.URL.Query().Get("user_id")
		breakoutService.RemoveUser(channel.(string), userID)
	})
}
