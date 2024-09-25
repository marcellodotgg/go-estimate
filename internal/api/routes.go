package api

import (
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gomarchy/estimate/internal/api/controller"
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
		PUT("", controller.Create)
}
