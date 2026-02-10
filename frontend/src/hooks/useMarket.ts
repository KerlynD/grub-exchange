"use client";

import { useState, useEffect, useCallback } from "react";
import { StockListItem } from "@/types";
import * as api from "@/lib/api";

export function useMarket() {
  const [stocks, setStocks] = useState<StockListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStocks = useCallback(async () => {
    try {
      setLoading(true);
      const data = await api.getStocks();
      setStocks(data.stocks || []);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load stocks");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStocks();
  }, [fetchStocks]);

  return { stocks, loading, error, refresh: fetchStocks };
}
