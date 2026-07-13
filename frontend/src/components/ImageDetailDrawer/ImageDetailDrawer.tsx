import { useAtom } from "jotai";
import { X } from "lucide-react";
import { useEffect } from "react";
import type { ImageRecord } from "@/lib/types";
import { cn } from "@/lib/utils";
import { detailImageAtom } from "@/store/uiAtoms";

type StatusTone = "green" | "yellow" | "red" | "muted";

const DOT_TONE: Record<StatusTone, string> = {
  green: "bg-success-foreground",
  yellow: "bg-warning-foreground",
  red: "bg-destructive-foreground",
  muted: "bg-muted-foreground",
};

function Row({ label, value }: { label: string; value?: string }) {
  return (
    <div className="flex items-center justify-between py-1 text-xs">
      <span className="text-muted-foreground">{label}</span>
      <span className="max-w-[300px] break-all text-right font-mono text-[11px] text-foreground">
        {value || "-"}
      </span>
    </div>
  );
}

function Status({
  tone,
  label,
  description,
}: {
  tone: StatusTone;
  label: string;
  description: string;
}) {
  return (
    <div className="flex items-start gap-2 py-2.5">
      <div className={cn("mt-1 size-2 shrink-0 rounded-full", DOT_TONE[tone])} />
      <div className="text-xs">
        <div>{label}</div>
        <div className="mt-px text-[11px] text-muted-foreground">{description}</div>
      </div>
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="border-border border-b px-6 py-4 last:border-b-0">
      <h4 className="mb-2.5 font-semibold text-[10px] text-muted-foreground uppercase tracking-wider">
        {title}
      </h4>
      {children}
    </div>
  );
}

function DetailContent({ img }: { img: ImageRecord }) {
  return (
    <>
      <div className="border-border border-b px-6 pt-5 pb-4">
        <h3 className="break-all font-mono font-semibold text-sm leading-snug">{img.image}</h3>
        {img.digest ? (
          <div className="mt-1 break-all font-mono text-[11px] text-muted-foreground">
            {img.digest}
          </div>
        ) : null}
      </div>

      <Section title="Workload">
        <Row label="Kind" value={img.workload.kind} />
        <Row label="Name" value={img.workload.name} />
        <Row label="Namespace" value={img.namespace} />
      </Section>

      <Section title="Signature Verification">
        {img.signature ? (
          <>
            <Status
              tone={img.signature.verified ? "green" : img.signature.signed ? "yellow" : "red"}
              label={
                img.signature.verified
                  ? "Verified"
                  : img.signature.signed
                    ? "Signed (unverified)"
                    : "No signature found"
              }
              description={
                img.signature.verified
                  ? "Signature exists and has been verified against the configured public key."
                  : img.signature.signed
                    ? "A cosign signature exists but no public key was configured for verification."
                    : "No cosign signature was found attached to this image."
              }
            />
            {img.signature.error ? <Row label="Error" value={img.signature.error} /> : null}
          </>
        ) : (
          <Status
            tone="muted"
            label="Not checked"
            description="Signature verification was disabled for this collection run."
          />
        )}
      </Section>

      <Section title="SLSA Provenance">
        {img.provenance?.hasProvenance ? (
          <>
            <Status
              tone="green"
              label="Provenance attestation found"
              description="This image has a SLSA provenance attestation attached via OCI referrers."
            />
            <Row label="Predicate Type" value={img.provenance.predicateType} />
          </>
        ) : img.provenance ? (
          <Status
            tone="muted"
            label="No provenance"
            description="No SLSA provenance attestation was found in the OCI referrers index."
          />
        ) : (
          <Status
            tone="muted"
            label="Not checked"
            description="Provenance detection was disabled for this collection run."
          />
        )}
      </Section>

      <Section title="Software Bill of Materials">
        {img.sbom?.hasSBOM ? (
          <>
            <Status
              tone="green"
              label="SBOM attached"
              description="An SBOM attestation was found attached to this image."
            />
            <Row label="Format" value={(img.sbom.format ?? "").toUpperCase()} />
          </>
        ) : img.sbom ? (
          <Status
            tone="muted"
            label="No SBOM"
            description="No SBOM attestation (SPDX or CycloneDX) was found attached to this image."
          />
        ) : (
          <Status
            tone="muted"
            label="Not checked"
            description="SBOM detection was disabled for this collection run."
          />
        )}
      </Section>

      <Section title="Version Status">
        {img.update ? (
          img.update.updateAvailable ? (
            <>
              <Status
                tone="yellow"
                label="Update available"
                description="A newer version exists in the registry."
              />
              <Row label="Current Tag" value={img.update.currentTag} />
              {img.update.latestInMajor ? (
                <Row label="Latest (same major)" value={img.update.latestInMajor} />
              ) : null}
              {img.update.newestAvailable ? (
                <Row label="Newest available" value={img.update.newestAvailable} />
              ) : null}
            </>
          ) : (
            <>
              <Status
                tone="green"
                label="Up to date"
                description="This image is running the latest available version (at the configured update level)."
              />
              <Row label="Current Tag" value={img.update.currentTag} />
            </>
          )
        ) : (
          <Status
            tone="muted"
            label="Not checked"
            description="Update checking was disabled for this collection run."
          />
        )}
      </Section>
    </>
  );
}

export function ImageDetailDrawer() {
  const [img, setImg] = useAtom(detailImageAtom);
  const open = img != null;

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setImg(null);
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [open, setImg]);

  return (
    <>
      <div
        aria-hidden
        onClick={() => setImg(null)}
        className={cn(
          "fixed inset-0 z-40 bg-scrim transition-opacity duration-200",
          open ? "opacity-100" : "pointer-events-none opacity-0",
        )}
      />
      <aside
        role="dialog"
        aria-modal="true"
        aria-label="Image details"
        className={cn(
          "fixed inset-y-0 right-0 z-50 w-[520px] max-w-full overflow-y-auto border-border border-l bg-card transition-transform duration-[250ms] ease-[--ease-standard]",
          open ? "translate-x-0" : "translate-x-full",
        )}
      >
        <button
          type="button"
          aria-label="Close"
          onClick={() => setImg(null)}
          className="absolute top-3 right-4 rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground"
        >
          <X className="size-4" />
        </button>
        {img ? <DetailContent img={img} /> : null}
      </aside>
    </>
  );
}
