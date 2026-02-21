# Rizon Test Task

Backend and client for a small feedback flow: Go GraphQL API (gqlgen) and a React Native (Expo) app that talks to it. Auth is email magic-link; the API can enqueue email and Slack notifications via RabbitMQ, and background workers send them.

---

## What’s in the repo

- **Backend**: Go server with GraphQL (99designs/gqlgen), PostgreSQL (GORM), Redis (rate limiting / in-memory state), RabbitMQ (email and Slack jobs).
- **Workers**: Single process that runs both an email worker (SMTP) and a Slack worker (mock that logs to console).
- **Client**: React Native app under `cmd/clients/rn-app` (Expo, Apollo Client, TypeScript). GraphQL types and hooks are generated from the live API.

All runnable steps are wired through [Task](https://taskfile.dev/). Use `task --list` to see available commands.

---

## Prerequisites

- **Go** (1.21+). The project uses Go modules.
- **Node.js** and **npm** (for the React Native app).
- **Docker** and **Docker Compose** (for PostgreSQL, Redis, RabbitMQ).
- **Task** (Taskfile runner):  
  `brew install go-task` on macOS, or see [taskfile.dev](https://taskfile.dev/).
- **gqlgen** (for backend GraphQL codegen):  
  `go install github.com/99designs/gqlgen/cmd/gqlgen@v0.17.86`  
  Ensure `$(go env GOPATH)/bin` is on your `PATH` so `task gql:gen` can find it.

---

## Quick start (backend only)

From the project root:

1. **Start dependencies** (Postgres, Redis, RabbitMQ):
   ```bash
   task deps:up
   ```

2. **Tidy Go modules** (if needed):
   ```bash
   task deps:tidy
   ```

3. **Run migrations**:
   ```bash
   task migrate:up
   ```

4. **Start the server**:
   ```bash
   task server:run
   ```

5. **Run the worker** (optional; processes email and Slack jobs from RabbitMQ). In a second terminal, set `EMAIL_USERNAME`, `EMAIL_PASSWORD`, and `EMAIL_FROM`, then:
   ```bash
   task worker:run
   ```
   See **Running the workers** below for details.

The GraphQL API is at `http://localhost:8080/graphql`. A playground is available at `http://localhost:8080/` (Apollo sandbox available at `http://localhost:8080/sandbox`) when the server is running.

Default config assumes the Docker stack above (see **Configuration** for overrides).

---

## Running the React Native app

1. Backend must be running (Quick start steps 1–4).

2. **Install frontend dependencies**:
   ```bash
   task rn:install
   ```

3. **Generate GraphQL client** (types and React hooks from the running API):
   ```bash
   task rn:codegen
   ```
   This uses `GRAPHQL_ENDPOINT` if set; otherwise `http://localhost:8080/graphql`. Start the server first so the schema is available.

4. **Start the dev server**:
   ```bash
   task rn:start
   ```
   Then run on a simulator/emulator or device (e.g. press `i` for iOS, `a` for Android in the Expo CLI), or use:
   - `task rn:ios`
   - `task rn:android`
   - `task rn:web` for web.

---

## Running the workers (email + Slack)

The worker process runs both the email and Slack consumers. It requires RabbitMQ (same as the server) and, for the email worker, SMTP settings.

1. **Start dependencies** (if not already):  
   `task deps:up`

2. **Set email env vars** (required for the worker to start):
   - `EMAIL_USERNAME` – SMTP username  
   - `EMAIL_PASSWORD` – SMTP password  
   - `EMAIL_FROM` – From address  

   Optional overrides: `EMAIL_SMTP_HOST` (default `smtp.gmail.com`), `EMAIL_SMTP_PORT` (default `587`).

3. **Run the worker**:
   ```bash
   task worker:run
   ```

Slack is implemented as a mock that logs to the console; no Slack API key is needed for local runs.

---

## GraphQL code generation

- **Backend (Go)**  
  Schema: `internal/graphql/schema.graphqls`. Generate resolvers and types with:
  ```bash
  task gql:gen
  ```
  Requires `gqlgen` installed (see Prerequisites).

- **Frontend (React Native)**  
  The app uses GraphQL Code Generator (see `cmd/clients/rn-app/codegen.config.js`). It introspects the running server and writes TypeScript types and React Apollo hooks to `cmd/clients/rn-app/graphql/generated/graphql.ts`. Run with:
  ```bash
  task rn:codegen
  ```
  Ensure the backend is up and reachable at the URL used by codegen (default `http://localhost:8080/graphql`, or set `GRAPHQL_ENDPOINT`).

---

## Configuration

The server and workers read from the environment; a `.env` in the project root is loaded if present (via `godotenv`). You can copy `.env.example` to `.env` and fill in any values you need to override (or set only the ones required for the worker).

**Server / API**

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP listen port |
| `BASE_URL` | `http://localhost:8080` | Base URL used for auth links etc. |

**Database (PostgreSQL)**  
Matches the Docker Compose defaults if not set:

| Variable | Default |
|----------|---------|
| `DB_HOST` | `localhost` |
| `DB_PORT` | `5432` |
| `DB_USER` | `rizon` |
| `DB_PASSWORD` | `rizon_dev_password` |
| `DB_NAME` | `rizon_db` |
| `DB_SSLMODE` | `disable` |

**Redis**  
Used for rate limiting and in-memory state. Defaults: `REDIS_HOST=localhost`, `REDIS_PORT=6379`, or set `REDIS_URL`.

**RabbitMQ**  
Used for email and Slack job queues. Defaults: host `localhost`, port `5672`, user `rizon`, password `rizon_dev_password`, vhost `/`, or set `RABBITMQ_URL`.

**Auth**  
Optional overrides: `EMAIL_AUTH_SALT`, `EMAIL_AUTH_ENDPOINT`, `JWT_SECRET`. Defaults are suitable only for local development.

**Rate limiting**  
`RATE_LIMIT_REQUESTS` (default 60), `RATE_LIMIT_WINDOW_SECONDS` (default 60).

**Email (worker)**  
Required when running the worker: `EMAIL_USERNAME`, `EMAIL_PASSWORD`, `EMAIL_FROM`. Optional: `EMAIL_SMTP_HOST`, `EMAIL_SMTP_PORT`.

---

## Task reference

All commands are run from the repo root via Task.

| Task | Description |
|------|-------------|
| `task deps:up` | Start PostgreSQL, Redis, RabbitMQ (Docker) |
| `task deps:down` | Stop dependency containers |
| `task deps:tidy` | `go mod tidy` |
| `task migrate:up` | Apply migrations |
| `task migrate:down` | Roll back last migration |
| `task server:run` | Run the GraphQL server |
| `task server:build` | Build server binary to `bin/server` |
| `task worker:run` | Run email + Slack workers (needs email env) |
| `task worker:build` | Build worker binary to `bin/worker` |
| `task gql:gen` | Generate Go GraphQL code (gqlgen) |
| `task rn:install` | npm install in React Native app |
| `task rn:codegen` | Generate GraphQL client for RN app |
| `task rn:start` | Start Expo dev server |
| `task rn:ios` / `rn:android` / `rn:web` | Run app on iOS / Android / web |
| `task mocks:gen` | Generate gomock mocks for tests |
| `task schema` | Print current DB schema (uses `tools/show-schema.sh`) |

---

## Project layout (high level)

- `cmd/server` – GraphQL HTTP server entrypoint  
- `cmd/worker` – Email and Slack worker entrypoint  
- `cmd/migration` – Database migrations (goose)  
- `cmd/clients/rn-app` – React Native (Expo) app and its GraphQL codegen config  
- `internal/` – Core logic: app (auth, feedback), config, database, GraphQL (schema, resolvers, generated), middleware, message broker, repositories, etc.  
- `internal/graphql/schema.graphqls` – Single GraphQL schema used by gqlgen  
- `docker-compose.deps.yml` – Postgres, Redis, RabbitMQ for local development  

---

## Tests

Run Go tests from the repo root:

```bash
go test ./...
```

Mock generation for tests: `task mocks:gen`.

---

## Stopping services

- Stop the Docker stack: `task deps:down`  
- To remove data as well: `task deps:clean`

Server and worker are stopped with Ctrl+C in their terminals.
