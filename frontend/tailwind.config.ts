import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        "grub-green": "#00C805",
        "grub-red": "#FF5000",
        "dark-bg": "#000000",
        "card-bg": "#1C1C1E",
        "card-hover": "#2C2C2E",
        "text-primary": "#FFFFFF",
        "text-secondary": "#8E8E93",
        "border-dark": "#38383A",
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
      },
    },
  },
  plugins: [],
};

export default config;
