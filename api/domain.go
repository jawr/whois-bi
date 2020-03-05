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

func (s Server) handleGetDomain() HandlerFunc {
	type Response struct {
		domain.Domain

		Records []domain.Record

		Whois domain.Whois
	}

	return func(u user.User, c *gin.Context) error {
		domainName := c.Param("domain")

		var response Response

		err := s.db.Model(&response.Domain).Where("domain = ? AND owner_id = ?", domainName, u.ID).Select()
		if err != nil {
			return errors.Wrap(err, "Select Domain")
		}

		if response.OwnerID != u.ID {
			return errors.New("Not allowed")
		}

		err = s.db.Model(&response.Records).Where("domain_id = ?", response.ID).Select()
		if err != nil {
			return errors.Wrap(err, "Select Records")
		}

		err = s.db.Model(&response.Whois).Where("domain_id = ?", response.ID).Order("updated_date DESC").Limit(1).Select()
		if err != nil {
			return errors.Wrap(err, "Select Records")
		}

		c.JSON(http.StatusOK, &response)

		return nil
	}
}
