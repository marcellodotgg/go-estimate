package service

import (
	"bytes"
	"text/template"

	"github.com/gomarchy/estimate/internal/domain"
	"github.com/gomarchy/estimate/internal/infrastructure/database"
	"github.com/gomarchy/estimate/internal/infrastructure/websocket"
	"github.com/google/uuid"
)

type BreakoutService interface {
	AddUser(breakoutID, userID string)
	RemoveUser(breakoutID, userID string)
	Create(userID string) (domain.Breakout, error)
	Exists(breakoutID string) bool
}

type breakoutService struct {
}

var tmpl = template.Must(template.ParseGlob("templates/**/*"))

func NewBreakoutService() BreakoutService {
	return breakoutService{}
}

func (s breakoutService) Exists(breakoutID string) bool {
	return database.DB.First(&domain.Breakout{}, "id = ?", breakoutID).Error == nil
}

func (s breakoutService) Create(userID string) (domain.Breakout, error) {
	breakout := domain.Breakout{
		Audit:     domain.Audit{ID: uuid.NewString()},
		CreatedBy: userID,
	}
	err := database.DB.Create(&breakout).Error
	return breakout, err
}

func (s breakoutService) AddUser(breakoutID, userID string) {
	// check if the user already exists
	if err := database.DB.First(&domain.User{}, "user_id = ? AND breakout_id = ?", userID, breakoutID).Error; err == nil {
		s.broadcast(breakoutID)
		return
	}

	user := domain.User{
		Name:       "Guest",
		UserID:     userID,
		BreakoutID: breakoutID,
		Vote:       "",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return
	}

	s.broadcast(breakoutID)
}

func (s breakoutService) RemoveUser(breakoutID, userID string) {
	if err := database.DB.Delete(domain.User{}, "breakout_id = ? AND user_id = ?", breakoutID, userID).Error; err != nil {
		return
	}

	s.broadcast(breakoutID)
}

func (s breakoutService) broadcast(breakoutID string) {
	var breakout domain.Breakout
	if err := database.DB.Preload("Users").First(&breakout, "id = ?", breakoutID).Error; err != nil {
		return
	}
	html, _ := s.renderTemplateToString("breakout/sample", breakout)
	websocket.UpdateChannel(breakout.ID, []byte(html))
}

func (s breakoutService) renderTemplateToString(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer

	err := tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
