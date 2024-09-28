package service

import (
	"github.com/gomarchy/estimate/internal/domain"
	"github.com/gomarchy/estimate/internal/infrastructure/database"
)

type VoteService interface {
	Vote(conection *domain.Connection, value string)
	ShowVotes(breakout *domain.Breakout)
	Reset(breakout *domain.Breakout)
}

type voteService struct {
	broadcast BroadcastService
}

func NewVoteService() VoteService {
	return voteService{
		broadcast: NewBroadcastService(),
	}
}

func (v voteService) Reset(breakout *domain.Breakout) {
	database.DB.Model(domain.Connection{}).Where("breakout_id = ?", breakout.ID).Update("vote", "")
	database.DB.Model(domain.Breakout{}).Where("id = ?", breakout.ID).Update("show_votes", false)

	v.broadcast.ResetVotes(breakout.ID)
	v.broadcast.Breakout(breakout.ID)
}

func (v voteService) ShowVotes(breakout *domain.Breakout) {
	database.DB.Model(domain.Breakout{}).Where("id = ?", breakout.ID).Update("show_votes", true)

	v.broadcast.ResetVotes(breakout.ID)
	v.broadcast.Breakout(breakout.ID)
}

func (v voteService) Vote(connection *domain.Connection, value string) {
	if value == connection.Vote {
		value = ""
	}
	database.DB.Model(connection).Update("vote", value)
	database.DB.First(connection)

	var breakout domain.Breakout
	database.DB.Preload("Connections", "is_online = ?", true).First(&breakout, "id = ?", connection.BreakoutID)

	if v.didEveryoneVote(breakout) {
		v.ShowVotes(&breakout)
		return
	}

	v.broadcast.Breakout(breakout.ID)
}

func (v voteService) didEveryoneVote(breakout domain.Breakout) bool {
	for _, connection := range breakout.Connections {
		if connection.Vote == "" {
			return false
		}
	}
	return true
}
