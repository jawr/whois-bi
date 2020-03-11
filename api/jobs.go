package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jawr/monere/job"
	"github.com/jawr/monere/user"
	"github.com/pkg/errors"
)

func (s Server) handleGetJobs() HandlerFunc {
	return func(u user.User, c *gin.Context) error {
		jobs := make([]job.Job, 0)

		err := s.db.Model(&jobs).
			Relation("Domain").
			Where("domain.owner_id = ?", u.ID).
			Order("job.created_at DESC").
			Select()
		if err != nil {
			return newApiError(http.StatusNotFound, "Not found", errors.Wrap(err, "Select"))
		}

		c.JSON(http.StatusOK, &jobs)

		return nil
	}
}
