package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mtvru/go-project-278/internal/errs"
	"github.com/mtvru/go-project-278/internal/schemas"
	"github.com/mtvru/go-project-278/internal/service"
)

type LinkHandler struct {
	links   *service.LinkService
	baseURL string
}

func NewLinkHandler(links *service.LinkService, baseURL string) *LinkHandler {
	return &LinkHandler{links: links, baseURL: baseURL}
}

func (h *LinkHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	total, err := h.links.Count(ctx)
	if err != nil {
		serverError(c, err)
		return
	}

	pr, ok := parseRange(c, total)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid range"})
		return
	}

	links, err := h.links.List(ctx, pr.offset, pr.limit)
	if err != nil {
		serverError(c, err)
		return
	}

	responses := make([]schemas.LinkResponse, 0, len(links))
	for _, l := range links {
		responses = append(responses, schemas.NewLinkResponse(l, h.baseURL))
	}

	setContentRange(c, "links", pr.offset, len(responses), total)
	c.JSON(http.StatusOK, responses)
}

func (h *LinkHandler) Create(c *gin.Context) {
	var payload schemas.CreateLinkPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondBindError(c, err)
		return
	}

	link, err := h.links.Create(c.Request.Context(), payload.OriginalURL, payload.ShortName)
	if errors.Is(err, errs.ErrShortNameTaken) {
		shortNameTaken(c)
		return
	}
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, schemas.NewLinkResponse(link, h.baseURL))
}

func (h *LinkHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	link, err := h.links.Get(c.Request.Context(), id)
	if errors.Is(err, errs.ErrNotFound) {
		notFound(c)
		return
	}
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, schemas.NewLinkResponse(link, h.baseURL))
}

func (h *LinkHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var payload schemas.UpdateLinkPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondBindError(c, err)
		return
	}

	link, err := h.links.Update(c.Request.Context(), id, payload.OriginalURL, payload.ShortName)
	if errors.Is(err, errs.ErrNotFound) {
		notFound(c)
		return
	}
	if errors.Is(err, errs.ErrShortNameTaken) {
		shortNameTaken(c)
		return
	}
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, schemas.NewLinkResponse(link, h.baseURL))
}

func (h *LinkHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	err := h.links.Delete(c.Request.Context(), id)
	if errors.Is(err, errs.ErrNotFound) {
		notFound(c)
		return
	}
	if err != nil {
		serverError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
