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

type DomainHandlerFunc func(domain.Domain, user.User, *gin.Context) error

func (s Server) handleDomain(fn DomainHandlerFunc) gin.HandlerFunc {
	return s.handleUser(func(u user.User, c *gin.Context) error {
		var d domain.Domain

		err := s.db.Model(&d).
			Where(
				"domain = ? AND owner_id = ?",
				c.Param("domain"),
				u.ID,
			).
			Select()
		if err != nil {
			return errors.Wrap(err, "Select Domain")
		}

		if d.OwnerID != u.ID {
			return errors.New("Not allowed")
		}

		return fn(d, u, c)
	})
}

func (s Server) handleGetDomain() DomainHandlerFunc {
	type Response struct {
		domain.Domain

		Records []domain.Record

		Whois domain.Whois
	}

	return func(d domain.Domain, u user.User, c *gin.Context) error {

		response := Response{
			Domain: d,
		}

		err := s.db.Model(&response.Records).
			Where("domain_id = ? AND removed_at IS NULL", response.ID).
			Order("added_at DESC").
			Select()
		if err != nil {
			return errors.Wrap(err, "Select Records")
		}

		err = s.db.Model(&response.Whois).
			Where("domain_id = ?", response.ID).
			Order("updated_date DESC").
			Limit(1).
			Select()
		if err != nil {
			return errors.Wrap(err, "Select Records")
		}

		c.JSON(http.StatusOK, &response)

		return nil
	}
}

func (s Server) handleGetDomainRecords() DomainHandlerFunc {
	return func(d domain.Domain, u user.User, c *gin.Context) error {
		var records []domain.Record
		err := s.db.Model(&records).Where("domain_id = ?", d.ID).Order("added_at DESC").Select()
		if err != nil {
			return errors.Wrap(err, "Select Records")
		}
		c.JSON(http.StatusOK, &records)
		return nil
	}
}

func (s Server) handleGetDomainWhois() DomainHandlerFunc {
	return func(d domain.Domain, u user.User, c *gin.Context) error {
		var whois []domain.Whois
		err := s.db.Model(&whois).Where("domain_id = ?", d.ID).Order("added_at DESC").Select()
		if err != nil {
			return errors.Wrap(err, "Select Whois")
		}
		c.JSON(http.StatusOK, &whois)
		return nil
	}
}
