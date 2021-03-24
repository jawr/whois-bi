package api

import (
	"os"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/jawr/whois-bi/pkg/internal/sender"
)

type Server struct {
	db      *pg.DB
	router  *gin.Engine
	emailer *sender.Sender
}

func NewServer(db *pg.DB, emailer *sender.Sender) *Server {
	router := gin.Default()

	store := cookie.NewStore([]byte(os.Getenv("HTTP_COOKIE_SECRET")))

	if os.Getenv("MODE") == "dev" {
		opts := sessions.Options{
			Secure: false,
		}
		store.Options(opts)
	}

	router.Use(sessions.Sessions(os.Getenv("HTTP_SESSION_ID"), store))

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
