package schemas

import (
	"time"

	"github.com/mtvru/go-project-278/internal/db"
)

type CreateLinkPayload struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	ShortName   string `json:"short_name" binding:"omitempty,min=3,max=32"`
}

type UpdateLinkPayload = CreateLinkPayload

type LinkResponse struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortURL    string `json:"short_url"`
}

func NewLinkResponse(l db.Link, baseURL string) LinkResponse {
	return LinkResponse{
		ID:          l.ID,
		OriginalURL: l.OriginalUrl,
		ShortName:   l.ShortName,
		ShortURL:    baseURL + "/r/" + l.ShortName,
	}
}

type VisitResponse struct {
	ID        int64     `json:"id"`
	LinkID    int64     `json:"link_id"`
	CreatedAt time.Time `json:"created_at"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
	Status    int32     `json:"status"`
}

func NewVisitResponse(v db.LinkVisit) VisitResponse {
	return VisitResponse{
		ID:        v.ID,
		LinkID:    v.LinkID,
		CreatedAt: v.CreatedAt,
		IP:        v.Ip,
		UserAgent: v.UserAgent,
		Referer:   v.Referer,
		Status:    v.Status,
	}
}
