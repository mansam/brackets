package handlers

import (
	"brackets/internal/db"
	"brackets/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type TournamentHandler struct{}

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
		"users": list,
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

}

func (r *TournamentHandler) AdvanceRound(c *gin.Context) {

}
