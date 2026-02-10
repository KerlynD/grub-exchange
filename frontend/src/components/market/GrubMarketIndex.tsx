"use client";

import { useEffect, useState, useRef, useCallback } from "react";
import { motion } from "framer-motion";
import Card from "@/components/ui/Card";
import ShootingStar from "@/components/icons/ShootingStar";
import { MarketOverview } from "@/types";
import * as api from "@/lib/api";
import { formatNumber } from "@/lib/utils";

export default function GrubMarketIndex() {
  const [data, setData] = useState<MarketOverview | null>(null);
  const [loading, setLoading] = useState(true);
  const initialFetch = useRef(true);

  const fetchData = useCallback(async () => {
    try {
      const overview = await api.getMarketOverview();
      setData(overview);
    } catch {
      // silently fail
    } finally {
      if (initialFetch.current) {
        setLoading(false);
        initialFetch.current = false;
      }
    }
  }, []);

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 15_000);
    return () => clearInterval(interval);
  }, [fetchData]);

  if (loading) {
    return <div className="bg-card-bg rounded-xl h-32 animate-pulse" />;
  }

  if (!data) return null;

  return (
    <Card className="bg-gradient-to-br from-card-bg to-grub-green/5 border border-grub-green/20">
      <div className="flex items-start justify-between">
        <div>
          <div className="flex items-center gap-2 mb-3">
            <ShootingStar size={20} className="text-grub-green" />
            <h2 className="text-lg font-bold text-white">The Grub Market</h2>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <div>
              <p className="text-text-secondary text-xs">Total Market Cap</p>
              <motion.p
                key={data.total_market_cap.toFixed(0)}
                initial={{ opacity: 0.5 }}
                animate={{ opacity: 1 }}
                className="text-white font-bold text-lg"
              >
                {formatNumber(data.total_market_cap)}
              </motion.p>
            </div>
            <div>
              <p className="text-text-secondary text-xs">Total Grub</p>
              <p className="text-white font-semibold">
                {formatNumber(data.total_grub)}
              </p>
            </div>
            <div>
              <p className="text-text-secondary text-xs">In Market</p>
              <p className="text-grub-green font-semibold">
                {formatNumber(data.total_invested)}
                <span className="text-text-secondary text-xs ml-1">
                  ({data.invested_percent.toFixed(1)}%)
                </span>
              </p>
            </div>
            <div>
              <p className="text-text-secondary text-xs">Held as Cash</p>
              <p className="text-white font-semibold">
                {formatNumber(data.total_cash)}
                <span className="text-text-secondary text-xs ml-1">
                  ({(100 - data.invested_percent).toFixed(1)}%)
                </span>
              </p>
            </div>
          </div>
        </div>

        <div className="hidden sm:block text-right">
          <p className="text-text-secondary text-xs mb-1">Active Stocks</p>
          <p className="text-white font-bold text-2xl">{data.total_stocks}</p>
        </div>
      </div>

      {/* Market flow bar */}
      <div className="mt-4">
        <div className="h-2 bg-card-hover rounded-full overflow-hidden">
          <motion.div
            initial={{ width: 0 }}
            animate={{ width: `${data.invested_percent}%` }}
            transition={{ duration: 1, ease: "easeOut" }}
            className="h-full bg-gradient-to-r from-grub-green to-grub-green/60 rounded-full"
          />
        </div>
        <div className="flex justify-between mt-1">
          <span className="text-[10px] text-grub-green">Invested</span>
          <span className="text-[10px] text-text-secondary">Cash</span>
        </div>
      </div>
    </Card>
  );
}
