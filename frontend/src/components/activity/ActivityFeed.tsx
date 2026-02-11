"use client";

import { useEffect, useState, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Notification } from "@/types";
import * as api from "@/lib/api";
import { ChevronDown, ChevronUp } from "lucide-react";

function timeAgo(dateStr: string): string {
  const now = new Date();
  const then = new Date(dateStr);
  const seconds = Math.floor((now.getTime() - then.getTime()) / 1000);

  if (seconds < 60) return "just now";
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  return `${Math.floor(seconds / 86400)}d ago`;
}

const COLLAPSED_COUNT = 5;

export default function ActivityFeed() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState(false);

  const fetchNotifications = useCallback(async () => {
    try {
      const data = await api.getNotifications();
      setNotifications(data.notifications || []);
      setUnreadCount(data.unread_count || 0);
    } catch {
      // silently fail
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchNotifications();
    // Poll every 15 seconds for new notifications
    const interval = setInterval(fetchNotifications, 15000);
    return () => clearInterval(interval);
  }, [fetchNotifications]);

  const handleMarkRead = async () => {
    try {
      await api.markNotificationsRead();
      setUnreadCount(0);
      setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
    } catch {
      // silently fail
    }
  };

  if (loading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3].map((i) => (
          <div key={i} className="bg-card-bg rounded-lg h-14 animate-pulse" />
        ))}
      </div>
    );
  }

  if (notifications.length === 0) {
    return (
      <div className="text-center py-8 text-text-secondary text-sm">
        No activity yet. When someone trades your stock, it&apos;ll show up
        here!
      </div>
    );
  }

  const displayedNotifications = expanded
    ? notifications
    : notifications.slice(0, COLLAPSED_COUNT);
  const hasMore = notifications.length > COLLAPSED_COUNT;

  return (
    <div className="space-y-2">
      {unreadCount > 0 && (
        <div className="flex items-center justify-between mb-2">
          <span className="text-xs text-grub-green font-medium">
            {unreadCount} new
          </span>
          <button
            onClick={handleMarkRead}
            className="text-xs text-text-secondary hover:text-white transition-colors"
          >
            Mark all read
          </button>
        </div>
      )}
      <AnimatePresence>
        {displayedNotifications.map((notif, i) => (
          <motion.div
            key={notif.id}
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -20 }}
            transition={{ delay: i * 0.05 }}
            className={`flex items-start gap-3 p-3 rounded-lg transition-colors ${
              notif.read
                ? "bg-card-bg/50"
                : "bg-card-bg border-l-2 border-grub-green"
            }`}
          >
            <div className="flex-shrink-0 mt-0.5">
              {notif.type === "trade_buy" ? (
                <div className="w-8 h-8 rounded-full bg-grub-green/20 flex items-center justify-center">
                  <span className="text-grub-green text-xs font-bold">B</span>
                </div>
              ) : notif.type === "trade_sell" ? (
                <div className="w-8 h-8 rounded-full bg-grub-red/20 flex items-center justify-center">
                  <span className="text-grub-red text-xs font-bold">S</span>
                </div>
              ) : (
                <div className="w-8 h-8 rounded-full bg-card-hover flex items-center justify-center">
                  <span className="text-text-secondary text-xs">!</span>
                </div>
              )}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-white text-sm leading-snug">
                {notif.message}
              </p>
              <p className="text-text-secondary text-xs mt-0.5">
                {timeAgo(notif.created_at)}
              </p>
            </div>
          </motion.div>
        ))}
      </AnimatePresence>

      {hasMore && (
        <button
          onClick={() => setExpanded(!expanded)}
          className="w-full py-2 text-xs font-medium text-text-secondary hover:text-white
            flex items-center justify-center gap-1 transition-colors"
        >
          {expanded ? (
            <>
              Show less <ChevronUp size={14} />
            </>
          ) : (
            <>
              Show {notifications.length - COLLAPSED_COUNT} more <ChevronDown size={14} />
            </>
          )}
        </button>
      )}
    </div>
  );
}
