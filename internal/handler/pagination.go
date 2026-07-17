package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type pageRange struct {
	offset int32
	limit  int32
}

func parseRange(c *gin.Context, total int64) (pageRange, bool) {
	raw := c.Query("range")
	if raw == "" {
		raw = c.GetHeader("Range")
	}
	if raw == "" {
		return pageRange{offset: 0, limit: int32(total)}, true
	}

	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimPrefix(trimmed, "[")
	trimmed = strings.TrimSuffix(trimmed, "]")
	parts := strings.Split(trimmed, ",")
	if len(parts) != 2 {
		return pageRange{}, false
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return pageRange{}, false
	}
	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return pageRange{}, false
	}
	if start < 0 || end < start {
		return pageRange{}, false
	}

	return pageRange{offset: int32(start), limit: int32(end - start)}, true
}

func setContentRange(c *gin.Context, resource string, offset int32, returned int, total int64) {
	c.Header("Content-Range", fmt.Sprintf("%s %d-%d/%d", resource, offset, int(offset)+returned, total))
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		notFound(c)
		return 0, false
	}
	return id, true
}
