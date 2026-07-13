import { useAtom, useSetAtom } from "jotai";
import { useMemo } from "react";
import { StatusBadge, type Tone } from "@/components/StatusBadge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  EMPTY_FILTERS,
  filterImages,
  hasActiveFilters,
  type ImageFilters,
  type SortColumn,
  sortImages,
} from "@/lib/images";
import type { ImageRecord } from "@/lib/types";
import {
  detailImageAtom,
  filtersAtom,
  pageAtom,
  pageSizeAtom,
  sortAtom,
  statFilterAtom,
} from "@/store/uiAtoms";

const PAGE_SIZE_ITEMS = [25, 50, 100].map((n) => ({ value: String(n), label: `${n}/page` }));

// Column header definitions: label + whether the column is sortable.
const COLUMNS: { key: SortColumn; label: string }[] = [
  { key: "image", label: "Image" },
  { key: "namespace", label: "NS" },
  { key: "workload", label: "Workload" },
  { key: "signature", label: "Sig" },
  { key: "provenance", label: "SLSA" },
  { key: "sbom", label: "SBOM" },
  { key: "update", label: "Update" },
];

function signatureBadge(img: ImageRecord) {
  const s = img.signature;
  if (!s) return <StatusBadge tone="muted">N/A</StatusBadge>;
  if (s.verified) return <StatusBadge tone="green">Verified</StatusBadge>;
  if (s.signed) return <StatusBadge tone="yellow">Signed</StatusBadge>;
  return <StatusBadge tone="red">Unsigned</StatusBadge>;
}

function toneBadge(tone: Tone, text: string) {
  return <StatusBadge tone={tone}>{text}</StatusBadge>;
}

function EmptyRow({ message }: { message: string }) {
  return (
    <TableRow>
      <TableCell colSpan={COLUMNS.length} className="py-10 text-center text-muted-foreground">
        {message}
      </TableCell>
    </TableRow>
  );
}

