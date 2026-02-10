"use client";

import { motion } from "framer-motion";

interface BalanceDisplayProps {
  totalPortfolioValue: number;
  totalPL?: number;
  totalPLPercent?: number;
}

export default function BalanceDisplay({
  totalPortfolioValue,
  totalPL,
  totalPLPercent,
}: BalanceDisplayProps) {
  return (
    <div>
      <p className="text-text-secondary text-sm mb-1">Total Portfolio Value</p>
      <motion.p
        key={totalPortfolioValue.toFixed(2)}
        initial={{ opacity: 0.5, y: -5 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-4xl font-bold text-white"
      >
        {totalPortfolioValue.toFixed(2)}{" "}
        <span className="text-lg font-normal text-text-secondary">GRUB</span>
      </motion.p>
      {totalPL !== undefined && (
        <p
          className={`text-sm mt-1 ${
            totalPL >= 0 ? "text-grub-green" : "text-grub-red"
          }`}
        >
          {totalPL >= 0 ? "+" : ""}
          {totalPL.toFixed(2)} ({totalPLPercent?.toFixed(2) || "0.00"}%) all
          time
        </p>
      )}
    </div>
  );
}
