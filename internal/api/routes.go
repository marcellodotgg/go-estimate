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
	setupRoutes()
	// Run the API
	router.Run()
}

func setupRoutes() {
	controller := controller.NewHomePageController()

	router.
		GET("", controller.HomePage)
}
