"use client";

import { useEffect, useState } from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { PortfolioSnapshot } from "@/types";
import * as api from "@/lib/api";

function CustomTooltip({ active, payload }: any) {
  if (active && payload && payload.length) {
    return (
      <div className="bg-card-bg border border-border-dark rounded-lg px-3 py-2 shadow-lg">
        <p className="text-white font-semibold text-sm">
          {payload[0].value.toFixed(2)} Grub
        </p>
        <p className="text-text-secondary text-xs">
          {new Date(payload[0].payload.timestamp).toLocaleDateString()}
        </p>
      </div>
    );
  }
  return null;
}

export default function PortfolioChart({ height = 200 }: { height?: number }) {
  const [data, setData] = useState<PortfolioSnapshot[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .getPortfolioGraph()
      .then((res) => setData(res.snapshots || []))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="animate-pulse bg-card-bg rounded-xl" style={{ height }} />
    );
  }

  if (data.length < 2) {
    return (
      <div
        className="flex items-center justify-center text-text-secondary text-sm rounded-xl bg-card-bg"
        style={{ height }}
      >
        Portfolio graph will appear as market activity occurs
      </div>
    );
  }

  const isPositive = data[data.length - 1].total_value >= data[0].total_value;
  const color = isPositive ? "#00C805" : "#FF5000";

  return (
    <ResponsiveContainer width="100%" height={height}>
      <AreaChart data={data}>
        <defs>
          <linearGradient id="portfolioGrad" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor={color} stopOpacity={0.3} />
            <stop offset="100%" stopColor={color} stopOpacity={0} />
          </linearGradient>
        </defs>
        <XAxis
          dataKey="timestamp"
          tickFormatter={(v) =>
            new Date(v).toLocaleDateString("en-US", {
              month: "short",
              day: "numeric",
            })
          }
          stroke="#8E8E93"
          tick={{ fontSize: 10 }}
          axisLine={false}
          tickLine={false}
        />
        <YAxis
          domain={["auto", "auto"]}
          stroke="#8E8E93"
          tick={{ fontSize: 10 }}
          axisLine={false}
          tickLine={false}
          width={45}
        />
        <Tooltip content={<CustomTooltip />} />
        <Area
          type="monotone"
          dataKey="total_value"
          stroke={color}
          strokeWidth={2}
          fill="url(#portfolioGrad)"
          isAnimationActive
          animationDuration={600}
        />
      </AreaChart>
    </ResponsiveContainer>
  );
}
