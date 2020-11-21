package models

import (
	"brackets/internal/db"
	"gorm.io/gorm"
)

type Faction struct {
	gorm.Model

	Name    string
	Players []Player
}

func (r *Faction) JSON() map[string]interface{} {
	return map[string]interface{}{
		"id":   r.ID,
		"name": r.Name,
	}
}

func NewFaction(name string) (*Faction, error) {
	f := Faction{Name: name}

	result := db.DB.Create(&f)
	if result.Error != nil {
		return nil, result.Error
	}
	return &f, nil
}