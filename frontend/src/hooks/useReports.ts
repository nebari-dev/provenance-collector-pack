import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ProvenanceReport, ReportEntry } from "@/lib/types";

/** GET /api/reports — the scan listing, newest-first (server-sorted). */
export function useReports() {
  return useQuery({
    queryKey: ["reports"],
    queryFn: () => api.get<ReportEntry[]>("/api/reports"),
  });
}

/**
 * GET /api/reports/{filename} — the full report for the active scan. Disabled
 * until a filename is known (the first render before the listing resolves).
 */
export function useReport(filename: string | null) {
  return useQuery({
    queryKey: ["report", filename],
    queryFn: () => api.get<ProvenanceReport>(`/api/reports/${encodeURIComponent(filename ?? "")}`),
    enabled: filename != null,
  });
}
