import { ChevronDown } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useRunScan } from "@/hooks/useRunScan";

function exportLinks(filename: string | null) {
  const suffix = filename ? `&filename=${encodeURIComponent(filename)}` : "";
  const jsonPath = filename
    ? `/api/reports/${encodeURIComponent(filename)}`
    : "/api/reports/latest";
  return {
    csv: `/api/export?format=csv${suffix}`,
    markdown: `/api/export?format=markdown${suffix}`,
    json: jsonPath,
  };
}

export function PageHeader({
  clusterName,
  generatedAt,
  activeFilename,
  canRunScan,
}: {
  clusterName?: string;
  generatedAt?: string;
  activeFilename: string | null;
  canRunScan: boolean;
}) {
  const { runScan, isScanning } = useRunScan();
  const links = exportLinks(activeFilename);
  const lastUpdated = generatedAt ? new Date(generatedAt).toLocaleString() : "";

  return (
    <div className="mb-5 flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 className="font-semibold text-[22px] tracking-tight">Provenance</h1>
        <p className="mt-0.5 max-w-[760px] text-muted-foreground text-sm">
          Container image provenance, signatures, SBOMs, and available updates discovered across
          your cluster.
        </p>
      </div>

      <div className="flex flex-wrap items-center gap-3">
        <div className="flex flex-col items-end text-[11px] text-muted-foreground leading-tight">
          <span>{clusterName || "Local"}</span>
          {lastUpdated ? <span>{lastUpdated}</span> : null}
        </div>

        <DropdownMenu modal={false}>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm" title="Download the selected report">
              Export
              <ChevronDown className="size-3.5" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-auto min-w-[140px]">
            <DropdownMenuItem asChild>
              <a href={links.csv} download>
                CSV
              </a>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <a href={links.markdown} download>
                Markdown
              </a>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <a href={links.json} download="provenance-report.json">
                JSON
              </a>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        {canRunScan ? (
          <Button
            data-testid="run-scan"
            variant="outline"
            size="sm"
            onClick={runScan}
            loading={isScanning}
            loadingText="Scan running"
            title="Trigger a manual provenance scan"
          >
            Run Scan
          </Button>
        ) : null}
      </div>
    </div>
  );
}
