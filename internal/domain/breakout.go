package domain

type Breakout struct {
	Audit
	CreatedBy string `gorm:"not null"`
	Users     []User
	ShowVotes bool   `gorm:"type:int; default:0"`
	Cards     []Card `gorm:"-"`
}

type Card struct {
	Value string
}
