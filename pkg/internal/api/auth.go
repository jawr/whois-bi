package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"golang.org/x/crypto/bcrypt"
)

const (
	SessionUserKey string = "monere-user-id"
)

func (s Server) handleGetStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Logged in"})
	}
}

func (s Server) handleGetLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(SessionUserKey)
		if userID == nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid session"},
			)
			return
		}

		session.Delete(SessionUserKey)
		if err := session.Save(); err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Failed to save sesion"},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Logged out"})
	}
}

func (s Server) handlePostRegister() gin.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		var request Request

		if c.ShouldBind(&request) != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Missing Email and/or Password"},
			)
			return
		}

		u, err := user.NewUser(request.Email, request.Password)
		if err != nil {
			log.Println(err)
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Unable to create User"},
			)
			return
		}

		if err := u.Insert(s.db); err != nil {
			log.Println(err)
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Unable to create User"},
			)
			return
		}

		body := fmt.Sprintf(
			`Thank you for registering with us. Please complete your registration by clicking <a href="https://%s/verify/%s">here<a/>`,
			os.Getenv("DOMAIN"),
			u.VerifiedCode,
		)

		if err := s.emailer.Send(u.Email, "Please verify your account", body); err != nil {
			log.Println(err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Unable to send verification email"},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Registration complete. Verification sent."})
	}
}

func (s Server) handlePostLogin() gin.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		var request Request

		if c.ShouldBind(&request) != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Missing Email and/or Password"},
			)
			return
		}

		// validate user
		var u user.User

		if err := s.db.Model(&u).Where("email = ? AND verified_at IS NOT NULL", request.Email).Select(); err != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Email or Password is invalid."},
			)
			return
		}

		// validate password
		if err := bcrypt.CompareHashAndPassword(u.Password, []byte(request.Password)); err != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Email or Password is invalid."},
			)
			return
		}

		// successful login! run database updates
		go func(u user.User) {
			// can handle ip logging if we want
			_, err := s.db.Model(&u).Set("last_login_at = now()").WherePK().Update()
			if err != nil {
				panic(err)
			}
		}(u)

		// set and save session
		session := sessions.Default(c)
		session.Set(SessionUserKey, u.ID)
		if err := session.Save(); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to save session"},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Logged in"})
	}
}

func (s Server) handlePostVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := user.VerifyUser(s.db, c.Param("code")); err != nil {
			log.Printf("Verify error: %s", err)
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Failed to verify using that code"},
			)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Verified. Please login."})
	}
}
