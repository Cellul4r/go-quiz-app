# ExecPlans

When writing complex features or significant refactors, use an ExecPlan (as described in .agent/PLANS.md) from design to implementation.

# Repository Guidelines

## Project Structure & Module Organization

This repository contains a Go quiz API with Supabase-backed persistence. The Go module is in `backend/`.

- `backend/app/`: entry point, HTTP setup and module wiring.
- `backend/domain/`: entities, interfaces, validation, and errors.
- `backend/internal/rest/`: Fiber handlers, DTOs, middleware, helpers, and REST tests.
- `backend/internal/repository/postgres/`: Postgres repositories and integration tests.
- `backend/{quiz,profile,favorite}/`: feature services and generated mocks.
- `supabase/migrations/`: SQL migrations and local Supabase config.
- `backend/swagger_doc/`: generated Swagger artifacts.
- `docs/images/`: project diagrams and documentation assets.

## Build, Test, and Development Commands

Run Go commands from `backend/` unless noted.

- `go run ./app`: start the API locally with the root `.env` configuration.
- `go test ./...`: run all backend tests.
- `go test ./internal/rest ./quiz ./profile ./favorite`: run faster service and REST tests.
- `go build ./...`: verify all backend packages compile.
- `supabase start`: start local Supabase from the root.
- `supabase migration up`: apply pending migrations locally.
- `supabase db reset`: reset local data and reapply migrations.
- `atlas migrate diff <name> --env gorm`: generate a migration from GORM models.
- `atlas migrate hash --env gorm`: repair Atlas migration checksums.

## Coding Style & Naming Conventions

Run `gofmt` or `go fmt ./...`. Keep package names short, lowercase, and aligned with directory names. Export only API, service, or domain symbols used across packages. Prefer names such as `QuizService`, `ProfileRepository`, and `FavoriteHandler`; keep request/response shapes in `internal/rest/dto`.

## Testing Guidelines

Tests use Go's `testing` package with `testify`, mocks, `sqlmock`, and Testcontainers for Postgres integration coverage. Name test files `*_test.go` beside the package under test. Use table-driven tests for handlers and services. Repository integration tests may require Docker and local environment values.

## Commit & Pull Request Guidelines

Use short Conventional Commit-style prefixes, for example `feat: add tests for favorite service` and `fix: change author_id to static in quiz_test.go`. Keep commits focused; use `feat:`, `fix:`, `test:`, `docs:`, or `chore:` where applicable.

Pull requests should include a summary, linked issue when available, test results (`go test ./...`), and notes for migrations or environment changes. 

## Security & Configuration Tips

Do not commit `.env`; update `.env.example` when configuration changes. Review generated migrations before applying them; never run `supabase db push` against production outside the approved deployment flow.
