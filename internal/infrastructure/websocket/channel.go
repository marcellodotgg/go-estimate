package websocket

import (
	"sync"

	"github.com/olahol/melody"
)

var Manager = melody.New()

type ChannelCounts struct {
	Counts map[string]int
	mux    sync.Mutex
}

// Increment the connection count for a channel
func (cc *ChannelCounts) Increment(channel string) {
	cc.mux.Lock()
	defer cc.mux.Unlock()
	cc.Counts[channel]++
}

// Decrement the connection count for a channel
func (cc *ChannelCounts) Decrement(channel string) {
	cc.mux.Lock()
	defer cc.mux.Unlock()
	cc.Counts[channel]--
}

// Get the connection count for a channel
func (cc *ChannelCounts) GetCount(channel string) int {
	cc.mux.Lock()
	defer cc.mux.Unlock()
	return cc.Counts[channel]
}

func UpdateChannel(channel string, json []byte) {
	Manager.BroadcastFilter(json, func(q *melody.Session) bool {
		qChannel, _ := q.Get("channel")
		return qChannel == channel
	})
}
