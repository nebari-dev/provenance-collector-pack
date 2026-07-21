import { useAtom, useSetAtom } from "jotai";
import { useEffect } from "react";
import { HelmTable } from "@/components/HelmTable";
import { ImageDetailDrawer } from "@/components/ImageDetailDrawer";
import { ImagesTable } from "@/components/ImagesTable";
import { PageHeader } from "@/components/PageHeader";
import { StatCards } from "@/components/StatCards";
import { Timeline } from "@/components/Timeline";
import { Toasts } from "@/components/Toasts";
import { Topbar } from "@/components/Topbar";
import { Spinner } from "@/components/ui/spinner";
import { useMe } from "@/hooks/useMe";
import { useReport, useReports } from "@/hooks/useReports";
import { EMPTY_FILTERS } from "@/lib/images";
import { cn } from "@/lib/utils";
import { activeReportFilenameAtom, filtersAtom, pageAtom, statFilterAtom } from "@/store/uiAtoms";

function Section({
  title,
  count,
  countTestId,
  className,
  children,
}: {
  title: string;
  count?: string;
  countTestId?: string;
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <section className={cn("mb-6", className)}>
      <header className="mb-3 flex items-baseline gap-2">
        <h2 className="font-semibold text-foreground text-sm tracking-tight">{title}</h2>
        {count ? (
          <span data-testid={countTestId} className="text-[13px] text-muted-foreground">
            ({count})
          </span>
        ) : null}
      </header>
      {children}
    </section>
  );
}

function CenteredMessage({ children }: { children: React.ReactNode }) {
  return <div className="py-16 text-center text-muted-foreground text-sm">{children}</div>;
}

export default function App() {
  const me = useMe();
  // The listing endpoint returns JSON `null` (not `[]`) when no scans exist yet,
  // and a destructuring default only fills in for `undefined` — so coalesce
  // explicitly, otherwise reports[0]/reports.length below throw on a fresh install.
  const { data, isLoading: reportsLoading, isError } = useReports();
  const reports = data ?? [];

  const [active] = useAtom(activeReportFilenameAtom);
  const activeFilename = active ?? reports[0]?.filename ?? null;
  const { data: report } = useReport(activeFilename);

  // Reset the images-table view whenever the active report changes, mirroring
  // the old loadReport() which cleared filters/stat/page on every switch.
  const setFilters = useSetAtom(filtersAtom);
  const setStatFilter = useSetAtom(statFilterAtom);
  const setPage = useSetAtom(pageAtom);
  // biome-ignore lint/correctness/useExhaustiveDependencies: activeFilename is the intended trigger — resetting the view each time the selected report changes.
  useEffect(() => {
    setFilters(EMPTY_FILTERS);
    setStatFilter("");
    setPage(0);
  }, [activeFilename, setFilters, setStatFilter, setPage]);

  return (
    <div className="min-h-full">
      <Topbar />
      <main className="w-full px-10 pt-6 pb-10">
        <PageHeader
          clusterName={report?.metadata.clusterName}
          generatedAt={report?.metadata.generatedAt}
          activeFilename={activeFilename}
          canRunScan={me.canRunScan}
        />

        {reportsLoading ? (
          <CenteredMessage>
            <Spinner className="mx-auto" />
          </CenteredMessage>
        ) : isError ? (
          <CenteredMessage>Failed to load reports</CenteredMessage>
        ) : reports.length === 0 ? (
          <Section title="Timeline">
            <Timeline reports={[]} activeFilename={null} />
          </Section>
        ) : (
          <>
            {report ? <StatCards summary={report.summary} /> : null}

            <Section title="Timeline" className="mt-10">
              <Timeline reports={reports} activeFilename={activeFilename} />
            </Section>

            <Section
              title="Container Images"
              count={report ? String(report.images.length) : undefined}
              countTestId="images-total"
            >
              {report ? (
                <ImagesTable images={report.images} />
              ) : (
                <CenteredMessage>
                  <Spinner className="mx-auto" />
                </CenteredMessage>
              )}
            </Section>

            <Section
              title="Helm Releases"
              count={report?.helmReleases?.length ? String(report.helmReleases.length) : undefined}
            >
              {report ? <HelmTable releases={report.helmReleases ?? []} /> : null}
            </Section>
          </>
        )}
      </main>

      <ImageDetailDrawer />
      <Toasts />
    </div>
  );
}
