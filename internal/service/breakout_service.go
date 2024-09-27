package service

import (
	"bytes"
	"text/template"

	"github.com/gomarchy/estimate/internal/domain"
	"github.com/gomarchy/estimate/internal/infrastructure/websocket"
)

var (
	Channels = make(map[string]*domain.Breakouts)
)

type BreakoutService interface {
	AddUser(channel, userID string)
	RemoveUser(channel, userID string)
	Create(channel, userID string) error
}

type breakoutService struct {
}

var tmpl = template.Must(template.ParseGlob("templates/**/*"))

func NewBreakoutService() BreakoutService {
	return breakoutService{}
}

func (s breakoutService) Create(channel, userID string) error {
	Channels[channel] = &domain.Breakouts{
		Breakout: domain.Breakout{
			Users:   make(map[string]domain.User),
			OwnerID: userID,
		},
	}
	return nil
}

func (s breakoutService) AddUser(channel, userID string) {
	Channels[channel].Mux.Lock()
	defer Channels[channel].Mux.Unlock()
	Channels[channel].Breakout.Users[userID] = domain.User{
		Name: "Guest",
		Vote: "",
	}
	s.updateChannel(channel)
}

func (s breakoutService) RemoveUser(channel, userID string) {
	Channels[channel].Mux.Lock()
	defer Channels[channel].Mux.Unlock()
	delete(Channels[channel].Breakout.Users, userID)
	s.updateChannel(channel)
}

func (s breakoutService) updateChannel(channel string) {
	html, _ := s.renderTemplateToString("breakout/sample", Channels[channel].Breakout)
	websocket.UpdateChannel(channel, []byte(html))
}

func (s breakoutService) renderTemplateToString(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer

	err := tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
