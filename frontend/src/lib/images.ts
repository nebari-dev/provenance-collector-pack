// Pure, framework-free helpers for the Container Images table: filtering,
// sorting, and the stat-card tiering. Kept side-effect free so they can be unit
// tested directly and reused by the table, stat cards, and stat→filter toggle.

import type { ImageRecord } from "./types";

export type Tone = "green" | "yellow" | "red";

/** The filter selection driving the images table. Empty string = "All". */
export interface ImageFilters {
  search: string;
  namespace: string;
  signature: "" | "verified" | "signed" | "unsigned";
  sbom: "" | "yes" | "no";
  provenance: "" | "yes" | "no";
  update: "" | "yes" | "no";
}

export const EMPTY_FILTERS: ImageFilters = {
  search: "",
  namespace: "",
  signature: "",
  sbom: "",
  provenance: "",
  update: "",
};

export type SortColumn =
  | ""
  | "image"
  | "namespace"
  | "workload"
  | "signature"
  | "provenance"
  | "sbom"
  | "update";

export interface SortState {
  col: SortColumn;
  asc: boolean;
}

/** True when any filter is active (used for the "N of M match" line). */
export function hasActiveFilters(f: ImageFilters): boolean {
  return Object.values(f).some((v) => v !== "");
}

/** Filter images by search text and the signature/sbom/slsa/update selects. */
export function filterImages(images: ImageRecord[], f: ImageFilters): ImageRecord[] {
  return images.filter((img) => {
    if (f.search) {
      const hay =
        `${img.image} ${img.namespace} ${img.workload.kind}/${img.workload.name}`.toLowerCase();
      if (!hay.includes(f.search.toLowerCase())) return false;
    }
    if (f.namespace && img.namespace !== f.namespace) return false;
    if (f.signature) {
      const sig = img.signature;
      if (f.signature === "verified" && !sig?.verified) return false;
      if (f.signature === "signed" && !sig?.signed) return false;
      if (f.signature === "unsigned" && sig?.signed) return false;
    }
    if (f.sbom) {
      const has = !!img.sbom?.hasSBOM;
      if (f.sbom === "yes" && !has) return false;
      if (f.sbom === "no" && has) return false;
    }
    if (f.provenance) {
      const has = !!img.provenance?.hasProvenance;
      if (f.provenance === "yes" && !has) return false;
      if (f.provenance === "no" && has) return false;
    }
    if (f.update) {
      const has = !!img.update?.updateAvailable;
      if (f.update === "yes" && !has) return false;
      if (f.update === "no" && has) return false;
    }
    return true;
  });
}

// Rank helpers so the status columns sort in a meaningful order (higher = more
// assurance). Absent signature sorts below "unsigned" (-1), matching the old UI.
function signatureRank(img: ImageRecord): number {
  const s = img.signature;
  if (!s) return -1;
  return s.verified ? 2 : s.signed ? 1 : 0;
}

/** Return a new array sorted by the given column/direction (stable no-op if unset). */
export function sortImages(images: ImageRecord[], { col, asc }: SortState): ImageRecord[] {
  if (!col) return images;
  const dir = asc ? 1 : -1;
  const sorted = [...images];
  sorted.sort((a, b) => {
    switch (col) {
      case "image":
        return a.image.localeCompare(b.image) * dir;
      case "namespace":
        return a.namespace.localeCompare(b.namespace) * dir;
      case "workload":
        return (
          `${a.workload.kind}/${a.workload.name}`.localeCompare(
            `${b.workload.kind}/${b.workload.name}`,
          ) * dir
        );
      case "signature":
        return (signatureRank(a) - signatureRank(b)) * dir;
      case "provenance":
        return (
          ((a.provenance?.hasProvenance ? 1 : 0) - (b.provenance?.hasProvenance ? 1 : 0)) * dir
        );
      case "sbom":
        return ((a.sbom?.hasSBOM ? 1 : 0) - (b.sbom?.hasSBOM ? 1 : 0)) * dir;
      case "update":
        return ((a.update?.updateAvailable ? 1 : 0) - (b.update?.updateAvailable ? 1 : 0)) * dir;
      default:
        return 0;
    }
  });
  return sorted;
}

/**
 * Tier a ratio-style stat (signed / verified / slsa / sbom of the unique-image
 * total) into a colour. A zero count is neutral ("" — no positive signal yet),
 * not red: a vanilla cluster pulling unsigned public images is expected to have
 * none. Red is reserved for when some images carry the attribute but the
 * percentage is poor.
 */
export function tier(n: number, total: number, lo: number, hi: number): Tone | "" {
  if (n === 0) return "";
  const pct = total ? Math.round((n / total) * 100) : 0;
  if (pct >= hi) return "green";
  if (pct >= lo) return "yellow";
  return "red";
}
