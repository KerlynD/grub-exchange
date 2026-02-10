# Deploying Grub Exchange

**Stack**: Frontend on Vercel, Backend on Fly.io, Database on Supabase (PostgreSQL)

---

## 1. Set Up Supabase (Database)

1. Go to [supabase.com](https://supabase.com) and create a new project
2. Choose a region close to your Fly.io region (e.g., **US East** if using `ewr`)
3. Set a strong database password — save it
4. Once the project is created, go to **Settings → Database**
5. Copy the **Connection string (URI)** — it looks like:
   ```
   postgresql://postgres.[ref]:[password]@aws-0-us-east-1.pooler.supabase.com:6543/postgres
   ```
6. **Important**: Use the **Session mode** connection string (port `5432`) not the pooler (`6543`) for Fly.io, since the Go backend manages its own connection pool. Or use the pooler in **Transaction mode** (port `6543`) — either works.

That's it. The backend automatically creates all tables on first startup via migrations.

---

## 2. Deploy Backend to Fly.io

### Install Fly CLI

```bash
# macOS
brew install flyctl

# Windows
powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"

# Linux
curl -L https://fly.io/install.sh | sh
```

### Login & Launch

```bash
cd backend

# Login to Fly.io
fly auth login

# Launch the app (first time only)
fly launch --no-deploy
# When prompted: choose a name like "grub-exchange-api", pick region "ewr" (or closest to your Supabase)
# Say NO to databases — we're using Supabase
```

### Set Secrets (Environment Variables)

```bash
# Your Supabase connection string
fly secrets set DATABASE_URL="postgresql://postgres.[ref]:[password]@aws-0-us-east-1.pooler.supabase.com:5432/postgres"

# Generate a strong random secret for JWT
fly secrets set JWT_SECRET="$(openssl rand -base64 32)"

# Your Vercel frontend URL (update after deploying frontend)
fly secrets set FRONTEND_URL="https://grub-exchange.vercel.app"
```

### Deploy

```bash
fly deploy
```

Your backend will be live at `https://grub-exchange-api.fly.dev` (or whatever name you chose).

### Verify

```bash
curl https://grub-exchange-api.fly.dev/api/stocks
```

Should return `{"stocks": []}` or similar if the DB is empty.

---

## 3. Deploy Frontend to Vercel

### Install Vercel CLI (optional — you can also use the dashboard)

```bash
npm install -g vercel
```

### Option A: Vercel Dashboard (Recommended)

1. Push your code to GitHub
2. Go to [vercel.com](https://vercel.com) → **New Project**
3. Import your GitHub repo
4. Set the **Root Directory** to `frontend`
5. Set the **Framework Preset** to `Next.js`
6. Add this **Environment Variable**:
   | Key | Value |
   |-----|-------|
   | `NEXT_PUBLIC_API_URL` | `https://grub-exchange-api.fly.dev` |
7. Click **Deploy**

### Option B: Vercel CLI

```bash
cd frontend
vercel --prod
# When prompted:
#   - Set root directory: ./
#   - Framework: Next.js
#   - Add env var NEXT_PUBLIC_API_URL=https://grub-exchange-api.fly.dev
```

Your frontend will be live at `https://grub-exchange.vercel.app` (or your custom domain).

---

## 4. Update CORS (Important!)

After deploying the frontend, go back and update the Fly.io secret with the actual Vercel URL:

```bash
cd backend
fly secrets set FRONTEND_URL="https://your-actual-vercel-url.vercel.app"
```

This ensures the backend allows requests from your frontend domain.

---

## 5. Seed Data (Optional)

To seed Ivan (the test user), run the seed script locally pointing at your Supabase DB:

```bash
cd backend
DATABASE_URL="postgresql://postgres.[ref]:[password]@aws-0-us-east-1.pooler.supabase.com:5432/postgres" go run scripts/seed.go
```

---

## 6. Custom Domain (Optional)

### Vercel (Frontend)
1. Go to your Vercel project → **Settings → Domains**
2. Add your domain (e.g., `grub.exchange`)
3. Update DNS as instructed

### Fly.io (Backend)
```bash
fly certs add api.grub.exchange
```
Then add a CNAME record for `api.grub.exchange` → `grub-exchange-api.fly.dev`

After setting custom domains, update:
- `FRONTEND_URL` on Fly.io to your frontend domain
- `NEXT_PUBLIC_API_URL` on Vercel to your backend domain

---

## Environment Variables Summary

### Fly.io (Backend)
| Variable | Example | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgresql://postgres...` | Supabase connection string |
| `JWT_SECRET` | `random-32-char-string` | Secret for signing auth tokens |
| `FRONTEND_URL` | `https://grub-exchange.vercel.app` | CORS allowed origin |
| `PORT` | `8080` | Server port (set in fly.toml) |

### Vercel (Frontend)
| Variable | Example | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | `https://grub-exchange-api.fly.dev` | Backend API URL |

---

## Architecture Notes

- **No CGo dependency**: The PostgreSQL driver (`lib/pq`) is pure Go — no GCC needed, builds anywhere
- **Connection pooling**: Backend uses 10 max connections, Supabase free tier supports 50
- **Migrations**: Run automatically on backend startup — no manual SQL needed
- **Market maker**: Runs as a goroutine inside the backend — no separate worker process
- **Scheduled jobs**: Daily decay and dividends run inside the backend process too
- **Auth**: JWT stored in httpOnly cookies. Make sure both frontend and backend are on HTTPS in production for cookies to work cross-origin

---

## Troubleshooting

**CORS errors**: Make sure `FRONTEND_URL` on Fly.io exactly matches your Vercel URL (including `https://`, no trailing slash).

**Cookies not working**: httpOnly cookies require both frontend and backend on HTTPS. If using custom domains, ensure SSL is configured on both.

**Database connection refused**: Check that the Supabase connection string is correct and you're using the right port (5432 for direct, 6543 for pooler).

**Backend logs**:
```bash
fly logs --app grub-exchange-api
```
