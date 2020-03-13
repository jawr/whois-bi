package api

import (
	"os"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	store := cookie.NewStore([]byte("1kEetoDbop4$%3lSF,xvmBpekREK3#$"))

	if os.Getenv("MONERE_ENV") == "dev" {
		opts := sessions.Options{
			Secure: false,
		}
		store.Options(opts)
	}

	router.Use(sessions.Sessions("monere", store))

	server := Server{
		db:      db,
		router:  router,
		emailer: emailer,
	}

	return &server
}

func (s *Server) Run(addr string) error {
	s.setupRoutes()
	endless.DefaultHammerTime = time.Second * 20
	return endless.ListenAndServe(addr, s.router)
}
