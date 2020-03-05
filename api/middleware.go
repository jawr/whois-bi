package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jawr/monere/user"
)

func handleAuth(c *gin.Context) {
	userID := sessions.Default(c).Get(SessionUserKey)

	if userID == nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"Error": "Not Authorized"},
		)
		return
	}

	c.Next()
}

func (s Server) handleUser(fn HandlerFunc) gin.HandlerFunc {
	// user cache
	return func(c *gin.Context) {
		userID := sessions.Default(c).Get(SessionUserKey)

		// use a pool for user objects
		var u user.User
		if err := s.db.Model(&u).Where("id = ?", userID).Select(); err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"Error": "Not Authorized"},
			)
			return
		}

		fn(u, c)
	}
}

func handleError(fn HandlerFunc) HandlerFunc {
	// depending on the user we want to be able to
	// offer different errors
	return func(u user.User, c *gin.Context) error {
		if err := fn(u, c); err != nil {
			log.Printf("Error: %s", err)

			c.JSON(
				http.StatusInternalServerError,
				gin.H{"Error": err.Error()},
			)
		}

		return nil
	}
}
