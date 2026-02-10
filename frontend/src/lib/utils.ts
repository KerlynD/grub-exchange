export function formatGrub(amount: number): string {
  return amount.toFixed(2);
}

export function formatPercent(percent: number): string {
  const sign = percent >= 0 ? "+" : "";
  return `${sign}${percent.toFixed(2)}%`;
}

export function formatNumber(num: number): string {
  if (num >= 1_000_000) return `${(num / 1_000_000).toFixed(1)}M`;
  if (num >= 1_000) return `${(num / 1_000).toFixed(1)}K`;
  return num.toFixed(2);
}

export function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function cn(...classes: (string | boolean | undefined | null)[]): string {
  return classes.filter(Boolean).join(" ");
}

export function getPLColor(value: number): string {
  if (value > 0) return "text-grub-green";
  if (value < 0) return "text-grub-red";
  return "text-text-secondary";
}

export function getPLBgColor(value: number): string {
  if (value > 0) return "bg-grub-green/10 text-grub-green";
  if (value < 0) return "bg-grub-red/10 text-grub-red";
  return "bg-card-bg text-text-secondary";
}
