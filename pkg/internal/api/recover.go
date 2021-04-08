package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jawr/whois-bi/pkg/internal/user"
)

func (s Server) handlePostRecover() gin.HandlerFunc {
	type Request struct {
		Email string `json:"email"`
	}

	return func(c *gin.Context) {
		var request Request

		if c.ShouldBind(&request) != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Missing Email"},
			)
			return
		}

		var u user.User
		if err := s.db.Model(&u).Where("email = ?", request.Email).Select(); err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Unknown email"},
			)
			return
		}

		rec := user.NewRecover(u)

		if err := rec.Insert(s.db); err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Recovery already in progress"},
			)
			return
		}

		body := fmt.Sprintf(
			`Account recover requested, if this was you please continue to <a href="https://%s/recover/%s">reset your password<a/>.\nIf you did not request a password reset, please ignore this email.`,
			os.Getenv("DOMAIN"),
			rec.Code,
		)

		if err := s.emailer.Send(u.Email, "Account Recovery", body); err != nil {
			log.Println(err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Unable to send recovery email"},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Recovery process started."})
	}
}

func (s Server) handlePostRecoverCode() gin.HandlerFunc {
	type Request struct {
		Code            string `json:"code"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"password_confirm"`
	}

	return func(c *gin.Context) {
		var request Request

		if c.ShouldBind(&request) != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Missing Code"},
			)
			return
		}

		// validate password
		if len(request.Password) == 0 || request.Password != request.PasswordConfirm {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid Password"},
			)
			return
		}

		// set new password
		newPassword, err := user.CreatePassword(request.Password)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			return
		}

		err = user.RecoverPassword(s.db, request.Code, newPassword)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Recovery complete. Please login with your new credentials!"})
	}
}
