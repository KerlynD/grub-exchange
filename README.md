# Grub Exchange

A Robinhood-style social trading app where users can buy and sell shares of their friends using "Grub" currency.

## Tech Stack

- **Frontend**: Next.js 14, TypeScript, Tailwind CSS, Framer Motion, Recharts
- **Backend**: Go (Gin framework)
- **Database**: SQLite (with WAL mode)
- **Auth**: JWT with httpOnly cookies

## Getting Started

### Prerequisites

- Node.js 18+
- Go 1.21+
- GCC (required for SQLite via CGo — install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) on Windows)

### Backend

```bash
cd backend
go mod tidy
go run ./cmd/server/
```

The API server starts on `http://localhost:8080`.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend starts on `http://localhost:3000`.

### Environment Variables

**Backend** (optional):
- `PORT` — Server port (default: `8080`)
- `DB_PATH` — SQLite database path (default: `grub_exchange.db`)
- `JWT_SECRET` — JWT signing secret (default: dev secret)
- `FRONTEND_URL` — Frontend origin for CORS (default: `http://localhost:3000`)
- `ENV` — Set to `production` for secure cookies

**Frontend** (optional):
- `NEXT_PUBLIC_API_URL` — Backend API URL (default: `http://localhost:8080`)

## Features

- User registration creates both a trading account AND a tradeable stock
- Users start with 100 Grub currency
- Daily login bonus: 10 Grub (claimable once per 24h)
- Buy/sell shares with market-forces pricing
- 2% appreciation earning when someone buys your shares
- 0.5% daily price decay for inactive stocks
- 1% daily dividends on portfolio value
- Price charts with multiple time ranges
- Leaderboard with 5 ranking categories
- Robinhood-inspired dark mode UI
