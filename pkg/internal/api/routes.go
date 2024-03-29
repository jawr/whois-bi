package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jawr/whois-bi/pkg/internal/user"
)

type HandlerFunc func(user.User, *gin.Context) error

func (s Server) setupRoutes() {

	base := s.router.Group("/api")

	// authentication
	base.POST("/register", s.handlePostRegister())
	base.POST("/login", s.handlePostLogin())
	base.GET("/logout", s.handleGetLogout())
	base.POST("/verify/:code", s.handlePostVerify())

	base.POST("/recover", s.handlePostRecover())
	base.POST("/recover/code", s.handlePostRecoverCode())

	// user routes
	user := base.Group("/user/")
	user.Use(handleAuth)

	user.GET("/status", s.handleGetStatus())

	// domain read
	user.GET("/domains", s.handleUser(s.handleGetDomains()))
	user.GET("/domain/:domain", s.handleDomain(s.handleGetDomain()))
	user.DELETE("/domain/:domain", s.handleDomain(s.handleDeleteDomain()))
	user.GET("/domain/:domain/records", s.handleDomain(s.handleGetDomainRecords()))
	user.GET("/domain/:domain/whois", s.handleDomain(s.handleGetDomainWhois()))
	user.PUT("/domain/:domain/batch", s.handleDomain(s.handlePutDomainBatch()))

	// lists
	user.GET("/lists", s.handleUser(s.handleGetMatches()))
	user.POST("/lists", s.handleUser(s.handlePostList()))
	user.DELETE("/lists/:id", s.handleUser(s.handleDeleteList()))

	// domain create
	user.POST("/domain", s.handleUser(s.handlePostDomain()))
	user.POST("/domain/:domain/record", s.handleDomain(s.handlePostRecord()))

	// job read
	user.GET("/jobs/:domain", s.handleUser(s.handleGetJobs()))
	user.POST("/jobs/:domain", s.handleUser(s.handlePostJob()))
}
