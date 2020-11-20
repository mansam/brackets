package main

import (
	"brackets/internal/db"
	"brackets/internal/handlers"
	"brackets/internal/middleware"
	"brackets/internal/models"
	"brackets/internal/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	db.Init()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	cookieStore := cookie.NewStore(util.GetRandomBytesOrDie(64), util.GetRandomBytesOrDie(32))
	cookieStore.Options(sessions.Options{
		Path:     "",
		Domain:   "",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	r.Use(sessions.Sessions("session", cookieStore))

	userHandler := handlers.UserHandler{}
	authHandler := handlers.AuthHandler{}
	err := db.DB.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}

	r.GET("/users/:id", userHandler.Get)
	r.GET("/users", userHandler.List)
	r.POST("/users", userHandler.Post)
	r.POST("/login", authHandler.Login)
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/login", "./web/login.html")

	loginRequired := r.Group("/")
	loginRequired.Use(middleware.AuthRequired)
	loginRequired.GET("/logout", authHandler.Logout)
	loginRequired.GET("/whoami", authHandler.WhoAmI)
	err = r.Run()
	if err != nil {
		panic(err)
	}
}
