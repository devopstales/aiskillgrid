#!/usr/bin/env node
/**
 * PRD Kanban — single-file local server for `.skillgrid/prd/PRD*.md` (YAML `status`).
 * No npm in repo — Node 18+ only. UI loads marked + DOMPurify from jsDelivr for markdown preview (offline: escaped plain text).
 * Front matter: expects classic `---` … `---` blocks with simple `title:` / `status:` lines.
 *
 * Run:  node prd-kanban.mjs
 * Open: http://127.0.0.1:8787
 *
 * Env: PRD_KANBAN_PORT, PRD_KANBAN_HOST, PRD_KANBAN_DIR  |  Flags: --port --host --prd-dir
 */
import http from 'node:http';
import fs from 'node:fs/promises';
import { existsSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import crypto from 'node:crypto';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const STATUSES = new Set(['draft', 'todo', 'inprogress', 'devdone', 'done']);
const PRD_NAME_RE = /^PRD\d{2}_.+\.md$/;

const STATIC = {
  '/index.html': `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>PRD Kanban</title>
  <link rel="stylesheet" href="/styles.css" />
</head>
<body>
  <header class="top">
    <h1>PRD Kanban</h1>
    <p class="hint">Local only — cards map to <code>.skillgrid/prd/PRD*.md</code> <code>status</code> in frontmatter.</p>
    <p id="banner" class="banner" hidden></p>
    <button type="button" id="refresh" class="btn">Refresh</button>
  </header>
  <main id="board" class="board" aria-live="polite"></main>
  <div id="ticket-overlay" class="ticket-overlay" hidden>
    <div class="ticket-backdrop" id="ticket-backdrop"></div>
    <div class="ticket-panel" id="ticket-panel" role="dialog" aria-modal="true" aria-labelledby="ticket-summary" aria-hidden="true">
      <div class="ticket-header">
        <span class="ticket-key" id="ticket-key"></span>
        <span class="ticket-status badge" id="ticket-status-badge"></span>
        <button type="button" class="ticket-close" id="ticket-close" aria-label="Close ticket detail">×</button>
      </div>
      <div class="ticket-body-wrap">
        <h2 class="ticket-summary" id="ticket-summary"></h2>
        <dl class="ticket-meta"><dt>Revision</dt><dd id="ticket-rev"></dd></dl>
        <h3 class="ticket-desc-label">Description</h3>
        <div class="ticket-desc markdown-body" id="ticket-desc"></div>
      </div>
    </div>
  </div>
  <script src="https://cdn.jsdelivr.net/npm/marked@12.0.2/marked.min.js" crossorigin="anonymous"></script>
  <script src="https://cdn.jsdelivr.net/npm/dompurify@3.1.6/dist/purify.min.js" crossorigin="anonymous"></script>
  <script src="/app.js" defer></script>
</body>
</html>
`,
  '/styles.css': `:root {
  font-family: system-ui, sans-serif;
  line-height: 1.4;
  color: #1a1a1a;
  background: #f4f4f5;
}

.top {
  max-width: 120rem;
  margin: 0 auto;
  padding: 0.75rem 1rem;
}

.top h1 {
  font-size: 1.25rem;
  margin: 0 0 0.25rem;
}

.hint {
  margin: 0 0 0.5rem;
  color: #52525b;
  font-size: 0.9rem;
}

.banner {
  margin: 0.5rem 0;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fbbf24;
}

.banner.error {
  background: #fee2e2;
  color: #991b1b;
  border-color: #f87171;
}

.btn {
  font: inherit;
  padding: 0.35rem 0.75rem;
  border-radius: 6px;
  border: 1px solid #d4d4d8;
  background: #fff;
  cursor: pointer;
}

.btn:hover {
  background: #fafafa;
}

.board {
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
  max-width: 120rem;
  margin: 0 auto;
  padding: 0 1rem 2rem;
  overflow-x: auto;
}

.column {
  flex: 1 1 12rem;
  min-width: 11rem;
  background: #e4e4e7;
  border-radius: 8px;
  padding: 0.5rem;
}

.column h2 {
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin: 0 0 0.5rem;
  color: #3f3f46;
}

.column-body {
  min-height: 4rem;
  border-radius: 6px;
  padding: 0.35rem;
  background: #fafafa;
  border: 2px dashed transparent;
}

.column-body.drag-over {
  border-color: #6366f1;
  background: #eef2ff;
}

.card {
  background: #fff;
  border: 1px solid #d4d4d8;
  border-radius: 6px;
  padding: 0.5rem 0.6rem;
  margin-bottom: 0.45rem;
  cursor: grab;
  box-shadow: 0 1px 2px rgb(0 0 0 / 0.05);
}

.card:active {
  cursor: grabbing;
}

.card-title {
  display: inline-block;
  font-weight: 600;
  font-size: 0.9rem;
  margin: 0 0 0.35rem;
  word-break: break-word;
  color: #0052cc;
  text-decoration: none;
  cursor: pointer;
  background: none;
  border: none;
  padding: 0;
  font: inherit;
  text-align: left;
  width: 100%;
}

.card-title:hover,
.card-title:focus-visible {
  text-decoration: underline;
  outline: none;
}

.card-meta {
  font-size: 0.75rem;
  color: #71717a;
  margin: 0 0 0.35rem;
}

.card-move {
  flex: 1;
  min-width: 0;
  font: inherit;
  font-size: 0.8rem;
  padding: 0.25rem;
  border-radius: 4px;
  border: 1px solid #d4d4d8;
  background: #fff;
}

.card-actions {
  display: flex;
  gap: 0.35rem;
  margin-top: 0.35rem;
  align-items: center;
}

.ticket-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 2rem 1rem;
  box-sizing: border-box;
}

.ticket-overlay[hidden] {
  display: none !important;
}

.ticket-backdrop {
  position: absolute;
  inset: 0;
  background: rgb(9 30 66 / 0.54);
}

.ticket-panel {
  position: relative;
  width: 100%;
  max-width: 44rem;
  max-height: calc(100vh - 4rem);
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgb(0 0 0 / 0.2);
  overflow: hidden;
  margin-top: 2vh;
}

.ticket-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
  padding: 0.65rem 1rem;
  background: linear-gradient(135deg, #0052cc 0%, #0747a6 100%);
  color: #fff;
}

.ticket-key {
  font-family: ui-monospace, monospace;
  font-size: 0.8rem;
  font-weight: 600;
  opacity: 0.95;
}

.ticket-status.badge {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  background: rgb(255 255 255 / 0.25);
}

.ticket-close {
  margin-left: auto;
  font: inherit;
  font-size: 1.35rem;
  line-height: 1;
  padding: 0.15rem 0.5rem;
  border: none;
  border-radius: 4px;
  background: rgb(255 255 255 / 0.15);
  color: #fff;
  cursor: pointer;
}

.ticket-close:hover {
  background: rgb(255 255 255 / 0.3);
}

.ticket-body-wrap {
  padding: 1rem 1.25rem;
  overflow: auto;
  flex: 1;
}

.ticket-summary {
  margin: 0 0 0.75rem;
  font-size: 1.15rem;
  font-weight: 600;
}

.ticket-meta {
  margin: 0 0 1rem;
  font-size: 0.8rem;
  color: #5e6c84;
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 0.25rem 0.75rem;
}

.ticket-meta dt {
  margin: 0;
  font-weight: 600;
}

.ticket-meta dd {
  margin: 0;
  font-family: ui-monospace, monospace;
  word-break: break-all;
}

.ticket-desc-label {
  margin: 0 0 0.35rem;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #5e6c84;
}

.ticket-desc {
  margin: 0;
  padding: 0.85rem 1rem;
  background: #f7f8f9;
  border: 1px solid #dfe1e6;
  border-radius: 4px;
  font-size: 0.875rem;
  line-height: 1.55;
  overflow: auto;
  max-height: min(50vh, 28rem);
  color: #172b4d;
}

.ticket-desc.markdown-body :first-child {
  margin-top: 0;
}

.ticket-desc.markdown-body :last-child {
  margin-bottom: 0;
}

.ticket-desc.markdown-body h1,
.ticket-desc.markdown-body h2,
.ticket-desc.markdown-body h3 {
  margin: 1em 0 0.5em;
  line-height: 1.25;
  font-weight: 600;
}

.ticket-desc.markdown-body h1 { font-size: 1.25rem; }
.ticket-desc.markdown-body h2 { font-size: 1.1rem; }
.ticket-desc.markdown-body h3 { font-size: 1rem; }

.ticket-desc.markdown-body p {
  margin: 0 0 0.75em;
}

.ticket-desc.markdown-body ul,
.ticket-desc.markdown-body ol {
  margin: 0 0 0.75em;
  padding-left: 1.35em;
}

.ticket-desc.markdown-body li {
  margin: 0.2em 0;
}

.ticket-desc.markdown-body code {
  font-family: ui-monospace, monospace;
  font-size: 0.88em;
  padding: 0.12em 0.35em;
  background: #ebecf0;
  border-radius: 3px;
}

.ticket-desc.markdown-body pre {
  margin: 0.75em 0;
  padding: 0.65rem 0.75rem;
  background: #fff;
  border: 1px solid #dfe1e6;
  border-radius: 4px;
  overflow: auto;
  font-size: 0.82rem;
  line-height: 1.4;
}

.ticket-desc.markdown-body pre code {
  padding: 0;
  background: none;
  border-radius: 0;
}

.ticket-desc.markdown-body blockquote {
  margin: 0.75em 0;
  padding-left: 0.85em;
  border-left: 3px solid #dfe1e6;
  color: #5e6c84;
}

.ticket-desc.markdown-body a {
  color: #0052cc;
}

.ticket-desc.markdown-body table {
  border-collapse: collapse;
  width: 100%;
  margin: 0.75em 0;
  font-size: 0.9em;
}

.ticket-desc.markdown-body th,
.ticket-desc.markdown-body td {
  border: 1px solid #dfe1e6;
  padding: 0.35em 0.5em;
  text-align: left;
}

.ticket-desc.markdown-body th {
  background: #ebecf0;
}

.ticket-desc .ticket-desc-fallback {
  margin: 0;
  white-space: pre-wrap;
  font-family: ui-monospace, monospace;
  font-size: 0.82rem;
  line-height: 1.45;
  background: transparent;
  border: none;
  padding: 0;
}
`,
  '/app.js': `const STATUSES = ['draft', 'todo', 'inprogress', 'devdone', 'done'];

if (typeof marked !== 'undefined') {
  marked.setOptions({ gfm: true, breaks: true });
}

function escapeHtml(s) {
  return String(s)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

/** Render PRD body as sanitized HTML when marked + DOMPurify load from CDN; else escaped pre. */
function renderMarkdownBody(md) {
  const raw = String(md ?? '');
  if (typeof marked !== 'undefined' && typeof DOMPurify !== 'undefined') {
    try {
      return DOMPurify.sanitize(marked.parse(raw));
    } catch {
      /* fall through */
    }
  }
  return \`<pre class="ticket-desc-fallback">\${escapeHtml(raw)}</pre>\`;
}

const boardEl = document.getElementById('board');
const bannerEl = document.getElementById('banner');
const refreshBtn = document.getElementById('refresh');
const ticketOverlay = document.getElementById('ticket-overlay');
const ticketBackdrop = document.getElementById('ticket-backdrop');
const ticketPanel = document.getElementById('ticket-panel');
const ticketKey = document.getElementById('ticket-key');
const ticketStatusBadge = document.getElementById('ticket-status-badge');
const ticketSummary = document.getElementById('ticket-summary');
const ticketRev = document.getElementById('ticket-rev');
const ticketDesc = document.getElementById('ticket-desc');
const ticketClose = document.getElementById('ticket-close');

let cards = [];

function showBanner(message, isError) {
  bannerEl.textContent = message;
  bannerEl.hidden = !message;
  bannerEl.classList.toggle('error', Boolean(isError));
}

function closeTicket() {
  ticketOverlay.hidden = true;
  ticketPanel.setAttribute('aria-hidden', 'true');
  document.body.style.overflow = '';
  ticketDesc.innerHTML = '';
}

function onTicketEscape(e) {
  if (e.key === 'Escape' && !ticketOverlay.hidden) {
    e.preventDefault();
    closeTicket();
  }
}

async function openTicket(file) {
  showBanner('');
  const res = await fetch(\`/api/prds/\${encodeURIComponent(file)}\`);
  if (!res.ok) {
    showBanner(\`Could not load PRD (\${res.status})\`, true);
    return;
  }
  const d = await res.json();
  ticketKey.textContent = d.file;
  ticketStatusBadge.textContent = d.status;
  ticketSummary.textContent = d.title;
  ticketRev.textContent = d.revision.slice(0, 12) + '…';
  ticketRev.title = d.revision;
  ticketDesc.innerHTML = renderMarkdownBody(d.body);
  ticketOverlay.hidden = false;
  ticketPanel.setAttribute('aria-hidden', 'false');
  document.body.style.overflow = 'hidden';
  ticketClose.focus();
}

ticketClose.addEventListener('click', closeTicket);
ticketBackdrop.addEventListener('click', closeTicket);
document.addEventListener('keydown', onTicketEscape);

async function loadPrds() {
  showBanner('');
  const res = await fetch('/api/prds');
  if (!res.ok) {
    showBanner(\`Failed to load PRDs (\${res.status})\`, true);
    return;
  }
  cards = await res.json();
  render();
}

function render() {
  boardEl.innerHTML = '';
  for (const status of STATUSES) {
    const col = document.createElement('section');
    col.className = 'column';
    col.dataset.status = status;

    const h = document.createElement('h2');
    h.textContent = status;
    col.appendChild(h);

    const body = document.createElement('div');
    body.className = 'column-body';
    body.dataset.status = status;

    body.addEventListener('dragover', (e) => {
      e.preventDefault();
      body.classList.add('drag-over');
    });
    body.addEventListener('dragleave', () => body.classList.remove('drag-over'));
    body.addEventListener('drop', async (e) => {
      e.preventDefault();
      body.classList.remove('drag-over');
      const file = e.dataTransfer.getData('text/prd-file');
      if (!file) return;
      await patchStatus(file, status);
    });

    const inCol = cards.filter((c) => c.status === status);
    for (const c of inCol) {
      body.appendChild(renderCard(c));
    }
    col.appendChild(body);
    boardEl.appendChild(col);
  }
}

function renderCard(c) {
  const el = document.createElement('article');
  el.className = 'card';
  el.draggable = true;
  el.dataset.file = c.file;
  el.addEventListener('dragstart', (e) => {
    e.dataTransfer.setData('text/prd-file', c.file);
    e.dataTransfer.effectAllowed = 'move';
  });

  const title = document.createElement('button');
  title.type = 'button';
  title.className = 'card-title';
  title.textContent = c.title;
  title.setAttribute('aria-label', \`Open PRD \${c.file}\`);
  title.addEventListener('click', (e) => {
    e.stopPropagation();
    openTicket(c.file);
  });
  el.appendChild(title);

  const meta = document.createElement('p');
  meta.className = 'card-meta';
  meta.textContent = c.file;
  el.appendChild(meta);

  const actions = document.createElement('div');
  actions.className = 'card-actions';

  const sel = document.createElement('select');
  sel.className = 'card-move';
  sel.setAttribute('aria-label', \`Move \${c.file}\`);
  for (const s of STATUSES) {
    const opt = document.createElement('option');
    opt.value = s;
    opt.textContent = s;
    if (s === c.status) opt.selected = true;
    sel.appendChild(opt);
  }
  sel.addEventListener('change', async () => {
    const next = sel.value;
    await patchStatus(c.file, next);
  });
  actions.appendChild(sel);
  el.appendChild(actions);

  return el;
}

async function patchStatus(file, status) {
  showBanner('');
  const card = cards.find((x) => x.file === file);
  if (!card) return;
  const res = await fetch(\`/api/prds/\${encodeURIComponent(file)}\`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status, revision: card.revision }),
  });
  if (res.status === 409) {
    showBanner('File changed on disk (409). Refreshing…', true);
    await loadPrds();
    return;
  }
  if (!res.ok) {
    let detail = '';
    try {
      const j = await res.json();
      detail = j.error ? \`: \${j.error}\` : '';
    } catch {
      /* ignore */
    }
    showBanner(\`Update failed (\${res.status})\${detail}\`, true);
    await loadPrds();
    return;
  }
  await loadPrds();
}

refreshBtn.addEventListener('click', () => loadPrds());

loadPrds().catch((err) => {
  showBanner(err.message || String(err), true);
});
`,
};

/** Walk up from startDir until a `.skillgrid/prd` directory exists. */
function findPrdDirUpwards(startDir) {
  let dir = path.resolve(startDir);
  for (;;) {
    const candidate = path.join(dir, '.skillgrid', 'prd');
    if (existsSync(candidate)) return candidate;
    const parent = path.dirname(dir);
    if (parent === dir) return null;
    dir = parent;
  }
}

/** Default when this file lives in `.skillgrid/scripts/`. */
function defaultPrdDir() {
  const sibling = path.join(__dirname, '..', 'prd');
  if (existsSync(sibling)) return path.resolve(sibling);
  return (
    findPrdDirUpwards(__dirname) ||
    path.resolve(process.cwd(), '.skillgrid', 'prd')
  );
}

function parseArgs(argv) {
  let port = Number(process.env.PRD_KANBAN_PORT) || 8787;
  let host = process.env.PRD_KANBAN_HOST || '127.0.0.1';
  let prdDirOpt = process.env.PRD_KANBAN_DIR || null;

  for (let i = 2; i < argv.length; i++) {
    const a = argv[i];
    if (a === '--port' && argv[i + 1]) {
      port = Number(argv[++i]);
    } else if (a === '--host' && argv[i + 1]) {
      host = argv[++i];
    } else if (a === '--prd-dir' && argv[i + 1]) {
      prdDirOpt = argv[++i];
    } else if (a === '--help' || a === '-h') {
      console.log(`Usage: node prd-kanban.mjs [--port N] [--host addr] [--prd-dir path]

Default PRD dir: ../prd next to this script (Skillgrid layout), else walk up for .skillgrid/prd, else <cwd>/.skillgrid/prd
Explicit --prd-dir / PRD_KANBAN_DIR: absolute or relative to cwd

Requires: Node.js 18+ (no repo packages; browser uses CDN for markdown preview)
Env: PRD_KANBAN_PORT, PRD_KANBAN_HOST, PRD_KANBAN_DIR`);
      process.exit(0);
    }
  }
  if (!Number.isFinite(port) || port < 1 || port > 65535) {
    console.error('Invalid --port');
    process.exit(1);
  }

  let prdDir;
  if (prdDirOpt != null && String(prdDirOpt).trim() !== '') {
    prdDir = path.isAbsolute(prdDirOpt)
      ? prdDirOpt
      : path.resolve(process.cwd(), prdDirOpt);
  } else {
    prdDir = defaultPrdDir();
  }
  return { port, host, prdDir };
}

function sha256Hex(buf) {
  return crypto.createHash('sha256').update(buf).digest('hex');
}

function safeResolvedPath(prdRoot, fileParam) {
  const name = path.basename(decodeURIComponent(fileParam || ''));
  if (!PRD_NAME_RE.test(name)) return null;
  const resolved = path.resolve(prdRoot, name);
  const root = path.resolve(prdRoot);
  if (path.dirname(resolved) !== root) return null;
  return resolved;
}

function firstH1Title(body) {
  const m = String(body).match(/^\s*#\s+(.+)$/m);
  return m ? m[1].trim() : null;
}

/** Split `---` / `---` front matter from markdown body (simple PRD-style files). */
const FRONT_MATTER_RE = /^---\r?\n([\s\S]*?)\r?\n---\r?\n([\s\S]*)$/;

function splitFrontMatter(rawStr) {
  const raw = String(rawStr);
  const m = raw.match(FRONT_MATTER_RE);
  if (!m) return { front: null, body: raw };
  return { front: m[1], body: m[2] };
}

function yamlScalarValue(raw) {
  let v = String(raw ?? '').trim();
  if ((v.startsWith('"') && v.endsWith('"')) || (v.startsWith("'") && v.endsWith("'")))
    v = v.slice(1, -1);
  return v;
}

function getYamlField(front, key) {
  if (front == null) return null;
  const re = new RegExp(`^${key}:\\s*(.*)$`, 'm');
  const mm = String(front).match(re);
  if (!mm) return null;
  return yamlScalarValue(mm[1]);
}

/** Emit a scalar suitable for a single-line YAML value. */
function yamlScalarLine(value) {
  const v = String(value);
  if (v === '' || /^[\w.-]+$/.test(v)) return v;
  return JSON.stringify(v);
}

function setYamlField(front, key, value) {
  const line = `${key}: ${yamlScalarLine(value)}`;
  if (front == null || !String(front).trim()) return line;
  const lines = String(front).split(/\r?\n/);
  const keyRe = new RegExp(`^${key}:\\s*.*$`);
  let found = false;
  const out = lines.map((ln) => {
    if (keyRe.test(ln)) {
      found = true;
      return line;
    }
    return ln;
  });
  if (!found) out.push(line);
  return out.join('\n');
}

function composeWithFrontMatter(front, body) {
  return `---\n${front}\n---\n${body}`;
}

function effectiveStatusString(statusRaw) {
  const s = statusRaw != null ? String(statusRaw).trim() : '';
  if (s && STATUSES.has(s)) return s;
  return 'draft';
}

/** Update only `status:` in front matter; preserve body bytes as string content after split. */
function applyStatusToDocument(rawBuf, newStatus) {
  const rawStr = rawBuf.toString('utf8');
  const { front, body } = splitFrontMatter(rawStr);
  if (front === null) {
    return composeWithFrontMatter(setYamlField(null, 'status', newStatus), body);
  }
  return composeWithFrontMatter(setYamlField(front, 'status', newStatus), body);
}

function extractPrdFromRaw(rawBuf, fileName, opts = {}) {
  const includeBody = Boolean(opts.includeBody);
  const rawStr = rawBuf.toString('utf8');
  const revision = sha256Hex(rawBuf);
  const { front, body } = splitFrontMatter(rawStr);
  const bodyOut = body;
  const st = getYamlField(front, 'status');
  const status = effectiveStatusString(st);
  let title = fileName;
  const t = getYamlField(front, 'title');
  if (t) title = t;
  else {
    const h = firstH1Title(bodyOut);
    if (h) title = h;
  }
  const out = { file: fileName, title, status, revision };
  if (includeBody) out.body = bodyOut;
  return out;
}

async function getPrdDetail(prdRoot, fileParam) {
  const decoded = decodeURIComponent(fileParam || '');
  if (!decoded || decoded !== path.basename(decoded)) {
    const err = new Error('not_found');
    err.code = 'not_found';
    throw err;
  }
  const full = safeResolvedPath(prdRoot, fileParam);
  if (!full || !existsSync(full)) {
    const err = new Error('not_found');
    err.code = 'not_found';
    throw err;
  }
  const raw = await fs.readFile(full);
  return extractPrdFromRaw(raw, path.basename(full), { includeBody: true });
}

async function listPrds(prdRoot) {
  const entries = await fs.readdir(prdRoot, { withFileTypes: true });
  const out = [];
  for (const ent of entries) {
    if (!ent.isFile()) continue;
    if (!PRD_NAME_RE.test(ent.name)) continue;
    const full = path.join(prdRoot, ent.name);
    const raw = await fs.readFile(full);
    try {
      const meta = extractPrdFromRaw(raw, ent.name);
      out.push({ file: meta.file, title: meta.title, status: meta.status, revision: meta.revision });
    } catch {
      out.push({
        file: ent.name,
        title: ent.name,
        status: 'draft',
        revision: sha256Hex(raw),
      });
    }
  }
  out.sort((a, b) => a.file.localeCompare(b.file));
  return out;
}

async function updateStatus(prdRoot, fileParam, body, currentRevision) {
  const full = safeResolvedPath(prdRoot, fileParam);
  if (!full) {
    const err = new Error('not_found');
    err.code = 'not_found';
    throw err;
  }
  if (!existsSync(full)) {
    const err = new Error('not_found');
    err.code = 'not_found';
    throw err;
  }
  const raw = await fs.readFile(full);
  const nowRev = sha256Hex(raw);
  if (nowRev !== currentRevision) {
    const err = new Error('conflict');
    err.code = 'conflict';
    throw err;
  }
  if (!STATUSES.has(body.status)) {
    const err = new Error('invalid_status');
    err.code = 'invalid_status';
    throw err;
  }
  const next = applyStatusToDocument(raw, body.status);
  const dir = path.dirname(full);
  const tmp = path.join(dir, `.${path.basename(full)}.${process.pid}.${Date.now()}.tmp`);
  await fs.writeFile(tmp, next, 'utf8');
  await fs.rename(tmp, full);
}

function json(res, status, obj) {
  const data = JSON.stringify(obj);
  res.writeHead(status, {
    'Content-Type': 'application/json; charset=utf-8',
    'Content-Length': Buffer.byteLength(data),
  });
  res.end(data);
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    const chunks = [];
    req.on('data', (c) => chunks.push(c));
    req.on('end', () => resolve(Buffer.concat(chunks)));
    req.on('error', reject);
  });
}

