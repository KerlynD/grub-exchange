# Grub Exchange

A Robinhood-style social trading platform where users can buy and sell shares of their friends using **Grub** currency. When you sign up, you become a tradeable stock — and so does everyone else.

## How It Works

- **Sign up and become a stock.** Your first name becomes your ticker symbol, and your share price starts at 10 GRUB.
- **Trade your friends.** Buy and sell shares of other users. Prices move based on supply and demand — more buyers push the price up, more sellers push it down.
- **Earn Grub.** Collect a daily login bonus, earn 2% appreciation when someone buys your stock, and receive daily dividends on your portfolio value.
- **Post news and influence the market.** Write bullish or bearish takes on any stock. Community sentiment (likes/dislikes) influences the market maker's trading behavior.

## Features

- **Market-forces pricing** — AMM-style execution with slippage protection to prevent arbitrage
- **Live price charts** — Real-time candlestick-style area charts with 1D, 1W, 1M, and ALL time ranges
- **Portfolio tracking** — P&L per holding, total portfolio value over time, and a historical portfolio graph
- **Leaderboard** — Rankings for most valuable stocks, biggest gainers/losers, richest traders, and best portfolio performance
- **Achievements** — Unlock badges like First Trade, Diamond Hands, Day Trader, and Whale
- **Activity feed** — Real-time notifications when someone trades your stock
- **News & sentiment** — Post and vote on stock news; sentiment drives AI market maker behavior
- **Market maker** — Background bot that trades every 60 seconds with a bullish bias, keeping the market alive
- **Daily claim** — 10 free GRUB every 24 hours
- **Daily dividends** — 1% of your portfolio value paid out daily

## Screenshots

<!-- Add screenshots here -->

## Tech Stack

- **Frontend:** Next.js, TypeScript, Tailwind CSS, Framer Motion, Recharts
- **Backend:** Go (Gin framework)
- **Database:** PostgreSQL (Supabase)
- **Auth:** JWT with httpOnly cookies
