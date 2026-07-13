import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { MeResponse } from "@/lib/types";

// Fail-quiet default mirroring the old UI: scan hidden, no email, features off.
const DEFAULT_ME: MeResponse = {
  authEnabled: false,
  canRunScan: false,
  features: { timelineDeltas: false },
};

/**
 * GET /api/me — resolves scan-button gating, the profile display, and the
 * opt-in feature flags. Falls back to the all-off default if the call errors so
 * the dashboard still renders (matches the old init()'s fail-quiet behaviour).
 */
export function useMe(): MeResponse {
  const { data } = useQuery({
    queryKey: ["me"],
    queryFn: () => api.get<MeResponse>("/api/me"),
  });
  return data ?? DEFAULT_ME;
}
