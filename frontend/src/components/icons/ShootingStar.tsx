"use client";

interface ShootingStarProps {
  size?: number;
  className?: string;
}

export default function ShootingStar({ size = 24, className = "" }: ShootingStarProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
    >
      {/* Star */}
      <path
        d="M12 2L14.09 8.26L20.18 8.63L15.55 12.74L17.09 18.77L12 15.27L6.91 18.77L8.45 12.74L3.82 8.63L9.91 8.26L12 2Z"
        fill="currentColor"
        opacity="0.9"
      />
      {/* Shooting trail */}
      <path
        d="M4 20L8 16"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        opacity="0.7"
      />
      <path
        d="M2 17L5 14.5"
        stroke="currentColor"
        strokeWidth="1"
        strokeLinecap="round"
        opacity="0.5"
      />
      <path
        d="M5 22L7.5 19"
        stroke="currentColor"
        strokeWidth="1"
        strokeLinecap="round"
        opacity="0.4"
      />
    </svg>
  );
}
