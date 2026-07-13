import { StatusBadge } from "@/components/StatusBadge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { HelmRecord } from "@/lib/types";

export function HelmTable({ releases }: { releases: HelmRecord[] }) {
  if (releases.length === 0) {
    return <div className="py-10 text-center text-muted-foreground text-sm">No Helm releases</div>;
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="px-4">Release</TableHead>
          <TableHead className="px-4">NS</TableHead>
          <TableHead className="px-4">Chart</TableHead>
          <TableHead className="px-4">Version</TableHead>
          <TableHead className="px-4">App</TableHead>
          <TableHead className="px-4">Status</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {releases.map((hr) => (
          <TableRow key={`${hr.namespace}/${hr.releaseName}`}>
            <TableCell className="text-xs">{hr.releaseName}</TableCell>
            <TableCell className="text-xs">{hr.namespace}</TableCell>
            <TableCell className="font-mono text-xs">{hr.chart}</TableCell>
            <TableCell className="font-mono text-xs">{hr.version}</TableCell>
            <TableCell className="font-mono text-xs">{hr.appVersion}</TableCell>
            <TableCell>
              {hr.status === "deployed" ? (
                <StatusBadge tone="green">Deployed</StatusBadge>
              ) : (
                <StatusBadge tone="yellow">{hr.status}</StatusBadge>
              )}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
