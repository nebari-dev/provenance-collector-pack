// TypeScript mirror of the dashboard's JSON contract. The report shapes come
// from internal/report/types.go; the /api/me and /api/scan shapes come from the
// handler structs in internal/dashboard/{server,scan}.go. Fields the Go side
// marks `omitempty` on a pointer are optional here — their absence is the
// "Not checked" state the UI renders distinctly from an explicit negative.

export interface ReportSummary {
  totalImages: number;
  uniqueImages: number;
  signedImages: number;
  verifiedImages: number;
  imagesWithSBOM: number;
  imagesWithProvenance: number;
  imagesWithUpdates: number;
  totalHelmReleases: number;
  helmReleasesWithUpdates: number;
}

export interface ReportMetadata {
  generatedAt: string;
  collectorVersion: string;
  clusterName?: string;
  namespacesScanned: string[];
}

export interface WorkloadRef {
  kind: string;
  name: string;
}

export interface SignatureInfo {
  signed: boolean;
  verified: boolean;
  error?: string;
}

export interface SBOMInfo {
  hasSBOM: boolean;
  format?: string;
}

export interface ProvenanceInfo {
  hasProvenance: boolean;
  predicateType?: string;
}

export interface UpdateInfo {
  currentTag: string;
  latestInMajor?: string;
  newestAvailable?: string;
  updateAvailable: boolean;
}

export interface ImageRecord {
  image: string;
  digest?: string;
  namespace: string;
  workload: WorkloadRef;
  // Absent (undefined) means the corresponding check was disabled for the run.
  signature?: SignatureInfo;
  sbom?: SBOMInfo;
  provenance?: ProvenanceInfo;
  update?: UpdateInfo;
}

export interface HelmRecord {
  releaseName: string;
  namespace: string;
  chart: string;
  version: string;
  appVersion: string;
  status: string;
  update?: UpdateInfo;
}

export interface ProvenanceReport {
  metadata: ReportMetadata;
  images: ImageRecord[];
  helmReleases?: HelmRecord[];
  summary: ReportSummary;
}

/** One entry in the GET /api/reports listing (newest-first). */
export interface ReportEntry {
  filename: string;
  generatedAt: string;
  summary: ReportSummary;
  clusterName?: string;
}

/** Opt-in UI feature flags advertised by /api/me. */
export interface MeFeatures {
  timelineDeltas: boolean;
}

/** GET /api/me — drives the profile menu, scan gating, and feature flags. */
export interface MeResponse {
  authEnabled: boolean;
  email?: string;
  groups?: string[];
  canRunScan: boolean;
  features: MeFeatures;
}

/** POST /api/scan success body. */
export interface ScanResponse {
  jobName: string;
  namespace: string;
}
