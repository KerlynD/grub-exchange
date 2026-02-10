"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import {
  LayoutDashboard,
  TrendingUp,
  Briefcase,
  Trophy,
  User,
  LogOut,
} from "lucide-react";
import ShootingStar from "@/components/icons/ShootingStar";

const navItems = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/market", label: "Market", icon: TrendingUp },
  { href: "/portfolio", label: "Portfolio", icon: Briefcase },
  { href: "/leaderboard", label: "Leaderboard", icon: Trophy },
  { href: "/profile", label: "Profile", icon: User },
];

export default function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuth();

  return (
    <>
      {/* Desktop Sidebar */}
      <aside className="hidden md:flex flex-col w-64 bg-card-bg border-r border-border-dark h-screen fixed left-0 top-0">
        <div className="p-6 border-b border-border-dark">
          <h1 className="text-xl font-bold text-white flex items-center gap-2">
            <ShootingStar size={22} className="text-grub-green" />
            <span><span className="text-grub-green">Grub</span> Exchange</span>
          </h1>
          {user && (
            <p className="text-text-secondary text-xs mt-1">
              ${user.ticker} &middot; {user.username}
            </p>
          )}
        </div>

        <nav className="flex-1 p-4 space-y-1">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-colors text-sm font-medium ${
                  isActive
                    ? "bg-grub-green/10 text-grub-green"
                    : "text-text-secondary hover:text-white hover:bg-card-hover"
                }`}
              >
                <Icon size={20} />
                {item.label}
              </Link>
            );
          })}
        </nav>

        <div className="p-4 border-t border-border-dark">
          <button
            onClick={logout}
            className="flex items-center gap-3 px-4 py-3 rounded-lg transition-colors text-sm font-medium text-text-secondary hover:text-grub-red hover:bg-card-hover w-full"
          >
            <LogOut size={20} />
            Log Out
          </button>
        </div>
      </aside>

      {/* Mobile Bottom Nav */}
      <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-card-bg border-t border-border-dark z-30">
        <div className="flex justify-around py-2">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`flex flex-col items-center gap-1 px-3 py-1.5 rounded-lg text-xs ${
                  isActive ? "text-grub-green" : "text-text-secondary"
                }`}
              >
                <Icon size={20} />
                {item.label}
              </Link>
            );
          })}
        </div>
      </nav>
    </>
  );
}
