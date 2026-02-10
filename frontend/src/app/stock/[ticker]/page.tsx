"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams } from "next/navigation";
import { motion } from "framer-motion";
import AppLayout from "@/components/layout/AppLayout";
import PriceChart from "@/components/charts/PriceChart";
import TradePanel from "@/components/trading/TradePanel";
import TransactionHistory from "@/components/trading/TransactionHistory";
import Card from "@/components/ui/Card";
import { useAuth } from "@/contexts/AuthContext";
import { StockDetail, PriceHistory } from "@/types";
import * as api from "@/lib/api";
import {
  formatGrub,
  formatPercent,
  formatNumber,
  getPLColor,
} from "@/lib/utils";

type TimeRange = "1D" | "1W" | "1M" | "ALL";

export default function StockDetailPage() {
  const params = useParams();
  const ticker = params.ticker as string;
  const { user, refreshUser } = useAuth();
  const [stock, setStock] = useState<StockDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [timeRange, setTimeRange] = useState<TimeRange>("1M");
  const [userShares, setUserShares] = useState(0);

  const fetchData = useCallback(async () => {
    try {
      setLoading(true);
      const data = await api.getStockDetail(ticker);
      setStock(data);

      // Get user's shares of this stock
      const portfolio = await api.getPortfolio();
      const holding = portfolio.holdings?.find((h) => h.ticker === ticker);
      setUserShares(holding ? holding.num_shares : 0);
    } catch {
      // handle error
    } finally {
      setLoading(false);
    }
  }, [ticker]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleTradeComplete = () => {
    fetchData();
    refreshUser();
  };

  const filterByTimeRange = (data: PriceHistory[]): PriceHistory[] => {
    if (!data) return [];
    const now = new Date();
    let since: Date;
    switch (timeRange) {
      case "1D":
        since = new Date(now.getTime() - 24 * 60 * 60 * 1000);
        break;
      case "1W":
        since = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
        break;
      case "1M":
        since = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
        break;
      default:
        return data;
    }
    return data.filter((p) => new Date(p.timestamp) >= since);
  };

  if (loading) {
    return (
      <AppLayout>
        <div className="space-y-4">
          <div className="bg-card-bg rounded-xl h-12 w-48 animate-pulse" />
          <div className="bg-card-bg rounded-xl h-80 animate-pulse" />
          <div className="bg-card-bg rounded-xl h-64 animate-pulse" />
        </div>
      </AppLayout>
    );
  }

  if (!stock) {
    return (
      <AppLayout>
        <div className="text-center py-20 text-text-secondary">
          Stock not found
        </div>
      </AppLayout>
    );
  }

  const filteredHistory = filterByTimeRange(stock.price_history);

  return (
    <AppLayout>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="space-y-6"
      >
        {/* Header */}
        <div>
          <div className="flex items-center gap-3 mb-2">
            <div className="w-12 h-12 rounded-full bg-grub-green/20 flex items-center justify-center">
              <span className="text-grub-green font-bold text-lg">
                {stock.user.ticker.charAt(0)}
              </span>
            </div>
            <div>
              <h1 className="text-2xl font-bold text-white">
                {stock.user.ticker}
              </h1>
              <p className="text-text-secondary text-sm">
                {stock.user.username}
              </p>
            </div>
          </div>
          <div className="flex items-baseline gap-3">
            <p className="text-3xl font-bold text-white">
              {formatGrub(stock.user.current_share_price)}
              <span className="text-lg font-normal text-text-secondary ml-1">
                Grub
              </span>
            </p>
            <span
              className={`text-sm font-semibold px-2 py-0.5 rounded ${
                stock.change_24h_percent >= 0
                  ? "bg-grub-green/10 text-grub-green"
                  : "bg-grub-red/10 text-grub-red"
              }`}
            >
              {formatPercent(stock.change_24h_percent)} today
            </span>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Chart Section */}
          <div className="lg:col-span-2 space-y-4">
            {/* Time Range Selector */}
            <div className="flex gap-2">
              {(["1D", "1W", "1M", "ALL"] as TimeRange[]).map((range) => (
                <button
                  key={range}
                  onClick={() => setTimeRange(range)}
                  className={`px-3 py-1.5 text-xs font-medium rounded-lg transition-colors ${
                    timeRange === range
                      ? "bg-grub-green text-black"
                      : "bg-card-bg text-text-secondary hover:text-white"
                  }`}
                >
                  {range}
                </button>
              ))}
            </div>

            {/* Price Chart */}
            <Card className="p-2">
              <PriceChart data={filteredHistory} height={350} />
            </Card>

            {/* Market Stats */}
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
              {[
                {
                  label: "Market Cap",
                  value: formatNumber(stock.market_cap),
                },
                {
                  label: "24h Volume",
                  value: formatNumber(stock.volume_24h),
                },
                {
                  label: "All-Time High",
                  value: formatGrub(stock.all_time_high),
                },
                {
                  label: "All-Time Low",
                  value: formatGrub(stock.all_time_low),
                },
              ].map((stat) => (
                <Card key={stat.label}>
                  <p className="text-text-secondary text-xs mb-1">
                    {stat.label}
                  </p>
                  <p className="text-white font-semibold">{stat.value}</p>
                </Card>
              ))}
            </div>

            {/* About Section */}
            {stock.user.bio && (
              <Card>
                <h3 className="text-white font-semibold mb-2">
                  About {stock.user.ticker}
                </h3>
                <p className="text-text-secondary text-sm leading-relaxed">
                  {stock.user.bio}
                </p>
              </Card>
            )}

            {/* Recent Trades */}
            <Card>
              <h3 className="text-white font-semibold mb-3">Recent Trades</h3>
              <TransactionHistory transactions={stock.recent_trades} />
            </Card>
          </div>

          {/* Trade Panel */}
          <div className="lg:col-span-1">
            <div className="lg:sticky lg:top-8">
              {user && stock.user.ticker === user.ticker ? (
                <Card className="text-center py-8">
                  <div className="w-12 h-12 rounded-full bg-card-hover flex items-center justify-center mx-auto mb-3">
                    <span className="text-text-secondary text-xl">$</span>
                  </div>
                  <p className="text-white font-semibold mb-1">
                    This is your stock!
                  </p>
                  <p className="text-text-secondary text-sm">
                    You can&apos;t trade your own shares.
                    <br />
                    Share your ticker so others can invest in you.
                  </p>
                </Card>
              ) : (
                <TradePanel
                  ticker={stock.user.ticker}
                  currentPrice={stock.user.current_share_price}
                  userBalance={user?.grub_balance || 0}
                  userShares={userShares}
                  onTradeComplete={handleTradeComplete}
                />
              )}
            </div>
          </div>
        </div>
      </motion.div>
    </AppLayout>
  );
}
