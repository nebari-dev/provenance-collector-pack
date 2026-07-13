import { useAtom, useSetAtom } from "jotai";
import { EMPTY_FILTERS, type ImageFilters, type Tone, tier } from "@/lib/images";
import type { ReportSummary } from "@/lib/types";
import { cn } from "@/lib/utils";
import { filtersAtom, pageAtom, statFilterAtom } from "@/store/uiAtoms";

interface StatDef {
  id: string;
  value: string;
  label: string;
  tone: Tone | "";
  title?: string;
  /** Filter fields applied when this card is activated. */
  filter?: Partial<ImageFilters>;
}

const VALUE_TONE: Record<Tone, string> = {
  green: "text-success-foreground",
  yellow: "text-warning-foreground",
  red: "text-destructive-foreground",
};

function buildStats(s: ReportSummary): StatDef[] {
  const total = s.uniqueImages || 0;
  return [
    { id: "all", value: String(s.uniqueImages), label: "Images", tone: "" },
    {
      id: "signed",
      value: String(s.signedImages),
      label: "Signed",
      tone: tier(s.signedImages, total, 50, 80),
      title:
        s.signedImages === 0
          ? "No signatures found. Common on clusters pulling unsigned public images."
          : undefined,
      filter: { signature: "signed" },
    },
    {
      id: "verified",
      value: String(s.verifiedImages),
      label: "Verified",
      tone: tier(s.verifiedImages, total, 50, 80),
      filter: { signature: "verified" },
    },
    {
      id: "provenance",
      value: String(s.imagesWithProvenance),
      label: "SLSA",
      tone: tier(s.imagesWithProvenance, total, 1, 50),
      filter: { provenance: "yes" },
    },
    {
      id: "sbom",
      value: String(s.imagesWithSBOM),
      label: "SBOM",
      tone: tier(s.imagesWithSBOM, total, 1, 50),
      filter: { sbom: "yes" },
    },
    {
      id: "updates",
      value: String(s.imagesWithUpdates),
      label: "Updates",
      tone: s.imagesWithUpdates > 0 ? "yellow" : "green",
      filter: { update: "yes" },
    },
    { id: "helm", value: String(s.totalHelmReleases), label: "Helm", tone: "" },
  ];
}

export function StatCards({ summary }: { summary: ReportSummary }) {
  const [statFilter, setStatFilter] = useAtom(statFilterAtom);
  const [filters, setFilters] = useAtom(filtersAtom);
  const setPage = useSetAtom(pageAtom);

  const stats = buildStats(summary);

  function toggle(stat: StatDef) {
    setPage(0);
    if (statFilter === stat.id) {
      // Toggling off — clear the stat and its derived selects, keep search/ns.
      setStatFilter("");
      setFilters({ ...EMPTY_FILTERS, search: filters.search, namespace: filters.namespace });
      return;
    }
    setStatFilter(stat.id);
    setFilters({
      ...EMPTY_FILTERS,
      search: filters.search,
      namespace: filters.namespace,
      ...(stat.filter ?? {}),
    });
  }

  return (
    <div className="mt-8 mb-5 grid grid-cols-[repeat(auto-fit,minmax(120px,1fr))] gap-2.5">
      {stats.map((stat) => {
        const active = statFilter === stat.id;
        return (
          <button
            key={stat.id}
            type="button"
            data-testid={`stat-${stat.id}`}
            title={stat.title}
            aria-pressed={active}
            onClick={() => toggle(stat)}
            className={cn(
              "rounded-md border bg-card px-4 py-3.5 text-left transition-colors hover:border-primary",
              active ? "border-primary bg-primary/8" : "border-border",
            )}
          >
            <div
              data-testid="stat-value"
              className={cn(
                "font-bold text-2xl tracking-tight",
                stat.tone ? VALUE_TONE[stat.tone] : "text-foreground",
              )}
            >
              {stat.value}
            </div>
            <div className="mt-0.5 font-medium text-[10px] text-muted-foreground uppercase tracking-wider">
              {stat.label}
            </div>
          </button>
        );
      })}
    </div>
  );
}
