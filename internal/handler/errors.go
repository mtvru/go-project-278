package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func serverError(c *gin.Context, err error) {
	slog.ErrorContext(c.Request.Context(), "request failed",
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
		slog.Any("error", err),
	)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func shortNameTaken(c *gin.Context) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"short_name": "short name already in use"}})
}

func respondBindError(c *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		fields := make(map[string]string, len(validationErrs))
		for _, fe := range validationErrs {
			fields[fe.Field()] = fe.Error()
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": fields})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
}
