package models

import (
	"brackets/internal/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Tournament struct {
	gorm.Model

	Name        string `gorm:"unique"`
	MaxRounds   int
	MaxPlayers  int
	OrganizerID int
	Organizer   User
	Players     []Player `gorm:"many2many:tournament_players;"`
	Pairings    []Pairing
}

func (r *Tournament) JSON() map[string]interface{} {
	players := make([]gin.H, 0)
	for _, p := range r.Players {
		players = append(players, p.JSON())
	}

	return map[string]interface{}{
		"id":         r.ID,
		"name":       r.Name,
		"maxplayers": r.MaxPlayers,
		"maxrounds":  r.MaxRounds,
		"players":    players,
	}
}

func NewTournament(name string, maxRounds int, maxPlayers int, organizer User) (*Tournament, error) {
	t := Tournament{
		Name:       name,
		MaxRounds:  maxRounds,
		MaxPlayers: maxPlayers,
		Organizer:  organizer,
	}

	result := db.DB.Create(&t)
	if result.Error != nil {
		return nil, result.Error
	}
	return &t, nil
}

type Pairing struct {
	gorm.Model

	Round        int
	PlayerA      Player
	PlayerAID    int
	ScoreA       int
	PlayerB      Player
	PlayerBID    int
	ScoreB       int
	TournamentID int
}

func (r *Pairing) JSON() map[string]interface{} {
	return map[string]interface{}{
		"id":      r.ID,
		"round":   r.Round,
		"playera": r.PlayerA.JSON(),
		"scorea":  r.ScoreA,
		"playerb": r.PlayerB.JSON(),
		"scoreb":  r.ScoreB,
	}
}
