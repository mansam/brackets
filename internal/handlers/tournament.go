package handlers

import (
	"brackets/internal/db"
	"brackets/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type TournamentHandler struct{}

type TournamentForm struct {
	Name       string `form:"name" json:"name" binding:"required,min=4"`
	MaxRounds  int    `form:"maxRounds" json:"maxRounds" binding:"required,min=3"`
	MaxPlayers int    `form:"maxPlayers" json:"maxPlayers" binding:"required,min=4"`
}

type PlayerForm struct {
	Name         string `form:"name" json:"name" binding:"required,alphanumunicode,min=3"`
	Email        string `form:"email" json:"email" binding:"required,email"`
	FactionID    int    `form:"factionid" json:"factionid" binding:"required"`
	SelfScorePin string `form:"pin" json:"pin" binding:"required"`
}

func (r *TournamentHandler) List(c *gin.Context) {
	var tournaments []models.Tournament
	result := db.DB.Find(&tournaments)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	list := make([]gin.H, 0)
	for _, t := range tournaments {
		list = append(list, t.JSON())
	}
	c.JSON(http.StatusOK, gin.H{
		"tournaments": list,
	})
}

func (r *TournamentHandler) Get(c *gin.Context) {
	t := models.Tournament{}
	tid := c.Param("id")
	result := db.DB.First(&t, "id = ?", tid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "not found",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	}
	c.JSON(http.StatusOK, t.JSON())
}

func (r *TournamentHandler) Create(c *gin.Context) {
	var form TournamentForm

	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organizer, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Validate form input
	if strings.Trim(form.Name, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name can't be empty"})
		return
	}

	result := db.DB.First(&models.Tournament{}, "name = ?", form.Name)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if result.Row() != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "name already in use",
		})
		return
	}

	t, err := models.NewTournament(form.Name, form.MaxRounds, form.MaxPlayers, organizer.(models.User))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't create tournament",
		})
		return
	}

	c.JSON(http.StatusOK, t.JSON())
}

func (r *TournamentHandler) AddPlayer(c *gin.Context) {
	var form PlayerForm
	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate form input
	if strings.Trim(form.Name, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name can't be empty"})
		return
	}

	result := db.DB.First(&models.Player{}, "email = ?", form.Email, "tournament_id = ?", tid)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if result.Row() != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "player already registered",
		})
		return
	}

	p, err := models.NewPlayer(form.Name, form.Email, form.SelfScorePin, form.FactionID, tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't add player",
		})
		return
	}

	c.JSON(http.StatusOK, p.JSON())
}

func (r *TournamentHandler) AdvanceRound(c *gin.Context) {

}
