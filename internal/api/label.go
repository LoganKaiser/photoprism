package api

import (
	"net/http"

	"github.com/photoprism/photoprism/pkg/sanitize"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/acl"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/i18n"
	"github.com/photoprism/photoprism/internal/query"
	"github.com/photoprism/photoprism/pkg/txt"
)

// UpdateLabel updates label properties.
//
// PUT /api/v1/labels/:uid
func UpdateLabel(router *gin.RouterGroup) {
	router.PUT("/labels/:uid", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourceLabels, acl.ActionUpdate)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		id := sanitize.IdString(c.Param("uid"))
		m, err := query.LabelByUID(id)

		if err != nil {
			Abort(c, http.StatusNotFound, i18n.ErrLabelNotFound)
			return
		}

		f, err := form.NewLabel(m)

		if err != nil {
			log.Error(err)
			AbortSaveFailed(c)
			return
		}

		if err := c.BindJSON(&f); err != nil {
			log.Error(err)
			AbortBadRequest(c)
			return
		}

		if err := m.SaveForm(f); err != nil {
			log.Error(err)
			AbortSaveFailed(c)
			return
		}

		event.SuccessMsg(i18n.MsgLabelSaved)

		PublishLabelEvent(EntityUpdated, id, c)

		c.JSON(http.StatusOK, m)
	})
}

// LikeLabel flags a label as favorite.
//
// POST /api/v1/labels/:uid/like
//
// Parameters:
//   uid: string Label UID
func LikeLabel(router *gin.RouterGroup) {
	router.POST("/labels/:uid/like", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourceLabels, acl.ActionUpdate)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		id := sanitize.IdString(c.Param("uid"))
		label, err := query.LabelByUID(id)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": txt.UcFirst(err.Error())})
			return
		}

		if err := label.Update("LabelFavorite", true); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": txt.UcFirst(err.Error())})
			return
		}

		if label.LabelPriority < 0 {
			event.Publish("count.labels", event.Data{
				"count": 1,
			})
		}

		PublishLabelEvent(EntityUpdated, id, c)

		c.JSON(http.StatusOK, http.Response{})
	})
}

// DislikeLabel removes the favorite flag from a label.
//
// DELETE /api/v1/labels/:uid/like
//
// Parameters:
//   uid: string Label UID
func DislikeLabel(router *gin.RouterGroup) {
	router.DELETE("/labels/:uid/like", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourceLabels, acl.ActionUpdate)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		id := sanitize.IdString(c.Param("uid"))
		label, err := query.LabelByUID(id)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": txt.UcFirst(err.Error())})
			return
		}

		if err := label.Update("LabelFavorite", false); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": txt.UcFirst(err.Error())})
			return
		}

		if label.LabelPriority < 0 {
			event.Publish("count.labels", event.Data{
				"count": -1,
			})
		}

		PublishLabelEvent(EntityUpdated, id, c)

		c.JSON(http.StatusOK, http.Response{})
	})
}
