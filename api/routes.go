package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jawr/monere/user"
)

type HandlerFunc func(user.User, *gin.Context) error

func (s Server) setupRoutes() {
	// authentication
	s.router.POST("/register", s.handlePostRegister())
	s.router.POST("/login", s.handlePostLogin())
	s.router.GET("/logout", s.handleGetLogout())
	s.router.POST("verify/:code", s.handlePostVerify())

	// user routes
	user := s.router.Group("/user/")
	user.Use(handleAuth)

	user.GET("/status", s.handleGetStatus())

	// domain read
	user.GET("/domains", s.handleUser(s.handleGetDomains()))
	user.GET("/domain/:domain", s.handleDomain(s.handleGetDomain()))
	user.GET("/domain/:domain/records", s.handleDomain(s.handleGetDomainRecords()))
	user.GET("/domain/:domain/whois", s.handleDomain(s.handleGetDomainWhois()))

	// domain create
	user.POST("/domain", s.handleUser(s.handlePostDomain()))
	user.POST("/domain/:domain/record", s.handleDomain(s.handlePostRecord()))
}
