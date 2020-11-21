package models

import (
	"brackets/internal/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model

	Name         string
	Email        string
	SelfScorePin []byte
	Faction      Faction
	FactionID    int
	Tournament   Tournament
	TournamentID int
}

func (r *Player) JSON() map[string]interface{} {
	return map[string]interface{}{
		"id":         r.ID,
		"name":       r.Name,
		"faction":    r.Faction.Name,
		"tournament": r.Tournament.Name,
	}
}

func NewPlayer(name, email, selfScorePin string, factionID, tournamentID int) (*Player, error) {
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(selfScorePin), BcryptCost)
	if err != nil {
		return nil, err
	}

	p := Player{
		Name:         name,
		Email:        email,
		SelfScorePin: hashedPin,
		FactionID:    factionID,
		TournamentID: tournamentID,
	}

	result := db.DB.Create(&p)
	if result.Error != nil {
		return nil, result.Error
	}
	return &p, nil
}