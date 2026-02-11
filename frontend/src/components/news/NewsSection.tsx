"use client";

import { useEffect, useState, useCallback, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Card from "@/components/ui/Card";
import { StockPost } from "@/types";
import * as api from "@/lib/api";
import { ThumbsUp, ThumbsDown, Send, Newspaper } from "lucide-react";

interface NewsSectionProps {
  ticker: string;
}

function timeAgo(dateStr: string): string {
  const now = new Date();
  const date = new Date(dateStr);
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (seconds < 60) return "just now";
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  return `${Math.floor(seconds / 86400)}d ago`;
}

export default function NewsSection({ ticker }: NewsSectionProps) {
  const [posts, setPosts] = useState<StockPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [content, setContent] = useState("");
  const [posting, setPosting] = useState(false);
  const [votingIds, setVotingIds] = useState<Set<number>>(new Set());
  const initialFetch = useRef(true);

  const fetchPosts = useCallback(async () => {
    try {
      const data = await api.getStockPosts(ticker);
      setPosts(data.posts || []);
    } catch {
      // silently fail
    } finally {
      if (initialFetch.current) {
        setLoading(false);
        initialFetch.current = false;
      }
    }
  }, [ticker]);

  useEffect(() => {
    fetchPosts();
    const interval = setInterval(fetchPosts, 15_000);
    return () => clearInterval(interval);
  }, [fetchPosts]);

  const handlePost = async () => {
    if (!content.trim() || posting) return;
    setPosting(true);
    try {
      await api.createStockPost(ticker, content.trim());
      setContent("");
      await fetchPosts();
    } catch {
      // handle error
    } finally {
      setPosting(false);
    }
  };

  const handleVote = async (postId: number, voteType: 1 | -1) => {
    if (votingIds.has(postId)) return;
    setVotingIds((prev) => new Set(prev).add(postId));

    // Optimistic update
    setPosts((prev) =>
      prev.map((p) => {
        if (p.id !== postId) return p;
        const wasLiked = p.user_vote === 1;
        const wasDisliked = p.user_vote === -1;

        let newLikes = p.likes;
        let newDislikes = p.dislikes;
        let newVote: number = voteType;

        if (voteType === 1) {
          if (wasLiked) {
            // Toggle off
            newLikes--;
            newVote = 0;
          } else {
            newLikes++;
            if (wasDisliked) newDislikes--;
          }
        } else {
          if (wasDisliked) {
            // Toggle off
            newDislikes--;
            newVote = 0;
          } else {
            newDislikes++;
            if (wasLiked) newLikes--;
          }
        }

        return {
          ...p,
          likes: Math.max(0, newLikes),
          dislikes: Math.max(0, newDislikes),
          user_vote: newVote,
        };
      })
    );

    try {
      await api.votePost(postId, voteType);
    } catch {
      // Revert on error
      await fetchPosts();
    } finally {
      setVotingIds((prev) => {
        const next = new Set(prev);
        next.delete(postId);
        return next;
      });
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handlePost();
    }
  };

  return (
    <Card>
      <div className="flex items-center gap-2 mb-4">
        <Newspaper size={18} className="text-grub-green" />
        <h3 className="text-white font-semibold">News</h3>
        {posts.length > 0 && (
          <span className="text-text-secondary text-xs bg-card-hover px-2 py-0.5 rounded-full">
            {posts.length}
          </span>
        )}
      </div>

      {/* Post input */}
      <div className="flex gap-2 mb-4">
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value.slice(0, 500))}
          onKeyDown={handleKeyDown}
          placeholder={`What's happening with $${ticker}?`}
          rows={2}
          className="flex-1 bg-dark-bg border border-border-dark rounded-lg px-3 py-2 text-white text-sm
            placeholder-text-secondary resize-none focus:outline-none focus:ring-1 focus:ring-grub-green/50
            focus:border-grub-green transition-all"
        />
        <button
          onClick={handlePost}
          disabled={!content.trim() || posting}
          className="self-end px-3 py-2 bg-grub-green text-black rounded-lg font-medium text-sm
            disabled:opacity-40 disabled:cursor-not-allowed hover:bg-grub-green/90 transition-colors"
        >
          <Send size={16} />
        </button>
      </div>
      <p className="text-text-secondary text-[10px] mb-3 -mt-2">
        News affects market sentiment. Likes push the stock up, dislikes push it down.
      </p>

      {/* Posts list */}
      {loading ? (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-dark-bg rounded-lg h-16 animate-pulse" />
          ))}
        </div>
      ) : posts.length === 0 ? (
        <div className="text-text-secondary text-center py-6 text-sm">
          No news yet. Be the first to post about ${ticker}!
        </div>
      ) : (
        <div className="space-y-3 max-h-[400px] overflow-y-auto pr-1">
          <AnimatePresence initial={false}>
            {posts.map((post) => {
              const netSentiment = post.likes - post.dislikes;
              return (
                <motion.div
                  key={post.id}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  className="bg-dark-bg rounded-lg p-3"
                >
                  <div className="flex items-start justify-between gap-2">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-white text-sm font-medium">
                          {post.author_username}
                        </span>
                        <span className="text-text-secondary text-[10px]">
                          {timeAgo(post.created_at)}
                        </span>
                        {netSentiment !== 0 && (
                          <span
                            className={`text-[10px] px-1.5 py-0.5 rounded ${
                              netSentiment > 0
                                ? "bg-grub-green/10 text-grub-green"
                                : "bg-grub-red/10 text-grub-red"
                            }`}
                          >
                            {netSentiment > 0 ? "Bullish" : "Bearish"}
                          </span>
                        )}
                      </div>
                      <p className="text-text-secondary text-sm leading-relaxed break-words">
                        {post.content}
                      </p>
                    </div>
                  </div>

                  {/* Vote buttons */}
                  <div className="flex items-center gap-3 mt-2">
                    <button
                      onClick={() => handleVote(post.id, 1)}
                      className={`flex items-center gap-1 text-xs transition-colors ${
                        post.user_vote === 1
                          ? "text-grub-green"
                          : "text-text-secondary hover:text-grub-green"
                      }`}
                    >
                      <ThumbsUp size={14} fill={post.user_vote === 1 ? "currentColor" : "none"} />
                      {post.likes > 0 && <span>{post.likes}</span>}
                    </button>
                    <button
                      onClick={() => handleVote(post.id, -1)}
                      className={`flex items-center gap-1 text-xs transition-colors ${
                        post.user_vote === -1
                          ? "text-grub-red"
                          : "text-text-secondary hover:text-grub-red"
                      }`}
                    >
                      <ThumbsDown size={14} fill={post.user_vote === -1 ? "currentColor" : "none"} />
                      {post.dislikes > 0 && <span>{post.dislikes}</span>}
                    </button>
                  </div>
                </motion.div>
              );
            })}
          </AnimatePresence>
        </div>
      )}
    </Card>
  );
}
