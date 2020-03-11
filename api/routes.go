package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jawr/monere/user"
)

type HandlerFunc func(user.User, *gin.Context) error

func (s Server) setupRoutes() {

	base := s.router.Group("/")
	if os.Getenv("MONERE_ENV") == "dev" {
		base = s.router.Group("/api")
	}

	// authentication
	base.POST("/register", s.handlePostRegister())
	base.POST("/login", s.handlePostLogin())
	base.GET("/logout", s.handleGetLogout())
	base.POST("verify/:code", s.handlePostVerify())

	// user routes
	user := base.Group("/user/")
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

	// job read
	user.GET("/jobs", s.handleUser(s.handleGetJobs()))
}
