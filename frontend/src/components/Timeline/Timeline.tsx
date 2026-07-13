import { useSetAtom } from "jotai";
import { cardVariants } from "@/components/ui/card";
import type { ReportEntry } from "@/lib/types";
import { cn } from "@/lib/utils";
import { activeReportFilenameAtom, pageAtom, statFilterAtom } from "@/store/uiAtoms";

export function Timeline({
  reports,
  activeFilename,
}: {
  reports: ReportEntry[];
  activeFilename: string | null;
}) {
  const setActive = useSetAtom(activeReportFilenameAtom);
  const setStatFilter = useSetAtom(statFilterAtom);
  const setPage = useSetAtom(pageAtom);

  if (reports.length === 0) {
    return (
      <div className="p-10 text-center text-muted-foreground text-sm">
        No reports yet. Run the collector to generate a provenance report.
      </div>
    );
  }

  function select(filename: string) {
    setActive(filename);
    setStatFilter("");
    setPage(0);
  }

  return (
    <div className="flex gap-2 overflow-x-auto pb-1">
      {reports.map((r) => {
        const d = new Date(r.generatedAt);
        const active = r.filename === activeFilename;
        return (
          <button
            key={r.filename}
            type="button"
            data-testid="timeline-item"
            aria-pressed={active}
            onClick={() => select(r.filename)}
            // Reuse the Nebari Card surface (bg-card, border, radius) but compact
            // the spacing and make it an interactive, selectable chip.
            className={cn(
              cardVariants(),
              "min-w-[120px] shrink-0 cursor-pointer gap-0 p-2.5 text-left transition-colors hover:border-primary",
              active ? "border-primary bg-primary/8" : "border-border",
            )}
          >
            <div className="font-semibold text-xs">{d.toLocaleDateString()}</div>
            <div className="text-[11px] text-muted-foreground">{d.toLocaleTimeString()}</div>
            <div className="mt-0.5 text-[10px] text-muted-foreground">
              {r.summary.totalImages} images
            </div>
          </button>
        );
      })}
    </div>
  );
}
