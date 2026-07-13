import { screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";
import { StatCards } from "@/components/StatCards";
import type { ImageRecord, ReportSummary } from "@/lib/types";
import { renderWithProviders } from "@/test/render";
import { ImagesTable } from "./ImagesTable";

function img(overrides: Partial<ImageRecord> & Pick<ImageRecord, "image">): ImageRecord {
  return {
    namespace: "default",
    workload: { kind: "Deployment", name: "app" },
    ...overrides,
  };
}

const IMAGES: ImageRecord[] = [
  img({ image: "nginx:1.25", namespace: "web", signature: { signed: true, verified: true } }),
  img({ image: "redis:7", namespace: "cache", signature: { signed: true, verified: false } }),
  img({ image: "busybox:latest", namespace: "web", signature: { signed: false, verified: false } }),
];

// tbody rows only (excludes the header row).
function bodyRowCount() {
  const rowgroups = screen.getAllByRole("rowgroup");
  const tbody = rowgroups[rowgroups.length - 1];
  return within(tbody).getAllByRole("row").length;
}

describe("ImagesTable", () => {
  it("renders a row per image", () => {
    renderWithProviders(<ImagesTable images={IMAGES} />);
    expect(screen.getByText("nginx:1.25")).toBeInTheDocument();
    expect(screen.getByText("redis:7")).toBeInTheDocument();
    expect(bodyRowCount()).toBe(3);
  });

  it("narrows rows via the search box", async () => {
    const user = userEvent.setup();
    renderWithProviders(<ImagesTable images={IMAGES} />);
    await user.type(screen.getByPlaceholderText("Search image, workload..."), "redis");
    expect(screen.getByText("redis:7")).toBeInTheDocument();
    expect(screen.queryByText("nginx:1.25")).not.toBeInTheDocument();
    expect(screen.getByText(/1 of 3 match/)).toBeInTheDocument();
  });

  it("paginates when there are more images than the page size", () => {
    const many = Array.from({ length: 30 }, (_, i) => img({ image: `img-${i}:1` }));
    renderWithProviders(<ImagesTable images={many} />);
    // 25 default page size → 25 body rows.
    expect(bodyRowCount()).toBe(25);
    expect(screen.getByText("1-25 of 30")).toBeInTheDocument();
  });

  it("clicking a stat card applies the matching table filter", async () => {
    const user = userEvent.setup();
    const summary: ReportSummary = {
      totalImages: 3,
      uniqueImages: 3,
      signedImages: 2,
      verifiedImages: 1,
      imagesWithSBOM: 0,
      imagesWithProvenance: 0,
      imagesWithUpdates: 0,
      totalHelmReleases: 0,
      helmReleasesWithUpdates: 0,
    };
    // Render both under one provider so they share the Jotai filter atoms.
    renderWithProviders(
      <>
        <StatCards summary={summary} />
        <ImagesTable images={IMAGES} />
      </>,
    );

    await user.click(screen.getByRole("button", { name: /Signed/ }));
    // Only the two signed images remain (nginx + redis), busybox is dropped.
    expect(screen.queryByText("busybox:latest")).not.toBeInTheDocument();
    expect(screen.getByText(/2 of 3 match/)).toBeInTheDocument();
  });
});
