"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { StockListItem } from "@/types";
import * as api from "@/lib/api";

const POLL_INTERVAL = 10_000; // 10 seconds

export function useMarket() {
  const [stocks, setStocks] = useState<StockListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const initialFetch = useRef(true);

  const fetchStocks = useCallback(async () => {
    try {
      // Only show loading spinner on the very first fetch
      if (initialFetch.current) setLoading(true);
      const data = await api.getStocks();
      setStocks(data.stocks || []);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load stocks");
    } finally {
      setLoading(false);
      initialFetch.current = false;
    }
  }, []);

  useEffect(() => {
    fetchStocks();
    const interval = setInterval(fetchStocks, POLL_INTERVAL);
    return () => clearInterval(interval);
  }, [fetchStocks]);

  return { stocks, loading, error, refresh: fetchStocks };
}
