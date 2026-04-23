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
<link href="https://fonts.googleapis.com/css2?family=DM+Sans:wght@300;400;500;600;700&family=DM+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
  :root {
    --bg: #0d0d14;
    --surface: #17171f;
    --surface-2: #1e1e28;
    --border: #26262f;
    --border-subtle: #1e1e28;
    --text: #f0f0f8;
    --text-secondary: #c8c8d8;
    --muted: #8888a0;
    --faint: #555568;
    --purple: #a78bfa;
    --purple-dark: #7c3aed;
    --purple-bg: rgba(124,58,237,0.15);
    --green: #10b981;
    --yellow: #f59e0b;
    --red: #ef4444;
    --info: #0b70e0;
    --font: 'DM Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    --mono: 'DM Mono', 'SF Mono', 'Consolas', monospace;
    --radius: 12px;
    --radius-sm: 8px;
  }
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { font-family: var(--font); background: var(--bg); color: var(--text); line-height: 1.6; -webkit-font-smoothing: antialiased; font-size: 14px; }

  /* Nav */
  .nav { position: sticky; top: 0; z-index: 100; background: rgba(13,13,20,0.88); backdrop-filter: blur(12px); border-bottom: 1px solid var(--border); }
  .nav-inner { max-width: 1280px; margin: 0 auto; padding: 0 24px; display: flex; align-items: center; justify-content: space-between; height: 52px; }
  .nav-brand { display: flex; align-items: center; gap: 10px; font-weight: 600; font-size: 14px; color: var(--text); }
  .nav-brand img { height: 22px; }
  .nav-brand .sep { color: var(--faint); font-weight: 300; margin: 0 2px; }
  .nav-meta { font-size: 12px; color: var(--muted); display: flex; gap: 16px; align-items: center; }

  .container { max-width: 1280px; margin: 0 auto; padding: 20px 24px; }

  /* Stats */
  .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(120px, 1fr)); gap: 10px; margin-bottom: 20px; }
  .stat { background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius); padding: 14px 16px; cursor: pointer; transition: all 0.15s ease; }
  .stat:hover { border-color: var(--purple-dark); }
  .stat.active { border-color: var(--purple-dark); background: var(--purple-bg); }
  .stat .value { font-size: 22px; font-weight: 700; letter-spacing: -0.02em; }
  .stat .label { font-size: 10px; color: var(--muted); margin-top: 1px; text-transform: uppercase; letter-spacing: 0.06em; font-weight: 500; }
  .stat.green .value { color: var(--green); }
  .stat.yellow .value { color: var(--yellow); }
  .stat.red .value { color: var(--red); }

  /* Panels */
  .panel { background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius); margin-bottom: 16px; overflow: hidden; }
  .panel-header { padding: 12px 20px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
  .panel-header h2 { font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: var(--muted); }
  .panel-body { padding: 0; }

  /* Filters */
  .filters { padding: 10px 20px; border-bottom: 1px solid var(--border); display: flex; gap: 8px; flex-wrap: wrap; align-items: center; }
  .filters input[type="text"] {
    background: var(--surface-2); border: 1px solid var(--border); border-radius: var(--radius-sm);
    padding: 5px 10px; color: var(--text); font-size: 12px; font-family: var(--font);
    min-width: 180px; outline: none; transition: all 0.15s ease;
  }
  .filters input[type="text"]::placeholder { color: var(--faint); }
  .filters input[type="text"]:focus { border-color: var(--purple-dark); box-shadow: 0 0 0 2px rgba(124,58,237,0.15); }
  .filters select {
    background: var(--surface-2); border: 1px solid var(--border); border-radius: var(--radius-sm);
    padding: 5px 8px; color: var(--text); font-size: 12px; font-family: var(--font);
    outline: none; cursor: pointer; transition: all 0.15s ease;
  }
  .filters select:focus { border-color: var(--purple-dark); box-shadow: 0 0 0 2px rgba(124,58,237,0.15); }
  .filter-label { font-size: 10px; color: var(--faint); text-transform: uppercase; letter-spacing: 0.04em; font-weight: 500; }
  .filter-reset { background: none; border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 4px 10px; color: var(--muted); font-size: 11px; font-family: var(--font); cursor: pointer; margin-left: auto; transition: all 0.15s ease; }
  .filter-reset:hover { border-color: var(--purple-dark); color: var(--text); }

  /* Tables */
  table { width: 100%; border-collapse: collapse; font-size: 13px; }
  th { text-align: left; padding: 8px 20px; color: var(--faint); font-weight: 500; font-size: 10px; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid var(--border); cursor: pointer; user-select: none; transition: color 0.15s; }
  th:hover { color: var(--muted); }
  th .sort-arrow { font-size: 9px; margin-left: 3px; }
  td { padding: 8px 20px; border-bottom: 1px solid var(--border); font-size: 12px; }
  tr:last-child td { border-bottom: none; }
  tr:hover { background: rgba(124,58,237,0.04); }

  /* Badges */
  .badge { display: inline-block; padding: 1px 7px; border-radius: 4px; font-size: 10px; font-weight: 500; letter-spacing: 0.02em; }
  .badge-green { background: rgba(16,185,129,0.12); color: var(--green); }
  .badge-yellow { background: rgba(245,158,11,0.12); color: var(--yellow); }
  .badge-red { background: rgba(239,68,68,0.12); color: var(--red); }
  .badge-muted { background: rgba(136,136,160,0.08); color: var(--faint); }
  .badge-purple { background: var(--purple-bg); color: var(--purple); }

  /* Timeline */
  .timeline { display: flex; gap: 8px; overflow-x: auto; padding: 12px 20px; }
  .timeline-item { min-width: 120px; padding: 10px; border: 1px solid var(--border); border-radius: var(--radius-sm); cursor: pointer; transition: all 0.15s ease; flex-shrink: 0; }
  .timeline-item:hover { border-color: var(--purple-dark); }
  .timeline-item.active { border-color: var(--purple-dark); background: var(--purple-bg); }
  .timeline-item .date { font-size: 12px; font-weight: 600; }
  .timeline-item .time { font-size: 11px; color: var(--muted); }
  .timeline-item .count { font-size: 10px; color: var(--faint); margin-top: 2px; }

  /* Export buttons */
  .export-group { display: flex; gap: 4px; }
  .export-btn { background: var(--surface-2); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 4px 10px; color: var(--muted); font-size: 11px; font-family: var(--font); cursor: pointer; transition: all 0.15s ease; text-decoration: none; display: inline-flex; align-items: center; gap: 4px; }
  .export-btn:hover { border-color: var(--purple-dark); color: var(--text); }

  /* Pagination */
  .pagination { display: flex; align-items: center; justify-content: space-between; padding: 8px 20px; border-top: 1px solid var(--border); font-size: 11px; color: var(--muted); }
  .pagination .page-info { }
  .pagination .page-controls { display: flex; gap: 4px; }
  .pagination button { background: var(--surface-2); border: 1px solid var(--border); border-radius: 4px; padding: 3px 10px; color: var(--text); font-size: 11px; font-family: var(--font); cursor: pointer; transition: all 0.15s ease; }
  .pagination button:hover:not(:disabled) { border-color: var(--purple-dark); }
  .pagination button:disabled { opacity: 0.35; cursor: default; }
  .pagination select { background: var(--surface-2); border: 1px solid var(--border); border-radius: 4px; padding: 3px 6px; color: var(--text); font-size: 11px; font-family: var(--font); }

  /* Detail panel */
  .detail-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 200; opacity: 0; pointer-events: none; transition: opacity 0.2s; }
  .detail-overlay.open { opacity: 1; pointer-events: auto; }
  .detail-panel { position: fixed; top: 0; right: 0; bottom: 0; width: 520px; max-width: 100vw; background: var(--surface); border-left: 1px solid var(--border); z-index: 201; transform: translateX(100%); transition: transform 0.25s ease; overflow-y: auto; }
  .detail-panel.open { transform: translateX(0); }
  .detail-close { position: absolute; top: 12px; right: 16px; background: none; border: none; color: var(--muted); font-size: 18px; cursor: pointer; padding: 4px 8px; border-radius: 4px; transition: all 0.15s; }
  .detail-close:hover { color: var(--text); background: var(--surface-2); }
  .detail-header { padding: 20px 24px 16px; border-bottom: 1px solid var(--border); }
  .detail-header h3 { font-size: 13px; font-weight: 600; word-break: break-all; font-family: var(--mono); line-height: 1.4; }
  .detail-header .detail-digest { font-size: 11px; color: var(--faint); font-family: var(--mono); margin-top: 4px; word-break: break-all; }
  .detail-section { padding: 16px 24px; border-bottom: 1px solid var(--border); }
  .detail-section:last-child { border-bottom: none; }
  .detail-section h4 { font-size: 10px; text-transform: uppercase; letter-spacing: 0.06em; color: var(--faint); font-weight: 600; margin-bottom: 10px; }
  .detail-row { display: flex; justify-content: space-between; align-items: center; padding: 4px 0; font-size: 12px; }
  .detail-row .label { color: var(--muted); }
  .detail-row .value { color: var(--text); font-family: var(--mono); font-size: 11px; text-align: right; max-width: 300px; word-break: break-all; }
  .detail-status { display: flex; align-items: center; gap: 8px; padding: 10px 0; }
  .detail-status .indicator { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
  .detail-status .indicator.green { background: var(--green); box-shadow: 0 0 6px rgba(16,185,129,0.4); }
  .detail-status .indicator.yellow { background: var(--yellow); box-shadow: 0 0 6px rgba(245,158,11,0.4); }
  .detail-status .indicator.red { background: var(--red); box-shadow: 0 0 6px rgba(239,68,68,0.4); }
  .detail-status .indicator.muted { background: var(--faint); }
  .detail-status-text { font-size: 12px; }
  .detail-status-text .sub { font-size: 11px; color: var(--faint); margin-top: 1px; }

  /* Misc */
  .empty { text-align: center; padding: 40px 20px; color: var(--muted); }
  .empty p { margin-top: 4px; font-size: 12px; }
  .loading { text-align: center; padding: 40px; color: var(--muted); font-size: 12px; }
  .mono { font-family: var(--mono); font-size: 11px; }
  .text-muted { color: var(--muted); }
  .result-count { font-size: 11px; color: var(--faint); padding: 6px 20px; border-bottom: 1px solid var(--border); }
  tr.clickable { cursor: pointer; }
</style>
</head>
<body>
<nav class="nav">
  <div class="nav-inner">
    <div class="nav-brand">
      <img src="https://raw.githubusercontent.com/nebari-dev/nebari-design/main/logo-mark/horizontal/Nebari-Logo-Horizontal-Lockup-White-text.svg" alt="Nebari"
           onerror="this.style.display='none'" style="height:20px">
      <span class="sep">/</span>
      <span>Provenance</span>
    </div>
    <div class="nav-meta">
      <span id="cluster-name"></span>
      <span id="last-updated"></span>
      <div class="export-group">
        <a class="export-btn" href="/api/export?format=csv" download>CSV</a>
        <a class="export-btn" href="/api/export?format=markdown" download>Markdown</a>
        <a class="export-btn" href="/api/reports/latest" download="provenance-report.json">JSON</a>
      </div>
    </div>
  </div>
</nav>

<div class="container">
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

async function init() {
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
    return '<div class="timeline-item' + (i === 0 ? ' active' : '') + '" onclick="selectReport(' + i + ')" id="tl-' + i + '">' +
      '<div class="date">' + d.toLocaleDateString() + '</div>' +
      '<div class="time">' + d.toLocaleTimeString() + '</div>' +
      '<div class="count">' + r.summary.totalImages + ' images</div></div>';
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
  const d = new Date(currentReport.metadata.generatedAt);
  document.getElementById('last-updated').textContent = d.toLocaleString();
  document.getElementById('cluster-name').textContent = currentReport.metadata.clusterName || '';
}

function resetFilters() {
  imageFilters = { search: '', namespace: '', signature: '', sbom: '', provenance: '', update: '' };
  statFilter = '';
}

function renderStats(r) {
  const s = r.summary;
  const pct = (n) => s.uniqueImages ? Math.round(n / s.uniqueImages * 100) : 0;
  const pS = pct(s.signedImages), pB = pct(s.imagesWithSBOM), pP = pct(s.imagesWithProvenance || 0);
  document.getElementById('stats').innerHTML =
    statCard('all', s.uniqueImages, 'Images', '') +
    statCard('signed', s.signedImages, 'Signed ' + pS + '%', pS > 80 ? 'green' : pS > 50 ? 'yellow' : 'red') +
    statCard('verified', s.verifiedImages, 'Verified', s.verifiedImages > 0 ? 'green' : '') +
    statCard('provenance', s.imagesWithProvenance || 0, 'SLSA ' + pP + '%', pP > 0 ? 'green' : '') +
    statCard('sbom', s.imagesWithSBOM, 'SBOM ' + pB + '%', pB > 50 ? 'green' : pB > 0 ? 'yellow' : '') +
    statCard('updates', s.imagesWithUpdates, 'Updates', s.imagesWithUpdates > 0 ? 'yellow' : 'green') +
    statCard('helm', s.totalHelmReleases, 'Helm', '');
}

function statCard(id, value, label, cls) {
  return '<div class="stat ' + cls + (statFilter === id ? ' active' : '') + '" onclick="toggleStatFilter(\'' + id + '\')" id="stat-' + id + '"><div class="value">' + value + '</div><div class="label">' + label + '</div></div>';
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
      '<td style="font-size:11px">' + esc(img.workload.kind) + '/' + esc(img.workload.name) + '</td>' +
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
