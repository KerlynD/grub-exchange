"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { PortfolioResponse, TransactionWithDetails } from "@/types";
import * as api from "@/lib/api";

const POLL_INTERVAL = 15_000; // 15 seconds

export function usePortfolio() {
  const [portfolio, setPortfolio] = useState<PortfolioResponse | null>(null);
  const [history, setHistory] = useState<TransactionWithDetails[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const initialFetch = useRef(true);

  const fetchPortfolio = useCallback(async () => {
    try {
      if (initialFetch.current) setLoading(true);
      const data = await api.getPortfolio();
      setPortfolio(data);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load portfolio");
    } finally {
      setLoading(false);
      initialFetch.current = false;
    }
  }, []);

  const fetchHistory = useCallback(async () => {
    try {
      const data = await api.getPortfolioHistory();
      setHistory(data.transactions || []);
    } catch {
      // silently fail
    }
  }, []);

  useEffect(() => {
    fetchPortfolio();
    fetchHistory();
    const interval = setInterval(fetchPortfolio, POLL_INTERVAL);
    return () => clearInterval(interval);
  }, [fetchPortfolio, fetchHistory]);

  return { portfolio, history, loading, error, refresh: fetchPortfolio };
}
