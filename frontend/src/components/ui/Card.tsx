"use client";

import { motion } from "framer-motion";
import { cn } from "@/lib/utils";

interface CardProps {
  children: React.ReactNode;
  className?: string;
  hover?: boolean;
  onClick?: () => void;
}

export default function Card({ children, className, hover = false, onClick }: CardProps) {
  return (
    <motion.div
      whileHover={hover ? { scale: 1.01, y: -2 } : undefined}
      className={cn(
        "bg-card-bg rounded-xl p-4 border border-border-dark/50",
        hover && "cursor-pointer transition-colors hover:border-border-dark",
        className
      )}
      onClick={onClick}
    >
      {children}
    </motion.div>
  );
}
