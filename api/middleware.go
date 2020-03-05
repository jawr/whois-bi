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

		handleError(fn)(u, c)
	}
}

type apiError struct {
	statusCode int
	friendly   string
	err        error
}

func newApiError(statusCode int, friendly string, err error) error {
	return &apiError{
		statusCode: statusCode,
		friendly:   friendly,
		err:        err,
	}
}

func (e *apiError) Error() string {
	return e.err.Error()
}

func handleError(fn HandlerFunc) HandlerFunc {
	// depending on the user we want to be able to
	// offer different errors
	return func(u user.User, c *gin.Context) error {
		if err := fn(u, c); err != nil {

			if aerr, ok := err.(*apiError); ok {
				log.Printf("API Error: %d - %s: %s", aerr.statusCode, aerr.friendly, aerr.err)
				c.JSON(
					aerr.statusCode,
					gin.H{"Error": aerr.friendly},
				)
			} else {
				log.Printf("Error: %s", err)
				c.JSON(
					http.StatusInternalServerError,
					gin.H{"Error": err.Error()},
				)
			}
		}

		return nil
	}
}
