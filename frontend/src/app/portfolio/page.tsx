"use client";

import { motion } from "framer-motion";
import AppLayout from "@/components/layout/AppLayout";
import BalanceDisplay from "@/components/portfolio/BalanceDisplay";
import PortfolioCard from "@/components/portfolio/PortfolioCard";
import TransactionHistory from "@/components/trading/TransactionHistory";
import Card from "@/components/ui/Card";
import { usePortfolio } from "@/hooks/usePortfolio";

export default function PortfolioPage() {
  const { portfolio, history, loading } = usePortfolio();

  if (loading) {
    return (
      <AppLayout>
        <div className="space-y-4">
          <div className="bg-card-bg rounded-xl h-20 animate-pulse" />
          <div className="bg-card-bg rounded-xl h-16 animate-pulse" />
          <div className="bg-card-bg rounded-xl h-16 animate-pulse" />
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="space-y-8"
      >
        {/* Header */}
        <div>
          <h1 className="text-2xl font-bold text-white mb-4">Portfolio</h1>
          {portfolio && (
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <Card>
                <p className="text-text-secondary text-xs mb-1">Grub Balance</p>
                <p className="text-white font-bold text-xl">
                  {portfolio.grub_balance.toFixed(2)}
                </p>
              </Card>
              <Card>
                <p className="text-text-secondary text-xs mb-1">
                  Portfolio Value
                </p>
                <p className="text-white font-bold text-xl">
                  {portfolio.total_value.toFixed(2)}
                </p>
              </Card>
              <Card>
                <p className="text-text-secondary text-xs mb-1">
                  Total P&L
                </p>
                <p
                  className={`font-bold text-xl ${
                    portfolio.total_pl >= 0
                      ? "text-grub-green"
                      : "text-grub-red"
                  }`}
                >
                  {portfolio.total_pl >= 0 ? "+" : ""}
                  {portfolio.total_pl.toFixed(2)} (
                  {portfolio.total_pl_percent.toFixed(2)}%)
                </p>
              </Card>
            </div>
          )}
        </div>

        {/* Holdings */}
        <div>
          <h2 className="text-lg font-bold text-white mb-3">Holdings</h2>
          {portfolio?.holdings && portfolio.holdings.length > 0 ? (
            <div className="space-y-2">
              {portfolio.holdings.map((h) => (
                <PortfolioCard key={h.stock_user_id} holding={h} />
              ))}
            </div>
          ) : (
            <Card>
              <p className="text-text-secondary text-center py-6 text-sm">
                No holdings yet. Start trading to build your portfolio!
              </p>
            </Card>
          )}
        </div>

        {/* Transaction History */}
        <div>
          <h2 className="text-lg font-bold text-white mb-3">
            Transaction History
          </h2>
          <Card>
            <TransactionHistory transactions={history} />
          </Card>
        </div>
      </motion.div>
    </AppLayout>
  );
}
