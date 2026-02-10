"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Button from "@/components/ui/Button";
import * as api from "@/lib/api";

interface DailyClaimButtonProps {
  canClaim: boolean;
  onClaimed: () => void;
}

export default function DailyClaimButton({ canClaim, onClaimed }: DailyClaimButtonProps) {
  const [loading, setLoading] = useState(false);
  const [claimed, setClaimed] = useState(false);

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
    <div className="relative">
      <Button
        variant={canClaim ? "success" : "secondary"}
        size="md"
        onClick={handleClaim}
        loading={loading}
        disabled={!canClaim}
      >
        {canClaim ? "Claim Daily +10 Grub" : "Already Claimed Today"}
      </Button>
      <AnimatePresence>
        {claimed && (
          <motion.div
            initial={{ opacity: 0, y: 0 }}
            animate={{ opacity: 1, y: -30 }}
            exit={{ opacity: 0 }}
            className="absolute -top-2 left-1/2 -translate-x-1/2 text-grub-green font-bold text-sm"
          >
            +10 Grub!
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
