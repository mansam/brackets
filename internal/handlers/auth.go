package handlers

import (
	"brackets/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const CookieKeyUser = "user"

type AuthHandler struct{}

func (r *AuthHandler) Login(c *gin.Context) {
	session := sessions.Default(c)
	name := c.PostForm("name")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(name, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	userid, valid, err := models.CheckPassword(name, password)
	if !valid || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	session.Set(CookieKeyUser, userid)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session."})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged in"})
}

func (r *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(CookieKeyUser)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(CookieKeyUser)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func (r *AuthHandler) WhoAmI(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(CookieKeyUser)
	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"status": "logged in",
	})
}
