import { describe, expect, it } from "vitest";
import {
  EMPTY_FILTERS,
  filterImages,
  hasActiveFilters,
  type ImageFilters,
  sortImages,
  tier,
} from "./images";
import type { ImageRecord } from "./types";

function img(overrides: Partial<ImageRecord> & Pick<ImageRecord, "image">): ImageRecord {
  return {
    namespace: "default",
    workload: { kind: "Deployment", name: "app" },
    ...overrides,
  };
}

const images: ImageRecord[] = [
  img({
    image: "nginx:1.25",
    namespace: "web",
    signature: { signed: true, verified: true },
    sbom: { hasSBOM: true, format: "spdx" },
    provenance: { hasProvenance: true, predicateType: "slsa" },
    update: { currentTag: "1.25", updateAvailable: false },
  }),
  img({
    image: "redis:7",
    namespace: "cache",
    signature: { signed: true, verified: false },
    update: { currentTag: "7", newestAvailable: "7.2", updateAvailable: true },
  }),
  img({
    image: "busybox:latest",
    namespace: "web",
    signature: { signed: false, verified: false },
  }),
  img({ image: "mystery:1" }), // no signature/sbom/provenance/update = "Not checked"
];

function withFilters(patch: Partial<ImageFilters>): ImageFilters {
  return { ...EMPTY_FILTERS, ...patch };
}

describe("filterImages", () => {
  it("returns all images with empty filters", () => {
    expect(filterImages(images, EMPTY_FILTERS)).toHaveLength(4);
  });

  it("matches search across image, namespace, and workload", () => {
    expect(filterImages(images, withFilters({ search: "redis" })).map((i) => i.image)).toEqual([
      "redis:7",
    ]);
    expect(filterImages(images, withFilters({ search: "cache" })).map((i) => i.image)).toEqual([
      "redis:7",
    ]);
  });

  it("filters by namespace", () => {
    expect(filterImages(images, withFilters({ namespace: "web" }))).toHaveLength(2);
  });

  it("filters by signature state", () => {
    expect(
      filterImages(images, withFilters({ signature: "verified" })).map((i) => i.image),
    ).toEqual(["nginx:1.25"]);
    // "unsigned" excludes signed images but includes the not-checked one.
    expect(
      filterImages(images, withFilters({ signature: "unsigned" })).map((i) => i.image),
    ).toEqual(["busybox:latest", "mystery:1"]);
  });

  it("filters by sbom, provenance, and update presence", () => {
    expect(filterImages(images, withFilters({ sbom: "yes" }))).toHaveLength(1);
    expect(filterImages(images, withFilters({ provenance: "yes" }))).toHaveLength(1);
    expect(filterImages(images, withFilters({ update: "yes" })).map((i) => i.image)).toEqual([
      "redis:7",
    ]);
  });
});

describe("hasActiveFilters", () => {
  it("is false for the empty filter set", () => {
    expect(hasActiveFilters(EMPTY_FILTERS)).toBe(false);
  });
  it("is true when any field is set", () => {
    expect(hasActiveFilters(withFilters({ namespace: "web" }))).toBe(true);
  });
});

describe("sortImages", () => {
  it("is a no-op when no column is set", () => {
    expect(sortImages(images, { col: "", asc: true })).toBe(images);
  });

  it("sorts by image name ascending and descending", () => {
    expect(sortImages(images, { col: "image", asc: true }).map((i) => i.image)).toEqual([
      "busybox:latest",
      "mystery:1",
      "nginx:1.25",
      "redis:7",
    ]);
    expect(sortImages(images, { col: "image", asc: false }).map((i) => i.image)[0]).toBe("redis:7");
  });

  it("ranks signature: verified > signed > unsigned > not-checked", () => {
    expect(sortImages(images, { col: "signature", asc: false }).map((i) => i.image)).toEqual([
      "nginx:1.25", // verified
      "redis:7", // signed
      "busybox:latest", // unsigned
      "mystery:1", // absent
    ]);
  });

  it("does not mutate the input array", () => {
    const before = images.map((i) => i.image);
    sortImages(images, { col: "image", asc: true });
    expect(images.map((i) => i.image)).toEqual(before);
  });
});

describe("tier", () => {
  it("is neutral for a zero count regardless of total", () => {
    expect(tier(0, 10, 50, 80)).toBe("");
  });
  it("greens at/above the high threshold", () => {
    expect(tier(8, 10, 50, 80)).toBe("green");
  });
  it("yellows between low and high", () => {
    expect(tier(6, 10, 50, 80)).toBe("yellow");
  });
  it("reds below the low threshold when some signal exists", () => {
    expect(tier(1, 10, 50, 80)).toBe("red");
  });
});
