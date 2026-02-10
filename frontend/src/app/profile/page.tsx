"use client";

import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import AppLayout from "@/components/layout/AppLayout";
import Button from "@/components/ui/Button";
import Card from "@/components/ui/Card";
import { useAuth } from "@/contexts/AuthContext";
import * as api from "@/lib/api";
import { formatGrub } from "@/lib/utils";

export default function ProfilePage() {
  const { user, refreshUser } = useAuth();
  const [bio, setBio] = useState("");
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    if (user) {
      setBio(user.bio || "");
    }
  }, [user]);

  const handleSave = async () => {
    setSaving(true);
    try {
      await api.updateProfile({ bio });
      await refreshUser();
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch {
      // handle error
    } finally {
      setSaving(false);
    }
  };

  if (!user) return null;

  return (
    <AppLayout>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="space-y-8 max-w-2xl mx-auto"
      >
        <div>
          <h1 className="text-2xl font-bold text-white mb-1">Your Profile</h1>
          <p className="text-text-secondary text-sm">
            This info appears on your stock&apos;s ticker page
          </p>
        </div>

        {/* Profile Card */}
        <Card>
          <div className="flex items-center gap-4 mb-6">
            <div className="w-16 h-16 rounded-full bg-grub-green/20 flex items-center justify-center">
              <span className="text-grub-green font-bold text-2xl">
                {user.ticker.charAt(0)}
              </span>
            </div>
            <div>
              <p className="text-white font-bold text-xl">{user.username}</p>
              <p className="text-grub-green font-semibold">${user.ticker}</p>
              <p className="text-text-secondary text-xs">
                Share price: {formatGrub(user.current_share_price)} Grub
              </p>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4 mb-6">
            <div className="bg-dark-bg rounded-lg p-3">
              <p className="text-text-secondary text-xs">Grub Balance</p>
              <p className="text-white font-bold">
                {formatGrub(user.grub_balance)}
              </p>
            </div>
            <div className="bg-dark-bg rounded-lg p-3">
              <p className="text-text-secondary text-xs">Member Since</p>
              <p className="text-white font-bold">
                {new Date(user.created_at).toLocaleDateString("en-US", {
                  month: "short",
                  year: "numeric",
                })}
              </p>
            </div>
          </div>
        </Card>

        {/* Bio Editor */}
        <Card>
          <h2 className="text-white font-semibold mb-3">About You</h2>
          <p className="text-text-secondary text-xs mb-3">
            Write something about yourself. This will show up in the
            &quot;About&quot; section on your stock page.
          </p>
          <textarea
            value={bio}
            onChange={(e) => setBio(e.target.value)}
            maxLength={500}
            rows={4}
            placeholder="Tell people about yourself... What makes your stock worth buying?"
            className="w-full bg-dark-bg border border-border-dark rounded-lg px-4 py-3 text-white text-sm
              placeholder-text-secondary focus:outline-none focus:ring-2 focus:ring-grub-green/50 focus:border-grub-green
              transition-all duration-150 resize-none"
          />
          <div className="flex items-center justify-between mt-3">
            <p className="text-text-secondary text-xs">{bio.length}/500</p>
            <div className="flex items-center gap-3">
              {saved && (
                <motion.span
                  initial={{ opacity: 0, x: 10 }}
                  animate={{ opacity: 1, x: 0 }}
                  className="text-grub-green text-sm"
                >
                  Saved!
                </motion.span>
              )}
              <Button
                variant="success"
                size="sm"
                onClick={handleSave}
                loading={saving}
              >
                Save Profile
              </Button>
            </div>
          </div>
        </Card>
      </motion.div>
    </AppLayout>
  );
}
