package middleware

import (
	"brackets/internal/db"
	"brackets/internal/models"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("user")
	if uid == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	u := models.User{}
	result := db.DB.First(&u, "id = ?", uid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			session.Clear()
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			session.Clear()
			return
		}
	}

	c.Set("user", u)
	// Continue down the chain to handler etc
	c.Next()
}

func AdminOnly(c *gin.Context) {
	u, ok := c.Get("user")
	if !ok || u.(models.User).Name != "admin" {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "must be admin"})
		return
	}

	c.Next()
}

func OrganizerOnly(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("user")
	if uid == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	t := models.Tournament{}
	result := db.DB.First(&t, "id = ?", c.Param("id"))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "tournament not found",
			})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	if uid != t.OrganizerID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "tournament does not belong to you",
		})
		return
	}

	// Continue down the chain to handler etc
	c.Next()
}
