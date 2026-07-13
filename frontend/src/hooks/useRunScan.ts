import { useQueryClient } from "@tanstack/react-query";
import { useSetAtom } from "jotai";
import { useCallback, useEffect, useRef, useState } from "react";
import { ApiError, api } from "@/lib/api";
import type { ReportEntry, ScanResponse } from "@/lib/types";
import { activeReportFilenameAtom, addToastAtom } from "@/store/uiAtoms";

const SCAN_POLL_MS = 5000;
const SCAN_POLL_MAX = 60; // ~5 minutes worst case

function sleep(ms: number, signal: AbortSignal): Promise<void> {
  return new Promise((resolve, reject) => {
    const id = setTimeout(resolve, ms);
    signal.addEventListener(
      "abort",
      () => {
        clearTimeout(id);
        reject(new DOMException("aborted", "AbortError"));
      },
      { once: true },
    );
  });
}

/**
 * Triggers a manual scan (POST /api/scan) then polls GET /api/reports until a
 * new report appears — at which point it switches the view to it — or until the
 * ~5-minute budget elapses. Surfaces progress via toasts. Mirrors the old
 * runScan()/pollForNewReport() flow.
 */
export function useRunScan() {
  const queryClient = useQueryClient();
  const addToast = useSetAtom(addToastAtom);
  const setActiveReport = useSetAtom(activeReportFilenameAtom);
  const [isScanning, setIsScanning] = useState(false);
  const abortRef = useRef<AbortController | null>(null);

  // Cancel any in-flight poll loop if the component unmounts.
  useEffect(() => () => abortRef.current?.abort(), []);

  const runScan = useCallback(async () => {
    if (isScanning) return;
    setIsScanning(true);
    const controller = new AbortController();
    abortRef.current = controller;

    const baseline = (queryClient.getQueryData<ReportEntry[]>(["reports"]) ?? []).length;

    try {
      const res = await api.post<ScanResponse>("/api/scan");
      addToast({
        kind: "success",
        title: "Scan started",
        sub: res.jobName ? `Job: ${res.jobName}` : undefined,
      });
    } catch (err) {
      const detail = err instanceof ApiError ? `${err.status} ${err.message}` : String(err);
      addToast({ kind: "error", title: "Scan request failed", sub: detail });
      setIsScanning(false);
      return;
    }

    try {
      for (let tick = 0; tick < SCAN_POLL_MAX; tick += 1) {
        await sleep(SCAN_POLL_MS, controller.signal);
        try {
          const latest = await api.get<ReportEntry[]>("/api/reports");
          if (latest.length > baseline) {
            queryClient.setQueryData(["reports"], latest);
            setActiveReport(latest[0].filename);
            addToast({
              kind: "success",
              title: "New report available",
              sub: latest[0].filename,
            });
            return;
          }
        } catch {
          // Transient error — keep polling.
        }
      }
      addToast({
        kind: "error",
        title: "Scan still running",
        sub: "Stopped watching after 5 min — refresh later",
      });
    } catch {
      // Aborted (unmount) — nothing to report.
    } finally {
      setIsScanning(false);
    }
  }, [isScanning, queryClient, addToast, setActiveReport]);

  return { runScan, isScanning };
}
