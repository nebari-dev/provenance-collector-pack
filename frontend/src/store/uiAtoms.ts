// Jotai UI state shared across the dashboard. Server data lives in TanStack
// Query; this holds the view state the old vanilla UI kept in module globals:
// the selected report, the images-table filters/sort/pagination, the open row
// detail, and the toast queue.

import { atom } from "jotai";
import { EMPTY_FILTERS, type ImageFilters, type SortState } from "@/lib/images";
import type { ImageRecord } from "@/lib/types";

/**
 * Filename of the report the user is viewing. `null` means "follow the newest"
 * — App resolves it to reports[0] so the dashboard opens on the latest scan.
 */
export const activeReportFilenameAtom = atom<string | null>(null);

/** Which stat card is active (drives the highlight + the derived table filter). */
export const statFilterAtom = atom<string>("");

export const filtersAtom = atom<ImageFilters>(EMPTY_FILTERS);
export const sortAtom = atom<SortState>({ col: "", asc: true });
export const pageAtom = atom<number>(0);
export const pageSizeAtom = atom<number>(25);

/** The image whose detail drawer is open, or null when closed. */
export const detailImageAtom = atom<ImageRecord | null>(null);

export interface Toast {
  id: number;
  kind: "success" | "error";
  title: string;
  sub?: string;
}

export const toastsAtom = atom<Toast[]>([]);

let nextToastId = 0;

/** Write-only: append a toast. Returns nothing; read toastsAtom to render. */
export const addToastAtom = atom(null, (get, set, toast: Omit<Toast, "id">) => {
  nextToastId += 1;
  set(toastsAtom, [...get(toastsAtom), { ...toast, id: nextToastId }]);
});

/** Write-only: remove a toast by id (called when it auto-dismisses). */
export const removeToastAtom = atom(null, (get, set, id: number) => {
  set(
    toastsAtom,
    get(toastsAtom).filter((t) => t.id !== id),
  );
});
