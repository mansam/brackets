package models

import "gorm.io/gorm"

type Tournament struct {
	gorm.Model

	Name       string `gorm:"unique"`
	MaxRounds  int
	MaxPlayers int
	Players    []Player `gorm:"many2many:tournament_players;"`
	Pairings   []Pairing
}

func (r *Tournament) JSON() map[string]interface{} {
	return map[string]interface{}{
		"id":         r.ID,
		"name":       r.Name,
		"maxplayers": r.MaxPlayers,
		"maxrounds":  r.MaxRounds,
	}
}

type Pairing struct {
	gorm.Model

	Round   int
	PlayerA Player
	PlayerB Player
}

func (r *Pairing) JSON() map[string]interface{} {
	return map[string]interface{}{
		"id":      r.ID,
		"round":   r.Round,
		"playera": r.PlayerA.JSON(),
		"playerb": r.PlayerB.JSON(),
	}
}
