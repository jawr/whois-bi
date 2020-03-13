package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jawr/monere/domain"
	"github.com/jawr/monere/job"
	"github.com/jawr/monere/user"
	"github.com/pkg/errors"
)

func (s Server) handleGetJobs() HandlerFunc {
	return func(u user.User, c *gin.Context) error {
		jobs := make([]job.Job, 0)

		err := s.db.Model(&jobs).
			Relation("Domain").
			Where("domain.owner_id = ? AND domain.domain = ?", u.ID, c.Param("domain")).
			Order("job.created_at DESC").
			Select()
		if err != nil {
			return newApiError(http.StatusNotFound, "Not found", errors.Wrap(err, "Select"))
		}

		c.JSON(http.StatusOK, &jobs)

		return nil
	}
}

func (s Server) handlePostJob() HandlerFunc {
	return func(u user.User, c *gin.Context) error {
		var d domain.Domain

		err := s.db.Model(&d).Where("domain = ? AND owner_id = ?", c.Param("domain"), u.ID).Select()
		if err != nil {
			return newApiError(http.StatusNotFound, "Not found", errors.Wrap(err, "Select"))
		}

		j := job.NewJob(d)
		if err := j.Insert(s.db); err != nil {
			return newApiError(http.StatusInternalServerError, "Job already queued", errors.Wrap(err, "Insert"))
		}

		c.JSON(http.StatusOK, &j)

		return nil
	}
}
