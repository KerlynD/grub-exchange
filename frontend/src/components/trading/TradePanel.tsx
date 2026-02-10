"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import Button from "@/components/ui/Button";
import Modal from "@/components/ui/Modal";
import { formatGrub } from "@/lib/utils";
import * as api from "@/lib/api";

interface TradePanelProps {
  ticker: string;
  currentPrice: number;
  userBalance: number;
  userShares: number;
  onTradeComplete: () => void;
}

type InputMode = "shares" | "grub";

export default function TradePanel({
  ticker,
  currentPrice,
  userBalance,
  userShares,
  onTradeComplete,
}: TradePanelProps) {
  const [tab, setTab] = useState<"buy" | "sell">("buy");
  const [inputMode, setInputMode] = useState<InputMode>("grub");
  const [inputValue, setInputValue] = useState("");
  const [showReview, setShowReview] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const rawValue = parseFloat(inputValue) || 0;

  // Calculate shares and cost based on input mode
  const numShares =
    inputMode === "shares" ? rawValue : currentPrice > 0 ? rawValue / currentPrice : 0;
  const estimatedCost =
    inputMode === "grub" ? rawValue : rawValue * currentPrice;

  const canBuy = tab === "buy" && rawValue > 0 && estimatedCost <= userBalance;
  const canSell = tab === "sell" && rawValue > 0 && numShares <= userShares;

  const handleSubmit = async () => {
    setLoading(true);
    setError(null);
    try {
      const payload =
        inputMode === "grub"
          ? { stock_ticker: ticker, grub_amount: rawValue }
          : { stock_ticker: ticker, num_shares: rawValue };

      if (tab === "buy") {
        await api.buyStock(payload);
      } else {
        await api.sellStock(payload);
      }
      setSuccess(true);
      setInputValue("");
      onTradeComplete();
      setTimeout(() => {
        setSuccess(false);
        setShowReview(false);
      }, 1500);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Trade failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-card-bg rounded-xl border border-border-dark p-4">
      {/* Buy / Sell Tab */}
      <div className="flex bg-dark-bg rounded-lg p-1 mb-4">
        {(["buy", "sell"] as const).map((t) => (
          <button
            key={t}
            onClick={() => {
              setTab(t);
              setInputValue("");
              setError(null);
            }}
            className={`flex-1 py-2 text-sm font-semibold rounded-md transition-colors ${
              tab === t
                ? t === "buy"
                  ? "bg-grub-green text-black"
                  : "bg-grub-red text-white"
                : "text-text-secondary hover:text-white"
            }`}
          >
            {t === "buy" ? "Buy" : "Sell"}
          </button>
        ))}
      </div>

      {/* Input Mode Toggle */}
      <div className="flex bg-dark-bg rounded-lg p-1 mb-4">
        {(["grub", "shares"] as const).map((mode) => (
          <button
            key={mode}
            onClick={() => {
              setInputMode(mode);
              setInputValue("");
            }}
            className={`flex-1 py-1.5 text-xs font-medium rounded-md transition-colors ${
              inputMode === mode
                ? "bg-card-hover text-white"
                : "text-text-secondary hover:text-white"
            }`}
          >
            {mode === "grub" ? "Grub Amount" : "Shares"}
          </button>
        ))}
      </div>

      {/* Input */}
      <div className="mb-4">
        <label className="text-text-secondary text-sm mb-1 block">
          {inputMode === "grub" ? "Grub to spend" : "Number of shares"}
        </label>
        <div className="relative">
          <input
            type="number"
            min="0"
            step={inputMode === "grub" ? "0.01" : "0.0001"}
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            placeholder="0"
            className="w-full bg-dark-bg border border-border-dark rounded-lg px-4 py-3 text-white text-lg font-semibold
              focus:outline-none focus:ring-2 focus:ring-grub-green/50 focus:border-grub-green pr-16"
          />
          <span className="absolute right-4 top-1/2 -translate-y-1/2 text-text-secondary text-sm">
            {inputMode === "grub" ? "GRUB" : "shares"}
          </span>
        </div>
      </div>

      {/* Quick amount buttons for grub mode */}
      {inputMode === "grub" && tab === "buy" && (
        <div className="flex gap-2 mb-4">
          {[5, 10, 25, 50].map((amt) => (
            <button
              key={amt}
              onClick={() => setInputValue(String(Math.min(amt, userBalance)))}
              className="flex-1 py-1.5 text-xs font-medium bg-dark-bg border border-border-dark rounded-lg
                text-text-secondary hover:text-white hover:border-grub-green/50 transition-colors"
            >
              {amt}G
            </button>
          ))}
        </div>
      )}

      {/* Price Info */}
      <div className="space-y-2 mb-4">
        <div className="flex justify-between text-sm">
          <span className="text-text-secondary">Market Price</span>
          <span className="text-white">{formatGrub(currentPrice)} Grub</span>
        </div>
        <div className="flex justify-between text-sm">
          <span className="text-text-secondary">Shares</span>
          <span className="text-white">
            {numShares > 0 ? numShares.toFixed(4) : "0"}
          </span>
        </div>
        <div className="flex justify-between text-sm font-semibold">
          <span className="text-text-secondary">
            {tab === "buy" ? "Est. Cost" : "Est. Proceeds"}
          </span>
          <span className="text-white">{formatGrub(estimatedCost)} Grub</span>
        </div>
        <div className="border-t border-border-dark/50 pt-2 flex justify-between text-sm">
          <span className="text-text-secondary">
            {tab === "buy" ? "Available" : "Shares Owned"}
          </span>
          <span className="text-white">
            {tab === "buy"
              ? `${formatGrub(userBalance)} Grub`
              : `${userShares.toFixed(4)}`}
          </span>
        </div>
      </div>

      {error && <p className="text-grub-red text-sm mb-3">{error}</p>}

      <Button
        variant={tab === "buy" ? "success" : "danger"}
        size="lg"
        className="w-full"
        onClick={() => {
          setError(null);
          setShowReview(true);
        }}
        disabled={tab === "buy" ? !canBuy : !canSell}
      >
        Review Order
      </Button>

      {/* Review Modal */}
      <Modal
        isOpen={showReview}
        onClose={() => {
          if (!loading) setShowReview(false);
        }}
        title="Review Order"
      >
        {success ? (
          <motion.div
            initial={{ scale: 0.8, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            className="text-center py-6"
          >
            <div className="w-16 h-16 rounded-full bg-grub-green/20 flex items-center justify-center mx-auto mb-3">
              <svg className="w-8 h-8 text-grub-green" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <p className="text-white font-semibold text-lg">Order Executed!</p>
          </motion.div>
        ) : (
          <div>
            <div className="bg-dark-bg rounded-lg p-4 mb-4 space-y-2">
              <p className="text-white">
                {tab === "buy" ? "Buying" : "Selling"}{" "}
                <span className="font-bold">{numShares.toFixed(4)}</span> shares
                of <span className="font-bold text-grub-green">{ticker}</span>
              </p>
              <div className="flex justify-between text-sm">
                <span className="text-text-secondary">Price per share</span>
                <span className="text-white">{formatGrub(currentPrice)}</span>
              </div>
              <div className="flex justify-between text-sm font-semibold border-t border-border-dark pt-2">
                <span className="text-text-secondary">
                  {tab === "buy" ? "Total Cost" : "Total Proceeds"}
                </span>
                <span className="text-white">
                  {formatGrub(estimatedCost)} Grub
                </span>
              </div>
            </div>
            {error && <p className="text-grub-red text-sm mb-3">{error}</p>}
            <div className="flex gap-3">
              <Button
                variant="secondary"
                className="flex-1"
                onClick={() => setShowReview(false)}
                disabled={loading}
              >
                Cancel
              </Button>
              <Button
                variant={tab === "buy" ? "success" : "danger"}
                className="flex-1"
                onClick={handleSubmit}
                loading={loading}
              >
                Confirm {tab === "buy" ? "Buy" : "Sell"}
              </Button>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}
