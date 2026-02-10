-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    ticker TEXT UNIQUE NOT NULL,
    bio TEXT DEFAULT '',
    current_share_price DOUBLE PRECISION DEFAULT 10.0,
    shares_outstanding INTEGER DEFAULT 1000,
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Balances table
CREATE TABLE IF NOT EXISTS balances (
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    grub_balance DOUBLE PRECISION DEFAULT 100.0,
    last_daily_claim TIMESTAMPTZ
);

-- Portfolios table (who owns shares of whom)
CREATE TABLE IF NOT EXISTS portfolios (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users(id),
    stock_user_id INTEGER NOT NULL REFERENCES users(id),
    num_shares DOUBLE PRECISION NOT NULL,
    avg_purchase_price DOUBLE PRECISION NOT NULL,
    UNIQUE(owner_id, stock_user_id)
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    buyer_id INTEGER NOT NULL REFERENCES users(id),
    stock_user_id INTEGER NOT NULL REFERENCES users(id),
    transaction_type TEXT NOT NULL,
    num_shares DOUBLE PRECISION NOT NULL,
    price_per_share DOUBLE PRECISION NOT NULL,
    total_grub DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

-- Price history for charting
CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    price DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

-- Portfolio value snapshots for the dashboard graph
CREATE TABLE IF NOT EXISTS portfolio_snapshots (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    total_value DOUBLE PRECISION NOT NULL,
    grub_balance DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_portfolio_snapshots_user_time ON portfolio_snapshots(user_id, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_portfolios_owner ON portfolios(owner_id);
CREATE INDEX IF NOT EXISTS idx_portfolios_stock ON portfolios(stock_user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_stock ON transactions(stock_user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_buyer ON transactions(buyer_id);
CREATE INDEX IF NOT EXISTS idx_price_history_user_time ON price_history(user_id, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_users_ticker ON users(ticker);