export function ImagesTable({ images }: { images: ImageRecord[] }) {
  const [filters, setFilters] = useAtom(filtersAtom);
  const [sort, setSort] = useAtom(sortAtom);
  const [page, setPage] = useAtom(pageAtom);
  const [pageSize, setPageSize] = useAtom(pageSizeAtom);
  const setStatFilter = useSetAtom(statFilterAtom);
  const setDetailImage = useSetAtom(detailImageAtom);

  const namespaces = useMemo(() => [...new Set(images.map((i) => i.namespace))].sort(), [images]);

  const filtered = useMemo(
    () => sortImages(filterImages(images, filters), sort),
    [images, filters, sort],
  );

  const total = images.length;
  const totalPages = Math.max(1, Math.ceil(filtered.length / pageSize));
  const clampedPage = Math.min(page, totalPages - 1);
  const start = clampedPage * pageSize;
  const pageRows = filtered.slice(start, start + pageSize);

  // A filter change clears the active stat card and returns to the first page.
  function patchFilter(patch: Partial<ImageFilters>) {
    setFilters({ ...filters, ...patch });
    setStatFilter("");
    setPage(0);
  }

  function clearFilters() {
    setFilters(EMPTY_FILTERS);
    setStatFilter("");
    setPage(0);
  }

  function onSortClick(col: SortColumn) {
    setSort((s) => (s.col === col ? { col, asc: !s.asc } : { col, asc: true }));
  }

  function sortArrow(col: SortColumn) {
    if (sort.col !== col) return null;
    return <span className="text-[9px]">{sort.asc ? "▲" : "▼"}</span>;
  }

  if (total === 0) {
    return (
      <div className="py-10 text-center text-muted-foreground text-sm">No images discovered</div>
    );
  }

  return (
    <div>
      {/* Toolbar: search on the left, filters + clear grouped on the right. The
          Nebari Input renders inside a w-full wrapper, so we cap its width with a
          sized container rather than classes on the input itself. */}
      <div className="mb-3 flex flex-wrap items-center gap-3">
        <div className="w-[320px] max-w-full">
          <Input
            type="text"
            data-testid="image-search"
            value={filters.search}
            placeholder="Search image, workload..."
            onChange={(e) => patchFilter({ search: e.target.value })}
            className="h-8"
          />
        </div>
        <div className="ml-auto flex flex-wrap items-center gap-2">
          <FilterSelect
            label="NS"
            value={filters.namespace}
            onChange={(v) => patchFilter({ namespace: v })}
            options={[
              { value: "", label: "All" },
              ...namespaces.map((n) => ({ value: n, label: n })),
            ]}
          />
          <FilterSelect
            label="Sig"
            value={filters.signature}
            onChange={(v) => patchFilter({ signature: v as ImageFilters["signature"] })}
            options={[
              { value: "", label: "All" },
              { value: "verified", label: "Verified" },
              { value: "signed", label: "Signed" },
              { value: "unsigned", label: "Unsigned" },
            ]}
          />
          <FilterSelect
            label="SBOM"
            value={filters.sbom}
            onChange={(v) => patchFilter({ sbom: v as ImageFilters["sbom"] })}
            options={YES_NO}
          />
          <FilterSelect
            label="SLSA"
            value={filters.provenance}
            onChange={(v) => patchFilter({ provenance: v as ImageFilters["provenance"] })}
            options={YES_NO}
          />
          <FilterSelect
            label="Update"
            value={filters.update}
            onChange={(v) => patchFilter({ update: v as ImageFilters["update"] })}
            options={YES_NO}
          />
          <Button variant="outline" size="sm" onClick={clearFilters} className="h-8">
            Clear
          </Button>
        </div>
      </div>

      {hasActiveFilters(filters) ? (
        <div data-testid="images-match" className="mb-2 px-0.5 text-[11px] text-muted-foreground">
          {filtered.length} of {total} match
        </div>
      ) : null}

      <Table>
        <TableHeader>
          <TableRow>
            {COLUMNS.map((c) => (
              <TableHead key={c.key} onClick={() => onSortClick(c.key)}>
                <span className="flex items-center gap-1">
                  {c.label}
                  {sortArrow(c.key)}
                </span>
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {pageRows.length === 0 ? (
            <EmptyRow message="No images match filters" />
          ) : (
            pageRows.map((img) => {
              const update = img.update;
              return (
                <TableRow
                  key={`${img.namespace}/${img.workload.kind}/${img.workload.name}/${img.image}`}
                  className="cursor-pointer"
                  onClick={() => setDetailImage(img)}
                >
                  <TableCell>
                    <span className="font-mono text-xs">{img.image}</span>
                    {img.digest ? (
                      <div className="font-mono text-[10px] text-muted-foreground">
                        {img.digest.substring(7, 19)}
                      </div>
                    ) : null}
                  </TableCell>
                  <TableCell className="text-xs">{img.namespace}</TableCell>
                  <TableCell
                    className="max-w-[260px] truncate text-xs"
                    title={`${img.workload.kind}/${img.workload.name}`}
                  >
                    {img.workload.kind}/{img.workload.name}
                  </TableCell>
                  <TableCell>{signatureBadge(img)}</TableCell>
                  <TableCell>
                    {img.provenance?.hasProvenance
                      ? toneBadge("purple", "SLSA")
                      : toneBadge("muted", "None")}
                  </TableCell>
                  <TableCell>
                    {img.sbom?.hasSBOM
                      ? toneBadge("green", (img.sbom.format ?? "").toUpperCase() || "SBOM")
                      : toneBadge("muted", "None")}
                  </TableCell>
                  <TableCell>
                    {update?.updateAvailable
                      ? toneBadge(
                          "yellow",
                          update.latestInMajor || update.newestAvailable || "Update",
                        )
                      : toneBadge("green", "Current")}
                  </TableCell>
                </TableRow>
              );
            })
          )}
        </TableBody>
      </Table>

      {filtered.length > pageSize ? (
        <div className="flex items-center justify-between px-0.5 pt-3 text-[11px] text-muted-foreground">
          <div>
            {start + 1}-{Math.min(start + pageSize, filtered.length)} of {filtered.length}
          </div>
          <div className="flex items-center gap-1">
            <Button
              variant="outline"
              size="sm"
              disabled={clampedPage === 0}
              onClick={() => setPage(0)}
            >
              «
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={clampedPage === 0}
              onClick={() => setPage(clampedPage - 1)}
            >
              ‹
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={clampedPage >= totalPages - 1}
              onClick={() => setPage(clampedPage + 1)}
            >
              ›
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={clampedPage >= totalPages - 1}
              onClick={() => setPage(totalPages - 1)}
            >
              »
            </Button>
            <Select
              items={PAGE_SIZE_ITEMS}
              value={String(pageSize)}
              onValueChange={(v) => {
                setPageSize(Number((v as string | null) ?? pageSize));
                setPage(0);
              }}
            >
              <SelectTrigger className="h-7 w-auto min-w-[92px]">
                <SelectValue>
                  {(v) => PAGE_SIZE_ITEMS.find((i) => i.value === v)?.label ?? `${pageSize}/page`}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                {PAGE_SIZE_ITEMS.map((i) => (
                  <SelectItem key={i.value} value={i.value}>
                    {i.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
      ) : null}
    </div>
  );
}

const YES_NO = [
  { value: "", label: "All" },
  { value: "yes", label: "Yes" },
  { value: "no", label: "No" },
];

function FilterSelect({
  label,
  value,
  onChange,
  options,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  options: { value: string; label: string }[];
}) {
  return (
    <div className="flex items-center gap-1.5">
      <span className="font-medium text-[10px] text-muted-foreground uppercase tracking-wide">
        {label}
      </span>
      <Select
        items={options}
        value={value}
        onValueChange={(v) => onChange((v as string | null) ?? "")}
      >
        <SelectTrigger className="h-8 w-auto min-w-[92px]">
          <SelectValue>{(v) => options.find((o) => o.value === v)?.label ?? "All"}</SelectValue>
        </SelectTrigger>
        <SelectContent>
          {options.map((o) => (
            <SelectItem key={o.value} value={o.value}>
              {o.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
