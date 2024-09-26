package api

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gomarchy/estimate/internal/api/controller"
	"github.com/gomarchy/estimate/internal/infrastructure/websocket"
	"github.com/olahol/melody"
)

var tmpl *template.Template
var router = gin.Default()
var channelCounts = &websocket.ChannelCounts{
	Counts: make(map[string]int),
}

func Start() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	// Load templates
	router.LoadHTMLGlob("templates/**/*")
	tmpl = template.Must(template.ParseGlob("templates/**/*"))
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

	router.
		Group("ws").
		GET(":channel", controller.Handle)

	websocket.Manager.HandleConnect(func(s *melody.Session) {
		channel, _ := s.Get("channel")
		channelCounts.Increment(channel.(string))
		updateChannel(channel.(string))
	})

	websocket.Manager.HandleDisconnect(func(s *melody.Session) {
		channel, _ := s.Get("channel")
		channelCounts.Decrement(channel.(string))
		updateChannel(channel.(string))
	})

	websocket.Manager.HandleMessage(func(s *melody.Session, msg []byte) {
		channel, _ := s.Get("channel")
		updateChannel(channel.(string))
	})
}

func updateChannel(channel string) {
	html, err := renderTemplateToString("breakout/sample", gin.H{})
	if err != nil {
		return
	}
	websocket.UpdateChannel(channel, []byte(html))
}

func renderTemplateToString(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer

	err := tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
