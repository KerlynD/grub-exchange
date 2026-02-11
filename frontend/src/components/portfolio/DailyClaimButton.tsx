"use client";

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Button from "@/components/ui/Button";
import * as api from "@/lib/api";

interface DailyClaimButtonProps {
  canClaim: boolean;
  lastDailyClaim?: string;
  onClaimed: () => void;
}

function formatCountdown(ms: number): string {
  if (ms <= 0) return "00:00:00";
  const totalSeconds = Math.floor(ms / 1000);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  return `${hours.toString().padStart(2, "0")}:${minutes
    .toString()
    .padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
}

export default function DailyClaimButton({
  canClaim,
  lastDailyClaim,
  onClaimed,
}: DailyClaimButtonProps) {
  const [loading, setLoading] = useState(false);
  const [claimed, setClaimed] = useState(false);
  const [countdown, setCountdown] = useState<string | null>(null);

  useEffect(() => {
    if (canClaim || !lastDailyClaim) {
      setCountdown(null);
      return;
    }

    const claimTime = new Date(lastDailyClaim).getTime();
    const nextClaimTime = claimTime + 24 * 60 * 60 * 1000; // 24 hours later

    const updateCountdown = () => {
      const now = Date.now();
      const remaining = nextClaimTime - now;
      if (remaining <= 0) {
        setCountdown(null);
        onClaimed(); // Trigger refresh when timer reaches 0
      } else {
        setCountdown(formatCountdown(remaining));
      }
    };

    updateCountdown();
    const interval = setInterval(updateCountdown, 1000);
    return () => clearInterval(interval);
  }, [canClaim, lastDailyClaim, onClaimed]);

  const handleClaim = async () => {
    if (!canClaim || loading) return;
    setLoading(true);
    try {
      await api.claimDaily();
      setClaimed(true);
      onClaimed();
      setTimeout(() => setClaimed(false), 3000);
    } catch {
      // already claimed
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative flex flex-col items-center gap-1">
      <Button
        variant={canClaim ? "success" : "secondary"}
        size="md"
        onClick={handleClaim}
        loading={loading}
        disabled={!canClaim}
      >
        {canClaim ? "Claim Daily Grub" : "Already Claimed Today"}
      </Button>
      {countdown && (
        <span className="text-xs text-text-secondary">
          Next claim in: <span className="font-mono">{countdown}</span>
        </span>
      )}
      <AnimatePresence>
        {claimed && (
          <motion.div
            initial={{ opacity: 0, y: 0 }}
            animate={{ opacity: 1, y: -30 }}
            exit={{ opacity: 0 }}
            className="absolute -top-2 left-1/2 -translate-x-1/2 text-grub-green font-bold text-sm"
          >
            +Grub!
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
