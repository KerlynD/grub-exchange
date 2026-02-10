export interface User {
  id: number;
  username: string;
  email: string;
  ticker: string;
  bio: string;
  current_share_price: number;
  shares_outstanding: number;
  grub_balance: number;
  last_login?: string;
  created_at: string;
}

export interface StockListItem {
  id: number;
  username: string;
  ticker: string;
  current_share_price: number;
  change_24h_percent: number;
  sparkline_data: number[];
}

export interface StockDetail {
  user: User;
  price_history: PriceHistory[];
  recent_trades: TransactionWithDetails[];
  change_24h: number;
  change_24h_percent: number;
  market_cap: number;
  volume_24h: number;
  all_time_high: number;
  all_time_low: number;
}

export interface PriceHistory {
  id: number;
  user_id: number;
  price: number;
  timestamp: string;
}

export interface PortfolioHolding {
  ticker: string;
  username: string;
  stock_user_id: number;
  num_shares: number;
  avg_purchase_price: number;
  current_price: number;
  total_value: number;
  profit_loss: number;
  profit_loss_percent: number;
}

export interface PortfolioResponse {
  grub_balance: number;
  total_value: number;
  total_pl: number;
  total_pl_percent: number;
  holdings: PortfolioHolding[] | null;
  can_claim_daily: boolean;
}

export interface TransactionWithDetails {
  id: number;
  buyer_username: string;
  stock_ticker: string;
  transaction_type: "BUY" | "SELL";
  num_shares: number;
  price_per_share: number;
  total_grub: number;
  timestamp: string;
}

export interface LeaderboardEntry {
  rank: number;
  user_id: number;
  username: string;
  ticker: string;
  value: number;
  change?: number;
}

export interface LeaderboardData {
  most_valuable: LeaderboardEntry[];
  biggest_gainers: LeaderboardEntry[];
  biggest_losers: LeaderboardEntry[];
  richest_traders: LeaderboardEntry[];
  best_performance: LeaderboardEntry[];
}

export interface TradeRequest {
  stock_ticker: string;
  num_shares?: number;
  grub_amount?: number;
}

export interface MarketSnapshot {
  id: number;
  total_market_cap: number;
  total_invested: number;
  total_cash: number;
  total_grub: number;
  timestamp: string;
}

export interface MarketOverview {
  total_market_cap: number;
  total_grub: number;
  total_invested: number;
  total_cash: number;
  invested_percent: number;
  total_stocks: number;
  history: MarketSnapshot[];
}

export interface PortfolioSnapshot {
  id: number;
  user_id: number;
  total_value: number;
  grub_balance: number;
  timestamp: string;
}

export interface Notification {
  id: number;
  user_id: number;
  type: string;
  message: string;
  actor_username: string;
  stock_ticker: string;
  num_shares: number;
  read: boolean;
  created_at: string;
}

export interface Achievement {
  id: string;
  name: string;
  description: string;
  icon: string;
}

export interface UserAchievement {
  id: number;
  user_id: number;
  achievement_id: string;
  name: string;
  description: string;
  icon: string;
  earned_at: string;
}
