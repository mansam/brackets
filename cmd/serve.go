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
	tournamentHandler := handlers.TournamentHandler{}
	factionHandler := handlers.FactionHandler{}
	err := db.DB.AutoMigrate(&models.User{}, &models.Tournament{}, &models.Player{}, &models.Pairing{})
	if err != nil {
		panic(err)
	}

	r.GET("/users/:id", userHandler.Get)
	r.GET("/users", userHandler.List)
	r.POST("/users", userHandler.Create)
	r.POST("/login", authHandler.Login)
	r.GET("/tournaments/:id", tournamentHandler.Get)
	r.GET("/tournaments", tournamentHandler.List)
	r.GET("/factions", factionHandler.List)
	r.GET("/factions/:id", factionHandler.Get)

	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/login", "./web/login.html")

	loginRequired := r.Group("/")
	loginRequired.Use(middleware.AuthRequired)
	loginRequired.GET("/logout", authHandler.Logout)
	loginRequired.GET("/whoami", authHandler.WhoAmI)
	loginRequired.POST("/tournaments", tournamentHandler.Create)
	loginRequired.StaticFile("/tournament", "./web/tournament.html")

	organizerOnly := r.Group("/")
	organizerOnly.Use(middleware.AuthRequired)
	organizerOnly.Use(middleware.OrganizerOnly)
	organizerOnly.POST("/tournaments/:id/players", tournamentHandler.AddPlayer)

	adminOnly := r.Group("/admin")
	adminOnly.Use(middleware.AuthRequired)
	adminOnly.Use(middleware.AdminOnly)
	adminOnly.POST("/factions", factionHandler.Create)
	adminOnly.StaticFile("/faction", "./web/faction.html")

	err = r.Run()
	if err != nil {
		panic(err)
	}
}
