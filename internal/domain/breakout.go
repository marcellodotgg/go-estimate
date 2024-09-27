package domain

import "sync"

type Breakouts struct {
	Breakout Breakout
	Mux      sync.Mutex
}

type Breakout struct {
	Users   map[string]User
	OwnerID string
}

type User struct {
	Name string
	Vote string
}
