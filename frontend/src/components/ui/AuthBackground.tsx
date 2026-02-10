"use client";

import { motion } from "framer-motion";

const circles = [
  { size: 350, x: "5%", y: "15%", delay: 0, duration: 18, opacity: 0.15 },
  { size: 250, x: "70%", y: "5%", delay: 1.5, duration: 22, opacity: 0.12 },
  { size: 200, x: "80%", y: "65%", delay: 3, duration: 16, opacity: 0.18 },
  { size: 300, x: "15%", y: "70%", delay: 0.5, duration: 20, opacity: 0.14 },
  { size: 220, x: "45%", y: "45%", delay: 2, duration: 24, opacity: 0.1 },
  { size: 180, x: "55%", y: "80%", delay: 4, duration: 19, opacity: 0.16 },
  { size: 400, x: "30%", y: "10%", delay: 1, duration: 26, opacity: 0.08 },
];

export default function AuthBackground() {
  return (
    <div className="fixed inset-0 overflow-hidden pointer-events-none">
      {circles.map((circle, i) => (
        <motion.div
          key={i}
          className="absolute rounded-full"
          style={{
            width: circle.size,
            height: circle.size,
            left: circle.x,
            top: circle.y,
            background: `radial-gradient(circle, rgba(0, 200, 5, ${circle.opacity}) 0%, rgba(0, 200, 5, ${circle.opacity * 0.3}) 50%, transparent 70%)`,
            filter: "blur(40px)",
          }}
          initial={{ opacity: 0, scale: 0.6 }}
          animate={{
            opacity: [0.5, 1, 0.5],
            scale: [1, 1.15, 1],
            x: [0, 40, -30, 0],
            y: [0, -35, 20, 0],
          }}
          transition={{
            duration: circle.duration,
            delay: circle.delay,
            repeat: Infinity,
            ease: "easeInOut",
          }}
        />
      ))}

      {/* Subtle grid overlay */}
      <div
        className="absolute inset-0 opacity-[0.04]"
        style={{
          backgroundImage:
            "linear-gradient(rgba(255,255,255,.15) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,.15) 1px, transparent 1px)",
          backgroundSize: "60px 60px",
        }}
      />
    </div>
  );
}
