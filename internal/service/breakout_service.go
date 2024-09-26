package service

import (
	"bytes"
	"text/template"

	"github.com/gomarchy/estimate/internal/domain"
	"github.com/gomarchy/estimate/internal/infrastructure/websocket"
)

var (
	channels = make(map[string]*domain.Breakouts)
)

type BreakoutService interface {
	JoinAs(channel, userID, displayName string)
	RemoveUser(channel, userID string)
}

type breakoutService struct {
}

var tmpl = template.Must(template.ParseGlob("templates/**/*"))

func NewBreakoutService() BreakoutService {
	channels["123"] = &domain.Breakouts{
		Breakout: domain.Breakout{
			Users: make(map[string]domain.User),
		},
	}
	return breakoutService{}
}

func (s breakoutService) JoinAs(channel, userID, displayName string) {
	channels[channel].Mux.Lock()
	defer channels[channel].Mux.Unlock()
	channels[channel].Breakout.Users[userID] = domain.User{
		Name: displayName,
		Vote: "",
	}
	s.updateChannel(channel)
}

func (s breakoutService) RemoveUser(channel, userID string) {
	channels[channel].Mux.Lock()
	defer channels[channel].Mux.Unlock()
	delete(channels[channel].Breakout.Users, userID)
	s.updateChannel(channel)
}

func (s breakoutService) updateChannel(channel string) {
	html, err := s.renderTemplateToString("breakout/sample", channels[channel].Breakout)
	if err != nil {
		return
	}
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
