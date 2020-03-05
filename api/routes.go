package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jawr/monere/user"
)

type HandlerFunc func(user.User, *gin.Context) error

func (s Server) setupRoutes() {
	// authentication
	s.router.POST("/login", s.handlePostLogin())
	s.router.GET("/logout", s.handleGetLogout())

	auth := s.router.Group("/auth/")
	auth.Use(handleAuth)

	// authentication required
	auth.GET("/status", s.handleGetStatus())
	auth.GET("/domains", s.handleUser(handleError(s.handleGetDomains())))
}
