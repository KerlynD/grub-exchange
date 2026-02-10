"use client";

import { motion } from "framer-motion";
import AppLayout from "@/components/layout/AppLayout";
import BalanceDisplay from "@/components/portfolio/BalanceDisplay";
import DailyClaimButton from "@/components/portfolio/DailyClaimButton";
import PortfolioCard from "@/components/portfolio/PortfolioCard";
import StockCard from "@/components/market/StockCard";
import PortfolioChart from "@/components/charts/PortfolioChart";
import ActivityFeed from "@/components/activity/ActivityFeed";
import BadgeGrid from "@/components/achievements/BadgeGrid";
import Card from "@/components/ui/Card";
import { useAuth } from "@/contexts/AuthContext";
import { usePortfolio } from "@/hooks/usePortfolio";
import { useMarket } from "@/hooks/useMarket";

export default function DashboardPage() {
  const { refreshUser } = useAuth();
  const {
    portfolio,
    loading: pLoading,
    refresh: refreshPortfolio,
  } = usePortfolio();
  const { stocks, loading: mLoading } = useMarket();

  const handleClaimed = () => {
    refreshPortfolio();
    refreshUser();
  };

  const sorted = [...stocks].sort(
    (a, b) => b.change_24h_percent - a.change_24h_percent
  );
  const topGainers = sorted.slice(0, 5);

  const totalPortfolioValue = portfolio
    ? portfolio.grub_balance + portfolio.total_value
    : 0;

  return (
    <AppLayout>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="space-y-8"
      >
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
          {portfolio && (
            <BalanceDisplay
              totalPortfolioValue={totalPortfolioValue}
              totalPL={portfolio.total_pl}
              totalPLPercent={portfolio.total_pl_percent}
            />
          )}
          {portfolio && (
            <DailyClaimButton
              canClaim={portfolio.can_claim_daily}
              onClaimed={handleClaimed}
            />
          )}
        </div>

        {/* Main content grid: left column (charts/holdings) + right column (activity/badges) */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left column */}
          <div className="lg:col-span-2 space-y-8">
            {/* Portfolio Value Graph */}
            <Card className="p-2">
              <div className="px-3 pt-2 pb-1">
                <p className="text-text-secondary text-xs font-medium">
                  Investing
                </p>
              </div>
              <PortfolioChart height={220} />
              {/* Buying Power */}
              {portfolio && (
                <div className="px-3 pb-3 pt-2 border-t border-border-dark mt-2">
                  <div className="flex items-center justify-between">
                    <p className="text-text-secondary text-sm">Buying Power</p>
                    <p className="text-white font-semibold text-sm">
                      {portfolio.grub_balance.toFixed(2)} GRUB
                    </p>
                  </div>
                </div>
              )}
            </Card>

            {/* My Holdings */}
            <div>
              <h2 className="text-lg font-bold text-white mb-3">
                Your Holdings
              </h2>
              {pLoading ? (
                <div className="space-y-2">
                  {[1, 2, 3].map((i) => (
                    <div
                      key={i}
                      className="bg-card-bg rounded-xl h-16 animate-pulse"
                    />
                  ))}
                </div>
              ) : portfolio?.holdings && portfolio.holdings.length > 0 ? (
                <div className="space-y-2">
                  {portfolio.holdings.map((h) => (
                    <PortfolioCard key={h.stock_user_id} holding={h} />
                  ))}
                </div>
              ) : (
                <Card>
                  <p className="text-text-secondary text-center py-4 text-sm">
                    You don&apos;t own any stocks yet. Head to the{" "}
                    <span className="text-grub-green">Market</span> to start
                    trading!
                  </p>
                </Card>
              )}
            </div>

            {/* Top Stocks */}
            <div>
              <h2 className="text-lg font-bold text-white mb-3">Top Movers</h2>
              {mLoading ? (
                <div className="space-y-2">
                  {[1, 2, 3].map((i) => (
                    <div
                      key={i}
                      className="bg-card-bg rounded-xl h-16 animate-pulse"
                    />
                  ))}
                </div>
              ) : topGainers.length > 0 ? (
                <div className="space-y-2">
                  {topGainers.map((s) => (
                    <StockCard key={s.id} stock={s} />
                  ))}
                </div>
              ) : (
                <Card>
                  <p className="text-text-secondary text-center py-4 text-sm">
                    No stocks available yet
                  </p>
                </Card>
              )}
            </div>
          </div>

          {/* Right column: Activity Feed + Achievements */}
          <div className="lg:col-span-1 space-y-6">
            {/* Activity Feed */}
            <Card>
              <h3 className="text-white font-semibold mb-3 flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-grub-green animate-pulse" />
                Activity
              </h3>
              <ActivityFeed />
            </Card>

            {/* Achievements */}
            <Card>
              <h3 className="text-white font-semibold mb-3">Achievements</h3>
              <BadgeGrid />
            </Card>
          </div>
        </div>
      </motion.div>
    </AppLayout>
  );
}
