#!/usr/bin/env node
/**
 * Skillgrid Dashboard — single-file local server for PRD Kanban, workflow events, previews, and graphify output.
 * No npm in repo — Node 18+ only. UI loads marked + DOMPurify from jsDelivr for markdown preview (offline: escaped plain text).
 * Front matter: supports classic `---` … `---` blocks with simple `title:` / `status:` lines and Skillgrid title-block metadata.
 *
 * Run:  node skillgrid-ui.mjs
 * Open: http://127.0.0.1:8787
 *
 * Env: SKILLGRID_UI_PORT, SKILLGRID_UI_HOST, SKILLGRID_UI_DIR
 * Legacy env: PRD_KANBAN_PORT, PRD_KANBAN_HOST, PRD_KANBAN_DIR
 * Flags: --port --host --prd-dir
 */
import http from 'node:http';
import fs from 'node:fs/promises';
import { existsSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import crypto from 'node:crypto';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const DEFAULT_WORKFLOW_PHASES = ['brainstorm', 'design', 'plan', 'breakdown', 'apply', 'test', 'security', 'validate', 'finish'];
const DEFAULT_PRD_WORKFLOW = Object.freeze({
  source: 'preset',
  preset: 'skillgrid-default',
  fallbackStatus: 'draft',
  statuses: Object.freeze([
    Object.freeze({ id: 'draft', label: 'Draft' }),
    Object.freeze({ id: 'todo', label: 'Todo' }),
    Object.freeze({ id: 'inprogress', label: 'In Progress' }),
    Object.freeze({ id: 'devdone', label: 'Dev Done' }),
    Object.freeze({ id: 'done', label: 'Done' }),
  ]),
  phaseStatusMap: Object.freeze({
    plan: 'draft',
    breakdown: 'todo',
    apply: 'inprogress',
    validate: 'devdone',
    finish: 'done',
  }),
  providerMapping: Object.freeze({}),
});
const STATUS_ID_RE = /^[A-Za-z0-9][A-Za-z0-9._-]*$/;
const PRD_NAME_RE = /^PRD\d{2}_.+\.md$/;

const STATIC = {
  '/index.html': `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Skillgrid Dashboard</title>
  <link rel="stylesheet" href="/styles.css" />
</head>
<body>
  <header class="top">
    <h1>Skillgrid Dashboard</h1>
    <p class="hint">Local only — PRDs, handoff state, previews, graphify output, and JSONL workflow events from this repo.</p>
    <nav class="tabs" aria-label="Dashboard views">
      <button type="button" class="tab is-active" data-view="board">Board</button>
      <button type="button" class="tab" data-view="workflow">Workflow</button>
      <button type="button" class="tab" data-view="agents">Subagents</button>
      <a id="graph-link" class="tab tab-link" href="/graphify/graph.html" target="_blank" rel="noreferrer" hidden>Graph</a>
    </nav>
    <p id="banner" class="banner" hidden></p>
    <div id="system-status" class="system-status" aria-live="polite"></div>
    <button type="button" id="refresh" class="btn">Refresh</button>
  </header>
  <main class="dashboard-shell">
    <section id="board-view" class="view is-active" aria-live="polite">
      <div id="board" class="board"></div>
    </section>
    <section id="workflow-view" class="view workflow-view" aria-live="polite" hidden>
      <aside class="workflow-list" aria-label="PRDs">
        <h2>Changes</h2>
        <div id="workflow-list"></div>
      </aside>
      <section class="workflow-detail">
        <div id="workflow-detail" class="workflow-empty">Select a PRD/change to inspect workflow events and artifacts.</div>
      </section>
    </section>
    <section id="agents-view" class="view agents-view" aria-live="polite" hidden>
      <div class="agents-header">
        <h2>Subagent Actions</h2>
        <p class="hint">Read-only view over workflow events that include <code>agent</code>, <code>subagent</code>, <code>role</code>, or subagent/review nodes.</p>
      </div>
      <div id="agents-list" class="agents-list"></div>
    </section>
  </main>
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

.system-status {
  display: flex;
  gap: 0.4rem;
  flex-wrap: wrap;
  margin: 0.5rem 0;
}

.status-chip {
  display: inline-flex;
  gap: 0.25rem;
  align-items: center;
  padding: 0.25rem 0.55rem;
  border-radius: 999px;
  border: 1px solid #d4d4d8;
  background: #fff;
  color: #3f3f46;
  font-size: 0.82rem;
}

.status-chip strong {
  color: #18181b;
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

.tabs {
  display: flex;
  gap: 0.4rem;
  flex-wrap: wrap;
  margin: 0.6rem 0;
}

.tab {
  font: inherit;
  padding: 0.35rem 0.75rem;
  border-radius: 999px;
  border: 1px solid #d4d4d8;
  background: #fff;
  color: #27272a;
  cursor: pointer;
  text-decoration: none;
}

.tab:hover,
.tab:focus-visible {
  background: #f8fafc;
}

.tab.is-active {
  border-color: #6366f1;
  background: #eef2ff;
  color: #3730a3;
}

.dashboard-shell {
  max-width: 120rem;
  margin: 0 auto;
  padding: 0 1rem 2rem;
}

.view[hidden] {
  display: none !important;
}

.board {
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
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

.card-badges,
.card-links {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
  margin: 0.35rem 0;
}

.badge,
.pill-link {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  max-width: 100%;
  padding: 0.15rem 0.4rem;
  border-radius: 999px;
  font-size: 0.68rem;
  line-height: 1.2;
  border: 1px solid #d4d4d8;
  background: #f8fafc;
  color: #3f3f46;
  text-decoration: none;
  overflow: hidden;
  text-overflow: ellipsis;
}

.badge.blocked {
  border-color: #fca5a5;
  background: #fee2e2;
  color: #991b1b;
}

.badge.passed,
.badge.completed {
  border-color: #86efac;
  background: #dcfce7;
  color: #166534;
}

.badge.failed {
  border-color: #fca5a5;
  background: #fee2e2;
  color: #991b1b;
}

.pill-link {
  color: #0052cc;
}

.pill-link:hover {
  background: #eef2ff;
  text-decoration: underline;
}

.workflow-view {
  display: grid;
  grid-template-columns: minmax(14rem, 22rem) minmax(0, 1fr);
  gap: 1rem;
  align-items: start;
}

.agents-view {
  display: grid;
  gap: 0.75rem;
}

.agents-header {
  background: #fff;
  border: 1px solid #d4d4d8;
  border-radius: 8px;
  padding: 0.85rem 1rem;
  box-shadow: 0 1px 2px rgb(0 0 0 / 0.05);
}

.agents-header h2 {
  margin: 0 0 0.25rem;
  font-size: 1rem;
}

.agents-list {
  display: grid;
  gap: 0.65rem;
}

.agent-action {
  background: #fff;
  border: 1px solid #d4d4d8;
  border-radius: 8px;
  padding: 0.8rem;
  box-shadow: 0 1px 2px rgb(0 0 0 / 0.05);
}

.agent-action-head {
  display: flex;
  gap: 0.45rem;
  flex-wrap: wrap;
  align-items: center;
  margin-bottom: 0.35rem;
}

.agent-name {
  font-weight: 700;
  color: #18181b;
}

.agent-summary {
  margin: 0.35rem 0;
}

.agent-meta {
  display: flex;
  gap: 0.35rem;
  flex-wrap: wrap;
  color: #71717a;
  font-size: 0.74rem;
}

.workflow-list,
.workflow-detail {
  background: #fff;
  border: 1px solid #d4d4d8;
  border-radius: 8px;
  box-shadow: 0 1px 2px rgb(0 0 0 / 0.05);
}

.workflow-list {
  padding: 0.75rem;
}

.workflow-list h2 {
  margin: 0 0 0.6rem;
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: #52525b;
}

.workflow-item {
  width: 100%;
  text-align: left;
  font: inherit;
  padding: 0.55rem;
  border: 1px solid transparent;
  border-radius: 6px;
  background: transparent;
  cursor: pointer;
}

.workflow-item:hover,
.workflow-item.is-active {
  background: #eef2ff;
  border-color: #c7d2fe;
}

.workflow-item-title {
  display: block;
  font-weight: 600;
  color: #18181b;
}

.workflow-item-meta {
  display: block;
  margin-top: 0.15rem;
  font-size: 0.72rem;
  color: #71717a;
}

.workflow-detail {
  padding: 1rem;
}

.workflow-empty {
  color: #71717a;
}

.workflow-header {
  display: grid;
  gap: 0.4rem;
  margin-bottom: 1rem;
}

.workflow-header h2 {
  margin: 0;
  font-size: 1.15rem;
}

.workflow-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.workflow-card {
  border: 1px solid #e4e4e7;
  border-radius: 6px;
  padding: 0.7rem;
  background: #fafafa;
}

.workflow-card h3 {
  margin: 0 0 0.4rem;
  font-size: 0.8rem;
  color: #52525b;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.phase-list {
  display: grid;
  gap: 0.35rem;
  margin-bottom: 1rem;
}

.phase-row {
  display: grid;
  grid-template-columns: 7rem 7rem minmax(0, 1fr);
  gap: 0.5rem;
  align-items: center;
  padding: 0.45rem 0.55rem;
  border: 1px solid #e4e4e7;
  border-radius: 6px;
  background: #fff;
}

.phase-node {
  font-weight: 600;
  color: #27272a;
}

.phase-summary {
  min-width: 0;
  color: #52525b;
  font-size: 0.78rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.event-list {
  display: grid;
  gap: 0.5rem;
}

.event {
  border-left: 3px solid #c7d2fe;
  padding: 0.45rem 0 0.45rem 0.65rem;
}

.event-time {
  color: #71717a;
  font-size: 0.72rem;
}

.event-summary {
  margin: 0.15rem 0;
  color: #18181b;
}

.artifact-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.3rem;
}

@media (max-width: 760px) {
  .workflow-view {
    grid-template-columns: 1fr;
  }

  .phase-row {
    grid-template-columns: 1fr;
  }
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
  '/app.js': `const DEFAULT_STATUS_COLUMNS = [
  { id: 'draft', label: 'Draft' },
  { id: 'todo', label: 'Todo' },
  { id: 'inprogress', label: 'In Progress' },
  { id: 'devdone', label: 'Dev Done' },
  { id: 'done', label: 'Done' },
];
const DEFAULT_PHASES = ['brainstorm', 'design', 'plan', 'breakdown', 'apply', 'test', 'security', 'validate', 'finish'];

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
const boardView = document.getElementById('board-view');
const workflowView = document.getElementById('workflow-view');
const workflowList = document.getElementById('workflow-list');
const workflowDetail = document.getElementById('workflow-detail');
const agentsView = document.getElementById('agents-view');
const agentsList = document.getElementById('agents-list');
const bannerEl = document.getElementById('banner');
const systemStatusEl = document.getElementById('system-status');
const refreshBtn = document.getElementById('refresh');
const graphLink = document.getElementById('graph-link');
const tabEls = Array.from(document.querySelectorAll('[data-view]'));
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
let dashboard = { prds: [], graph: { available: false } };
let activeView = 'board';
let selectedWorkflowFile = null;
let selectedChangeId = null;

function selectWorkflowCard(card) {
  selectedWorkflowFile = card ? card.file : null;
  selectedChangeId = card ? card.changeId : null;
}

function selectedWorkflowCard() {
  if (selectedWorkflowFile) {
    const byFile = cards.find((c) => c.file === selectedWorkflowFile);
    if (byFile) return byFile;
  }
  if (selectedChangeId) {
    const byChange = cards.find((c) => c.changeId === selectedChangeId);
    if (byChange) return byChange;
  }
  return cards[0] || null;
}

function ensureWorkflowSelection() {
  const card = selectedWorkflowCard();
  selectWorkflowCard(card);
  return card;
}

function statusColumns() {
  const configured = dashboard.prdWorkflow && dashboard.prdWorkflow.statuses;
  return Array.isArray(configured) && configured.length ? configured : DEFAULT_STATUS_COLUMNS;
}

function statusLabel(statusId) {
  const match = statusColumns().find((status) => status.id === statusId);
  return match ? match.label : statusId;
}

function workflowPhases() {
  return Array.isArray(dashboard.workflowPhases) && dashboard.workflowPhases.length
    ? dashboard.workflowPhases
    : DEFAULT_PHASES;
}

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

function setView(view) {
  activeView = view;
  boardView.hidden = view !== 'board';
  workflowView.hidden = view !== 'workflow';
  agentsView.hidden = view !== 'agents';
  boardView.classList.toggle('is-active', view === 'board');
  workflowView.classList.toggle('is-active', view === 'workflow');
  agentsView.classList.toggle('is-active', view === 'agents');
  for (const tab of tabEls) {
    tab.classList.toggle('is-active', tab.dataset.view === view);
  }
  if (view === 'workflow' && cards.length) {
    ensureWorkflowSelection();
  }
  renderWorkflow();
  renderAgents();
}

ticketClose.addEventListener('click', closeTicket);
ticketBackdrop.addEventListener('click', closeTicket);
document.addEventListener('keydown', onTicketEscape);
for (const tab of tabEls) {
  tab.addEventListener('click', () => setView(tab.dataset.view));
}

async function loadPrds() {
  showBanner('');
  const res = await fetch('/api/dashboard');
  if (!res.ok) {
    showBanner('Failed to load dashboard (' + res.status + ')', true);
    return;
  }
  dashboard = await res.json();
  cards = dashboard.prds || [];
  graphLink.hidden = !(dashboard.graph && dashboard.graph.available);
  renderSystemStatus();
  ensureWorkflowSelection();
  render();
  renderWorkflow();
  renderAgents();
}

function renderSystemStatus() {
  systemStatusEl.innerHTML = '';
  const memory = dashboard.memory || {};
  const engram = memory.engram || {};
  const registry = memory.skillRegistry || {};

  const engramText = engram.available
    ? \`Engram: \${engram.memories || 0} memories / \${engram.sessions || 0} sessions\`
    : 'Engram: no .engram manifest';
  systemStatusEl.appendChild(renderStatusChip(engramText));

  const registryText = registry.available
    ? \`Skill registry: \${registry.skills || 0} skills\`
    : 'Skill registry: missing';
  systemStatusEl.appendChild(renderStatusChip(registryText));
}

function renderStatusChip(text) {
  const chip = document.createElement('span');
  chip.className = 'status-chip';
  const [label, rest] = String(text).split(/:\s*/, 2);
  if (rest !== undefined) {
    const strong = document.createElement('strong');
    strong.textContent = label + ':';
    chip.appendChild(strong);
    chip.appendChild(document.createTextNode(' ' + rest));
  } else {
    chip.textContent = text;
  }
  return chip;
}

function render() {
  boardEl.innerHTML = '';
  for (const statusDef of statusColumns()) {
    const status = statusDef.id;
    const col = document.createElement('section');
    col.className = 'column';
    col.dataset.status = status;

    const h = document.createElement('h2');
    h.textContent = statusDef.label || status;
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
  meta.textContent = c.file + (c.changeId ? ' / ' + c.changeId : '');
  el.appendChild(meta);

  const badges = document.createElement('div');
  badges.className = 'card-badges';
  if (c.handoff && c.handoff.phase) badges.appendChild(renderBadge('phase: ' + c.handoff.phase));
  if (c.latestEvent && c.latestEvent.status) badges.appendChild(renderBadge(c.latestEvent.status));
  if (c.blocked) badges.appendChild(renderBadge('blocked', 'blocked'));
  if (badges.children.length) el.appendChild(badges);

  const links = document.createElement('div');
  links.className = 'card-links';
  if (c.previews && c.previews.length) {
    links.appendChild(renderPillLink('Preview', c.previews[0].url, true));
  }
  if (c.external && c.external.url) {
    links.appendChild(renderPillLink(c.external.label || 'External', c.external.url, true));
  } else if (c.external && c.external.label) {
    links.appendChild(renderBadge(c.external.label));
  }
  const workflowButton = document.createElement('button');
  workflowButton.type = 'button';
  workflowButton.className = 'pill-link';
  workflowButton.textContent = 'Workflow';
  workflowButton.addEventListener('click', (e) => {
    e.stopPropagation();
    selectWorkflowCard(c);
    setView('workflow');
  });
  links.appendChild(workflowButton);
  el.appendChild(links);

  const actions = document.createElement('div');
  actions.className = 'card-actions';

  const sel = document.createElement('select');
  sel.className = 'card-move';
  sel.setAttribute('aria-label', \`Move \${c.file}\`);
  for (const statusDef of statusColumns()) {
    const s = statusDef.id;
    const opt = document.createElement('option');
    opt.value = s;
    opt.textContent = statusDef.label || s;
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

function renderBadge(text, extraClass) {
  const badge = document.createElement('span');
  badge.className = 'badge' + (extraClass ? ' ' + extraClass : '');
  const normalized = String(text || '').toLowerCase();
  if (['passed', 'completed', 'failed', 'blocked'].includes(normalized)) {
    badge.classList.add(normalized);
  }
  badge.textContent = text;
  return badge;
}

function renderPillLink(label, href, external) {
  const a = document.createElement('a');
  a.className = 'pill-link';
  a.href = href;
  a.textContent = label;
  if (external) {
    a.target = '_blank';
    a.rel = 'noreferrer';
  }
  return a;
}

function renderWorkflow() {
  const card = ensureWorkflowSelection();
  renderWorkflowList();
  if (!card) {
    workflowDetail.className = 'workflow-empty';
    workflowDetail.textContent = 'No PRDs found in .skillgrid/prd.';
    return;
  }
  renderWorkflowDetail(card);
}

function renderWorkflowList() {
  workflowList.innerHTML = '';
  for (const c of cards) {
    const item = document.createElement('button');
    item.type = 'button';
    item.className = 'workflow-item';
    item.classList.toggle('is-active', c.file === selectedWorkflowFile);
    item.addEventListener('click', () => {
      selectWorkflowCard(c);
      renderWorkflow();
    });
    const title = document.createElement('span');
    title.className = 'workflow-item-title';
    title.textContent = c.title;
    item.appendChild(title);
    const meta = document.createElement('span');
    meta.className = 'workflow-item-meta';
    meta.textContent = statusLabel(c.status) + ' / ' + (c.changeId || c.file);
    item.appendChild(meta);
    workflowList.appendChild(item);
  }
}

function renderWorkflowDetail(c) {
  workflowDetail.className = '';
  workflowDetail.innerHTML = '';

  const header = document.createElement('div');
  header.className = 'workflow-header';
  const h = document.createElement('h2');
  h.textContent = c.title;
  header.appendChild(h);
  const meta = document.createElement('div');
  meta.className = 'card-meta';
  meta.textContent = 'Change: ' + c.changeId + ' / PRD: ' + c.file + ' / Status: ' + statusLabel(c.status);
  header.appendChild(meta);
  const headerLinks = document.createElement('div');
  headerLinks.className = 'artifact-list';
  headerLinks.appendChild(renderPillLink('Open PRD detail', '#', false));
  headerLinks.lastChild.addEventListener('click', (e) => {
    e.preventDefault();
    openTicket(c.file);
  });
  for (const artifact of c.artifacts || []) {
    if (artifact.url) headerLinks.appendChild(renderPillLink(artifact.label, artifact.url, true));
  }
  header.appendChild(headerLinks);
  workflowDetail.appendChild(header);

  const grid = document.createElement('div');
  grid.className = 'workflow-grid';
  grid.appendChild(renderInfoCard('Current phase', (c.handoff && c.handoff.phase) || 'unknown'));
  grid.appendChild(renderInfoCard('Next action', (c.handoff && c.handoff.nextAction) || 'not recorded'));
  grid.appendChild(renderInfoCard('HITL blockers', listOrNone(c.handoff && c.handoff.hitlBlockers)));
  grid.appendChild(renderInfoCard('AFK-ready work', listOrNone(c.handoff && c.handoff.afkReady)));
  workflowDetail.appendChild(grid);

  const phaseTitle = document.createElement('h3');
  phaseTitle.textContent = 'Phase status';
  workflowDetail.appendChild(phaseTitle);
  const phaseList = document.createElement('div');
  phaseList.className = 'phase-list';
  const phaseMap = (c.workflow && c.workflow.phases) || {};
  for (const phase of workflowPhases()) {
    const latest = phaseMap[phase];
    const row = document.createElement('div');
    row.className = 'phase-row';
    const node = document.createElement('span');
    node.className = 'phase-node';
    node.textContent = phase;
    row.appendChild(node);
    row.appendChild(renderBadge(latest ? latest.status : 'no events'));
    const summary = document.createElement('span');
    summary.className = 'phase-summary';
    summary.textContent = latest ? latest.summary || '' : '';
    row.appendChild(summary);
    phaseList.appendChild(row);
  }
  workflowDetail.appendChild(phaseList);

  const eventTitle = document.createElement('h3');
  eventTitle.textContent = 'Handoff event log';
  workflowDetail.appendChild(eventTitle);
  const events = (c.workflow && c.workflow.events) || [];
  const eventList = document.createElement('div');
  eventList.className = 'event-list';
  if (!events.length) {
    eventList.textContent = 'No events yet. Append JSONL lines under .skillgrid/tasks/events/' + c.changeId + '.jsonl.';
  } else {
    for (const ev of events) {
      eventList.appendChild(renderEvent(ev));
    }
  }
  workflowDetail.appendChild(eventList);
}

function renderInfoCard(title, value) {
  const card = document.createElement('div');
  card.className = 'workflow-card';
  const h = document.createElement('h3');
  h.textContent = title;
  card.appendChild(h);
  const p = document.createElement('div');
  p.textContent = value;
  card.appendChild(p);
  return card;
}

function listOrNone(items) {
  if (!items || !items.length) return 'none';
  return items.join('; ');
}

function renderEvent(ev) {
  const el = document.createElement('article');
  el.className = 'event';
  const time = document.createElement('div');
  time.className = 'event-time';
  time.textContent = [ev.time, ev.phase || ev.node, ev.status].filter(Boolean).join(' / ');
  el.appendChild(time);
  const summary = document.createElement('p');
  summary.className = 'event-summary';
  summary.textContent = ev.summary || '(no summary)';
  el.appendChild(summary);
  if (ev.blocker) el.appendChild(renderBadge('blocked: ' + ev.blocker, 'blocked'));
  if (ev.artifactLinks && ev.artifactLinks.length) {
    const artifacts = document.createElement('div');
    artifacts.className = 'artifact-list';
    for (const artifact of ev.artifactLinks) {
      if (artifact.url) artifacts.appendChild(renderPillLink(artifact.label, artifact.url, true));
      else artifacts.appendChild(renderBadge(artifact.label));
    }
    el.appendChild(artifacts);
  }
  return el;
}

function renderAgents() {
  agentsList.innerHTML = '';
  const actions = dashboard.agentActions || [];
  if (!actions.length) {
    const empty = document.createElement('div');
    empty.className = 'workflow-empty';
    empty.textContent = 'No subagent actions found yet. Add event log entries with agent, subagent, role, or subagent/review node fields.';
    agentsList.appendChild(empty);
    return;
  }
  for (const action of actions) {
    agentsList.appendChild(renderAgentAction(action));
  }
}

function renderAgentAction(action) {
  const el = document.createElement('article');
  el.className = 'agent-action';

  const head = document.createElement('div');
  head.className = 'agent-action-head';
  const name = document.createElement('span');
  name.className = 'agent-name';
  name.textContent = action.agent || action.role || 'subagent';
  head.appendChild(name);
  if (action.role && action.role !== action.agent) head.appendChild(renderBadge(action.role));
  if (action.status) head.appendChild(renderBadge(action.status));
  el.appendChild(head);

  const summary = document.createElement('p');
  summary.className = 'agent-summary';
  summary.textContent = action.summary || action.task || '(no summary)';
  el.appendChild(summary);

  const meta = document.createElement('div');
  meta.className = 'agent-meta';
  meta.appendChild(document.createTextNode([action.time, action.changeId, action.prd].filter(Boolean).join(' / ')));
  el.appendChild(meta);

  if (action.blocker) el.appendChild(renderBadge('blocked: ' + action.blocker, 'blocked'));
  if (action.artifactLinks && action.artifactLinks.length) {
    const artifacts = document.createElement('div');
    artifacts.className = 'artifact-list';
    for (const artifact of action.artifactLinks) {
      if (artifact.url) artifacts.appendChild(renderPillLink(artifact.label, artifact.url, true));
      else artifacts.appendChild(renderBadge(artifact.label));
    }
    el.appendChild(artifacts);
  }

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

function skillgridRootFromPrdRoot(prdRoot) {
  return path.dirname(path.resolve(prdRoot));
}

function repoRootFromPrdRoot(prdRoot) {
  return path.dirname(skillgridRootFromPrdRoot(prdRoot));
}

function configError(message) {
  const err = new Error(`invalid .skillgrid/config.json prdWorkflow: ${message}`);
  err.code = 'invalid_config';
  return err;
}

function cloneDefaultPrdWorkflow() {
  return {
    source: DEFAULT_PRD_WORKFLOW.source,
    preset: DEFAULT_PRD_WORKFLOW.preset,
    fallbackStatus: DEFAULT_PRD_WORKFLOW.fallbackStatus,
    statuses: DEFAULT_PRD_WORKFLOW.statuses.map((status) => ({ ...status })),
    phaseStatusMap: { ...DEFAULT_PRD_WORKFLOW.phaseStatusMap },
    providerMapping: {},
  };
}

function normalizeStatusId(value, fieldName) {
  const id = String(value ?? '').trim();
  if (!STATUS_ID_RE.test(id)) {
    throw configError(`${fieldName} must be a single token using letters, numbers, dots, underscores, or hyphens`);
  }
  return id;
}

function statusIdSet(prdWorkflow) {
  return new Set(prdWorkflow.statuses.map((status) => status.id));
}

function normalizePrdWorkflow(config = {}) {
  const raw = config && typeof config === 'object' ? config.prdWorkflow : null;
  if (raw == null) return cloneDefaultPrdWorkflow();
  if (typeof raw !== 'object' || Array.isArray(raw)) {
    throw configError('must be an object when present');
  }

  const statusesRaw = raw.statuses;
  if (!Array.isArray(statusesRaw) || statusesRaw.length === 0) {
    throw configError('statuses must be a non-empty array');
  }

  const seen = new Set();
  const statuses = statusesRaw.map((entry, index) => {
    const item = typeof entry === 'string' ? { id: entry } : entry;
    if (!item || typeof item !== 'object' || Array.isArray(item)) {
      throw configError(`statuses[${index}] must be a string or object`);
    }
    const id = normalizeStatusId(item.id, `statuses[${index}].id`);
    if (seen.has(id)) throw configError(`duplicate status id "${id}"`);
    seen.add(id);
    const label = String(item.label ?? id).trim() || id;
    return { id, label };
  });

  const ids = new Set(statuses.map((status) => status.id));
  const fallbackStatus = normalizeStatusId(raw.fallbackStatus ?? statuses[0].id, 'fallbackStatus');
  if (!ids.has(fallbackStatus)) {
    throw configError(`fallbackStatus "${fallbackStatus}" is not listed in statuses`);
  }

  const phaseStatusMap = {};
  const rawPhaseMap = raw.phaseStatusMap && typeof raw.phaseStatusMap === 'object' && !Array.isArray(raw.phaseStatusMap)
    ? raw.phaseStatusMap
    : {};
  for (const [phase, status] of Object.entries(rawPhaseMap)) {
    const phaseId = normalizeStatusId(phase, 'phaseStatusMap key');
    const statusId = normalizeStatusId(status, `phaseStatusMap.${phaseId}`);
    if (!ids.has(statusId)) {
      throw configError(`phaseStatusMap.${phaseId} points to unknown status "${statusId}"`);
    }
    phaseStatusMap[phaseId] = statusId;
  }

  const providerMapping = raw.providerMapping && typeof raw.providerMapping === 'object' && !Array.isArray(raw.providerMapping)
    ? raw.providerMapping
    : {};

  return {
    source: typeof raw.source === 'string' && raw.source.trim() ? raw.source.trim() : 'custom',
    preset: typeof raw.preset === 'string' && raw.preset.trim() ? raw.preset.trim() : null,
    fallbackStatus,
    statuses,
    phaseStatusMap,
    providerMapping,
  };
}

async function readSkillgridConfig(prdRoot) {
  const configPath = path.join(skillgridRootFromPrdRoot(prdRoot), 'config.json');
  if (!existsSync(configPath)) return {};
  try {
    return JSON.parse(await fs.readFile(configPath, 'utf8'));
  } catch (e) {
    const err = new Error(`invalid .skillgrid/config.json: ${e.message || e}`);
    err.code = 'invalid_config';
    throw err;
  }
}

async function loadPrdWorkflow(prdRoot) {
  return normalizePrdWorkflow(await readSkillgridConfig(prdRoot));
}

function envWithLegacy(currentName, legacyName) {
  const current = process.env[currentName];
  if (current != null && current !== '') return current;
  return process.env[legacyName];
}

function parseArgs(argv) {
  let port = Number(envWithLegacy('SKILLGRID_UI_PORT', 'PRD_KANBAN_PORT')) || 8787;
  let host = envWithLegacy('SKILLGRID_UI_HOST', 'PRD_KANBAN_HOST') || '127.0.0.1';
  let prdDirOpt = envWithLegacy('SKILLGRID_UI_DIR', 'PRD_KANBAN_DIR') || null;

  for (let i = 2; i < argv.length; i++) {
    const a = argv[i];
    if (a === '--port' && argv[i + 1]) {
      port = Number(argv[++i]);
    } else if (a === '--host' && argv[i + 1]) {
      host = argv[++i];
    } else if (a === '--prd-dir' && argv[i + 1]) {
      prdDirOpt = argv[++i];
    } else if (a === '--help' || a === '-h') {
      console.log(`Usage: node skillgrid-ui.mjs [--port N] [--host addr] [--prd-dir path]

Default PRD dir: ../prd next to this script (Skillgrid layout), else walk up for .skillgrid/prd, else <cwd>/.skillgrid/prd
Explicit --prd-dir / SKILLGRID_UI_DIR: absolute or relative to cwd

Requires: Node.js 18+ (no repo packages; browser uses CDN for markdown preview)
Env: SKILLGRID_UI_PORT, SKILLGRID_UI_HOST, SKILLGRID_UI_DIR
Deprecated compatibility env: PRD_KANBAN_PORT, PRD_KANBAN_HOST, PRD_KANBAN_DIR`);
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

function safeFileUnder(rootDir, fileParam, allowedExts = null) {
  const name = path.basename(decodeURIComponent(fileParam || ''));
  if (!name || name === '.' || name === '..') return null;
  if (allowedExts && !allowedExts.has(path.extname(name).toLowerCase())) return null;
  const root = path.resolve(rootDir);
  const resolved = path.resolve(root, name);
  if (path.dirname(resolved) !== root) return null;
  return resolved;
}

function toPosixPath(p) {
  return String(p || '').replace(/\\/g, '/');
}

function urlForArtifact(repoRoot, artifactPath) {
  const normalized = toPosixPath(artifactPath);
  const previewPrefix = '.skillgrid/preview/';
  if (normalized.startsWith(previewPrefix)) {
    return '/preview/' + encodeURIComponent(path.basename(normalized));
  }
  if (normalized === 'graphify-out/graph.html') return '/graphify/graph.html';
  return null;
}

function firstH1Title(body) {
  const m = String(body).match(/^\s*#\s+(.+)$/m);
  return m ? m[1].trim() : null;
}

function firstPrdTitle(body) {
  const m = String(body).match(/^\s*#{1,3}\s+PRD:\s+(.+)$/mi);
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

function getMarkdownMetaField(body, label) {
  const escaped = String(label).replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  const re = new RegExp(`^-\\s*\\*\\*${escaped}:\\*\\*\\s*(.+)$`, 'mi');
  const mm = String(body || '').match(re);
  if (!mm) return null;
  return mm[1].trim().replace(/^`|`$/g, '');
}

function extractChangeIdFromPath(value) {
  const raw = toPosixPath(String(value || '').trim().replace(/^`|`$/g, ''));
  const mm = raw.match(/openspec\/changes\/([^/\s`]+)/);
  return mm ? mm[1] : null;
}

