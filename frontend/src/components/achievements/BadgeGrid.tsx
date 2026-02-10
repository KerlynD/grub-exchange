"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Achievement, UserAchievement } from "@/types";
import * as api from "@/lib/api";

export default function BadgeGrid() {
  const [allAchievements, setAllAchievements] = useState<Achievement[]>([]);
  const [earned, setEarned] = useState<UserAchievement[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .getAchievements()
      .then((data) => {
        setAllAchievements(data.all || []);
        setEarned(data.earned || []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="grid grid-cols-2 gap-3">
        {[1, 2, 3, 4].map((i) => (
          <div
            key={i}
            className="bg-card-bg rounded-xl h-24 animate-pulse"
          />
        ))}
      </div>
    );
  }

  const earnedIds = new Set(earned.map((e) => e.achievement_id));

  return (
    <div className="grid grid-cols-2 gap-3">
      {allAchievements.map((achievement, i) => {
        const isEarned = earnedIds.has(achievement.id);
        const earnedData = earned.find(
          (e) => e.achievement_id === achievement.id
        );

        return (
          <motion.div
            key={achievement.id}
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: i * 0.1 }}
            className={`relative rounded-xl p-3 transition-all ${
              isEarned
                ? "bg-grub-green/10 border border-grub-green/30"
                : "bg-card-bg border border-border-dark opacity-50"
            }`}
          >
            <div className="text-2xl mb-1.5">{achievement.icon}</div>
            <p
              className={`text-sm font-semibold ${
                isEarned ? "text-white" : "text-text-secondary"
              }`}
            >
              {achievement.name}
            </p>
            <p className="text-text-secondary text-xs mt-0.5 leading-snug">
              {achievement.description}
            </p>
            {isEarned && earnedData && (
              <p className="text-grub-green text-[10px] mt-1.5 font-medium">
                Earned{" "}
                {new Date(earnedData.earned_at).toLocaleDateString("en-US", {
                  month: "short",
                  day: "numeric",
                })}
              </p>
            )}
            {!isEarned && (
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="w-6 h-6 rounded-full bg-card-hover/80 flex items-center justify-center">
                  <span className="text-text-secondary text-xs">?</span>
                </div>
              </div>
            )}
          </motion.div>
        );
      })}
    </div>
  );
}
