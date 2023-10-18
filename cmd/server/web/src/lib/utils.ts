import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function timeLeftUntil(date: Date): string {
  const now = new Date();
  let diffInSeconds = Math.floor((date.getTime() - now.getTime()) / 1000);

  if (diffInSeconds < 0) {
    return "The date has already passed!";
  }

  const months = Math.floor(diffInSeconds / 2592000);
  diffInSeconds -= months * 2592000;

  const days = Math.floor(diffInSeconds / 86400);
  diffInSeconds -= days * 86400;

  const hours = Math.floor(diffInSeconds / 3600);
  diffInSeconds -= hours * 3600;

  const minutes = Math.floor(diffInSeconds / 60);
  diffInSeconds -= minutes * 60;

  const seconds = diffInSeconds;

  if (months > 0) return months === 1 ? `${months} month` : `${months} months`;
  if (days > 0) return days === 1 ? `${days} day` : `${days} days`;
  if (hours > 0) return hours === 1 ? `${hours} hour` : `${hours} hours`;
  if (minutes > 0) return minutes === 1 ? `${minutes} minute` : `${minutes} minutes`;

  return seconds === 1 ? `${seconds} second` : `${seconds} seconds`;
}

export function parseDate(dateString: string): Date | null {
  const date = new Date(dateString);

  // Check if the date is valid
  if (isNaN(date.getTime())) {
    console.error('Invalid date string');
    return null;
  }

  return date;
}

export function addTwoHours(): Date {
  const now = new Date();
  now.setHours(now.getHours() + 2);
  return now;
}
