# e-book

Simple e‑book website.

## Tech Stack

| Layer         | Stack                                                        |
| ------------- | ------------------------------------------------------------ |
| **Frontend**  | Next.js (React) + TypeScript + Tailwind CSS + Tiptap         |
| **Backend**   | Go (Fiber) + PostgreSQL + Redis + S3‑compatible (MinIO)       |
| **Search**    | OpenSearch                                                   |
| **Auth**      | JWT access tokens + refresh tokens                           |
<!-- | **Payments**  | Omise (PromptPay, Cards)                                     | -->
| **DevOps**    | Docker Compose + Traefik + Prometheus (optional: Grafana)    |


Env example
```env

POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=appdb

# Host port to bind the container's 5432
HOST_DB_PORT=5433

# Convenience DSN for the app or psql
DATABASE_URL=postgres://postgres:postgres@localhost:5433/appdb?sslmode=disable

 - MAILTRAP_HOST=live.smtp.mailtrap.io
          - MAILTRAP_PORT=587
          - MAILTRAP_USERNAME=...
          - MAILTRAP_PASSWORD=...
          - MAIL_FROM=no-reply@example.com

```
