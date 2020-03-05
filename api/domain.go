package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jawr/monere/domain"
	"github.com/jawr/monere/user"
	"github.com/pkg/errors"
)

func (s Server) handleGetDomains() HandlerFunc {
	return func(u user.User, c *gin.Context) error {
		var domains []domain.Domain

		err := s.db.Model(&domains).Where("owner_id = ?", u.ID).Select()
		if err != nil {
			return errors.Wrap(err, "Select")
		}

		c.JSON(http.StatusOK, &domains)

		return nil
	}
}
