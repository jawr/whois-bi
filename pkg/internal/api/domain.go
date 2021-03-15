package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

func (s Server) handleGetDomains() HandlerFunc {
	return func(u user.User, c *gin.Context) error {
		domains := make([]domain.DisplayDomain, 0)

		err := s.db.Model(&domains).
			ColumnExpr("domain.*, coalesce(count(distinct record.id), 0) as records, coalesce(count(distinct whois.id), 0) as whois").
			Join("left join records as record on domain.id = record.domain_id left join whois on domain.id = whois.domain_id").
			Where("domain.owner_id = ?", u.ID).
			Group("domain.id").
			Order("domain.domain").
			Select()
		if err != nil {
			return newApiError(http.StatusNotFound, "Not found", errors.Wrap(err, "Select"))
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
			return newApiError(http.StatusNotFound, "Not found", errors.Wrap(err, "Select Domain"))
		}

		if d.OwnerID != u.ID {
			return newApiError(http.StatusNotFound, "Not found", errors.New("Not allowed"))
		}

		return fn(d, u, c)
	})
}

func (s Server) handleGetDomain() DomainHandlerFunc {
	return func(d domain.Domain, u user.User, c *gin.Context) error {

		c.JSON(http.StatusOK, &d)

		return nil
	}
}

func (s Server) handleGetDomainRecords() DomainHandlerFunc {
	return func(d domain.Domain, u user.User, c *gin.Context) error {
		records := make([]domain.Record, 0)
		err := s.db.Model(&records).Where("domain_id = ?", d.ID).Order("added_at DESC").Select()
		if err != nil {
			return newApiError(http.StatusNotFound, "Not found", errors.Wrap(err, "Select"))
		}
		c.JSON(http.StatusOK, &records)
		return nil
	}
}

func (s Server) handleGetDomainWhois() DomainHandlerFunc {
	return func(d domain.Domain, u user.User, c *gin.Context) error {
		whois := make([]domain.Whois, 0)
		err := s.db.Model(&whois).Where("domain_id = ?", d.ID).Order("added_at DESC").Select()
		if err != nil {
			return newApiError(http.StatusNotFound, "Not found", errors.Wrap(err, "Select"))
		}
		c.JSON(http.StatusOK, &whois)
		return nil
	}
}

func (s Server) handlePostDomain() HandlerFunc {
	type Request struct {
		Domain string
	}

	return func(u user.User, c *gin.Context) error {
		var request Request
		if err := c.ShouldBind(&request); err != nil {
			return newApiError(http.StatusBadRequest, "Bad Request", errors.Wrap(err, "ShouldBind"))
		}

		d := domain.NewDomain(request.Domain, u)

		if len(d.Domain) == 0 || !strings.Contains(d.Domain, ".") {
			return newApiError(
				http.StatusBadRequest,
				fmt.Sprintf("Invalid Domain: '%s'", d.Domain),
				errors.Errorf("Invalid Domain: '%s'", d.Domain),
			)
		}

		_, err := s.db.Model(&d).
			OnConflict("(domain, owner_id) DO UPDATE").
			Set("deleted_at = NULL").
			Returning("*").
			Insert()
		if err != nil {
			return newApiError(http.StatusInternalServerError, "Inserting Domain", errors.Wrap(err, "Insert"))
		}

		dd := domain.DisplayDomain{Domain: d}
		c.JSON(http.StatusCreated, &dd)

		return nil
	}
}

func (s Server) handlePostRecord() DomainHandlerFunc {
	type Request struct {
		Raw string
	}

	type Response struct {
		Records []domain.Record `json:"records"`
		Errors  []string        `json:"errors"`
	}

	return func(d domain.Domain, u user.User, c *gin.Context) error {
		var request Request
		if c.ShouldBind(&request) != nil {
			return newApiError(http.StatusBadRequest, "Bad Request", errors.New("ShouldBind"))
		}

		zp := dns.NewZoneParser(strings.NewReader(request.Raw), "", "")

		response := Response{
			Records: make([]domain.Record, 0),
			Errors:  make([]string, 0),
		}

		for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
			record := domain.NewRecord(d, rr, domain.RecordSourceManual)

			if err := record.Insert(s.db); err != nil {
				// return newApiError(http.StatusInternalServerError, "Unable to save record", errors.Wrap(err, "Insert"))
				response.Errors = append(response.Errors, fmt.Sprintf("Unable to add '%s'", record.Raw))
				continue
			}

			response.Records = append(response.Records, record)
		}

		if err := zp.Err(); err != nil {
			return newApiError(http.StatusInternalServerError, "Unable to parse record(s)", errors.Wrap(err, "zoneParser Err"))
		}

		c.JSON(http.StatusCreated, &response)

		return nil
	}
}

func (s Server) handleDeleteDomain() DomainHandlerFunc {
	type Response struct{}
	return func(d domain.Domain, u user.User, c *gin.Context) error {

		response := Response{}

		_, err := s.db.Model(&d).
			Where("domain.id = ?", d.ID).
			Delete()
		if err != nil {
			return newApiError(http.StatusInternalServerError, "Unable to delete domain", errors.Wrap(err, "Select Whois"))
		}

		c.JSON(http.StatusOK, &response)

		return nil
	}
}
