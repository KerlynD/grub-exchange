"use client";

import { useRouter } from "next/navigation";
import Card from "@/components/ui/Card";
import Sparkline from "@/components/charts/Sparkline";
import { StockListItem } from "@/types";
import { formatGrub, formatPercent, getPLColor } from "@/lib/utils";

interface StockCardProps {
  stock: StockListItem;
}

export default function StockCard({ stock }: StockCardProps) {
  const router = useRouter();

  return (
    <Card
      hover
      onClick={() => router.push(`/stock/${stock.ticker}`)}
      className="flex items-center justify-between gap-4"
    >
      <div className="flex items-center gap-3 min-w-0">
        <div className="w-10 h-10 rounded-full bg-grub-green/20 flex-shrink-0 flex items-center justify-center">
          <span className="text-grub-green font-bold text-sm">
            {stock.ticker.charAt(0)}
          </span>
        </div>
        <div className="min-w-0">
          <p className="text-white font-semibold truncate">{stock.ticker}</p>
          <p className="text-text-secondary text-xs truncate">{stock.username}</p>
        </div>
      </div>

      <div className="w-20 flex-shrink-0">
        <Sparkline data={stock.sparkline_data} height={32} />
      </div>

      <div className="text-right flex-shrink-0">
        <p className="text-white font-semibold">
          {formatGrub(stock.current_share_price)}
        </p>
        <p className={`text-xs ${getPLColor(stock.change_24h_percent)}`}>
          {formatPercent(stock.change_24h_percent)}
        </p>
      </div>
    </Card>
  );
}
