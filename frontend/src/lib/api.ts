import {
  User,
  StockListItem,
  StockDetail,
  PortfolioResponse,
  TransactionWithDetails,
  LeaderboardData,
  TradeRequest,
  PortfolioSnapshot,
} from "@/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function fetchAPI<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  const res = await fetch(`${API_URL}${endpoint}`, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    ...options,
  });

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: "Request failed" }));
    throw new Error(error.error || `HTTP ${res.status}`);
  }

  return res.json();
}

// Auth
export async function register(data: {
  username: string;
  email: string;
  password: string;
  first_name: string;
}): Promise<{ user: User }> {
  return fetchAPI("/api/auth/register", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function login(data: {
  email: string;
  password: string;
}): Promise<{ user: User }> {
  return fetchAPI("/api/auth/login", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function logout(): Promise<void> {
  await fetchAPI("/api/auth/logout", { method: "POST" });
}

export async function getMe(): Promise<{ user: User }> {
  return fetchAPI("/api/auth/me");
}

// Stocks
export async function getStocks(): Promise<{ stocks: StockListItem[] }> {
  return fetchAPI("/api/stocks");
}

export async function getStockDetail(ticker: string): Promise<StockDetail> {
  return fetchAPI(`/api/stocks/${ticker}`);
}

// Trading
export async function buyStock(
  data: TradeRequest
): Promise<{ message: string; transaction: TransactionWithDetails }> {
  return fetchAPI("/api/trade/buy", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function sellStock(
  data: TradeRequest
): Promise<{ message: string; transaction: TransactionWithDetails }> {
  return fetchAPI("/api/trade/sell", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

// Portfolio
export async function getPortfolio(): Promise<PortfolioResponse> {
  return fetchAPI("/api/portfolio");
}

export async function claimDaily(): Promise<{
  message: string;
  new_balance: number;
}> {
  return fetchAPI("/api/portfolio/claim-daily", { method: "POST" });
}

export async function getPortfolioHistory(): Promise<{
  transactions: TransactionWithDetails[];
}> {
  return fetchAPI("/api/portfolio/history");
}

export async function getPortfolioGraph(): Promise<{
  snapshots: PortfolioSnapshot[];
}> {
  return fetchAPI("/api/portfolio/graph");
}

// Profile
export async function getProfile(): Promise<{ user: User }> {
  return fetchAPI("/api/profile");
}

export async function updateProfile(data: {
  bio: string;
}): Promise<{ user: User }> {
  return fetchAPI("/api/profile", {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

// Market
export async function getLeaderboard(): Promise<LeaderboardData> {
  return fetchAPI("/api/leaderboard");
}

export async function getRecentTransactions(): Promise<{
  transactions: TransactionWithDetails[];
}> {
  return fetchAPI("/api/transactions");
}
