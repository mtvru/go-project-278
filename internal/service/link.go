package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mtvru/go-project-278/internal/db"
	"github.com/mtvru/go-project-278/internal/errs"
	"github.com/mtvru/go-project-278/internal/shortname"
)

const maxShortNameAttempts = 5

type LinkStore interface {
	CreateLink(ctx context.Context, arg db.CreateLinkParams) (db.Link, error)
	GetLink(ctx context.Context, id int64) (db.Link, error)
	ListLinks(ctx context.Context, arg db.ListLinksParams) ([]db.Link, error)
	CountLinks(ctx context.Context) (int64, error)
	UpdateLink(ctx context.Context, arg db.UpdateLinkParams) (db.Link, error)
	DeleteLink(ctx context.Context, id int64) (int64, error)
}

type LinkService struct {
	store LinkStore
}

func NewLinkService(store LinkStore) *LinkService {
	return &LinkService{store: store}
}

func (s *LinkService) Count(ctx context.Context) (int64, error) {
	return s.store.CountLinks(ctx)
}

func (s *LinkService) List(ctx context.Context, offset, limit int32) ([]db.Link, error) {
	return s.store.ListLinks(ctx, db.ListLinksParams{Limit: limit, Offset: offset})
}

func (s *LinkService) Create(ctx context.Context, originalURL, shortName string) (db.Link, error) {
	generate := shortName == ""
	for attempt := 0; ; attempt++ {
		name := shortName
		if generate {
			generated, err := shortname.Generate()
			if err != nil {
				return db.Link{}, fmt.Errorf("generate short name: %w", err)
			}
			name = generated
		}

		link, err := s.store.CreateLink(ctx, db.CreateLinkParams{
			OriginalUrl: originalURL,
			ShortName:   name,
		})
		if err == nil {
			return link, nil
		}
		if errs.IsUniqueViolation(err) {
			if generate && attempt < maxShortNameAttempts {
				continue
			}
			return db.Link{}, errs.ErrShortNameTaken
		}
		return db.Link{}, fmt.Errorf("create link: %w", err)
	}
}

func (s *LinkService) Get(ctx context.Context, id int64) (db.Link, error) {
	link, err := s.store.GetLink(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return db.Link{}, errs.ErrNotFound
	}
	if err != nil {
		return db.Link{}, fmt.Errorf("get link: %w", err)
	}
	return link, nil
}

func (s *LinkService) Update(ctx context.Context, id int64, originalURL, shortName string) (db.Link, error) {
	name := shortName
	if name == "" {
		generated, err := shortname.Generate()
		if err != nil {
			return db.Link{}, fmt.Errorf("generate short name: %w", err)
		}
		name = generated
	}

	link, err := s.store.UpdateLink(ctx, db.UpdateLinkParams{
		ID:          id,
		OriginalUrl: originalURL,
		ShortName:   name,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return db.Link{}, errs.ErrNotFound
	}
	if errs.IsUniqueViolation(err) {
		return db.Link{}, errs.ErrShortNameTaken
	}
	if err != nil {
		return db.Link{}, fmt.Errorf("update link: %w", err)
	}
	return link, nil
}

func (s *LinkService) Delete(ctx context.Context, id int64) error {
	affected, err := s.store.DeleteLink(ctx, id)
	if err != nil {
		return fmt.Errorf("delete link: %w", err)
	}
	if affected == 0 {
		return errs.ErrNotFound
	}
	return nil
}
