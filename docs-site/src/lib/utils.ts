// Landed via `shadcn add @nebari/utils`. Source: nebari-design.
// Do not edit by hand — re-run the shadcn command to refresh.
import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
