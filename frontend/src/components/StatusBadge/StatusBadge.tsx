import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

export type Tone = "green" | "yellow" | "red" | "purple" | "muted";

// Map a semantic tone onto the Nebari status tokens. success/warning/destructive
// are light-bg + darker-fg pairs (and invert in dark mode), which reads as a
// chip; purple reuses the brand primary, muted the neutral surface.
const TONE_CLASS: Record<Tone, string> = {
  green: "bg-success text-success-foreground",
  yellow: "bg-warning text-warning-foreground",
  red: "bg-destructive text-destructive-foreground",
  purple: "bg-primary/12 text-primary",
  muted: "bg-muted text-muted-foreground",
};

/** A small status chip built on the Nebari Badge, coloured by tone. */
export function StatusBadge({ tone, children }: { tone: Tone; children: React.ReactNode }) {
  // Keep the Nebari Badge's pill shape (rounded-full) and sizing; only recolour
  // it per tone with the semantic status tokens.
  return <Badge className={cn(TONE_CLASS[tone])}>{children}</Badge>;
}
