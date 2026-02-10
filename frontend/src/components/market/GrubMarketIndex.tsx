"use client";

import { useEffect, useState, useRef, useCallback } from "react";
import { motion } from "framer-motion";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
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
    return <div className="bg-card-bg rounded-xl h-48 animate-pulse" />;
  }

  if (!data) return null;

  const chartData = (data.history || []).map((snap) => ({
    time: new Date(snap.timestamp).toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
    }),
    invested: snap.total_invested,
    cash: snap.total_cash,
    marketCap: snap.total_market_cap,
  }));

  const hasHistory = chartData.length > 1;

  return (
    <Card className="bg-gradient-to-br from-card-bg to-grub-green/5 border border-grub-green/20 p-0 overflow-hidden">
      <div className="p-4 pb-0">
        <div className="flex items-start justify-between mb-3">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <ShootingStar size={20} className="text-grub-green" />
              <h2 className="text-lg font-bold text-white">The Grub Market</h2>
            </div>
            <div className="flex items-baseline gap-2">
              <motion.p
                key={data.total_market_cap.toFixed(0)}
                initial={{ opacity: 0.5 }}
                animate={{ opacity: 1 }}
                className="text-2xl font-bold text-white"
              >
                {formatNumber(data.total_market_cap)}
              </motion.p>
              <span className="text-text-secondary text-xs">Total Market Cap</span>
            </div>
          </div>
          <div className="text-right hidden sm:block">
            <p className="text-text-secondary text-xs">Active Stocks</p>
            <p className="text-white font-bold text-xl">{data.total_stocks}</p>
          </div>
        </div>

        {/* Stats row */}
        <div className="grid grid-cols-3 gap-3 mb-3">
          <div>
            <p className="text-text-secondary text-[10px] uppercase tracking-wide">Total Grub</p>
            <p className="text-white font-semibold text-sm">{formatNumber(data.total_grub)}</p>
          </div>
          <div>
            <p className="text-text-secondary text-[10px] uppercase tracking-wide">In Market</p>
            <p className="text-grub-green font-semibold text-sm">
              {formatNumber(data.total_invested)}
              <span className="text-text-secondary text-[10px] ml-1">
                ({data.invested_percent.toFixed(1)}%)
              </span>
            </p>
          </div>
          <div>
            <p className="text-text-secondary text-[10px] uppercase tracking-wide">Held as Cash</p>
            <p className="text-white font-semibold text-sm">
              {formatNumber(data.total_cash)}
              <span className="text-text-secondary text-[10px] ml-1">
                ({(100 - data.invested_percent).toFixed(1)}%)
              </span>
            </p>
          </div>
        </div>
      </div>

      {/* Chart */}
      {hasHistory ? (
        <div className="h-40">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData} margin={{ top: 0, right: 0, left: 0, bottom: 0 }}>
              <defs>
                <linearGradient id="investedGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#00C805" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#00C805" stopOpacity={0} />
                </linearGradient>
                <linearGradient id="cashGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#6B7280" stopOpacity={0.2} />
                  <stop offset="95%" stopColor="#6B7280" stopOpacity={0} />
                </linearGradient>
              </defs>
              <XAxis
                dataKey="time"
                tick={{ fill: "#6B7280", fontSize: 10 }}
                axisLine={false}
                tickLine={false}
                interval="preserveStartEnd"
              />
              <YAxis hide domain={["dataMin", "dataMax"]} />
              <Tooltip
                contentStyle={{
                  background: "#1A1A2E",
                  border: "1px solid #2D2D3D",
                  borderRadius: "8px",
                  fontSize: "12px",
                }}
                labelStyle={{ color: "#6B7280" }}
                formatter={(value: number, name: string) => [
                  formatNumber(value),
                  name === "invested" ? "In Market" : name === "cash" ? "Cash" : "Market Cap",
                ]}
              />
              <Area
                type="monotone"
                dataKey="invested"
                stroke="#00C805"
                strokeWidth={2}
                fill="url(#investedGrad)"
                isAnimationActive={false}
              />
              <Area
                type="monotone"
                dataKey="cash"
                stroke="#6B7280"
                strokeWidth={1}
                fill="url(#cashGrad)"
                strokeDasharray="4 4"
                isAnimationActive={false}
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      ) : (
        <div className="px-4 pb-4">
          {/* Fallback bar when no history yet */}
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
          <p className="text-text-secondary text-[10px] text-center mt-2">
            Chart data will appear as snapshots are recorded
          </p>
        </div>
      )}
    </Card>
  );
}
