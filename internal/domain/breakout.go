package domain

import "sync"

type Breakouts struct {
	Breakout Breakout
	Mux      sync.Mutex
}

type Breakout struct {
	UsersCount int
}
