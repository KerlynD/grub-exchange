"use client";

import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import AppLayout from "@/components/layout/AppLayout";
import LeaderboardTable from "@/components/leaderboard/LeaderboardTable";
import Card from "@/components/ui/Card";
import { LeaderboardData } from "@/types";
import * as api from "@/lib/api";

const tabs = [
  { id: "valuable", label: "Most Valuable" },
  { id: "gainers", label: "Top Gainers" },
  { id: "losers", label: "Top Losers" },
  { id: "richest", label: "Richest" },
  { id: "performance", label: "Best P&L" },
];

export default function LeaderboardPage() {
  const [data, setData] = useState<LeaderboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("valuable");

  useEffect(() => {
    api
      .getLeaderboard()
      .then(setData)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const getTabContent = () => {
    if (!data) return null;

    switch (activeTab) {
      case "valuable":
        return (
          <LeaderboardTable
            entries={data.most_valuable}
            valueLabel="Price"
          />
        );
      case "gainers":
        return (
          <LeaderboardTable
            entries={data.biggest_gainers}
            valueLabel="Price"
            showChange
          />
        );
      case "losers":
        return (
          <LeaderboardTable
            entries={data.biggest_losers}
            valueLabel="Price"
            showChange
          />
        );
      case "richest":
        return (
          <LeaderboardTable
            entries={data.richest_traders}
            valueLabel="Balance"
          />
        );
      case "performance":
        return (
          <LeaderboardTable
            entries={data.best_performance}
            valueLabel="P&L %"
            isPercent
          />
        );
      default:
        return null;
    }
  };

  return (
    <AppLayout>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="space-y-6"
      >
        <div>
          <h1 className="text-2xl font-bold text-white mb-1">Leaderboard</h1>
          <p className="text-text-secondary text-sm">
            See who&apos;s on top of the Grub Exchange
          </p>
        </div>

        {/* Tabs */}
        <div className="flex gap-2 overflow-x-auto pb-2">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-2 text-sm font-medium rounded-lg whitespace-nowrap transition-colors ${
                activeTab === tab.id
                  ? "bg-grub-green text-black"
                  : "bg-card-bg text-text-secondary hover:text-white border border-border-dark"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Content */}
        <Card className="p-0 overflow-hidden">
          {loading ? (
            <div className="space-y-3 p-4">
              {[1, 2, 3, 4, 5].map((i) => (
                <div
                  key={i}
                  className="h-12 bg-dark-bg rounded animate-pulse"
                />
              ))}
            </div>
          ) : (
            <AnimatePresence mode="wait">
              <motion.div
                key={activeTab}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                transition={{ duration: 0.2 }}
                className="py-2"
              >
                {getTabContent()}
              </motion.div>
            </AnimatePresence>
          )}
        </Card>
      </motion.div>
    </AppLayout>
  );
}
