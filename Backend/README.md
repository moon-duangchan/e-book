# E-Book Backend

Go + Fiber REST API with Postgres, GORM, JWT auth, and email verification.

## Stack
- Go (Fiber, GORM)
- Postgres (Docker Compose for local dev)
- JWT (github.com/golang-jwt/jwt/v5)
- SMTP via Mailtrap (email verification)

## Quick Start
1. Install Go and Docker.
2. Copy `.env` (already present) and update values as needed.
3. Start Postgres:
   - `docker compose up -d`
4. Install deps and run the app:
   - `go mod tidy`
   - `go run main.go`
5. App listens on `http://localhost:${PORT}` (default `3001`).

## Environment
Key variables used by the app (see `.env`):
- `PORT` HTTP port (default `3001`)
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `HOST_DB_PORT` (Docker Postgres)
- `JWT_SECRET` secret for signing JWTs (required)
- `MAILTRAP_HOST`, `MAILTRAP_PORT`, `MAILTRAP_USERNAME`, `MAILTRAP_PASSWORD`, `MAIL_FROM`
- `APP_BASE_URL` base URL for email verification links (fallbacks to `http://localhost:$PORT`)

## Database
- Starts via `docker compose up -d` using `docker-compose.yml`.
- The app auto-migrates `Todo` and `User` on startup.

## Auth Flow
- Register: user signs up and receives a verification link by email.
- Verify: user clicks the link to verify the account.
- Login: on success, API returns a JWT.
- Use: include the JWT in `Authorization: Bearer <token>` for protected endpoints.

## API
- `POST /auth/register` — body: `{ "name": "...", "email": "...", "password": "..." }`
- `GET  /auth/verify?token=...` — verify email address
- `POST /auth/login` — body: `{ "email": "...", "password": "..." }` → returns `{ "token": "..." }`

### Todos (Protected)
All `/todos` routes require a valid JWT via `Authorization: Bearer <token>`.
- `GET    /todos`
- `GET    /todos/:id`
- `POST   /todos` — body: `{ "title": "...", "completed": false }`
- `PUT    /todos/:id` — body: same fields as `POST`
- `DELETE /todos/:id`

If you want only write operations protected (and allow public reads), replace the global `app.Use("/todos", auth.RequireAuth())` with the middleware on specific routes only.

## Examples
- Register:
```
curl -X POST http://localhost:3001/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Jane","email":"jane@example.com","password":"secret"}'
```
- Login:
```
TOKEN=$(curl -s -X POST http://localhost:3001/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"jane@example.com","password":"secret"}' | jq -r .token)
```
- Use token:
```
curl http://localhost:3001/todos -H "Authorization: Bearer $TOKEN"
```

## Troubleshooting
- Missing checksums: `go mod tidy`
- 401 Unauthorized: ensure `Authorization: Bearer <token>` header and `JWT_SECRET` set.
- Email not received: verify Mailtrap creds; for local testing, you can fetch the `verification_token` for the user from the DB and call `/auth/verify?token=...` directly.

## Project Structure
- `main.go` — app entry, routes and middleware wiring
- `controller/` — handlers for auth and todos
- `auth/` — email sender and JWT middleware
- `database/` — DB connection holder
- `models/` — GORM models
- `docker-compose.yml` — local Postgres

