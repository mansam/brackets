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

type UserHandler struct{}

func (r *UserHandler) List(c *gin.Context) {
	var users []models.User
	result := db.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	list := make([]gin.H, 0)
	for _, u := range users {
		list = append(list, u.JSON())
	}
	c.JSON(http.StatusOK, gin.H{
		"users": list,
	})
}

func (r *UserHandler) Post(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(name, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	result := db.DB.First(&models.User{}, "name = ?", name)
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

	u, err := models.NewUser(name, email, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, u.JSON())
}

func (r *UserHandler) Get(c *gin.Context) {
	u := models.User{}
	uid := c.Param("id")
	result := db.DB.First(&u, "id = ?", uid)
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
	c.JSON(http.StatusOK, u.JSON())
}
