import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { type RenderOptions, render } from "@testing-library/react";
import { Provider as JotaiProvider } from "jotai";
import type { ReactElement, ReactNode } from "react";

import { ThemeProvider } from "@/hooks/theme-provider";
import { THEME_STORAGE_KEY } from "@/lib/theme";

function TestProviders({ client, children }: { client: QueryClient; children: ReactNode }) {
  return (
    <QueryClientProvider client={client}>
      <JotaiProvider>
        <ThemeProvider storageKey={THEME_STORAGE_KEY}>{children}</ThemeProvider>
      </JotaiProvider>
    </QueryClientProvider>
  );
}

/**
 * Render a component wrapped in the app's providers (fresh QueryClient + Jotai
 * store per call, plus ThemeProvider) so tests are isolated from one another.
 */
export function renderWithProviders(ui: ReactElement, options?: Omit<RenderOptions, "wrapper">) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });

  return render(ui, {
    wrapper: ({ children }) => <TestProviders client={queryClient}>{children}</TestProviders>,
    ...options,
  });
}
