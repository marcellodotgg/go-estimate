package service

import (
	"github.com/gomarchy/estimate/internal/domain"
	"github.com/gomarchy/estimate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type BreakoutService interface {
	AddUser(breakoutID, userID string)
	RemoveUser(breakoutID, userID string)
	UpdateDisplayName(connection domain.Connection) error
	Create(userID string) (domain.Breakout, error)
	FindByID(breakoutID string) (domain.Breakout, error)
}

type breakoutService struct {
	broadcast BroadcastService
}

func NewBreakoutService() BreakoutService {
	return breakoutService{
		broadcast: NewBroadcastService(),
	}
}

func (s breakoutService) FindByID(breakoutID string) (domain.Breakout, error) {
	var breakout domain.Breakout
	if err := database.DB.First(&breakout, "id = ?", breakoutID).Error; err != nil {
		return breakout, err
	}

	breakout.Cards = FibonacciCards()

	return breakout, nil
}

func (s breakoutService) UpdateDisplayName(connection domain.Connection) error {
	if err := database.DB.Model(&connection).Update("display_name", connection.DisplayName).Error; err != nil {
		return err
	}
	s.broadcast.Breakout(connection.BreakoutID)
	return nil
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
	if err := database.DB.First(&domain.Connection{}, "user_id = ? AND breakout_id = ?", userID, breakoutID).Error; err == nil {
		database.DB.Where("breakout_id = ? AND user_id = ?", breakoutID, userID).Model(&domain.Connection{}).Update("is_online", true)
		s.broadcast.Breakout(breakoutID)
		return
	}

	user := domain.Connection{
		DisplayName: "Guest",
		UserID:      userID,
		BreakoutID:  breakoutID,
		Vote:        "",
		IsOnline:    true,
	}

	database.DB.Create(&user)
	s.broadcast.Breakout(breakoutID)
}

func (s breakoutService) RemoveUser(breakoutID, userID string) {
	database.DB.Where("breakout_id = ? AND user_id = ?", breakoutID, userID).Model(&domain.Connection{}).Update("is_online", false)
	s.broadcast.Breakout(breakoutID)
}
