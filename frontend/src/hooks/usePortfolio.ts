"use client";

import { useState, useEffect, useCallback } from "react";
import { PortfolioResponse, TransactionWithDetails } from "@/types";
import * as api from "@/lib/api";

export function usePortfolio() {
  const [portfolio, setPortfolio] = useState<PortfolioResponse | null>(null);
  const [history, setHistory] = useState<TransactionWithDetails[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPortfolio = useCallback(async () => {
    try {
      setLoading(true);
      const data = await api.getPortfolio();
      setPortfolio(data);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load portfolio");
    } finally {
      setLoading(false);
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
  }, [fetchPortfolio, fetchHistory]);

  return { portfolio, history, loading, error, refresh: fetchPortfolio };
}
