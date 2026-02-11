"use client";

import { useEffect, useState, useCallback, useRef } from "react";
import { useRouter } from "next/navigation";
import { StockPost } from "@/types";
import * as api from "@/lib/api";
import { ThumbsUp, ThumbsDown } from "lucide-react";

function timeAgo(dateStr: string): string {
  const now = new Date();
  const date = new Date(dateStr);
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (seconds < 60) return "just now";
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  return `${Math.floor(seconds / 86400)}d ago`;
}

export default function RecentNews() {
  const [posts, setPosts] = useState<StockPost[]>([]);
  const [loading, setLoading] = useState(true);
  const initialFetch = useRef(true);
  const router = useRouter();

  const fetchPosts = useCallback(async () => {
    try {
      const data = await api.getRecentPosts();
      setPosts(data.posts || []);
    } catch {
      // silently fail
    } finally {
      if (initialFetch.current) {
        setLoading(false);
        initialFetch.current = false;
      }
    }
  }, []);

  useEffect(() => {
    fetchPosts();
    const interval = setInterval(fetchPosts, 15_000);
    return () => clearInterval(interval);
  }, [fetchPosts]);

  if (loading) {
    return (
      <div className="space-y-2">
        {[1, 2, 3].map((i) => (
          <div key={i} className="bg-dark-bg rounded-lg h-14 animate-pulse" />
        ))}
      </div>
    );
  }

  if (posts.length === 0) {
    return (
      <p className="text-text-secondary text-center py-4 text-sm">
        No news yet. Visit a stock to post the first story!
      </p>
    );
  }

  return (
    <div className="space-y-2 max-h-[320px] overflow-y-auto pr-1">
      {posts.map((post) => {
        const net = post.likes - post.dislikes;
        return (
          <div
            key={post.id}
            onClick={() => router.push(`/stock/${post.stock_ticker}`)}
            className="bg-dark-bg rounded-lg p-2.5 cursor-pointer hover:bg-card-hover transition-colors"
          >
            <div className="flex items-center gap-2 mb-1">
              <span
                className="text-grub-green text-xs font-bold cursor-pointer"
              >
                ${post.stock_ticker}
              </span>
              <span className="text-text-secondary text-[10px]">
                {post.author_username} &middot; {timeAgo(post.created_at)}
              </span>
              {net !== 0 && (
                <span
                  className={`text-[10px] px-1 py-0.5 rounded ${
                    net > 0
                      ? "bg-grub-green/10 text-grub-green"
                      : "bg-grub-red/10 text-grub-red"
                  }`}
                >
                  {net > 0 ? "Bullish" : "Bearish"}
                </span>
              )}
            </div>
            <p className="text-text-secondary text-xs leading-relaxed line-clamp-2">
              {post.content}
            </p>
            <div className="flex items-center gap-3 mt-1.5">
              <span className="flex items-center gap-1 text-[10px] text-text-secondary">
                <ThumbsUp size={10} />
                {post.likes}
              </span>
              <span className="flex items-center gap-1 text-[10px] text-text-secondary">
                <ThumbsDown size={10} />
                {post.dislikes}
              </span>
            </div>
          </div>
        );
      })}
    </div>
  );
}