const MIME = {
  '.html': 'text/html; charset=utf-8',
  '.js': 'text/javascript; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
};

function serveStatic(res, urlPath) {
  let rel = String(urlPath || '/');
  if (rel.startsWith('/')) rel = rel.slice(1);
  rel = path.normalize(rel).replace(/^(\.\.(\/|\\|$))+/, '');
  if (!rel || rel === '.') rel = 'index.html';
  const key = '/' + rel.replace(/\\/g, '/');
  const body = STATIC[key];
  if (body === undefined) {
    res.writeHead(404);
    res.end('Not found');
    return;
  }
  const ext = path.extname(key);
  res.writeHead(200, { 'Content-Type': MIME[ext] || 'application/octet-stream' });
  res.end(body);
}

async function handle(req, res, prdRoot) {
  const u = new URL(req.url || '/', 'http://127.0.0.1');

  if (u.pathname === '/api/prds' && req.method === 'GET') {
    try {
      const list = await listPrds(prdRoot);
      json(res, 200, list);
    } catch (e) {
      json(res, 500, { error: String(e.message || e) });
    }
    return;
  }

  if (u.pathname.startsWith('/api/prds/') && req.method === 'GET') {
    const filePart = u.pathname.slice('/api/prds/'.length);
    try {
      const detail = await getPrdDetail(prdRoot, filePart);
      json(res, 200, detail);
    } catch (e) {
      if (e.code === 'not_found') json(res, 404, { error: 'not_found' });
      else json(res, 500, { error: String(e.message || e) });
    }
    return;
  }

  if (u.pathname.startsWith('/api/prds/') && req.method === 'PATCH') {
    const filePart = u.pathname.slice('/api/prds/'.length);
    let payload;
    try {
      const raw = await readBody(req);
      payload = JSON.parse(raw.toString('utf8') || '{}');
    } catch {
      json(res, 400, { error: 'invalid_json' });
      return;
    }
    if (typeof payload.status !== 'string' || typeof payload.revision !== 'string') {
      json(res, 400, { error: 'expected status and revision strings' });
      return;
    }
    try {
      await updateStatus(prdRoot, filePart, payload, payload.revision);
      json(res, 200, { ok: true });
    } catch (e) {
      if (e.code === 'not_found') json(res, 404, { error: 'not_found' });
      else if (e.code === 'conflict') json(res, 409, { error: 'conflict' });
      else if (e.code === 'invalid_status') json(res, 400, { error: 'invalid_status' });
      else json(res, 500, { error: String(e.message || e) });
    }
    return;
  }

  if (req.method === 'GET') {
    let p = u.pathname;
    if (p === '/') p = '/index.html';
    serveStatic(res, p);
    return;
  }

  res.writeHead(405);
  res.end('Method not allowed');
}

async function main() {
  const { port, host, prdDir } = parseArgs(process.argv);
  if (!existsSync(prdDir)) {
    console.error(`PRD directory not found: ${prdDir}`);
    console.error('Pass --prd-dir <path> or set PRD_KANBAN_DIR');
    process.exit(1);
  }

  const server = http.createServer((req, res) => {
    handle(req, res, prdDir).catch((e) => {
      console.error(e);
      if (!res.headersSent) {
        res.writeHead(500);
        res.end('Internal error');
      }
    });
  });

  server.listen(port, host, () => {
    console.log(`PRD Kanban serving http://${host}:${port}`);
    console.log(`PRD dir: ${prdDir}`);
  });
}

main();
