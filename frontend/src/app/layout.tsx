import type { Metadata } from "next";
import { AuthProvider } from "@/contexts/AuthContext";
import "./globals.css";

export const metadata: Metadata = {
  title: "Grub Exchange",
  description: "Buy and sell shares of your friends",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="bg-dark-bg min-h-screen">
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