function extractContextChangeId(value) {
  const raw = toPosixPath(String(value || '').trim().replace(/^`|`$/g, ''));
  const mm = raw.match(/context_([^/\s`]+)\.md$/);
  return mm ? mm[1] : null;
}

function prdSlugFromFile(fileName) {
  return String(fileName || '')
    .replace(/^PRD\d+_/, '')
    .replace(/\.md$/, '')
    .replace(/_/g, '-');
}

function fallbackChangeId(fileName) {
  return prdSlugFromFile(fileName).toLowerCase();
}

function parseExternal(value) {
  const raw = String(value || '').trim();
  if (!raw || raw.toLowerCase() === 'local') return null;
  const markdownLink = raw.match(/^\[([^\]]+)\]\(([^)]+)\)$/);
  if (markdownLink) return { label: markdownLink[1], url: markdownLink[2] };
  if (/^https?:\/\//i.test(raw)) return { label: raw, url: raw };
  return { label: raw, url: null };
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

function effectiveStatusString(statusRaw, prdWorkflow = DEFAULT_PRD_WORKFLOW) {
  const s = statusRaw != null ? String(statusRaw).trim().replace(/^`|`$/g, '') : '';
  if (s && statusIdSet(prdWorkflow).has(s)) return s;
  return prdWorkflow.fallbackStatus;
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
  const prdWorkflow = opts.prdWorkflow || DEFAULT_PRD_WORKFLOW;
  const rawStr = rawBuf.toString('utf8');
  const revision = sha256Hex(rawBuf);
  const { front, body } = splitFrontMatter(rawStr);
  const bodyOut = body;
  const st = getYamlField(front, 'status') || getMarkdownMetaField(bodyOut, 'Status');
  const status = effectiveStatusString(st, prdWorkflow);
  let title = fileName;
  const t = getYamlField(front, 'title');
  if (t) title = t;
  else {
    const h = firstPrdTitle(bodyOut) || firstH1Title(bodyOut);
    if (h) title = h;
  }
  const specPath = getYamlField(front, 'spec') || getMarkdownMetaField(bodyOut, 'Spec / change');
  const contextPath = getYamlField(front, 'context') || getMarkdownMetaField(bodyOut, 'Session context');
  const previewPath = getYamlField(front, 'preview') || getMarkdownMetaField(bodyOut, 'Preview');
  const externalRaw = getYamlField(front, 'external') || getMarkdownMetaField(bodyOut, 'External');
  const changeId =
    getYamlField(front, 'changeId') ||
    extractChangeIdFromPath(specPath) ||
    extractContextChangeId(contextPath) ||
    fallbackChangeId(fileName);
  const external = parseExternal(externalRaw);
  const out = {
    file: fileName,
    title,
    status,
    revision,
    changeId,
    specPath: specPath || null,
    contextPath: contextPath || `.skillgrid/tasks/context_${changeId}.md`,
    previewPath: previewPath || null,
    external,
  };
  if (includeBody) out.body = bodyOut;
  return out;
}

async function getPrdDetail(prdRoot, fileParam, prdWorkflow = null) {
  const workflow = prdWorkflow || await loadPrdWorkflow(prdRoot);
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
  return enrichPrd(prdRoot, extractPrdFromRaw(raw, path.basename(full), { includeBody: true, prdWorkflow: workflow }));
}

async function listPrds(prdRoot, prdWorkflow = DEFAULT_PRD_WORKFLOW) {
  const entries = await fs.readdir(prdRoot, { withFileTypes: true });
  const out = [];
  for (const ent of entries) {
    if (!ent.isFile()) continue;
    if (!PRD_NAME_RE.test(ent.name)) continue;
    const full = path.join(prdRoot, ent.name);
    const raw = await fs.readFile(full);
    try {
      const meta = extractPrdFromRaw(raw, ent.name, { prdWorkflow });
      out.push(meta);
    } catch {
      out.push({
        file: ent.name,
        title: ent.name,
        status: prdWorkflow.fallbackStatus,
        revision: sha256Hex(raw),
        changeId: fallbackChangeId(ent.name),
        specPath: null,
        contextPath: `.skillgrid/tasks/context_${fallbackChangeId(ent.name)}.md`,
        previewPath: null,
        external: null,
      });
    }
  }
  out.sort((a, b) => a.file.localeCompare(b.file));
  return out;
}

async function readTextIfExists(fullPath) {
  if (!fullPath || !existsSync(fullPath)) return null;
  return fs.readFile(fullPath, 'utf8');
}

function extractSectionList(markdown, heading) {
  const text = String(markdown || '');
  const escaped = String(heading).replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  const re = new RegExp(`^##\\s+${escaped}\\s*$([\\s\\S]*?)(?=^##\\s+|$)`, 'im');
  const mm = text.match(re);
  if (!mm) return [];
  return mm[1]
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter((line) => /^[-*]\s+/.test(line))
    .map((line) => line.replace(/^[-*]\s+/, '').replace(/^\[[ xX]\]\s+/, '').trim())
    .filter(Boolean);
}

function extractSingleLineValue(markdown, label) {
  const escaped = String(label).replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  const re = new RegExp(`^-\\s*\\*\\*${escaped}:\\*\\*\\s*(.+)$`, 'mi');
  const mm = String(markdown || '').match(re);
  return mm ? mm[1].trim() : null;
}

function parseHandoff(markdown) {
  if (!markdown) {
    return {
      exists: false,
      phase: null,
      status: null,
      nextAction: null,
      hitlBlockers: [],
      afkReady: [],
    };
  }
  const nextMatch = String(markdown).match(/^##\s+Next recommended action\s*$([\s\S]*?)(?=^##\s+|$)/im);
  const nextAction = nextMatch ? nextMatch[1].trim().split(/\r?\n/).find(Boolean) || null : null;
  return {
    exists: true,
    phase: extractSingleLineValue(markdown, 'Phase'),
    status: extractSingleLineValue(markdown, 'Status'),
    nextAction,
    hitlBlockers: extractSectionList(markdown, 'HITL blockers').concat(extractSectionList(markdown, 'HITL / human gates')),
    afkReady: extractSectionList(markdown, 'AFK-ready work'),
  };
}

async function listPreviewArtifacts(repoRoot, prd) {
  const previewRoot = path.join(repoRoot, '.skillgrid', 'preview');
  if (!existsSync(previewRoot)) return [];
  const candidates = new Map();
  const add = (fileName, explicit = false) => {
    if (!fileName) return;
    const base = path.basename(fileName);
    if (!base.toLowerCase().endsWith('.html')) return;
    candidates.set(base, {
      label: explicit ? `Preview: ${base}` : base,
      path: `.skillgrid/preview/${base}`,
      url: '/preview/' + encodeURIComponent(base),
    });
  };

  add(prd.previewPath, true);
  const entries = await fs.readdir(previewRoot, { withFileTypes: true }).catch(() => []);
  const needles = [prd.changeId, prdSlugFromFile(prd.file)]
    .filter(Boolean)
    .map((s) => String(s).toLowerCase());
  for (const ent of entries) {
    if (!ent.isFile()) continue;
    const lower = ent.name.toLowerCase();
    if (!lower.endsWith('.html')) continue;
    if (needles.some((needle) => lower.startsWith(needle))) add(ent.name);
  }
  return Array.from(candidates.values()).sort((a, b) => a.path.localeCompare(b.path));
}

function artifactLink(repoRoot, artifactPath) {
  const normalized = toPosixPath(String(artifactPath || '').trim().replace(/^`|`$/g, ''));
  if (!normalized) return null;
  const url = urlForArtifact(repoRoot, normalized);
  return {
    label: path.basename(normalized) || normalized,
    path: normalized,
    url,
  };
}

function withArtifactLinks(repoRoot, event) {
  const artifacts = Array.isArray(event.artifacts) ? event.artifacts : [];
  return {
    ...event,
    artifactLinks: artifacts.map((p) => artifactLink(repoRoot, p)).filter(Boolean),
  };
}

async function readWorkflowEvents(repoRoot, changeId) {
  if (!changeId) return [];
  const eventsPath = path.join(repoRoot, '.skillgrid', 'tasks', 'events', `${changeId}.jsonl`);
  const text = await readTextIfExists(eventsPath);
  if (!text) return [];
  return text
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => {
      try {
        return JSON.parse(line);
      } catch {
        return null;
      }
    })
    .filter((event) => event && typeof event === 'object')
    .map((event) => withArtifactLinks(repoRoot, event))
    .sort((a, b) => String(a.time || '').localeCompare(String(b.time || '')));
}

function phaseMapFromEvents(events) {
  const phases = {};
  for (const event of events) {
    const key = event.phase || event.node;
    if (!key) continue;
    phases[key] = event;
  }
  return phases;
}

function isSubagentEvent(event) {
  const values = [
    event.agent,
    event.subagent,
    event.role,
    event.node,
    event.phase,
    event.summary,
  ]
    .filter(Boolean)
    .map((value) => String(value).toLowerCase());
  return values.some((value) =>
    value.includes('subagent') ||
    value.includes('agent') ||
    value.includes('reviewer') ||
    value.includes('implementer') ||
    value.includes('critic') ||
    value.includes('auditor') ||
    value.includes('verifier')
  );
}

function agentActionFromEvent(prd, event) {
  const agent = event.agent || event.subagent || event.role || event.node || 'subagent';
  return {
    ...event,
    agent,
    role: event.role || event.node || null,
    task: event.task || null,
    changeId: event.changeId || prd.changeId,
    prd: event.prd || prd.file,
    prdTitle: prd.title,
  };
}

function collectAgentActions(prds) {
  return prds
    .flatMap((prd) =>
      ((prd.workflow && prd.workflow.events) || [])
        .filter(isSubagentEvent)
        .map((event) => agentActionFromEvent(prd, event))
    )
    .sort((a, b) => String(b.time || '').localeCompare(String(a.time || '')));
}

async function listResearchArtifacts(repoRoot, changeId) {
  const researchRoot = path.join(repoRoot, '.skillgrid', 'tasks', 'research', changeId || '');
  if (!changeId || !existsSync(researchRoot)) return [];
  const entries = await fs.readdir(researchRoot, { withFileTypes: true }).catch(() => []);
  return entries
    .filter((ent) => ent.isFile())
    .map((ent) => ({
      label: ent.name,
      path: `.skillgrid/tasks/research/${changeId}/${ent.name}`,
      url: null,
    }))
    .sort((a, b) => a.path.localeCompare(b.path));
}

async function readEngramManifest(repoRoot) {
  const manifestPath = path.join(repoRoot, '.engram', 'manifest.json');
  if (!existsSync(manifestPath)) {
    return {
      available: false,
      path: '.engram/manifest.json',
      chunks: 0,
      sessions: 0,
      memories: 0,
      prompts: 0,
    };
  }
  try {
    const manifest = JSON.parse(await fs.readFile(manifestPath, 'utf8'));
    const chunks = Array.isArray(manifest.chunks) ? manifest.chunks : [];
    return {
      available: true,
      path: '.engram/manifest.json',
      chunks: chunks.length,
      sessions: chunks.reduce((sum, chunk) => sum + Number(chunk.sessions || 0), 0),
      memories: chunks.reduce((sum, chunk) => sum + Number(chunk.memories || 0), 0),
      prompts: chunks.reduce((sum, chunk) => sum + Number(chunk.prompts || 0), 0),
    };
  } catch (e) {
    return {
      available: false,
      path: '.engram/manifest.json',
      error: String(e.message || e),
      chunks: 0,
      sessions: 0,
      memories: 0,
      prompts: 0,
    };
  }
}

async function readSkillRegistryMetadata(repoRoot) {
  const registryPath = path.join(repoRoot, '.skillgrid', 'project', 'SKILL_REGISTRY.md');
  if (!existsSync(registryPath)) {
    return {
      available: false,
      path: '.skillgrid/project/SKILL_REGISTRY.md',
      skills: 0,
      conventions: 0,
    };
  }
  const text = await fs.readFile(registryPath, 'utf8');
  const skills = (text.match(/^###\s+\S+/gm) || []).length;
  const conventionsSection = text.match(/^##\s+Project Conventions\s*$([\s\S]*?)(?=^##\s+|$)/im);
  const conventions = conventionsSection
    ? conventionsSection[1].split(/\r?\n/).filter((line) => /^\|\s*`?[^|-]/.test(line.trim())).length
    : 0;
  return {
    available: true,
    path: '.skillgrid/project/SKILL_REGISTRY.md',
    skills,
    conventions,
  };
}

async function readMemoryMetadata(repoRoot) {
  const [engram, skillRegistry] = await Promise.all([
    readEngramManifest(repoRoot),
    readSkillRegistryMetadata(repoRoot),
  ]);
  return { engram, skillRegistry };
}

async function enrichPrd(prdRoot, prd) {
  const repoRoot = repoRootFromPrdRoot(prdRoot);
  const changeId = prd.changeId || fallbackChangeId(prd.file);
  const contextRel = prd.contextPath || `.skillgrid/tasks/context_${changeId}.md`;
  const contextFull = path.join(repoRoot, contextRel);
  const handoffText = await readTextIfExists(contextFull);
  const handoff = parseHandoff(handoffText);
  const events = await readWorkflowEvents(repoRoot, changeId);
  const latestEvent = events.length ? events[events.length - 1] : null;
  const previews = await listPreviewArtifacts(repoRoot, prd);
  const research = await listResearchArtifacts(repoRoot, changeId);
  const artifacts = [
    { label: prd.file, path: `.skillgrid/prd/${prd.file}`, url: null },
    prd.specPath ? { label: path.basename(prd.specPath), path: prd.specPath, url: null } : null,
    { label: path.basename(contextRel), path: contextRel, url: null },
    ...previews,
    ...research,
  ].filter(Boolean);
  const blocked = Boolean(
    (latestEvent && latestEvent.status === 'blocked') ||
      (latestEvent && latestEvent.blocker) ||
      (handoff.status && String(handoff.status).toLowerCase().includes('blocked')) ||
      handoff.hitlBlockers.length
  );
  return {
    ...prd,
    changeId,
    contextPath: contextRel,
    handoff,
    previews,
    artifacts,
    latestEvent,
    blocked,
    workflow: {
      events,
      phases: phaseMapFromEvents(events),
    },
  };
}

async function buildDashboardData(prdRoot) {
  const repoRoot = repoRootFromPrdRoot(prdRoot);
  const prdWorkflow = await loadPrdWorkflow(prdRoot);
  const prds = await Promise.all((await listPrds(prdRoot, prdWorkflow)).map((prd) => enrichPrd(prdRoot, prd)));
  return {
    prds,
    prdWorkflow,
    agentActions: collectAgentActions(prds),
    memory: await readMemoryMetadata(repoRoot),
    workflowPhases: DEFAULT_WORKFLOW_PHASES,
    graph: {
      available: existsSync(path.join(repoRoot, 'graphify-out', 'graph.html')),
      path: 'graphify-out/graph.html',
      url: '/graphify/graph.html',
    },
  };
}

async function updateStatus(prdRoot, fileParam, body, currentRevision, prdWorkflow = null) {
  const workflow = prdWorkflow || await loadPrdWorkflow(prdRoot);
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
  if (!statusIdSet(workflow).has(body.status)) {
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
  '.json': 'application/json; charset=utf-8',
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

async function serveDiskFile(res, fullPath) {
  if (!fullPath || !existsSync(fullPath)) {
    res.writeHead(404);
    res.end('Not found');
    return;
  }
  const ext = path.extname(fullPath);
  const data = await fs.readFile(fullPath);
  res.writeHead(200, { 'Content-Type': MIME[ext] || 'application/octet-stream' });
  res.end(data);
}

function validChangeId(value) {
  return /^[a-zA-Z0-9._-]+$/.test(String(value || ''));
}

async function handle(req, res, prdRoot) {
  const u = new URL(req.url || '/', 'http://127.0.0.1');
  const repoRoot = repoRootFromPrdRoot(prdRoot);

  if (u.pathname === '/api/prds' && req.method === 'GET') {
    try {
      const data = await buildDashboardData(prdRoot);
      json(res, 200, data.prds);
    } catch (e) {
      json(res, 500, { error: String(e.message || e) });
    }
    return;
  }

  if (u.pathname === '/api/dashboard' && req.method === 'GET') {
    try {
      const data = await buildDashboardData(prdRoot);
      json(res, 200, data);
    } catch (e) {
      json(res, 500, { error: String(e.message || e) });
    }
    return;
  }

  if (u.pathname.startsWith('/api/changes/') && req.method === 'GET') {
    const rest = u.pathname.slice('/api/changes/'.length).split('/');
    const changeId = decodeURIComponent(rest[0] || '');
    const resource = rest[1] || '';
    if (!validChangeId(changeId)) {
      json(res, 404, { error: 'not_found' });
      return;
    }
    try {
      if (resource === 'events') {
        json(res, 200, await readWorkflowEvents(repoRoot, changeId));
        return;
      }
      if (resource === 'artifacts') {
        const data = await buildDashboardData(prdRoot);
        const prd = data.prds.find((item) => item.changeId === changeId);
        json(res, 200, prd ? prd.artifacts : []);
        return;
      }
    } catch (e) {
      json(res, 500, { error: String(e.message || e) });
      return;
    }
  }

  if (u.pathname === '/api/agents' && req.method === 'GET') {
    try {
      const data = await buildDashboardData(prdRoot);
      json(res, 200, data.agentActions);
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
    if (p.startsWith('/preview/')) {
      const filePart = p.slice('/preview/'.length);
      const full = safeFileUnder(path.join(repoRoot, '.skillgrid', 'preview'), filePart, new Set(['.html']));
      await serveDiskFile(res, full);
      return;
    }
    if (p === '/graphify/graph.html') {
      await serveDiskFile(res, path.join(repoRoot, 'graphify-out', 'graph.html'));
      return;
    }
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
    console.error('Pass --prd-dir <path> or set SKILLGRID_UI_DIR (PRD_KANBAN_DIR is deprecated but still accepted)');
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
    console.log(`Skillgrid Dashboard serving http://${host}:${port}`);
    console.log(`PRD dir: ${prdDir}`);
  });
}

export {
  DEFAULT_PRD_WORKFLOW,
  DEFAULT_WORKFLOW_PHASES,
  applyStatusToDocument,
  buildDashboardData,
  effectiveStatusString,
  extractPrdFromRaw,
  loadPrdWorkflow,
  normalizePrdWorkflow,
  updateStatus,
};

if (process.argv[1] && path.resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  main();
}
