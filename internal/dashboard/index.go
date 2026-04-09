package dashboard

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Provenance Collector</title>
<style>
  :root {
    --bg: #0f1117; --surface: #1a1d27; --border: #2a2d3a;
    --text: #e1e4ed; --muted: #8b8fa3; --accent: #6c8cff;
    --green: #4caf87; --yellow: #e5a84b; --red: #e5574b;
  }
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: var(--bg); color: var(--text); line-height: 1.6; }
  .container { max-width: 1200px; margin: 0 auto; padding: 24px; }
  header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 32px; }
  header h1 { font-size: 24px; font-weight: 600; }
  header .cluster { color: var(--muted); font-size: 14px; }

  .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 16px; margin-bottom: 32px; }
  .stat { background: var(--surface); border: 1px solid var(--border); border-radius: 8px; padding: 20px; }
  .stat .value { font-size: 28px; font-weight: 700; }
  .stat .label { font-size: 13px; color: var(--muted); margin-top: 4px; }
  .stat.green .value { color: var(--green); }
  .stat.yellow .value { color: var(--yellow); }
  .stat.red .value { color: var(--red); }

  .panel { background: var(--surface); border: 1px solid var(--border); border-radius: 8px; margin-bottom: 24px; }
  .panel-header { padding: 16px 20px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
  .panel-header h2 { font-size: 16px; font-weight: 600; }
  .panel-body { padding: 0; }

  table { width: 100%; border-collapse: collapse; font-size: 14px; }
  th { text-align: left; padding: 12px 20px; color: var(--muted); font-weight: 500; font-size: 12px; text-transform: uppercase; letter-spacing: 0.5px; border-bottom: 1px solid var(--border); }
  td { padding: 12px 20px; border-bottom: 1px solid var(--border); }
  tr:last-child td { border-bottom: none; }
  tr:hover { background: rgba(108,140,255,0.04); }

  .badge { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 12px; font-weight: 500; }
  .badge-green { background: rgba(76,175,135,0.15); color: var(--green); }
  .badge-yellow { background: rgba(229,168,75,0.15); color: var(--yellow); }
  .badge-red { background: rgba(229,87,75,0.15); color: var(--red); }
  .badge-muted { background: rgba(139,143,163,0.15); color: var(--muted); }

  .timeline { display: flex; gap: 12px; overflow-x: auto; padding: 16px 20px; }
  .timeline-item { min-width: 140px; padding: 12px; border: 1px solid var(--border); border-radius: 8px; cursor: pointer; transition: border-color 0.2s; flex-shrink: 0; }
  .timeline-item:hover, .timeline-item.active { border-color: var(--accent); }
  .timeline-item .date { font-size: 13px; font-weight: 600; }
  .timeline-item .time { font-size: 12px; color: var(--muted); }
  .timeline-item .count { font-size: 12px; color: var(--muted); margin-top: 4px; }

  .empty { text-align: center; padding: 60px 20px; color: var(--muted); }
  .empty p { margin-top: 8px; font-size: 14px; }
  .loading { text-align: center; padding: 60px; color: var(--muted); }

  .mono { font-family: "SF Mono", "Fira Code", monospace; font-size: 13px; }
  .text-muted { color: var(--muted); }
</style>
</head>
<body>
<div class="container">
  <header>
    <div>
      <h1>Provenance Collector</h1>
      <div class="cluster" id="cluster-name"></div>
    </div>
    <div class="text-muted" id="last-updated"></div>
  </header>

  <div class="stats" id="stats">
    <div class="loading">Loading...</div>
  </div>

  <div class="panel">
    <div class="panel-header">
      <h2>Report Timeline</h2>
    </div>
    <div class="timeline" id="timeline"></div>
  </div>

  <div class="panel">
    <div class="panel-header">
      <h2>Container Images</h2>
      <span class="text-muted" id="image-count"></span>
    </div>
    <div class="panel-body" id="images-table"></div>
  </div>

  <div class="panel">
    <div class="panel-header">
      <h2>Helm Releases</h2>
      <span class="text-muted" id="helm-count"></span>
    </div>
    <div class="panel-body" id="helm-table"></div>
  </div>
</div>

<script>
let reports = [];
let currentReport = null;

async function init() {
  try {
    const res = await fetch('/api/reports');
    reports = await res.json();
    if (!reports || reports.length === 0) {
      showEmpty();
      return;
    }
    renderTimeline();
    await loadReport(reports[0].filename);
  } catch (e) {
    document.getElementById('stats').innerHTML = '<div class="empty"><p>Failed to load reports</p></div>';
  }
}

function showEmpty() {
  document.getElementById('stats').innerHTML = '';
  document.getElementById('timeline').innerHTML = '<div class="empty"><p>No reports found. Run the collector to generate your first report.</p></div>';
  document.getElementById('images-table').innerHTML = '';
  document.getElementById('helm-table').innerHTML = '';
}

function renderTimeline() {
  const el = document.getElementById('timeline');
  el.innerHTML = reports.map((r, i) => {
    const d = new Date(r.generatedAt);
    return '<div class="timeline-item' + (i === 0 ? ' active' : '') + '" onclick="selectReport(' + i + ')" id="tl-' + i + '">' +
      '<div class="date">' + d.toLocaleDateString() + '</div>' +
      '<div class="time">' + d.toLocaleTimeString() + '</div>' +
      '<div class="count">' + r.summary.totalImages + ' images</div>' +
    '</div>';
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
  renderStats(currentReport);
  renderImages(currentReport);
  renderHelm(currentReport);

  const d = new Date(currentReport.metadata.generatedAt);
  document.getElementById('last-updated').textContent = 'Report: ' + d.toLocaleString();
  document.getElementById('cluster-name').textContent = currentReport.metadata.clusterName ? 'Cluster: ' + currentReport.metadata.clusterName : '';
}

function renderStats(r) {
  const s = r.summary;
  const pctSigned = s.uniqueImages ? Math.round(s.signedImages / s.uniqueImages * 100) : 0;
  const pctSbom = s.uniqueImages ? Math.round(s.imagesWithSBOM / s.uniqueImages * 100) : 0;

  document.getElementById('stats').innerHTML =
    stat(s.uniqueImages, 'Unique Images', '') +
    stat(s.signedImages, 'Signed (' + pctSigned + '%)', pctSigned > 80 ? 'green' : pctSigned > 50 ? 'yellow' : 'red') +
    stat(s.verifiedImages, 'Verified', s.verifiedImages > 0 ? 'green' : '') +
    stat(s.imagesWithSBOM, 'With SBOM (' + pctSbom + '%)', pctSbom > 50 ? 'green' : pctSbom > 20 ? 'yellow' : 'red') +
    stat(s.imagesWithUpdates, 'Updates Available', s.imagesWithUpdates > 0 ? 'yellow' : 'green') +
    stat(s.totalHelmReleases, 'Helm Releases', '');
}

function stat(value, label, cls) {
  return '<div class="stat ' + cls + '"><div class="value">' + value + '</div><div class="label">' + label + '</div></div>';
}

function renderImages(r) {
  if (!r.images || r.images.length === 0) {
    document.getElementById('images-table').innerHTML = '<div class="empty"><p>No images discovered</p></div>';
    document.getElementById('image-count').textContent = '';
    return;
  }
  document.getElementById('image-count').textContent = r.images.length + ' images';

  let html = '<table><thead><tr><th>Image</th><th>Namespace</th><th>Workload</th><th>Signature</th><th>SBOM</th><th>Update</th></tr></thead><tbody>';
  for (const img of r.images) {
    const sig = img.signature
      ? (img.signature.verified ? badge('Verified', 'green') : img.signature.signed ? badge('Signed', 'yellow') : badge('Unsigned', 'red'))
      : badge('N/A', 'muted');
    const sbom = img.sbom && img.sbom.hasSBOM ? badge(img.sbom.format.toUpperCase(), 'green') : badge('None', 'muted');
    const update = img.update && img.update.updateAvailable ? badge(img.update.latestInMajor || img.update.newestAvailable, 'yellow') : badge('Current', 'green');
    const digest = img.digest ? '<span class="text-muted mono">' + img.digest.substring(7, 19) + '</span>' : '';

    html += '<tr>' +
      '<td><span class="mono">' + esc(img.image) + '</span><br>' + digest + '</td>' +
      '<td>' + esc(img.namespace) + '</td>' +
      '<td>' + esc(img.workload.kind) + '/' + esc(img.workload.name) + '</td>' +
      '<td>' + sig + '</td>' +
      '<td>' + sbom + '</td>' +
      '<td>' + update + '</td>' +
    '</tr>';
  }
  html += '</tbody></table>';
  document.getElementById('images-table').innerHTML = html;
}

function renderHelm(r) {
  if (!r.helmReleases || r.helmReleases.length === 0) {
    document.getElementById('helm-table').innerHTML = '<div class="empty"><p>No Helm releases discovered</p></div>';
    document.getElementById('helm-count').textContent = '';
    return;
  }
  document.getElementById('helm-count').textContent = r.helmReleases.length + ' releases';

  let html = '<table><thead><tr><th>Release</th><th>Namespace</th><th>Chart</th><th>Version</th><th>App Version</th><th>Status</th></tr></thead><tbody>';
  for (const hr of r.helmReleases) {
    const status = hr.status === 'deployed' ? badge('Deployed', 'green') : badge(hr.status, 'yellow');
    html += '<tr>' +
      '<td>' + esc(hr.releaseName) + '</td>' +
      '<td>' + esc(hr.namespace) + '</td>' +
      '<td class="mono">' + esc(hr.chart) + '</td>' +
      '<td class="mono">' + esc(hr.version) + '</td>' +
      '<td class="mono">' + esc(hr.appVersion) + '</td>' +
      '<td>' + status + '</td>' +
    '</tr>';
  }
  html += '</tbody></table>';
  document.getElementById('helm-table').innerHTML = html;
}

function badge(text, cls) { return '<span class="badge badge-' + cls + '">' + esc(text) + '</span>'; }
function esc(s) { const d = document.createElement('div'); d.textContent = s || ''; return d.innerHTML; }

init();
</script>
</body>
</html>
` + ""
