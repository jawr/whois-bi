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

	// user routes
	user := s.router.Group("/user/")
	user.Use(handleAuth)

	user.GET("/status", s.handleGetStatus())
	user.GET("/domains", s.handleUser(s.handleGetDomains()))
	user.GET("/domain/:domain", s.handleDomain(s.handleGetDomain()))
	user.GET("/domain/:domain/records", s.handleDomain(s.handleGetDomainRecords()))
	user.GET("/domain/:domain/whois", s.handleDomain(s.handleGetDomainWhois()))
}
