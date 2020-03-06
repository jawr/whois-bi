package api

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/jawr/monere/sender"
)

type Server struct {
	db      *pg.DB
	router  *gin.Engine
	emailer *sender.Sender
}

func NewServer(db *pg.DB, emailer *sender.Sender) *Server {
	router := gin.Default()

	router.Use(
		sessions.Sessions(
			"monere",
			sessions.NewCookieStore(
				[]byte("1kEetoDbop4$%3lSF,xvmBpekREK3#$"),
			),
		),
	)

	server := Server{
		db:      db,
		router:  router,
		emailer: emailer,
	}

	return &server
}

func (s *Server) Run(addr string) error {
	s.setupRoutes()
	return s.router.Run(addr)
}
