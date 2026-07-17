package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mtvru/go-project-278/internal/errs"
	"github.com/mtvru/go-project-278/internal/schemas"
	"github.com/mtvru/go-project-278/internal/service"
)

type VisitHandler struct {
	visits *service.VisitService
}

func NewVisitHandler(visits *service.VisitService) *VisitHandler {
	return &VisitHandler{visits: visits}
}

func (h *VisitHandler) Redirect(c *gin.Context) {
	link, err := h.visits.Resolve(c.Request.Context(), c.Param("code"), service.VisitMeta{
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Referer:   c.GetHeader("Referer"),
	}, http.StatusFound)
	if errors.Is(err, errs.ErrNotFound) {
		notFound(c)
		return
	}
	if err != nil {
		serverError(c, err)
		return
	}
	c.Redirect(http.StatusFound, link.OriginalUrl)
}

func (h *VisitHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	total, err := h.visits.Count(ctx)
	if err != nil {
		serverError(c, err)
		return
	}

	pr, ok := parseRange(c, total)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid range"})
		return
	}

	visits, err := h.visits.List(ctx, pr.offset, pr.limit)
	if err != nil {
		serverError(c, err)
		return
	}

	responses := make([]schemas.VisitResponse, 0, len(visits))
	for _, v := range visits {
		responses = append(responses, schemas.NewVisitResponse(v))
	}

	setContentRange(c, "link_visits", pr.offset, len(responses), total)
	c.JSON(http.StatusOK, responses)
}
