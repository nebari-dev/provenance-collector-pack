import { useAtomValue, useSetAtom } from "jotai";
import { useEffect } from "react";
import { cn } from "@/lib/utils";
import { removeToastAtom, type Toast, toastsAtom } from "@/store/uiAtoms";

const TOAST_TTL_MS = 6000;

const ACCENT: Record<Toast["kind"], string> = {
  success: "border-l-success-foreground",
  error: "border-l-destructive-foreground",
};

function ToastCard({ toast }: { toast: Toast }) {
  const remove = useSetAtom(removeToastAtom);

  useEffect(() => {
    const id = setTimeout(() => remove(toast.id), TOAST_TTL_MS);
    return () => clearTimeout(id);
  }, [toast.id, remove]);

  return (
    <div
      role="status"
      data-testid={`toast-${toast.kind}`}
      className={cn(
        "animate-slide-up-fade rounded-md border border-border border-l-[3px] bg-card px-3.5 py-2.5 text-foreground text-xs shadow-lg",
        ACCENT[toast.kind],
      )}
    >
      <div>{toast.title}</div>
      {toast.sub ? (
        <div className="mt-0.5 font-mono text-[11px] text-muted-foreground">{toast.sub}</div>
      ) : null}
    </div>
  );
}

export function Toasts() {
  const toasts = useAtomValue(toastsAtom);
  return (
    <div
      data-testid="toasts"
      className="fixed right-5 bottom-5 z-[300] flex max-w-[420px] flex-col gap-2"
    >
      {toasts.map((t) => (
        <ToastCard key={t.id} toast={t} />
      ))}
    </div>
  );
}
