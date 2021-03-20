package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jawr/whois-bi/pkg/internal/list"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/pkg/errors"
)

func (s *Server) handleGetMatches() HandlerFunc {
	type response struct {
		Whitelists []list.List `json:"whitelists"`
		Blacklists []list.List `json:"blacklists"`
	}

	return func(u user.User, c *gin.Context) error {
		whitelists := make([]list.List, 0)
		blacklists := make([]list.List, 0)

		err := s.db.Model(&whitelists).Where("owner_id = ? AND list_type = ?", u.ID, list.Whitelist).Select()
		if err != nil {
			return newApiError(http.StatusInternalServerError, "Internal Server Error", errors.Wrap(err, "Select"))
		}

		err = s.db.Model(&blacklists).Where("owner_id = ? AND list_type = ?", u.ID, list.Blacklist).Select()
		if err != nil {
			return newApiError(http.StatusInternalServerError, "Internal Server Error", errors.Wrap(err, "Select"))
		}

		c.JSON(
			http.StatusOK,
			response{whitelists, blacklists},
		)

		return nil
	}
}

func (s Server) handlePostList() HandlerFunc {
	return func(u user.User, c *gin.Context) error {
		var l list.List

		if err := c.ShouldBind(&l); err != nil {
			return newApiError(http.StatusBadRequest, "Bad Request", errors.Wrap(err, "ShouldBind"))
		}

		if len(l.ListType) == 0 {
			return newApiError(http.StatusBadRequest, "Invalid List Type", errors.New("Invalid List Type"))
		}

		if err := l.Validate(); err != nil {
			return newApiError(http.StatusBadRequest, err.Error(), errors.Wrap(err, "ShouldBind"))
		}

		l.OwnerID = u.ID

		_, err := s.db.Model(&l).
			OnConflict("(list_type, domain, rr_type, record, owner_id) DO NOTHING").
			Insert()
		if err != nil {
			return newApiError(http.StatusInternalServerError, "Inserting", errors.Wrap(err, "Insert"))
		}

		l.AddedAt = time.Now()

		c.JSON(
			http.StatusCreated,
			&l,
		)

		return nil
	}
}
