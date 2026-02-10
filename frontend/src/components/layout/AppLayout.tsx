"use client";

import Sidebar from "./Sidebar";
import AuthGuard from "./AuthGuard";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthGuard>
      <div className="min-h-screen bg-dark-bg">
        <Sidebar />
        <main className="md:ml-64 pb-20 md:pb-0 min-h-screen">
          <div className="max-w-6xl mx-auto p-4 md:p-8">{children}</div>
        </main>
      </div>
    </AuthGuard>
  );
}
