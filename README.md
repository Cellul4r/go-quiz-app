# go-quiz-app

GoQuiz is a dynamic, community-driven learning platform featuring interactive solo quizzes, flashcard study modes, and customizable question banks.

## GoQuiz Architecture

Here is the database schema for Phase 1:

![GoQuiz Database Schema](./docs/images/GoQuiz%20DB%20diagram.png)

## Prerequisites

- [Go 1.26+](https://golang.org/dl/)
- [Docker](https://www.docker.com/products/docker-desktop)
- [Supabase CLI](https://supabase.com/docs/guides/cli/getting-started)
- [Atlas CLI](https://atlasgo.io/getting-started)

## Project Structure

go-quiz-app/
├── backend/ # Go Fiber backend
│ ├── app/ # Entry point
│ ├── domain/ # Entities + interfaces
│ ├── internal/
│ │ ├── atlas/ # Atlas schema generator
│ │ ├── config/ # App configuration
│ │ ├── database/ # GORM connection
│ │ ├── middleware/ # JWT auth middleware
│ │ ├── repository/ # DB implementations
│ │ └── rest/ # Fiber HTTP handlers
│ └── atlas.hcl # Atlas config
├── supabase/
│ ├── migrations/ # SQL migration files
│ └── seed.sql # Seed data
├── frontend/
├── .env # Single env file for entire project (never commit)
├── .env.example # Example env file (commit this)
└── docker-compose.yml

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Cellul4r/go-quiz-app.git
cd go-quiz-app
```

### 2. Setup environment variables

```bash
cp .env.example .env
```

Fill in the values:

```env
#App
APP_ENV=development # The application environment (e.g., development, production)
APP_DEBUG=true # Enable or disable debug mode

PORT=3000 # The port on which the application will run

# Database configuration
DB_HOST=
DB_PORT=6543
DB_NAME=postgres
DB_USER=
DB_PASSWORD=

```

### 3. Start local Supabase

```bash
supabase start
```

After starting, Supabase will print your local credentials:
API URL: http://localhost:54321
DB URL: postgresql://postgres:postgres@localhost:54322/postgres
Studio: http://localhost:54323
Anon Key: <your-local-anon-key>
JWT Secret: <your-local-jwt-secret>

```bash
supabase status
```

### 4. Apply database migrations

```bash
supabase migration up
```

To reset the database and reapply all migrations from scratch including seed data:

```bash
supabase db reset
```

### 5. Start the backend

```bash
# development with hot reload
docker compose -f docker-compose.yml -f docker-compose.dev.yml up

# or without docker
cd backend
go run ./app/main.go
```

### 6. Verify everything is running

Backend API: http://localhost:3000
Supabase Studio: http://localhost:54323
Swagger UI: http://localhost:3000/swagger/index.html

---

## Atlas Migration Guide

Atlas generates SQL migrations from your GORM models automatically.

### Install Atlas CLI

```bash
# macOS
brew install ariga/tap/atlas

# Linux / WSL
curl -sSf https://atlasgo.sh | sh
```

### Generate a new migration

```bash
cd backend
atlas migrate diff <migration_name> --env gorm
```

Example:

```bash
atlas migrate diff create_quizzes --env gorm
```

### Re-hash migration directory

Run this if Atlas reports a checksum error:

```bash
atlas migrate hash --env gorm
```

### Apply migrations locally

```bash
# at root
supabase migration up
```

### Full workflow when you change a domain struct

```bash
# 1. update your struct in domain/

# 2. generate migration (inside backend/)
cd backend
atlas migrate diff <migration_name> --env gorm

# 3. review the generated SQL in supabase/migrations/

# 4. apply locally (at root)
cd ..
supabase migration up

# 5. push to cloud when ready
supabase db push
```

> ⚠️ Only push to cloud via CI/CD in production. Never run `db push` manually against production.

---

## Stopping local Supabase

```bash
supabase stop
```

To stop and wipe all local data:

```bash
supabase stop --no-backup
```

---

## Common Commands Reference

| Command                                | Description                              | Run from   |
| -------------------------------------- | ---------------------------------------- | ---------- |
| `supabase start`                       | Start local Supabase                     | root       |
| `supabase stop`                        | Stop local Supabase                      | root       |
| `supabase status`                      | Show local URLs and keys                 | root       |
| `supabase migration up`                | Apply pending migrations                 | root       |
| `supabase migration new <name>`        | Create empty migration file              | root       |
| `supabase db reset`                    | Reset DB + reapply all migrations + seed | root       |
| `supabase db push`                     | Push migrations to cloud                 | root       |
| `atlas migrate diff <name> --env gorm` | Generate migration from GORM models      | `backend/` |
| `atlas migrate hash --env gorm`        | Fix checksum errors                      | `backend/` |
