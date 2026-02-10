"use client";

import { LineChart, Line, YAxis, ResponsiveContainer } from "recharts";

interface SparklineProps {
  data: number[];
  color?: string;
  height?: number;
}

export default function Sparkline({ data, color, height = 40 }: SparklineProps) {
  if (!data || data.length === 0) return null;

  const isPositive = data.length >= 2 ? data[data.length - 1] >= data[0] : true;
  const lineColor = color || (isPositive ? "#00C805" : "#FF5000");

  // Use only the last 30 data points for a tighter view
  const trimmed = data.length > 30 ? data.slice(-30) : data;
  const chartData = trimmed.map((value, index) => ({ value, index }));

  // Compute tight Y domain with 10% padding so movements are clearly visible
  const min = Math.min(...trimmed);
  const max = Math.max(...trimmed);
  const range = max - min || max * 0.01 || 0.1; // fallback if flat
  const padding = range * 0.1;
  const yDomain: [number, number] = [min - padding, max + padding];

  return (
    <ResponsiveContainer width="100%" height={height}>
      <LineChart data={chartData}>
        <YAxis domain={yDomain} hide />
        <Line
          type="monotone"
          dataKey="value"
          stroke={lineColor}
          strokeWidth={1.5}
          dot={false}
          isAnimationActive={false}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
