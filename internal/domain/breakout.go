package domain

type Breakout struct {
	Audit
	CreatedBy string `gorm:"not null"`
	Users     []User
}
