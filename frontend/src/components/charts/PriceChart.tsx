"use client";

import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { PriceHistory } from "@/types";

interface PriceChartProps {
  data: PriceHistory[];
  height?: number;
}

function formatTime(timestamp: string) {
  const d = new Date(timestamp);
  return d.toLocaleTimeString("en-US", { hour: "2-digit", minute: "2-digit" });
}

function CustomTooltip({ active, payload }: any) {
  if (active && payload && payload.length) {
    return (
      <div className="bg-card-bg border border-border-dark rounded-lg px-3 py-2 shadow-lg">
        <p className="text-white font-semibold text-sm">
          {payload[0].value.toFixed(2)} Grub
        </p>
        <p className="text-text-secondary text-xs">
          {new Date(payload[0].payload.timestamp).toLocaleString()}
        </p>
      </div>
    );
  }
  return null;
}

export default function PriceChart({ data, height = 300 }: PriceChartProps) {
  if (!data || data.length === 0) {
    return (
      <div
        className="flex items-center justify-center text-text-secondary"
        style={{ height }}
      >
        No price data available
      </div>
    );
  }

  const isPositive = data[data.length - 1].price >= data[0].price;
  const color = isPositive ? "#00C805" : "#FF5000";

  return (
    <ResponsiveContainer width="100%" height={height}>
      <AreaChart data={data}>
        <defs>
          <linearGradient id="priceGradient" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor={color} stopOpacity={0.3} />
            <stop offset="100%" stopColor={color} stopOpacity={0} />
          </linearGradient>
        </defs>
        <XAxis
          dataKey="timestamp"
          tickFormatter={formatTime}
          stroke="#8E8E93"
          tick={{ fontSize: 11 }}
          axisLine={false}
          tickLine={false}
        />
        <YAxis
          domain={["auto", "auto"]}
          stroke="#8E8E93"
          tick={{ fontSize: 11 }}
          axisLine={false}
          tickLine={false}
          width={50}
        />
        <Tooltip content={<CustomTooltip />} />
        <Area
          type="monotone"
          dataKey="price"
          stroke={color}
          strokeWidth={2}
          fill="url(#priceGradient)"
          isAnimationActive={true}
          animationDuration={600}
        />
      </AreaChart>
    </ResponsiveContainer>
  );
}
