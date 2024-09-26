package service

import (
	"bytes"
	"text/template"

	"github.com/gomarchy/estimate/internal/domain"
	"github.com/gomarchy/estimate/internal/infrastructure/websocket"
)

type BreakoutService interface {
	Add(id string)
	Remove(id string)
}

type breakoutService struct {
	channels map[string]*domain.Breakouts
}

var tmpl = template.Must(template.ParseGlob("templates/**/*"))

func NewBreakoutService() BreakoutService {
	return breakoutService{
		channels: make(map[string]*domain.Breakouts),
	}
}

func (s breakoutService) Add(channel string) {
	_, exists := s.channels[channel]

	if !exists {
		s.channels[channel] = &domain.Breakouts{}
	}

	s.channels[channel].Mux.Lock()
	defer s.channels[channel].Mux.Unlock()
	s.channels[channel].Breakout.UsersCount++
	s.UpdateChannel(channel)
}

func (s breakoutService) Remove(channel string) {
	s.channels[channel].Mux.Lock()
	defer s.channels[channel].Mux.Unlock()
	s.channels[channel].Breakout.UsersCount--
	s.UpdateChannel(channel)
}

func (s breakoutService) UpdateChannel(channel string) {
	html, err := s.renderTemplateToString("breakout/sample", s.channels[channel].Breakout)
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
