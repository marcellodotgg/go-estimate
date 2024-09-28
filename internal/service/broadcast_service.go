package service

import (
	"bytes"
	"html/template"

	"github.com/gomarchy/estimate/internal/domain"
	"github.com/gomarchy/estimate/internal/infrastructure/database"
	"github.com/gomarchy/estimate/internal/infrastructure/websocket"
)

type BroadcastService interface {
	Breakout(breakoutID string)
	ResetVotes(breakoutID string)
}

type broadcastService struct {
}

func NewBroadcastService() BroadcastService {
	return broadcastService{}
}

// Broadcasts the latest version of the breakout that matches the given `breakoutID`.
// If this errors, it is a no-op.
func (s broadcastService) Breakout(breakoutID string) {
	var breakout domain.Breakout
	if err := database.DB.Preload("Users", "is_online = ?", true).First(&breakout, "id = ?", breakoutID).Error; err != nil {
		return
	}
	breakout.Cards = []domain.Card{
		{Value: "1"},
		{Value: "2"},
		{Value: "3"},
		{Value: "5"},
		{Value: "8"},
		{Value: "13"},
		{Value: "?"},
	}
	html, _ := s.renderTemplateToString("breakout/user_panel", breakout)
	websocket.UpdateChannel(breakout.ID, []byte(html))
}

// Broadcasts a reset vote event
func (s broadcastService) ResetVotes(breakoutID string) {
	var breakout domain.Breakout
	if err := database.DB.Preload("Users", "is_online = ?", true).First(&breakout, "id = ?", breakoutID).Error; err != nil {
		return
	}
	breakout.Cards = []domain.Card{
		{Value: "1"},
		{Value: "2"},
		{Value: "3"},
		{Value: "5"},
		{Value: "8"},
		{Value: "13"},
		{Value: "?"},
	}
	html, _ := s.renderTemplateToString("breakout/cards", breakout)
	websocket.UpdateChannel(breakout.ID, []byte(html))
}

func (s broadcastService) renderTemplateToString(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	var tmpl = template.Must(template.ParseGlob("templates/**/*"))

	err := tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
