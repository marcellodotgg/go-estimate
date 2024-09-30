package websocket

import (
	"github.com/olahol/melody"
)

var Manager = melody.New()

func UpdateChannel(channel string, json []byte) {
	Manager.BroadcastFilter(json, func(q *melody.Session) bool {
		qChannel, _ := q.Get("channel")
		return qChannel == channel
	})
}
