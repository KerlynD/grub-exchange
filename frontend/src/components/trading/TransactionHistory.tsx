"use client";

import { TransactionWithDetails } from "@/types";
import { formatGrub, formatDate } from "@/lib/utils";

interface TransactionHistoryProps {
  transactions: TransactionWithDetails[];
}

export default function TransactionHistory({ transactions }: TransactionHistoryProps) {
  if (!transactions || transactions.length === 0) {
    return (
      <div className="text-text-secondary text-center py-8 text-sm">
        No transactions yet
      </div>
    );
  }

  return (
    <div className="space-y-2">
      {transactions.map((txn, i) => (
        <div
          key={txn.id || i}
          className="flex items-center justify-between py-3 border-b border-border-dark/50 last:border-0"
        >
          <div className="flex items-center gap-3">
            <div
              className={`w-8 h-8 rounded-full flex items-center justify-center ${
                txn.transaction_type === "BUY"
                  ? "bg-grub-green/20"
                  : "bg-grub-red/20"
              }`}
            >
              <span
                className={`text-xs font-bold ${
                  txn.transaction_type === "BUY"
                    ? "text-grub-green"
                    : "text-grub-red"
                }`}
              >
                {txn.transaction_type === "BUY" ? "B" : "S"}
              </span>
            </div>
            <div>
              <p className="text-white text-sm font-medium">
                {txn.buyer_username}{" "}
                <span className="text-text-secondary">
                  {txn.transaction_type === "BUY" ? "bought" : "sold"}
                </span>
              </p>
              <p className="text-text-secondary text-xs">
                {txn.num_shares.toFixed(2)} shares @ {formatGrub(txn.price_per_share)}
              </p>
            </div>
          </div>
          <div className="text-right">
            {/* Amounts are relative to the stock: BUY = money flowing in (+), SELL = money flowing out (-) */}
            <p
              className={`text-sm font-semibold ${
                txn.transaction_type === "BUY"
                  ? "text-grub-green"
                  : "text-grub-red"
              }`}
            >
              {txn.transaction_type === "BUY" ? "+" : "-"}
              {formatGrub(txn.total_grub)}
            </p>
            <p className="text-text-secondary text-xs">
              {formatDate(txn.timestamp)}
            </p>
          </div>
        </div>
      ))}
    </div>
  );
}
