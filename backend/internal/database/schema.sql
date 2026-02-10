-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    ticker TEXT UNIQUE NOT NULL,
    current_share_price REAL DEFAULT 10.0,
    shares_outstanding INTEGER DEFAULT 1000,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Balances table
CREATE TABLE IF NOT EXISTS balances (
    user_id INTEGER PRIMARY KEY,
    grub_balance REAL DEFAULT 100.0,
    last_daily_claim TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Portfolios table (who owns shares of whom)
CREATE TABLE IF NOT EXISTS portfolios (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    stock_user_id INTEGER NOT NULL,
    num_shares REAL NOT NULL,
    avg_purchase_price REAL NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (stock_user_id) REFERENCES users(id),
    UNIQUE(owner_id, stock_user_id)
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    buyer_id INTEGER NOT NULL,
    stock_user_id INTEGER NOT NULL,
    transaction_type TEXT NOT NULL,
    num_shares REAL NOT NULL,
    price_per_share REAL NOT NULL,
    total_grub REAL NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (buyer_id) REFERENCES users(id),
    FOREIGN KEY (stock_user_id) REFERENCES users(id)
);

-- Price history for charting
CREATE TABLE IF NOT EXISTS price_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    price REAL NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Portfolio value snapshots for the dashboard graph
CREATE TABLE IF NOT EXISTS portfolio_snapshots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    total_value REAL NOT NULL,
    grub_balance REAL NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
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
