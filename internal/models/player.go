package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model

	Name         string
	Email        string `gorm:"unique"`
	SelfScorePin []byte
	Tournaments  []Tournament `gorm:"many2many:tournament_players;"`
}

func (r *Player) JSON() map[string]interface{} {
	return map[string]interface{}{
		"id":   r.ID,
		"name": r.Name,
	}
}
