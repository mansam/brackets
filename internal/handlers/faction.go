package handlers

import (
	"brackets/internal/db"
	"brackets/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type FactionHandler struct{}

type FactionForm struct {
	Name string `form:"name" json:"name" binding:"required,min=1"`
}

func (r *FactionHandler) List(c *gin.Context) {
	var factions []models.Faction
	result := db.DB.Find(&factions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	list := make([]gin.H, 0)
	for _, f := range factions {
		list = append(list, f.JSON())
	}
	c.JSON(http.StatusOK, gin.H{
		"factions": list,
	})
}

func (r *FactionHandler) Get(c *gin.Context) {
	f := models.Faction{}
	fid := c.Param("id")
	result := db.DB.First(&f, "id = ?", fid)
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
	c.JSON(http.StatusOK, f.JSON())
}

func (r *FactionHandler) Create(c *gin.Context) {
	var form FactionForm

	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate form input
	if strings.Trim(form.Name, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name can't be empty"})
		return
	}

	result := db.DB.First(&models.Faction{}, "name = ?", form.Name)
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

	f, err := models.NewFaction(form.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't create tournament",
		})
		return
	}

	c.JSON(http.StatusOK, f.JSON())
}
