"use client";

import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { LeaderboardEntry } from "@/types";
import { formatGrub, formatPercent } from "@/lib/utils";

interface LeaderboardTableProps {
  entries: LeaderboardEntry[];
  valueLabel: string;
  showChange?: boolean;
  isPercent?: boolean;
}

export default function LeaderboardTable({
  entries,
  valueLabel,
  showChange = false,
  isPercent = false,
}: LeaderboardTableProps) {
  const router = useRouter();

  if (!entries || entries.length === 0) {
    return (
      <div className="text-text-secondary text-center py-8 text-sm">
        No data yet
      </div>
    );
  }

  return (
    <div className="space-y-1">
      {/* Header */}
      <div className="flex items-center px-4 py-2 text-xs text-text-secondary font-medium">
        <span className="w-8">#</span>
        <span className="flex-1">Name</span>
        <span className="w-24 text-right">{valueLabel}</span>
        {showChange && <span className="w-20 text-right">Change</span>}
      </div>

      {/* Rows */}
      {entries.map((entry, i) => (
        <motion.div
          key={entry.user_id}
          initial={{ opacity: 0, x: -10 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: i * 0.05 }}
          onClick={() => router.push(`/stock/${entry.ticker}`)}
          className="flex items-center px-4 py-3 rounded-lg hover:bg-card-hover cursor-pointer transition-colors"
        >
          <span className="w-8 text-text-secondary font-semibold text-sm">
            {entry.rank}
          </span>
          <div className="flex-1 flex items-center gap-2">
            <div className="w-8 h-8 rounded-full bg-grub-green/20 flex items-center justify-center">
              <span className="text-grub-green font-bold text-xs">
                {entry.ticker.charAt(0)}
              </span>
            </div>
            <div>
              <p className="text-white text-sm font-medium">{entry.ticker}</p>
              <p className="text-text-secondary text-xs">{entry.username}</p>
            </div>
          </div>
          <span className="w-24 text-right text-white font-semibold text-sm">
            {isPercent ? formatPercent(entry.value) : formatGrub(entry.value)}
          </span>
          {showChange && entry.change !== undefined && (
            <span
              className={`w-20 text-right text-sm font-medium ${
                entry.change >= 0 ? "text-grub-green" : "text-grub-red"
              }`}
            >
              {formatPercent(entry.change)}
            </span>
          )}
        </motion.div>
      ))}
    </div>
  );
}
