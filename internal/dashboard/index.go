package dashboard

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Provenance Collector</title>
<link rel="icon" href="https://raw.githubusercontent.com/nebari-dev/nebari-design/main/symbol/favicon.ico">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@300;400;500;600;700&family=IBM+Plex+Mono:wght@400;500&display=swap" rel="stylesheet">
<script>
  // Apply the persisted (or default) theme before first paint to avoid a flash
  // of the wrong theme. Default is "light" when nothing is stored — this app
  // intentionally opens light regardless of the OS preference.
  (function () {
    var KEY = 'provenance:themeMode';
    function stored() {
      try { var v = localStorage.getItem(KEY); if (v === 'light' || v === 'dark' || v === 'system') return v; } catch (e) {}
      return 'light';
    }
    function systemDark() {
      try { return window.matchMedia('(prefers-color-scheme: dark)').matches; } catch (e) { return false; }
    }
    function isDark(mode) { return mode === 'dark' || (mode === 'system' && systemDark()); }
    window.__themeKey = KEY;
    window.__themeStored = stored;
    window.__themeIsDark = isDark;
    window.__applyTheme = function (mode) {
      document.documentElement.classList.toggle('dark', isDark(mode));
    };
    window.__applyTheme(stored());
  })();
</script>
<style>
  /* Nebari design-system tokens (oklch), ported from the @nebari/theme shadcn
     registry via nebari-llm-serving-pack. Light lives on :root, dark on .dark.
     Treat these as upstream-managed; app-specific status / text tokens are
     layered on top of each block. Below the tokens, dashboard-facing aliases
     (--bg, --surface, --purple, --ok-*, ...) map onto these so the component
     rules read cleanly and stay theme-aware. */
  :root {
    --radius: 0.625rem;

    --background: oklch(100% 0 0);
    --foreground: oklch(26.94% 0.0037 286.15);

    --card: oklch(100% 0 0);
    --card-foreground: oklch(26.94% 0.0037 286.15);

    --popover: oklch(100% 0 0);
    --popover-foreground: oklch(26.94% 0.0037 286.15);

    --primary: oklch(55.06% 0.1886 311.45);
    --primary-foreground: oklch(100% 0 0);
    --primary-hover: oklch(47.01% 0.1577 311.26);

    --secondary: oklch(95.04% 0.0042 236.5);
    --secondary-foreground: oklch(32.95% 0.0209 254.12);

    --muted: oklch(94.94% 0.0013 286.37);
    --muted-foreground: oklch(54.86% 0.0154 285.88);
    --muted-foreground-strong: oklch(47.01% 0.0112 285.96);

    --accent: oklch(95.04% 0.0042 236.5);
    --accent-foreground: oklch(32.95% 0.0209 254.12);

    --destructive: oklch(95.06% 0.0247 29.93);
    --destructive-foreground: oklch(54.97% 0.2151 27.33);

    --warning: oklch(95.02% 0.0692 92.11);
    --warning-foreground: oklch(46.91% 0.096 91.9);

    --success: oklch(94.94% 0.0433 149.41);
    --success-foreground: oklch(47.01% 0.1313 149.41);

    --border: oklch(78.06% 0.0056 286.27);
    --input: oklch(69.88% 0.013 286.06);
    --border-strong: oklch(61.96% 0.0134 286);
    --ring: oklch(61.98% 0.2159 311.67);

    /* App-specific tokens (not part of the Nebari theme) */
    --body-background: #ffffff;
    --header-background: #f8f8f8;
    --text-secondary: #4e596a;

    --status-healthy-bg: #10b9811a;
    --status-healthy-fg: #047857;
    --status-healthy-dot: #10b981;

    --status-unhealthy-bg: #ef44441a;
    --status-unhealthy-fg: #b42318;
    --status-unhealthy-dot: #ef4444;

    --shadow: 0 4px 16px rgba(16,24,40,0.12);
  }

  .dark {
    --background: oklch(26.94% 0.0037 286.15);
    --foreground: oklch(97.91% 0 0);

    --card: oklch(33.01% 0.0052 286.11);
    --card-foreground: oklch(97.91% 0 0);

    --popover: oklch(26.94% 0.0037 286.15);
    --popover-foreground: oklch(97.91% 0 0);

    --primary: oklch(61.98% 0.2159 311.67);
    --primary-foreground: oklch(0% 0 0);
    --primary-hover: oklch(69.98% 0.1926 311.48);

    --secondary: oklch(40% 0.0269 250.57);
    --secondary-foreground: oklch(95.04% 0.0042 236.5);

    --muted: oklch(33.01% 0.0052 286.11);
    --muted-foreground: oklch(69.88% 0.013 286.06);
    --muted-foreground-strong: oklch(78.06% 0.0056 286.27);

    --accent: oklch(40% 0.0269 250.57);
    --accent-foreground: oklch(95.04% 0.0042 236.5);

    --destructive: oklch(33.1% 0.1332 27.42);
    --destructive-foreground: oklch(78.02% 0.1024 27.78);

    --warning: oklch(33.09% 0.0677 92.2);
    --warning-foreground: oklch(87.09% 0.1248 91.93);

    --success: oklch(32.88% 0.0887 149.73);
    --success-foreground: oklch(87.04% 0.0814 149.55);

    --border: oklch(47.01% 0.0112 285.96);
    --input: oklch(54.86% 0.0154 285.88);
    --border-strong: oklch(61.96% 0.0134 286);
    --ring: oklch(69.98% 0.1926 311.48);

    /* App-specific tokens (not part of the Nebari theme) */
    --body-background: #262628;
    --header-background: #353538;
    --text-secondary: #d0d5dd;

    --status-healthy-bg: #10b9811a;
    --status-healthy-fg: #6ee7b7;
    --status-healthy-dot: #34d399;

    --status-unhealthy-bg: #ef44441a;
    --status-unhealthy-fg: #fda29b;
    --status-unhealthy-dot: #f97066;

    --shadow: 0 4px 16px rgba(0,0,0,0.45);
  }

  :root, .dark {
    /* Dashboard-facing aliases onto the tokens above. */
    --bg: var(--body-background);
    --surface: var(--card);
    --surface-2: var(--muted);
    --text: var(--foreground);
    --text-muted: var(--muted-foreground);
    --faint: var(--muted-foreground);

    --purple: var(--primary);
    --purple-dark: var(--primary-hover);
    --purple-bg: color-mix(in oklch, var(--primary) 12%, transparent);
    --focus-ring: 0 0 0 2px color-mix(in oklch, var(--ring) 35%, transparent);
    --hover-tint: color-mix(in oklch, var(--primary) 6%, transparent);

    --ok-fg: var(--status-healthy-fg);
    --ok-bg: var(--status-healthy-bg);
    --ok-dot: var(--status-healthy-dot);
    --warn-fg: var(--warning-foreground);
    --warn-bg: var(--warning);
    --warn-dot: var(--warning-foreground);
    --bad-fg: var(--status-unhealthy-fg);
    --bad-bg: var(--status-unhealthy-bg);
    --bad-dot: var(--status-unhealthy-dot);

    --font: 'Geist', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    --mono: 'IBM Plex Mono', 'SF Mono', 'Consolas', monospace;
    --radius-sm: 8px;
  }

  * { margin: 0; padding: 0; box-sizing: border-box; }
  html, body { width: 100%; }
  body { font-family: var(--font); background: var(--bg); color: var(--text); line-height: 1.6; -webkit-font-smoothing: antialiased; font-size: 14px; }

  /* Header — full-width app bar: brand left, account menu right. */
  .appbar { display: flex; align-items: center; justify-content: space-between; height: 60px; padding: 0 40px; background: var(--header-background); border-bottom: 1px solid var(--border); }
  .brand { display: flex; align-items: center; gap: 10px; text-decoration: none; }
  .brand-logo { height: 28px; width: auto; color: var(--foreground); display: block; }
  .brand-sep { color: var(--border-strong); font-weight: 300; }
  .brand-app { font-weight: 600; font-size: 14px; color: var(--text); }

  /* Account / profile menu */
  .profile-menu { position: relative; }
  .profile-btn { display: inline-flex; align-items: center; gap: 10px; background: none; border: none; padding: 4px 6px; border-radius: var(--radius-sm); cursor: pointer; color: var(--text); font-family: var(--font); transition: background 0.15s ease; }
  .profile-btn:hover, .profile-btn[aria-expanded="true"] { background: var(--accent); }
  .profile-btn:focus-visible { outline: none; box-shadow: var(--focus-ring); }
  .avatar { display: inline-flex; align-items: center; justify-content: center; width: 32px; height: 32px; border-radius: 50%; background: var(--primary); color: var(--primary-foreground); font-size: 12px; font-weight: 600; flex-shrink: 0; }
  .profile-name { font-size: 13px; font-weight: 500; max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .profile-btn .chevron { color: var(--muted-foreground); transition: transform 0.15s ease; }
  .profile-btn[aria-expanded="true"] .chevron { transform: rotate(180deg); }
  .profile-popover { position: absolute; top: calc(100% + 6px); right: 0; width: 280px; background: var(--popover); color: var(--popover-foreground); border: 1px solid var(--border); border-radius: var(--radius); box-shadow: var(--shadow); padding: 6px; z-index: 220; }
  .profile-popover[hidden] { display: none; }
  .profile-head { padding: 8px 10px; border-bottom: 1px solid var(--border); margin-bottom: 6px; }
  .profile-head-name { font-size: 13px; font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .profile-head-email { font-size: 12px; color: var(--muted-foreground); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .theme-seg { display: flex; gap: 2px; background: var(--muted); border-radius: var(--radius-sm); padding: 3px; margin: 4px 4px 6px; }
  .theme-opt { flex: 1; display: inline-flex; align-items: center; justify-content: center; gap: 5px; padding: 6px 4px; border: none; background: none; border-radius: 6px; color: var(--muted-foreground); font-size: 12px; font-family: var(--font); cursor: pointer; transition: color 0.15s ease, background 0.15s ease; }
  .theme-opt:hover { color: var(--text); }
  .theme-opt svg { width: 15px; height: 15px; }
  .theme-opt.active { background: var(--background); color: var(--text); box-shadow: 0 1px 2px rgba(16,24,40,0.12); }
  .menu-sep { height: 1px; background: var(--border); margin: 6px 0; }
  .menu-item { display: block; width: 100%; text-align: left; padding: 8px 10px; border: none; background: none; border-radius: 6px; color: var(--text); font-size: 13px; font-family: var(--font); cursor: pointer; }
  .menu-item:hover:not(:disabled) { background: var(--accent); }
  .menu-item:disabled { color: var(--muted-foreground); opacity: 0.6; cursor: not-allowed; }

  /* Full-width content area */
  .container { width: 100%; padding: 24px 40px 40px; }

  /* Page title + toolbar */
  .page-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; flex-wrap: wrap; margin-bottom: 20px; }
  .page-title h1 { font-size: 22px; font-weight: 600; letter-spacing: -0.02em; }
  .page-sub { font-size: 13px; color: var(--muted-foreground); margin-top: 2px; max-width: 640px; }
  .page-actions { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
  .page-meta { display: flex; flex-direction: column; align-items: flex-end; font-size: 11px; color: var(--muted-foreground); line-height: 1.4; }

  /* Stats */
  .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(120px, 1fr)); gap: 10px; margin-bottom: 20px; }
  .stat { background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius); padding: 14px 16px; cursor: pointer; transition: all 0.15s ease; }
  .stat:hover { border-color: var(--purple); }
  .stat.active { border-color: var(--purple); background: var(--purple-bg); }
  .stat .value { font-size: 22px; font-weight: 700; letter-spacing: -0.02em; }
  .stat .label { font-size: 10px; color: var(--text-muted); margin-top: 1px; text-transform: uppercase; letter-spacing: 0.06em; font-weight: 500; }
  .stat.green .value { color: var(--ok-fg); }
  .stat.yellow .value { color: var(--warn-fg); }
  .stat.red .value { color: var(--bad-fg); }

  /* Panels */
  .panel { background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius); margin-bottom: 16px; overflow: hidden; }
  .panel-header { padding: 12px 20px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
  .panel-header h2 { font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: var(--text-muted); }
  .panel-body { padding: 0; }

  /* Filters */
  .filters { padding: 10px 20px; border-bottom: 1px solid var(--border); display: flex; gap: 8px; flex-wrap: wrap; align-items: center; }
  .filters input[type="text"] {
    background: var(--background); border: 1px solid var(--input); border-radius: var(--radius-sm);
    padding: 5px 10px; color: var(--text); font-size: 12px; font-family: var(--font);
    min-width: 180px; outline: none; transition: all 0.15s ease;
  }
  .filters input[type="text"]::placeholder { color: var(--faint); }
  .filters input[type="text"]:focus { border-color: var(--purple); box-shadow: var(--focus-ring); }
  .filters select {
    appearance: none; -webkit-appearance: none;
    /* Custom chevron (native arrow removed) so it isn't cramped against the
       value. Padding-right leaves room for the 12px chevron at right 10px. */
    background: var(--background) url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%238b8b9c' stroke-width='2.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'/%3E%3C/svg%3E") no-repeat right 10px center;
    border: 1px solid var(--input); border-radius: var(--radius-sm);
    padding: 5px 30px 5px 10px; color: var(--text); font-size: 12px; font-family: var(--font);
    min-width: 92px; outline: none; cursor: pointer; transition: border-color 0.15s ease, box-shadow 0.15s ease;
  }
  .filters select:hover { border-color: var(--border-strong); }
  .filters select:focus { border-color: var(--purple); box-shadow: var(--focus-ring); }
  .filter-label { font-size: 10px; color: var(--faint); text-transform: uppercase; letter-spacing: 0.04em; font-weight: 500; }
  .filter-reset { background: none; border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 4px 10px; color: var(--text-muted); font-size: 11px; font-family: var(--font); cursor: pointer; transition: all 0.15s ease; }
  .filter-reset:hover { border-color: var(--purple); color: var(--text); }

  /* Tables */
  table { width: 100%; border-collapse: collapse; font-size: 13px; }
  /* Sticky header — the appbar is not sticky, so the table header sticks to the
     top of the viewport once scrolled past. Surface background so rows scrolling
     under it don't bleed through. */
  thead th { position: sticky; top: 0; z-index: 50; background: var(--surface); }
  th { text-align: left; padding: 8px 20px; color: var(--faint); font-weight: 500; font-size: 10px; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid var(--border); cursor: pointer; user-select: none; transition: color 0.15s; }
  th:hover { color: var(--text-muted); }
  th .sort-arrow { font-size: 9px; margin-left: 3px; }
  td { padding: 8px 20px; border-bottom: 1px solid var(--border); font-size: 12px; }
  tr:last-child td { border-bottom: none; }
  tr:hover { background: var(--hover-tint); }
  /* Workload column truncation — keeps long ReplicaSet/StatefulSet names
     from wrapping mid-token. Hover the cell to see the full name. */
  td.workload { max-width: 260px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 11px; }

  /* Badges */
  .badge { display: inline-block; padding: 1px 7px; border-radius: 4px; font-size: 10px; font-weight: 500; letter-spacing: 0.02em; }
  .badge-green { background: var(--ok-bg); color: var(--ok-fg); }
  .badge-yellow { background: var(--warn-bg); color: var(--warn-fg); }
  .badge-red { background: var(--bad-bg); color: var(--bad-fg); }
  .badge-muted { background: var(--muted); color: var(--faint); }
  .badge-purple { background: var(--purple-bg); color: var(--purple); }

  /* Timeline */
  .timeline { display: flex; gap: 8px; overflow-x: auto; padding: 12px 20px; }
  .timeline-item { min-width: 120px; padding: 10px; background: var(--body-background); border: 1px solid var(--border); border-radius: var(--radius-sm); cursor: pointer; transition: all 0.15s ease; flex-shrink: 0; }
  .timeline-item:hover { border-color: var(--purple); }
  .timeline-item.active { border-color: var(--purple); background: var(--purple-bg); }
  .timeline-item .date { font-size: 12px; font-weight: 600; }
  .timeline-item .time { font-size: 11px; color: var(--text-muted); }
  .timeline-item .count { font-size: 10px; color: var(--faint); margin-top: 2px; }
  /* Delta vs the previous (older) scan — only rendered when a previous scan
     exists and the unique-image count differs. */
  .timeline-item .delta { display: inline-block; margin-top: 4px; padding: 1px 6px; border-radius: 3px; font-size: 10px; font-weight: 600; letter-spacing: 0.02em; }
  .timeline-item .delta-up { background: var(--ok-bg); color: var(--ok-fg); }
  .timeline-item .delta-down { background: var(--bad-bg); color: var(--bad-fg); }
  .timeline-item .delta-zero { background: var(--muted); color: var(--faint); }

  /* Pagination */
  .pagination { display: flex; align-items: center; justify-content: space-between; padding: 8px 20px; border-top: 1px solid var(--border); font-size: 11px; color: var(--text-muted); }
  .pagination .page-controls { display: flex; gap: 4px; }
  .pagination button { background: var(--background); border: 1px solid var(--border); border-radius: 4px; padding: 3px 10px; color: var(--text); font-size: 11px; font-family: var(--font); cursor: pointer; transition: all 0.15s ease; }
  .pagination button:hover:not(:disabled) { border-color: var(--purple); }
  .pagination button:disabled { opacity: 0.35; cursor: default; }
  .pagination select { background: var(--background); border: 1px solid var(--border); border-radius: 4px; padding: 3px 6px; color: var(--text); font-size: 11px; font-family: var(--font); }

  /* Detail panel */
  .detail-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 200; opacity: 0; pointer-events: none; transition: opacity 0.2s; }
  .detail-overlay.open { opacity: 1; pointer-events: auto; }
  .detail-panel { position: fixed; top: 0; right: 0; bottom: 0; width: 520px; max-width: 100vw; background: var(--surface); border-left: 1px solid var(--border); z-index: 201; transform: translateX(100%); transition: transform 0.25s ease; overflow-y: auto; }
  .detail-panel.open { transform: translateX(0); }
  .detail-close { position: absolute; top: 12px; right: 16px; background: none; border: none; color: var(--text-muted); font-size: 18px; cursor: pointer; padding: 4px 8px; border-radius: 4px; transition: all 0.15s; }
  .detail-close:hover { color: var(--text); background: var(--surface-2); }
  .detail-header { padding: 20px 24px 16px; border-bottom: 1px solid var(--border); }
  .detail-header h3 { font-size: 13px; font-weight: 600; word-break: break-all; font-family: var(--mono); line-height: 1.4; }
  .detail-header .detail-digest { font-size: 11px; color: var(--faint); font-family: var(--mono); margin-top: 4px; word-break: break-all; }
  .detail-section { padding: 16px 24px; border-bottom: 1px solid var(--border); }
  .detail-section:last-child { border-bottom: none; }
  .detail-section h4 { font-size: 10px; text-transform: uppercase; letter-spacing: 0.06em; color: var(--faint); font-weight: 600; margin-bottom: 10px; }
  .detail-row { display: flex; justify-content: space-between; align-items: center; padding: 4px 0; font-size: 12px; }
  .detail-row .label { color: var(--text-muted); }
  .detail-row .value { color: var(--text); font-family: var(--mono); font-size: 11px; text-align: right; max-width: 300px; word-break: break-all; }
  .detail-status { display: flex; align-items: center; gap: 8px; padding: 10px 0; }
  .detail-status .indicator { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
  .detail-status .indicator.green { background: var(--ok-dot); }
  .detail-status .indicator.yellow { background: var(--warn-dot); }
  .detail-status .indicator.red { background: var(--bad-dot); }
  .detail-status .indicator.muted { background: var(--faint); }
  .detail-status-text { font-size: 12px; }
  .detail-status-text .sub { font-size: 11px; color: var(--faint); margin-top: 1px; }

  /* Run Scan button */
  .btn-scan { background: var(--purple-bg); border: 1px solid var(--purple); border-radius: var(--radius-sm); padding: 5px 12px; color: var(--purple); font-size: 12px; font-family: var(--font); font-weight: 500; cursor: pointer; transition: all 0.15s ease; display: inline-flex; align-items: center; gap: 6px; }
  .btn-scan:hover:not(:disabled) { background: var(--purple); color: var(--primary-foreground); }
  .btn-scan:disabled { opacity: 0.5; cursor: progress; }
  .btn-scan .dot { width: 6px; height: 6px; border-radius: 50%; background: var(--purple); }
  .btn-scan:hover:not(:disabled) .dot { background: var(--primary-foreground); }

  /* Export menu */
  .export-menu { position: relative; display: inline-flex; }
  .export-btn { background: var(--background); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 4px 10px; color: var(--text-muted); font-size: 11px; font-family: var(--font); cursor: pointer; transition: all 0.15s ease; text-decoration: none; display: inline-flex; align-items: center; gap: 5px; }
  .export-btn:hover, .export-btn[aria-expanded="true"] { border-color: var(--purple); color: var(--text); }
  .export-btn .caret { font-size: 9px; transition: transform 0.15s ease; }
  .export-btn[aria-expanded="true"] .caret { transform: rotate(180deg); }
  .export-popover { position: absolute; top: calc(100% + 4px); right: 0; background: var(--popover); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 4px; min-width: 140px; box-shadow: var(--shadow); z-index: 150; display: flex; flex-direction: column; gap: 2px; }
  .export-popover[hidden] { display: none; }
  .export-item { padding: 6px 10px; font-size: 12px; color: var(--text-secondary); border-radius: 4px; text-decoration: none; cursor: pointer; transition: background 0.1s ease; }
  .export-item:hover, .export-item:focus { background: var(--accent); color: var(--text); outline: none; }

  /* Toast */
  .toast-wrap { position: fixed; bottom: 20px; right: 20px; z-index: 300; display: flex; flex-direction: column; gap: 8px; max-width: 420px; }
  .toast { background: var(--surface); border: 1px solid var(--border); border-left: 3px solid var(--purple); border-radius: var(--radius-sm); padding: 10px 14px; font-size: 12px; color: var(--text); box-shadow: var(--shadow); animation: slidein 0.2s ease; }
  .toast.error { border-left-color: var(--bad-dot); }
  .toast.success { border-left-color: var(--ok-dot); }
  .toast .sub { font-size: 11px; color: var(--text-muted); margin-top: 2px; font-family: var(--mono); }
  @keyframes slidein { from { transform: translateX(100%); opacity: 0; } to { transform: translateX(0); opacity: 1; } }

  /* Misc */
  .empty { text-align: center; padding: 40px 20px; color: var(--text-muted); }
  .empty p { margin-top: 4px; font-size: 12px; }
  .loading { text-align: center; padding: 40px; color: var(--text-muted); font-size: 12px; }
  .mono { font-family: var(--mono); font-size: 11px; }
  .text-muted { color: var(--text-muted); }
  .result-count { font-size: 11px; color: var(--faint); padding: 6px 20px; border-bottom: 1px solid var(--border); }
  tr.clickable { cursor: pointer; }

  @media (max-width: 640px) {
    .appbar { padding: 0 16px; }
    .container { padding: 16px; }
    .profile-name { display: none; }
  }
</style>
</head>
<body>
<header class="appbar">
  <a href="/" class="brand" aria-label="Provenance home">
    <svg class="brand-logo" viewBox="0 0 128 32" fill="none" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="Nebari">
      <g clip-path="url(#nb_clip)">
      <path d="M57.7327 24.8129H53.3465L46.0113 13.7068V24.8129H41.6252V6.80753H46.0113L53.3465 17.9657V6.80753H57.7327V24.8129Z" fill="currentColor"/>
      <path d="M74.3281 18.7103H64.4025C64.4713 19.5988 64.7571 20.2787 65.2618 20.75C65.7665 21.2214 66.386 21.455 67.1222 21.455C68.2172 21.455 68.9763 20.9941 69.4039 20.0701H74.0716C73.8317 21.0107 73.4 21.8575 72.7764 22.6104C72.1528 23.3634 71.3706 23.9536 70.43 24.3811C69.4894 24.8087 68.4382 25.0214 67.2744 25.0214C65.8729 25.0214 64.6236 24.7232 63.5307 24.1246C62.4358 23.526 61.5806 22.6709 60.9654 21.5593C60.3501 20.4476 60.0414 19.1483 60.0414 17.6612C60.0414 16.1741 60.3459 14.8748 60.9529 13.7631C61.5598 12.6515 62.4107 11.7964 63.5057 11.1978C64.6007 10.5992 65.8562 10.301 67.2765 10.301C68.6969 10.301 69.8919 10.5909 70.9702 11.1728C72.0485 11.7547 72.889 12.5827 73.4959 13.6609C74.1028 14.7392 74.4073 15.9948 74.4073 17.4318C74.4073 17.8427 74.3823 18.2702 74.3302 18.7145L74.3281 18.7103ZM69.917 16.2743C69.917 15.5213 69.6604 14.9228 69.1474 14.4785C68.6343 14.0343 67.994 13.8111 67.2244 13.8111C66.4548 13.8111 65.8687 14.0259 65.364 14.4535C64.8593 14.881 64.5464 15.488 64.4275 16.2743H69.917Z" fill="currentColor"/>
      <path d="M82.8688 10.9121C83.6384 10.5012 84.5185 10.2968 85.5113 10.2968C86.6918 10.2968 87.7596 10.5971 88.7169 11.1936C89.6742 11.7922 90.4313 12.6473 90.9861 13.759C91.5409 14.8706 91.8204 16.1616 91.8204 17.632C91.8204 19.1024 91.543 20.3976 90.9861 21.5176C90.4292 22.6375 89.6742 23.501 88.7169 24.1079C87.7596 24.7148 86.6918 25.0193 85.5113 25.0193C84.5018 25.0193 83.6217 24.8191 82.8688 24.4166C82.1159 24.0141 81.5256 23.4801 81.0981 22.8127V24.8129H76.712V5.83353H81.0981V12.5284C81.509 11.861 82.0992 11.3229 82.8688 10.9121ZM86.4477 15.0541C85.8408 14.4305 85.0921 14.1177 84.2036 14.1177C83.3151 14.1177 82.5914 14.4347 81.9845 15.0667C81.3776 15.6986 81.0731 16.5621 81.0731 17.657C81.0731 18.752 81.3755 19.6154 81.9845 20.2474C82.5914 20.8793 83.3318 21.1964 84.2036 21.1964C85.0754 21.1964 85.82 20.8752 86.4352 20.2349C87.0505 19.5946 87.3592 18.727 87.3592 17.632C87.3592 16.537 87.0547 15.6778 86.4477 15.0541Z" fill="currentColor"/>
      <path d="M94.0896 13.759C94.6444 12.6473 95.4014 11.7922 96.3588 11.1936C97.3161 10.5951 98.386 10.2968 99.5644 10.2968C100.574 10.2968 101.458 10.5012 102.219 10.9121C102.981 11.3229 103.567 11.861 103.976 12.5285V10.5012H108.362V24.8129H103.976V22.7856C103.548 23.453 102.954 23.9911 102.192 24.402C101.431 24.8129 100.547 25.0173 99.5373 25.0173C98.3735 25.0173 97.314 24.7149 96.3567 24.1058C95.3994 23.4989 94.6423 22.6355 94.0875 21.5155C93.5327 20.3955 93.2532 19.1003 93.2532 17.6299C93.2532 16.1596 93.5306 14.8685 94.0875 13.7569L94.0896 13.759ZM103.066 15.0667C102.459 14.4347 101.719 14.1177 100.847 14.1177C99.9753 14.1177 99.2349 14.4306 98.6279 15.0542C98.021 15.6778 97.7165 16.5371 97.7165 17.632C97.7165 18.727 98.0189 19.5946 98.6279 20.2349C99.2349 20.8773 99.9753 21.1964 100.847 21.1964C101.719 21.1964 102.459 20.8794 103.066 20.2474C103.673 19.6155 103.978 18.752 103.978 17.657C103.978 16.5621 103.673 15.6986 103.066 15.0667Z" fill="currentColor"/>
      <path d="M117.853 11.0268C118.623 10.5742 119.478 10.3469 120.419 10.3469V14.9895H119.213C118.118 14.9895 117.299 15.2252 116.75 15.6944C116.202 16.1658 115.928 16.9896 115.928 18.1701V24.8129H111.542V10.5012H115.928V12.8872C116.441 12.1009 117.082 11.4815 117.851 11.0268H117.853Z" fill="currentColor"/>
      <path d="M122.842 8.33422C122.354 7.88163 122.11 7.3206 122.11 6.65528C122.11 5.98996 122.354 5.40389 122.842 4.94922C123.33 4.49664 123.958 4.2693 124.728 4.2693C125.497 4.2693 126.1 4.49664 126.588 4.94922C127.076 5.40181 127.318 5.97119 127.318 6.65528C127.318 7.33937 127.074 7.88163 126.588 8.33422C126.1 8.7868 125.481 9.01414 124.728 9.01414C123.975 9.01414 123.33 8.7868 122.842 8.33422ZM126.907 10.5012V24.8129H122.521V10.5012H126.907Z" fill="currentColor"/>
      <path d="M33.5455 9.77332H27.5972C26.7978 9.77332 26.1498 10.4214 26.1498 11.2208C26.1498 12.0201 26.7978 12.6682 27.5972 12.6682H33.5455C34.3449 12.6682 34.9929 12.0201 34.9929 11.2208C34.9929 10.4214 34.3449 9.77332 33.5455 9.77332Z" fill="#20AAA1"/>
      <path d="M31.9041 17.7029H27.7057C26.9063 17.7029 26.2582 18.351 26.2582 19.1503C26.2582 19.9497 26.9063 20.5978 27.7057 20.5978H31.9041C32.7035 20.5978 33.3515 19.9497 33.3515 19.1503C33.3515 18.351 32.7035 17.7029 31.9041 17.7029Z" fill="#20AAA1"/>
      <path d="M22.1746 10.1237H6.04839C5.24899 10.1237 4.60095 10.7717 4.60095 11.5711C4.60095 12.3705 5.24899 13.0186 6.04839 13.0186H22.1746C22.974 13.0186 23.622 12.3705 23.622 11.5711C23.622 10.7717 22.974 10.1237 22.1746 10.1237Z" fill="#20AAA1"/>
      <path d="M15.609 6.41751H4.66142C3.86203 6.41751 3.21399 7.06555 3.21399 7.86495C3.21399 8.66434 3.86203 9.31238 4.66142 9.31238H15.609C16.4084 9.31238 17.0564 8.66434 17.0564 7.86495C17.0564 7.06555 16.4084 6.41751 15.609 6.41751Z" fill="#20AAA1"/>
      <path d="M5.07019 13.8257H1.44744C0.648039 13.8257 0 14.4738 0 15.2731C0 16.0725 0.648039 16.7206 1.44744 16.7206H5.07019C5.86959 16.7206 6.51763 16.0725 6.51763 15.2731C6.51763 14.4738 5.86959 13.8257 5.07019 13.8257Z" fill="#20AAA1"/>
      <path d="M18.4558 2.7009H13.2104C12.411 2.7009 11.763 3.34894 11.763 4.14833C11.763 4.94773 12.411 5.59577 13.2104 5.59577H18.4558C19.2552 5.59577 19.9033 4.94773 19.9033 4.14833C19.9033 3.34894 19.2552 2.7009 18.4558 2.7009Z" fill="#20AAA1"/>
      <path d="M29.9894 13.8153H16.756C15.9566 13.8153 15.3086 14.4633 15.3086 15.2627C15.3086 16.0621 15.9566 16.7101 16.756 16.7101H29.9894C30.7888 16.7101 31.4369 16.0621 31.4369 15.2627C31.4369 14.4633 30.7888 13.8153 29.9894 13.8153Z" fill="#20AAA1"/>
      <path d="M22.7502 28.5795C25.3364 27.9851 27.1029 26.9507 28.3126 25.3218C28.9883 24.4124 29.4305 23.3362 29.735 21.8554H26.3271C26.0455 23.0797 25.2968 24.0141 24.0975 24.6398C22.6876 25.3739 21.1359 25.5408 19.6363 25.7014C19.5195 25.7139 19.4027 25.7264 19.2859 25.7389C18.8459 25.7869 18.7666 25.4448 18.7291 25.2821C18.3599 23.695 17.4506 22.4957 16.0261 21.7178C15.9135 21.6573 15.7612 21.5551 15.7278 21.3841C15.6965 21.2172 15.7946 21.0733 15.8843 20.9669C16.5955 20.1139 17.0042 19.094 17.3025 18.1701H13.7673C13.3147 19.1232 12.5618 19.8136 11.4648 20.2808C9.86509 20.9607 8.19449 21.0483 6.49678 20.5394C5.21828 20.1577 4.44868 19.3631 4.20675 18.1784L0.863464 18.1722C1.28476 21.2089 3.07424 23.1277 6.18185 23.8764C7.66683 24.2331 9.21646 24.2101 10.7139 24.1893L10.9475 24.1851C11.1061 24.183 11.2625 24.1788 11.421 24.1747C11.6525 24.1684 11.8861 24.1622 12.1218 24.1622C12.6161 24.1622 13.1166 24.1893 13.6151 24.2998C14.8852 24.5814 15.5422 25.3802 15.5672 26.6774C15.5881 27.7703 15.1939 28.7735 14.3617 29.7475C13.7006 30.5213 12.8934 31.1324 12.1113 31.7226C11.9904 31.8144 11.8694 31.9061 11.7463 32L16.7873 31.9937C17.5027 31.1741 18.0033 30.3711 18.3203 29.5452C18.4767 29.1385 18.752 29.0676 19.0252 29.0509C20.1369 28.9862 21.4487 28.8778 22.7481 28.5795H22.7502Z" fill="#BA18DD"/>
      <path d="M26.1039 8.29459C28.4169 8.30502 30.296 6.48634 30.3148 4.21925C30.3336 1.90627 28.4878 0.0145932 26.1998 -6.30273e-06C23.8889 -0.0146058 21.9868 1.84787 21.9847 4.12331C21.9827 6.41334 23.8285 8.28416 26.1039 8.29459Z" fill="#EAB54E"/>
      </g>
      <defs><clipPath id="nb_clip"><rect width="127.318" height="32" fill="white"/></clipPath></defs>
    </svg>
    <span class="brand-sep">/</span>
    <span class="brand-app">Provenance</span>
  </a>

  <div class="profile-menu">
    <button id="profile-toggle" class="profile-btn" type="button" aria-haspopup="true" aria-expanded="false" aria-controls="profile-popover" onclick="toggleProfileMenu()">
      <span class="avatar" id="profile-initials">U</span>
      <span class="profile-name" id="profile-name">Account</span>
      <svg class="chevron" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg>
    </button>
    <div id="profile-popover" class="profile-popover" role="menu" hidden>
      <div class="profile-head">
        <div class="profile-head-name" id="profile-head-name">Signed in</div>
        <div class="profile-head-email" id="profile-head-email"></div>
      </div>
      <div class="theme-seg" role="group" aria-label="Theme">
        <button type="button" class="theme-opt" data-mode="light" onclick="setThemeMode('light')" title="Light mode" aria-label="Light mode">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="4"></circle><path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"></path></svg>
          <span>Light</span>
        </button>
        <button type="button" class="theme-opt" data-mode="dark" onclick="setThemeMode('dark')" title="Dark mode" aria-label="Dark mode">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
          <span>Dark</span>
        </button>
        <button type="button" class="theme-opt" data-mode="system" onclick="setThemeMode('system')" title="System theme" aria-label="System theme">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2"></rect><path d="M8 21h8M12 17v4"></path></svg>
          <span>System</span>
        </button>
      </div>
      <div class="menu-sep"></div>
      <button class="menu-item" type="button" disabled title="Sign out is not available in this deployment">Sign out</button>
    </div>
  </div>
</header>

<div class="container">
  <div class="page-head">
    <div class="page-title">
      <h1>Provenance</h1>
      <p class="page-sub">Container image provenance, signatures, SBOMs, and available updates discovered across your cluster.</p>
    </div>
    <div class="page-actions">
      <div class="page-meta">
        <span id="cluster-name"></span>
        <span id="last-updated"></span>
      </div>
      <div class="export-menu">
        <button id="export-toggle" class="export-btn" type="button" aria-haspopup="true" aria-expanded="false" aria-controls="export-popover" onclick="toggleExportMenu()" title="Download the selected report">
          <span>Export</span><span class="caret">&#9662;</span>
        </button>
        <div id="export-popover" class="export-popover" role="menu" hidden>
          <a id="export-csv" class="export-item" role="menuitem" href="/api/export?format=csv" download>CSV</a>
          <a id="export-md" class="export-item" role="menuitem" href="/api/export?format=markdown" download>Markdown</a>
          <a id="export-json" class="export-item" role="menuitem" href="/api/reports/latest" download="provenance-report.json">JSON</a>
        </div>
      </div>
      <button id="btn-scan" class="btn-scan" style="display:none" onclick="runScan()" title="Trigger a manual provenance scan">
        <span class="dot"></span><span>Run Scan</span>
      </button>
    </div>
  </div>

  <div class="stats" id="stats"><div class="loading">Loading...</div></div>

  <div class="panel">
    <div class="panel-header"><h2>Timeline</h2></div>
    <div class="timeline" id="timeline"></div>
  </div>

  <div class="panel">
    <div class="panel-header">
      <h2>Container Images</h2>
      <span class="text-muted" style="font-size:11px" id="image-count"></span>
    </div>
    <div class="filters" id="image-filters"></div>
    <div id="image-result-count"></div>
    <div class="panel-body" id="images-table"></div>
    <div class="pagination" id="images-pagination"></div>
  </div>

  <div class="panel">
    <div class="panel-header">
      <h2>Helm Releases</h2>
      <span class="text-muted" style="font-size:11px" id="helm-count"></span>
    </div>
    <div class="panel-body" id="helm-table"></div>
  </div>
</div>

<div class="toast-wrap" id="toast-wrap"></div>

<div class="detail-overlay" id="detail-overlay" onclick="closeDetail()"></div>
<div class="detail-panel" id="detail-panel">
  <button class="detail-close" onclick="closeDetail()">&times;</button>
  <div id="detail-content"></div>
</div>

<script>
let reports = [];
let currentReport = null;
let imageFilters = { search: '', namespace: '', signature: '', sbom: '', provenance: '', update: '' };
let statFilter = '';
let imageSortCol = '';
let imageSortAsc = true;
let pageSize = 25;
let currentPage = 0;
let lastFilteredImages = [];
// Feature flags from /api/me. Defaults are all-off to match the chart's
// opt-in posture; initConfig() overwrites this with the server's response.
let features = { timelineDeltas: false };

// --- Theme control -------------------------------------------------------
// The pre-paint script in <head> already applied the stored (or default-light)
// theme and exposed helpers on window. Here we wire the segmented control and
// keep "system" mode in sync with the OS preference.
function setThemeMode(mode) {
  try { localStorage.setItem(window.__themeKey, mode); } catch (e) {}
  window.__applyTheme(mode);
  syncThemeUI(mode);
}

function syncThemeUI(mode) {
  document.querySelectorAll('.theme-opt').forEach(function (el) {
    el.classList.toggle('active', el.getAttribute('data-mode') === mode);
  });
}

(function initThemeSync() {
  syncThemeUI(window.__themeStored());
  try {
    const mq = window.matchMedia('(prefers-color-scheme: dark)');
    mq.addEventListener('change', function () {
      // Only "system" mode should follow the OS; light/dark are explicit.
      if (window.__themeStored() === 'system') window.__applyTheme('system');
    });
  } catch (e) { /* matchMedia unsupported — explicit modes still work */ }
})();

// --- Profile menu --------------------------------------------------------
// initials returns up to two uppercase characters derived from the user's
// email (this deployment's /api/me has no name field). Falls back to "U".
function initials(email) {
  if (email) {
    const local = email.split('@')[0].replace(/[^a-zA-Z0-9]/g, '');
    if (local) return local.slice(0, 2).toUpperCase();
    return email.slice(0, 2).toUpperCase();
  }
  return 'U';
}

function setProfile(me) {
  const email = (me && me.email) || '';
  const display = email || 'Account';
  document.getElementById('profile-initials').textContent = initials(email);
  document.getElementById('profile-name').textContent = display;
  document.getElementById('profile-head-name').textContent = email ? 'Signed in' : 'Not signed in';
  document.getElementById('profile-head-email').textContent = email;
}

function toggleProfileMenu() {
  const btn = document.getElementById('profile-toggle');
  setProfileMenuOpen(btn.getAttribute('aria-expanded') !== 'true');
}

function setProfileMenuOpen(open) {
  const btn = document.getElementById('profile-toggle');
  const menu = document.getElementById('profile-popover');
  if (!btn || !menu) return;
  btn.setAttribute('aria-expanded', open ? 'true' : 'false');
  if (open) { menu.removeAttribute('hidden'); } else { menu.setAttribute('hidden', ''); }
}

async function init() {
  // Await the config fetch so renderTimeline() below sees the correct feature
  // flags on first paint. The call is fast and fail-quiet — defaults stand
  // if it errors.
  await initConfig();
  try {
    const res = await fetch('/api/reports');
    reports = await res.json();
    if (!reports || reports.length === 0) { showEmpty(); return; }
    renderTimeline();
    await loadReport(reports[0].filename);
  } catch (e) {
    document.getElementById('stats').innerHTML = '<div class="empty"><p>Failed to load reports</p></div>';
  }
}

// initConfig resolves the scan-button visibility (auth), the profile display,
// and the dashboard feature flags from /api/me in a single fetch. Fail-quiet:
// if /api/me errors, the scan button stays hidden, the profile shows the
// signed-out state, and feature flags stay at their defaults (all-off).
async function initConfig() {
  try {
    const res = await fetch('/api/me');
    if (!res.ok) return;
    const me = await res.json();
    setProfile(me);
    if (me && me.canRunScan) {
      document.getElementById('btn-scan').style.display = 'inline-flex';
    }
    if (me && me.features) {
      features = Object.assign(features, me.features);
    }
  } catch (e) { /* keep defaults */ }
}

let scanPollTimer = null;
const SCAN_POLL_MS = 5000;
const SCAN_POLL_MAX = 60; // 5 minutes worst case

async function runScan() {
  const btn = document.getElementById('btn-scan');
  btn.disabled = true;
  btn.querySelector('span:last-child').textContent = 'Starting...';
  try {
    const res = await fetch('/api/scan', { method: 'POST' });
    const text = await res.text();
    if (!res.ok) {
      showToast('error', 'Scan request failed', res.status + ' ' + (text || res.statusText));
      return;
    }
    let body = {};
    try { body = JSON.parse(text); } catch (e) {}
    showToast('success', 'Scan started', body.jobName ? 'Job: ' + body.jobName : '');
    btn.querySelector('span:last-child').textContent = 'Scan running';
    pollForNewReport();
  } catch (e) {
    showToast('error', 'Scan request failed', String(e));
  } finally {
    // Re-enable after the poller decides we're done; until then keep disabled.
    if (!scanPollTimer) {
      btn.disabled = false;
      btn.querySelector('span:last-child').textContent = 'Run Scan';
    }
  }
}

function pollForNewReport() {
  const baselineCount = reports.length;
  let ticks = 0;
  const tick = async () => {
    ticks++;
    try {
      const res = await fetch('/api/reports');
      const latest = await res.json();
      if (latest && latest.length > baselineCount) {
        // New report landed — refresh UI and stop polling.
        reports = latest;
        renderTimeline();
        await loadReport(reports[0].filename);
        showToast('success', 'New report available', reports[0].filename);
        finishScanPolling();
        return;
      }
    } catch (e) { /* keep polling */ }
    if (ticks >= SCAN_POLL_MAX) {
      showToast('error', 'Scan still running', 'Stopped watching after ' + (SCAN_POLL_MAX * SCAN_POLL_MS / 60000) + ' min — refresh later');
      finishScanPolling();
      return;
    }
    scanPollTimer = setTimeout(tick, SCAN_POLL_MS);
  };
  scanPollTimer = setTimeout(tick, SCAN_POLL_MS);
}

function finishScanPolling() {
  if (scanPollTimer) { clearTimeout(scanPollTimer); scanPollTimer = null; }
  const btn = document.getElementById('btn-scan');
  btn.disabled = false;
  btn.querySelector('span:last-child').textContent = 'Run Scan';
}

function showToast(kind, title, sub) {
  const wrap = document.getElementById('toast-wrap');
  const el = document.createElement('div');
  el.className = 'toast ' + (kind || '');
  el.innerHTML = '<div>' + esc(title) + '</div>' + (sub ? '<div class="sub">' + esc(sub) + '</div>' : '');
  wrap.appendChild(el);
  setTimeout(() => { el.style.transition = 'opacity 0.3s'; el.style.opacity = '0'; setTimeout(() => el.remove(), 300); }, 6000);
}

function showEmpty() {
  document.getElementById('stats').innerHTML = '';
  document.getElementById('timeline').innerHTML = '<div class="empty"><p>No reports yet. Run the collector to generate a provenance report.</p></div>';
  document.getElementById('images-table').innerHTML = '';
  document.getElementById('helm-table').innerHTML = '';
  document.getElementById('images-pagination').innerHTML = '';
}

function renderTimeline() {
  document.getElementById('timeline').innerHTML = reports.map((r, i) => {
    const d = new Date(r.generatedAt);
    // reports[] is sorted DESC by generatedAt, so the older scan to compare
    // against is the next index up (i + 1). The oldest card has no neighbor
    // to diff against; render no delta there. Also gated on the
    // timelineDeltas feature flag — off by default in the chart, enable via
    // webUI.features.timelineDeltas: true.
    const prev = reports[i + 1];
    let deltaHTML = '';
    if (features.timelineDeltas && prev) {
      const delta = (r.summary.uniqueImages || 0) - (prev.summary.uniqueImages || 0);
      const cls = delta > 0 ? 'delta-up' : delta < 0 ? 'delta-down' : 'delta-zero';
      const text = delta > 0 ? '+' + delta : String(delta);
      deltaHTML = '<div class="delta ' + cls + '" data-delta="' + delta + '" title="' + (
        delta > 0 ? delta + ' new unique image(s) vs previous scan' :
        delta < 0 ? Math.abs(delta) + ' image(s) gone vs previous scan' :
        'no change vs previous scan'
      ) + '">' + text + '</div>';
    }
    return '<div class="timeline-item' + (i === 0 ? ' active' : '') + '" onclick="selectReport(' + i + ')" id="tl-' + i + '">' +
      '<div class="date">' + d.toLocaleDateString() + '</div>' +
      '<div class="time">' + d.toLocaleTimeString() + '</div>' +
      '<div class="count">' + r.summary.totalImages + ' images</div>' +
      deltaHTML + '</div>';
  }).join('');
}

async function selectReport(idx) {
  document.querySelectorAll('.timeline-item').forEach(el => el.classList.remove('active'));
  document.getElementById('tl-' + idx).classList.add('active');
  await loadReport(reports[idx].filename);
}

async function loadReport(filename) {
  const res = await fetch('/api/reports/' + filename);
  currentReport = await res.json();
  resetFilters(); currentPage = 0;
  renderStats(currentReport);
  renderFilters(currentReport);
  renderImages(currentReport);
  renderHelm(currentReport);
  updateExportLinks(filename);
  const d = new Date(currentReport.metadata.generatedAt);
  document.getElementById('last-updated').textContent = d.toLocaleString();
  // Falls back to "Local" when clusterName is unset (typical for standalone /
  // dev installs) so the header doesn't show a blank slot before the timestamp.
  document.getElementById('cluster-name').textContent = currentReport.metadata.clusterName || 'Local';
}

// updateExportLinks repoints the Export menu items at whichever report is
// currently being viewed. Without this, the CSV / MD / JSON entries would
// always export the latest report regardless of the timeline selection.
function updateExportLinks(filename) {
  const enc = encodeURIComponent(filename);
  const csv = document.getElementById('export-csv');
  const md = document.getElementById('export-md');
  const json = document.getElementById('export-json');
  if (csv) csv.href = '/api/export?format=csv&filename=' + enc;
  if (md) md.href = '/api/export?format=markdown&filename=' + enc;
  if (json) json.href = '/api/reports/' + enc;
}

// toggleExportMenu opens / closes the Export dropdown. Closing is also wired
// to outside-clicks and Escape (see the listeners below).
function toggleExportMenu() {
  const btn = document.getElementById('export-toggle');
  const menu = document.getElementById('export-popover');
  if (!btn || !menu) return;
  const open = btn.getAttribute('aria-expanded') === 'true';
  setExportMenuOpen(!open);
}

function setExportMenuOpen(open) {
  const btn = document.getElementById('export-toggle');
  const menu = document.getElementById('export-popover');
  if (!btn || !menu) return;
  btn.setAttribute('aria-expanded', open ? 'true' : 'false');
  if (open) {
    menu.removeAttribute('hidden');
  } else {
    menu.setAttribute('hidden', '');
  }
}

document.addEventListener('click', (ev) => {
  const exportWrap = document.querySelector('.export-menu');
  if (exportWrap && !exportWrap.contains(ev.target)) setExportMenuOpen(false);
  const profileWrap = document.querySelector('.profile-menu');
  if (profileWrap && !profileWrap.contains(ev.target)) setProfileMenuOpen(false);
  // Clicking an item should close the popover after the browser starts the
  // download. The role=menuitem anchors live inside .export-menu so the
  // outside-click branch above doesn't fire; close explicitly here.
  if (ev.target && ev.target.classList && ev.target.classList.contains('export-item')) {
    setExportMenuOpen(false);
  }
});

document.addEventListener('keydown', (ev) => {
  if (ev.key === 'Escape') { setExportMenuOpen(false); setProfileMenuOpen(false); }
});

function resetFilters() {
  imageFilters = { search: '', namespace: '', signature: '', sbom: '', provenance: '', update: '' };
  statFilter = '';
}

function renderStats(r) {
  const s = r.summary;
  const total = s.uniqueImages || 0;

  // "N of M" everywhere with a denominator. Images and Helm are the absolutes,
  // the rest are fractions of the unique-image total. Keeps the 0-state from
  // looking like "0 / 0 / 0 / 0" without context.
  const ratio = (n) => total ? n + ' / ' + total : String(n);

  // Threshold helper for the four ratio-style cards. Painting red on a "0 of N"
  // when the cluster has *no* signed/provenance/SBOM artifacts at all is
  // alarmist — a vanilla cluster pulling public images is expected to be
  // unsigned. Reserve red for the case where some images carry the attribute
  // but the overall percentage is poor; otherwise stay muted.
  const tier = (n, lo, hi) => {
    if (n === 0) return '';                       // no positive signal — neutral, not "bad"
    const pct = total ? Math.round(n / total * 100) : 0;
    if (pct >= hi) return 'green';
    if (pct >= lo) return 'yellow';
    return 'red';
  };

  document.getElementById('stats').innerHTML =
    statCard('all',        String(s.uniqueImages),               'Images',   '') +
    statCard('signed',     ratio(s.signedImages),                'Signed',   tier(s.signedImages, 50, 80)) +
    statCard('verified',   ratio(s.verifiedImages),              'Verified', tier(s.verifiedImages, 50, 80)) +
    statCard('provenance', ratio(s.imagesWithProvenance || 0),   'SLSA',     tier(s.imagesWithProvenance || 0, 1, 50)) +
    statCard('sbom',       ratio(s.imagesWithSBOM),              'SBOM',     tier(s.imagesWithSBOM, 1, 50)) +
    statCard('updates',    ratio(s.imagesWithUpdates),           'Updates',  s.imagesWithUpdates > 0 ? 'yellow' : 'green') +
    statCard('helm',       String(s.totalHelmReleases),          'Helm',     '');
}

function statCard(id, value, label, cls) {
  const tip = id === 'signed' && value.indexOf('0 /') === 0
    ? ' title="No signatures found. Common on clusters pulling unsigned public images."'
    : '';
  return '<div class="stat ' + cls + (statFilter === id ? ' active' : '') + '" onclick="toggleStatFilter(\'' + id + '\')" id="stat-' + id + '"' + tip + '><div class="value">' + value + '</div><div class="label">' + label + '</div></div>';
}

function toggleStatFilter(id) {
  statFilter = statFilter === id ? '' : id;
  document.querySelectorAll('.stat').forEach(el => el.classList.remove('active'));
  if (statFilter) document.getElementById('stat-' + statFilter).classList.add('active');
  imageFilters = { search: imageFilters.search, namespace: imageFilters.namespace, signature: '', sbom: '', provenance: '', update: '' };
  switch (statFilter) {
    case 'signed': imageFilters.signature = 'signed'; break;
    case 'verified': imageFilters.signature = 'verified'; break;
    case 'sbom': imageFilters.sbom = 'yes'; break;
    case 'provenance': imageFilters.provenance = 'yes'; break;
    case 'updates': imageFilters.update = 'yes'; break;
  }
  currentPage = 0; syncFilterUI(); renderImages(currentReport);
}

function renderFilters(r) {
  const ns = [...new Set((r.images || []).map(i => i.namespace))].sort();
  document.getElementById('image-filters').innerHTML =
    '<input type="text" id="f-search" placeholder="Search image, workload..." oninput="onFilterChange()">' +
    '<span class="filter-label">NS</span><select id="f-namespace" onchange="onFilterChange()"><option value="">All</option>' + ns.map(n => '<option value="' + esc(n) + '">' + esc(n) + '</option>').join('') + '</select>' +
    '<span class="filter-label">Sig</span><select id="f-signature" onchange="onFilterChange()"><option value="">All</option><option value="verified">Verified</option><option value="signed">Signed</option><option value="unsigned">Unsigned</option></select>' +
    '<span class="filter-label">SBOM</span><select id="f-sbom" onchange="onFilterChange()"><option value="">All</option><option value="yes">Yes</option><option value="no">No</option></select>' +
    '<span class="filter-label">SLSA</span><select id="f-provenance" onchange="onFilterChange()"><option value="">All</option><option value="yes">Yes</option><option value="no">No</option></select>' +
    '<span class="filter-label">Update</span><select id="f-update" onchange="onFilterChange()"><option value="">All</option><option value="yes">Yes</option><option value="no">No</option></select>' +
    '<button class="filter-reset" onclick="clearFilters()">Clear</button>';
}

function syncFilterUI() {
  ['search','namespace','signature','sbom','provenance','update'].forEach(id => { const el = document.getElementById('f-' + id); if (el) el.value = imageFilters[id]; });
}

function onFilterChange() {
  imageFilters.search = (document.getElementById('f-search').value || '').toLowerCase();
  imageFilters.namespace = document.getElementById('f-namespace').value;
  imageFilters.signature = document.getElementById('f-signature').value;
  imageFilters.sbom = document.getElementById('f-sbom').value;
  imageFilters.provenance = document.getElementById('f-provenance').value;
  imageFilters.update = document.getElementById('f-update').value;
  statFilter = ''; currentPage = 0;
  document.querySelectorAll('.stat').forEach(el => el.classList.remove('active'));
  renderImages(currentReport);
}

function clearFilters() {
  resetFilters(); currentPage = 0; syncFilterUI();
  document.querySelectorAll('.stat').forEach(el => el.classList.remove('active'));
  renderImages(currentReport);
}

function filterImages(images) {
  return images.filter(img => {
    if (imageFilters.search) {
      const hay = (img.image + ' ' + img.namespace + ' ' + img.workload.kind + '/' + img.workload.name).toLowerCase();
      if (!hay.includes(imageFilters.search)) return false;
    }
    if (imageFilters.namespace && img.namespace !== imageFilters.namespace) return false;
    if (imageFilters.signature) {
      const sig = img.signature;
      if (imageFilters.signature === 'verified' && !(sig && sig.verified)) return false;
      if (imageFilters.signature === 'signed' && !(sig && sig.signed)) return false;
      if (imageFilters.signature === 'unsigned' && sig && sig.signed) return false;
    }
    if (imageFilters.sbom) { const h = img.sbom && img.sbom.hasSBOM; if (imageFilters.sbom === 'yes' && !h) return false; if (imageFilters.sbom === 'no' && h) return false; }
    if (imageFilters.provenance) { const h = img.provenance && img.provenance.hasProvenance; if (imageFilters.provenance === 'yes' && !h) return false; if (imageFilters.provenance === 'no' && h) return false; }
    if (imageFilters.update) { const h = img.update && img.update.updateAvailable; if (imageFilters.update === 'yes' && !h) return false; if (imageFilters.update === 'no' && h) return false; }
    return true;
  });
}

function sortImages(images) {
  if (!imageSortCol) return images;
  const sorted = [...images];
  sorted.sort((a, b) => {
    let va, vb;
    switch (imageSortCol) {
      case 'image': va = a.image; vb = b.image; break;
      case 'namespace': va = a.namespace; vb = b.namespace; break;
      case 'workload': va = a.workload.kind + '/' + a.workload.name; vb = b.workload.kind + '/' + b.workload.name; break;
      case 'signature': va = a.signature ? (a.signature.verified ? 2 : a.signature.signed ? 1 : 0) : -1; vb = b.signature ? (b.signature.verified ? 2 : b.signature.signed ? 1 : 0) : -1; return imageSortAsc ? va - vb : vb - va;
      case 'provenance': va = a.provenance && a.provenance.hasProvenance ? 1 : 0; vb = b.provenance && b.provenance.hasProvenance ? 1 : 0; return imageSortAsc ? va - vb : vb - va;
      case 'sbom': va = a.sbom && a.sbom.hasSBOM ? 1 : 0; vb = b.sbom && b.sbom.hasSBOM ? 1 : 0; return imageSortAsc ? va - vb : vb - va;
      case 'update': va = a.update && a.update.updateAvailable ? 1 : 0; vb = b.update && b.update.updateAvailable ? 1 : 0; return imageSortAsc ? va - vb : vb - va;
      default: return 0;
    }
    if (typeof va === 'string') return imageSortAsc ? va.localeCompare(vb) : vb.localeCompare(va);
    return 0;
  });
  return sorted;
}

function onSortClick(col) {
  if (imageSortCol === col) { imageSortAsc = !imageSortAsc; } else { imageSortCol = col; imageSortAsc = true; }
  renderImages(currentReport);
}
function sortArrow(col) { return imageSortCol !== col ? '' : '<span class="sort-arrow">' + (imageSortAsc ? ' &#9650;' : ' &#9660;') + '</span>'; }

function renderImages(r) {
  if (!r.images || r.images.length === 0) {
    document.getElementById('images-table').innerHTML = '<div class="empty"><p>No images discovered</p></div>';
    document.getElementById('image-count').textContent = '';
    document.getElementById('image-result-count').innerHTML = '';
    document.getElementById('images-pagination').innerHTML = '';
    return;
  }

  const all = sortImages(filterImages(r.images));
  lastFilteredImages = all;
  const total = r.images.length;
  const totalPages = Math.ceil(all.length / pageSize);
  if (currentPage >= totalPages) currentPage = Math.max(0, totalPages - 1);
  const start = currentPage * pageSize;
  const page = all.slice(start, start + pageSize);

  document.getElementById('image-count').textContent = total + ' total';
  const hasFilters = Object.values(imageFilters).some(v => v);
  document.getElementById('image-result-count').innerHTML = hasFilters ? '<div class="result-count">' + all.length + ' of ' + total + ' match</div>' : '';

  if (page.length === 0) {
    document.getElementById('images-table').innerHTML = '<div class="empty"><p>No images match filters</p></div>';
    document.getElementById('images-pagination').innerHTML = '';
    return;
  }

  let html = '<table><thead><tr>' +
    '<th onclick="onSortClick(\'image\')">Image' + sortArrow('image') + '</th>' +
    '<th onclick="onSortClick(\'namespace\')">NS' + sortArrow('namespace') + '</th>' +
    '<th onclick="onSortClick(\'workload\')">Workload' + sortArrow('workload') + '</th>' +
    '<th onclick="onSortClick(\'signature\')">Sig' + sortArrow('signature') + '</th>' +
    '<th onclick="onSortClick(\'provenance\')">SLSA' + sortArrow('provenance') + '</th>' +
    '<th onclick="onSortClick(\'sbom\')">SBOM' + sortArrow('sbom') + '</th>' +
    '<th onclick="onSortClick(\'update\')">Update' + sortArrow('update') + '</th>' +
    '</tr></thead><tbody>';

  for (let idx = 0; idx < page.length; idx++) {
    const img = page[idx];
    const globalIdx = start + idx;
    const sig = img.signature ? (img.signature.verified ? badge('Verified', 'green') : img.signature.signed ? badge('Signed', 'yellow') : badge('Unsigned', 'red')) : badge('N/A', 'muted');
    const prov = img.provenance && img.provenance.hasProvenance ? badge('SLSA', 'purple') : badge('None', 'muted');
    const sbom = img.sbom && img.sbom.hasSBOM ? badge(img.sbom.format.toUpperCase(), 'green') : badge('None', 'muted');
    const update = img.update && img.update.updateAvailable ? badge(img.update.latestInMajor || img.update.newestAvailable, 'yellow') : badge('Current', 'green');
    const digest = img.digest ? '<span class="text-muted" style="font-family:var(--mono);font-size:10px">' + img.digest.substring(7, 19) + '</span>' : '';

    html += '<tr class="clickable" onclick="openDetailIdx(' + globalIdx + ')">' +
      '<td><span class="mono">' + esc(img.image) + '</span><br>' + digest + '</td>' +
      '<td>' + esc(img.namespace) + '</td>' +
      '<td class="workload" title="' + esc(img.workload.kind + '/' + img.workload.name) + '">' + esc(img.workload.kind) + '/' + esc(img.workload.name) + '</td>' +
      '<td>' + sig + '</td><td>' + prov + '</td><td>' + sbom + '</td><td>' + update + '</td></tr>';
  }
  html += '</tbody></table>';
  document.getElementById('images-table').innerHTML = html;

  // Pagination
  const pEl = document.getElementById('images-pagination');
  if (all.length <= pageSize) { pEl.innerHTML = ''; return; }
  pEl.innerHTML =
    '<div class="page-info">' + (start + 1) + '-' + Math.min(start + pageSize, all.length) + ' of ' + all.length + '</div>' +
    '<div class="page-controls">' +
    '<button onclick="goPage(0)"' + (currentPage === 0 ? ' disabled' : '') + '>&laquo;</button>' +
    '<button onclick="goPage(' + (currentPage - 1) + ')"' + (currentPage === 0 ? ' disabled' : '') + '>&lsaquo;</button>' +
    '<button onclick="goPage(' + (currentPage + 1) + ')"' + (currentPage >= totalPages - 1 ? ' disabled' : '') + '>&rsaquo;</button>' +
    '<button onclick="goPage(' + (totalPages - 1) + ')"' + (currentPage >= totalPages - 1 ? ' disabled' : '') + '>&raquo;</button>' +
    '<select onchange="changePageSize(this.value)">' +
    [25,50,100].map(n => '<option value="' + n + '"' + (pageSize === n ? ' selected' : '') + '>' + n + '/page</option>').join('') +
    '</select></div>';
}

function goPage(p) { currentPage = p; renderImages(currentReport); }
function changePageSize(n) { pageSize = parseInt(n); currentPage = 0; renderImages(currentReport); }

function renderHelm(r) {
  if (!r.helmReleases || r.helmReleases.length === 0) {
    document.getElementById('helm-table').innerHTML = '<div class="empty"><p>No Helm releases</p></div>';
    document.getElementById('helm-count').textContent = ''; return;
  }
  document.getElementById('helm-count').textContent = r.helmReleases.length + ' releases';
  let html = '<table><thead><tr><th>Release</th><th>NS</th><th>Chart</th><th>Version</th><th>App</th><th>Status</th></tr></thead><tbody>';
  for (const hr of r.helmReleases) {
    html += '<tr><td>' + esc(hr.releaseName) + '</td><td>' + esc(hr.namespace) + '</td><td class="mono">' + esc(hr.chart) + '</td><td class="mono">' + esc(hr.version) + '</td><td class="mono">' + esc(hr.appVersion) + '</td><td>' + (hr.status === 'deployed' ? badge('Deployed', 'green') : badge(hr.status, 'yellow')) + '</td></tr>';
  }
  html += '</tbody></table>';
  document.getElementById('helm-table').innerHTML = html;
}

function openDetailIdx(idx) {
  openDetail(lastFilteredImages[idx]);
}

function openDetail(img) {
  let html = '<div class="detail-header"><h3>' + esc(img.image) + '</h3>';
  if (img.digest) html += '<div class="detail-digest">' + esc(img.digest) + '</div>';
  html += '</div>';

  // Workload
  html += '<div class="detail-section"><h4>Workload</h4>' +
    detailRow('Kind', img.workload.kind) +
    detailRow('Name', img.workload.name) +
    detailRow('Namespace', img.namespace) +
    '</div>';

  // Signature
  html += '<div class="detail-section"><h4>Signature Verification</h4>';
  if (img.signature) {
    const s = img.signature;
    const color = s.verified ? 'green' : s.signed ? 'yellow' : 'red';
    const label = s.verified ? 'Verified' : s.signed ? 'Signed (unverified)' : 'No signature found';
    const desc = s.verified ? 'Signature exists and has been verified against the configured public key.'
      : s.signed ? 'A cosign signature exists but no public key was configured for verification.'
      : 'No cosign signature was found attached to this image.';
    html += detailStatus(color, label, desc);
    if (s.error) html += detailRow('Error', s.error);
  } else {
    html += detailStatus('muted', 'Not checked', 'Signature verification was disabled for this collection run.');
  }
  html += '</div>';

  // SLSA Provenance
  html += '<div class="detail-section"><h4>SLSA Provenance</h4>';
  if (img.provenance && img.provenance.hasProvenance) {
    html += detailStatus('green', 'Provenance attestation found', 'This image has a SLSA provenance attestation attached via OCI referrers.');
    html += detailRow('Predicate Type', img.provenance.predicateType);
  } else if (img.provenance) {
    html += detailStatus('muted', 'No provenance', 'No SLSA provenance attestation was found in the OCI referrers index.');
  } else {
    html += detailStatus('muted', 'Not checked', 'Provenance detection was disabled for this collection run.');
  }
  html += '</div>';

  // SBOM
  html += '<div class="detail-section"><h4>Software Bill of Materials</h4>';
  if (img.sbom && img.sbom.hasSBOM) {
    html += detailStatus('green', 'SBOM attached', 'An SBOM attestation was found attached to this image.');
    html += detailRow('Format', img.sbom.format.toUpperCase());
  } else if (img.sbom) {
    html += detailStatus('muted', 'No SBOM', 'No SBOM attestation (SPDX or CycloneDX) was found attached to this image.');
  } else {
    html += detailStatus('muted', 'Not checked', 'SBOM detection was disabled for this collection run.');
  }
  html += '</div>';

  // Update
  html += '<div class="detail-section"><h4>Version Status</h4>';
  if (img.update) {
    const u = img.update;
    if (u.updateAvailable) {
      html += detailStatus('yellow', 'Update available', 'A newer version exists in the registry.');
      html += detailRow('Current Tag', u.currentTag);
      if (u.latestInMajor) html += detailRow('Latest (same major)', u.latestInMajor);
      if (u.newestAvailable) html += detailRow('Newest available', u.newestAvailable);
    } else {
      html += detailStatus('green', 'Up to date', 'This image is running the latest available version (at the configured update level).');
      html += detailRow('Current Tag', u.currentTag);
    }
  } else {
    html += detailStatus('muted', 'Not checked', 'Update checking was disabled for this collection run.');
  }
  html += '</div>';

  document.getElementById('detail-content').innerHTML = html;
  document.getElementById('detail-overlay').classList.add('open');
  document.getElementById('detail-panel').classList.add('open');
}

function closeDetail() {
  document.getElementById('detail-overlay').classList.remove('open');
  document.getElementById('detail-panel').classList.remove('open');
}

function detailRow(label, value) {
  return '<div class="detail-row"><span class="label">' + esc(label) + '</span><span class="value">' + esc(value || '-') + '</span></div>';
}

function detailStatus(color, label, description) {
  return '<div class="detail-status"><div class="indicator ' + color + '"></div><div class="detail-status-text"><div>' + esc(label) + '</div><div class="sub">' + esc(description) + '</div></div></div>';
}

// Close detail on Escape key
document.addEventListener('keydown', function(e) { if (e.key === 'Escape') closeDetail(); });

function badge(t, c) { return '<span class="badge badge-' + c + '">' + esc(t) + '</span>'; }
function esc(s) { const d = document.createElement('div'); d.textContent = s || ''; return d.innerHTML; }

init();
</script>
</body>
</html>
` + ""
