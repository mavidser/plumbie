package models

type Session struct {
	Key    string `xorm:"pk char(16)"`
	Data   string `xorm:"blob"`
	Expiry uint   `xorm:"not null"`
}
