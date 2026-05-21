# ⚽ League Simulation

A full-stack football league simulation built with **Go** and **Vue.js**. Simulates a 4-team Premier League-style season with match results, standings, championship predictions, and an admin panel.

🔗 **Live Demo:** [league-simulation-web.vercel.app](https://league-simulation-web.vercel.app/)  
📖 **API Docs:** [league-simulation-api.onrender.com/swagger/index.html](https://league-simulation-api.onrender.com/swagger/index.html)

---

## 🧱 Tech Stack

| Layer     | Technology                              |
| --------- | --------------------------------------- |
| Backend   | Go, Gin, GORM, JWT, Swagger, Zerolog    |
| Frontend  | Vue 3, Pinia, Chart.js, Tailwind CSS v4 |
| Database  | PostgreSQL (Supabase)                   |
| Cache     | Redis (Upstash)                         |
| Container | Docker, Docker Compose                  |
| CI/CD     | GitHub Actions                          |
| Deploy    | Render (backend) + Vercel (frontend)    |

---

## ✨ Features

- 🏆 4-team double round-robin fixture (6 weeks, 12 matches)
- ⚽ Match simulation engine using Poisson distribution with home advantage
- 📊 Premier League-style standings (PTS → GD → GF)
- 🎯 Championship predictions from Week 4 onwards
- 🔐 JWT-based admin authentication
- ✏️ Admin can edit match results — standings recalculate automatically
- ⚡ Redis caching for standings and predictions
- 📖 Swagger UI for interactive API documentation
- 📈 Season statistics with Chart.js visualizations
- 🧪 89.5% test coverage (unit + integration)

---

## 🚀 Quick Start (Docker)

### Prerequisites

- [Docker Desktop](https://www.docker.com/get-started) installed and running
- [Git](https://git-scm.com)

### 1 — Clone the repository

```bash
git clone https://github.com/yahyayesilyurt/league-simulation.git
cd league-simulation
```

### 2 — Set up environment variables

```bash
cp .env.example .env
```

The default `.env` is pre-configured for local Docker development. No changes needed to run locally.

### 3 — Start all services

```bash
docker compose up --build
```

This starts:

- **Go backend** on `http://localhost:8080`
- **Vue frontend** on `http://localhost:5173`
- **PostgreSQL** on `localhost:5432`
- **Redis** on `localhost:6379`
- **Nginx** reverse proxy on `http://localhost:80`

### 4 — Open the app

| Service    | URL                                      |
| ---------- | ---------------------------------------- |
| Frontend   | http://localhost:5173                    |
| API Health | http://localhost:8080/health             |
| Swagger UI | http://localhost:8080/swagger/index.html |

---

## 🔑 Admin Credentials

| Username | Password   |
| -------- | ---------- |
| `admin`  | `admin123` |

Used for editing match results and resetting the league.

---

## 📡 API Endpoints

### Public

| Method | Endpoint               | Description                       |
| ------ | ---------------------- | --------------------------------- |
| GET    | `/health`              | Health check                      |
| POST   | `/auth/login`          | Get JWT token                     |
| GET    | `/league/table`        | Current standings                 |
| GET    | `/league/fixtures`     | Full fixture list                 |
| GET    | `/league/week/:weekNo` | Matches for a specific week (1-6) |
| GET    | `/league/predictions`  | Championship predictions          |
| GET    | `/league/status`       | League status                     |
| POST   | `/league/next-week`    | Play the next week                |
| POST   | `/league/play-all`     | Play all remaining weeks          |

### Protected (Bearer token required)

| Method | Endpoint            | Description         |
| ------ | ------------------- | ------------------- |
| POST   | `/league/reset`     | Reset the league    |
| PUT    | `/match/:id/result` | Edit a match result |

---

## 🏗️ Project Structure

```
league-simulation/
├── backend/
│   ├── cmd/server/         # Entry point
│   ├── config/             # DB, Redis, JWT, Logger config
│   ├── docs/               # Swagger generated files
│   ├── internal/
│   │   ├── cache/          # Redis cache layer
│   │   ├── handler/        # HTTP handlers + router
│   │   ├── middleware/     # JWT auth middleware
│   │   ├── model/          # GORM models
│   │   ├── repository/     # DB queries (interface-based)
│   │   └── service/        # Business logic
│   ├── migrations/         # SQL schema
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/            # Axios services
│   │   ├── components/     # Vue components
│   │   ├── composables/    # Reusable logic
│   │   ├── stores/         # Pinia stores
│   │   └── views/          # Page views
│   └── Dockerfile
├── nginx/
│   └── nginx.conf
├── docker-compose.yml
└── .env.example
```

---

## 🗃️ Database Schema

```sql
teams      — id, name, strength, created_at, updated_at
matches    — id, week, home_team_id, away_team_id, home_goals,
             away_goals, played, created_at, updated_at
standings  — id, team_id, played, won, drawn, lost,
             goals_for, goals_against, goal_diff, points, updated_at
```

Full schema: [`backend/migrations/001_initial_schema.sql`](./backend/migrations/001_initial_schema.sql)

---

## 🧪 Running Tests

```bash
cd backend

# Run all tests
go test ./...

# With coverage report
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Verbose output
go test ./... -v
```

**Current coverage: 89.5%**

---

## 🔄 CI/CD Pipeline

GitHub Actions runs on every push to `main`.

Config: [`.github/workflows/ci.yml`](./.github/workflows/ci.yml)

---

## 🌍 Deployment

### Backend — Render.com

| Variable            | Value                        |
| ------------------- | ---------------------------- |
| `SUPABASE_URL`      | PostgreSQL connection string |
| `UPSTASH_REDIS_URL` | Redis connection string      |
| `JWT_SECRET`        | Your secret key              |
| `ADMIN_USERNAME`    | Admin username               |
| `ADMIN_PASSWORD`    | Admin password               |
| `APP_ENV`           | `production`                 |

### Frontend — Vercel

| Variable            | Value              |
| ------------------- | ------------------ |
| `VITE_API_BASE_URL` | Render backend URL |

---

## ⚙️ League Rules

- 4 teams with different strength ratings (75–90)
- Double round-robin: each team plays 6 matches (3 home, 3 away)
- **Win:** 3 pts · **Draw:** 1 pt · **Loss:** 0 pts
- Tiebreaker: Points → Goal Difference → Goals Scored
- Championship predictions unlock from **Week 4**
- Match simulation uses Poisson distribution with home advantage (+5 strength)

---

## 📦 Docker Commands

```bash
# Start all services
docker compose up -d

# Rebuild after code changes
docker compose up --build

# View backend logs
docker compose logs -f backend

# Stop all services
docker compose down

# Full reset (removes volumes)
docker compose down -v
```
