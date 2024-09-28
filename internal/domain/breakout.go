package domain

type Breakout struct {
	Audit
	CreatedBy   string `gorm:"not null"`
	Connections []Connection
	ShowVotes   bool   `gorm:"type:int; default:0"`
	Cards       []Card `gorm:"-"`
}

type Card struct {
	DisplayValue string
	Value        int8
}
