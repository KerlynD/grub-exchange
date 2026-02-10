"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { motion } from "framer-motion";
import Input from "@/components/ui/Input";
import Button from "@/components/ui/Button";
import AuthBackground from "@/components/ui/AuthBackground";
import { useAuth } from "@/contexts/AuthContext";

export default function RegisterPage() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [firstName, setFirstName] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { register, user } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (user) router.push("/dashboard");
  }, [user, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await register(username, email, password, firstName);
      router.push("/dashboard");
    } catch (e) {
      setError(e instanceof Error ? e.message : "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-dark-bg flex items-center justify-center p-4 relative">
      <AuthBackground />
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="w-full max-w-md relative z-10"
      >
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-white mb-2">
            <span className="text-grub-green">Grub</span> Exchange
          </h1>
          <p className="text-text-secondary">
            Create your account & become a tradeable stock
          </p>
        </div>

        <div className="bg-card-bg rounded-2xl border border-border-dark p-8">
          <h2 className="text-xl font-bold text-white mb-6">Sign Up</h2>

          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              label="First Name"
              type="text"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
              placeholder="Your first name"
              required
            />

            {firstName.length >= 2 && (
              <motion.div
                initial={{ opacity: 0, height: 0 }}
                animate={{ opacity: 1, height: "auto" }}
                className="bg-grub-green/10 border border-grub-green/30 rounded-lg px-4 py-2"
              >
                <p className="text-grub-green text-sm">
                  Your ticker will be:{" "}
                  <span className="font-bold">
                    ${firstName.toUpperCase()}
                  </span>
                </p>
              </motion.div>
            )}

            <Input
              label="Username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Choose a username"
              required
            />
            <Input
              label="Email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              required
            />
            <Input
              label="Password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="At least 6 characters"
              required
              minLength={6}
            />

            {error && (
              <p className="text-grub-red text-sm">{error}</p>
            )}

            <Button
              type="submit"
              variant="success"
              size="lg"
              className="w-full"
              loading={loading}
            >
              Create Account
            </Button>
          </form>

          <p className="text-text-secondary text-sm text-center mt-6">
            Already have an account?{" "}
            <Link href="/login" className="text-grub-green hover:underline">
              Log In
            </Link>
          </p>
        </div>

        <p className="text-text-secondary text-xs text-center mt-4">
          You&apos;ll start with 100 Grub and become a tradeable stock!
        </p>
      </motion.div>
    </div>
  );
}
