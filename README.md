# URL Shortener (Go)

### Hexlet tests and linter status:
[![Actions Status](https://github.com/mtvru/go-project-278/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/mtvru/go-project-278/actions)
[![CI](https://github.com/mtvru/go-project-278/actions/workflows/ci.yml/badge.svg)](https://github.com/mtvru/go-project-278/actions/workflows/ci.yml)

A URL shortener built with Go (Gin) and PostgreSQL, with redirects and visit analytics.

## Stack

Gin · PostgreSQL + sqlc · goose (migrations) · validator · Sentry · Caddy · Docker / Render

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/ping` | health check (`pong`) |
| GET | `/api/links` | list links (`range` pagination) |
| POST | `/api/links` | create a link |
| GET | `/api/links/:id` | get a link |
| PUT | `/api/links/:id` | update a link |
| DELETE | `/api/links/:id` | delete a link |
| GET | `/api/link_visits` | list visits (pagination) |
| GET | `/r/:code` | redirect by short name |

## Run locally

```bash
make db-up                                                                   # start PostgreSQL
export DATABASE_URL="postgres://postgres:password@localhost:5434/appdb?sslmode=disable"
make migrate                                                                 # apply migrations
npm install && make dev                                                      # run backend + frontend
```

Frontend: http://localhost:5173 · API: http://localhost:8080. See `.env.example` for env vars.
