"use client";

import { useRouter } from "next/navigation";
import Card from "@/components/ui/Card";
import { PortfolioHolding } from "@/types";
import { formatGrub, formatPercent, getPLColor } from "@/lib/utils";

interface PortfolioCardProps {
  holding: PortfolioHolding;
}

export default function PortfolioCard({ holding }: PortfolioCardProps) {
  const router = useRouter();

  return (
    <Card
      hover
      onClick={() => router.push(`/stock/${holding.ticker}`)}
      className="flex items-center justify-between"
    >
      <div className="flex items-center gap-3">
        <div className="w-10 h-10 rounded-full bg-grub-green/20 flex items-center justify-center">
          <span className="text-grub-green font-bold text-sm">
            {holding.ticker.charAt(0)}
          </span>
        </div>
        <div>
          <p className="text-white font-semibold">{holding.ticker}</p>
          <p className="text-text-secondary text-xs">
            {holding.num_shares} shares
          </p>
        </div>
      </div>
      <div className="text-right">
        <p className="text-white font-semibold">
          {formatGrub(holding.total_value)}
        </p>
        <p className={`text-xs ${getPLColor(holding.profit_loss)}`}>
          {holding.profit_loss >= 0 ? "+" : ""}
          {formatGrub(holding.profit_loss)} ({formatPercent(holding.profit_loss_percent)})
        </p>
      </div>
    </Card>
  );
}
