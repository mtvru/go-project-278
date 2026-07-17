package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/mtvru/go-project-278/internal/db"
	"github.com/mtvru/go-project-278/internal/errs"
)

type VisitMeta struct {
	IP        string
	UserAgent string
	Referer   string
}

type VisitStore interface {
	GetLinkByShortName(ctx context.Context, shortName string) (db.Link, error)
	CreateLinkVisit(ctx context.Context, arg db.CreateLinkVisitParams) (db.LinkVisit, error)
	ListLinkVisits(ctx context.Context, arg db.ListLinkVisitsParams) ([]db.LinkVisit, error)
	CountLinkVisits(ctx context.Context) (int64, error)
}

type VisitService struct {
	store VisitStore
}

func NewVisitService(store VisitStore) *VisitService {
	return &VisitService{store: store}
}

func (s *VisitService) Count(ctx context.Context) (int64, error) {
	return s.store.CountLinkVisits(ctx)
}

func (s *VisitService) List(ctx context.Context, offset, limit int32) ([]db.LinkVisit, error) {
	return s.store.ListLinkVisits(ctx, db.ListLinkVisitsParams{Limit: limit, Offset: offset})
}

func (s *VisitService) Resolve(ctx context.Context, code string, meta VisitMeta, status int32) (db.Link, error) {
	link, err := s.store.GetLinkByShortName(ctx, code)
	if errors.Is(err, pgx.ErrNoRows) {
		return db.Link{}, errs.ErrNotFound
	}
	if err != nil {
		return db.Link{}, fmt.Errorf("get link by short name: %w", err)
	}

	if _, err := s.store.CreateLinkVisit(ctx, db.CreateLinkVisitParams{
		LinkID:    link.ID,
		Ip:        meta.IP,
		UserAgent: meta.UserAgent,
		Referer:   meta.Referer,
		Status:    status,
	}); err != nil {
		slog.ErrorContext(ctx, "record visit failed",
			slog.Int64("link_id", link.ID),
			slog.Any("error", err),
		)
	}

	return link, nil
}
