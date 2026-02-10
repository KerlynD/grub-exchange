"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import AppLayout from "@/components/layout/AppLayout";
import StockCard from "@/components/market/StockCard";
import GrubMarketIndex from "@/components/market/GrubMarketIndex";
import { useMarket } from "@/hooks/useMarket";
import { Search } from "lucide-react";

export default function MarketPage() {
  const { stocks, loading } = useMarket();
  const [search, setSearch] = useState("");

  const filtered = stocks.filter(
    (s) =>
      s.ticker.toLowerCase().includes(search.toLowerCase()) ||
      s.username.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <AppLayout>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="space-y-6"
      >
        <div>
          <h1 className="text-2xl font-bold text-white mb-1">Market</h1>
          <p className="text-text-secondary text-sm">
            All tradeable stocks on Grub Exchange
          </p>
        </div>

        {/* Grub Market Index */}
        <GrubMarketIndex />

        {/* Search */}
        <div className="relative">
          <Search
            className="absolute left-3 top-1/2 -translate-y-1/2 text-text-secondary"
            size={18}
          />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search by name or ticker..."
            className="w-full bg-card-bg border border-border-dark rounded-lg pl-10 pr-4 py-2.5 text-white
              placeholder-text-secondary focus:outline-none focus:ring-2 focus:ring-grub-green/50 focus:border-grub-green
              transition-all duration-150"
          />
        </div>

        {/* Stock List */}
        {loading ? (
          <div className="space-y-2">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div
                key={i}
                className="bg-card-bg rounded-xl h-16 animate-pulse"
              />
            ))}
          </div>
        ) : filtered.length > 0 ? (
          <div className="space-y-2">
            {filtered.map((stock) => (
              <StockCard key={stock.id} stock={stock} />
            ))}
          </div>
        ) : (
          <div className="text-text-secondary text-center py-12">
            {search ? "No stocks match your search" : "No stocks available yet"}
          </div>
        )}
      </motion.div>
    </AppLayout>
  );
}
