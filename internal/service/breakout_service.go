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
	FindByID(breakoutID string) (domain.Breakout, error)
	Broadcast(breakoutID string)
}

type breakoutService struct {
}

func NewBreakoutService() BreakoutService {
	return breakoutService{}
}

func (s breakoutService) FindByID(breakoutID string) (domain.Breakout, error) {
	var breakout domain.Breakout
	if err := database.DB.First(&breakout, "id = ?", breakoutID).Error; err != nil {
		return breakout, err
	}
	return breakout, nil
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
	if err := database.DB.First(&domain.User{}, "user_id = ? AND breakout_id = ?", userID, breakoutID).Error; err == nil {
		database.DB.Where("breakout_id = ? AND user_id = ?", breakoutID, userID).Model(&domain.User{}).Update("is_online", true)
		s.Broadcast(breakoutID)
		return
	}

	user := domain.User{
		Name:       "Guest",
		UserID:     userID,
		BreakoutID: breakoutID,
		Vote:       "",
		IsOnline:   true,
	}

	database.DB.Create(&user)
	s.Broadcast(breakoutID)
}

func (s breakoutService) RemoveUser(breakoutID, userID string) {
	database.DB.Where("breakout_id = ? AND user_id = ?", breakoutID, userID).Model(&domain.User{}).Update("is_online", false)
	s.Broadcast(breakoutID)
}

// Broadcasts the latest version of the breakout that matches the given `breakoutID`.
// If this errors, it is a no-op.
func (s breakoutService) Broadcast(breakoutID string) {
	var breakout domain.Breakout
	if err := database.DB.Preload("Users", "is_online = ?", true).First(&breakout, "id = ?", breakoutID).Error; err != nil {
		return
	}
	html, _ := s.renderTemplateToString("breakout/sample", breakout)
	websocket.UpdateChannel(breakout.ID, []byte(html))
}

func (s breakoutService) renderTemplateToString(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	var tmpl = template.Must(template.ParseGlob("templates/**/*"))

	err := tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
