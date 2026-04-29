#!/usr/bin/env node
import assert from 'node:assert/strict';
import fs from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';

import {
  buildDashboardData,
  effectiveStatusString,
  extractPrdFromRaw,
  loadPrdWorkflow,
  normalizePrdWorkflow,
  updateStatus,
} from './skillgrid-ui.mjs';

async function makeFixture(config = null, prdBody = null) {
  const root = await fs.mkdtemp(path.join(os.tmpdir(), 'skillgrid-ui-'));
  const skillgrid = path.join(root, '.skillgrid');
  const prdRoot = path.join(skillgrid, 'prd');
  await fs.mkdir(prdRoot, { recursive: true });
  if (config) {
    await fs.writeFile(path.join(skillgrid, 'config.json'), JSON.stringify(config, null, 2), 'utf8');
  }
  const body = prdBody || `---
title: Example PRD
status: draft
---
# Example PRD
`;
  const prdPath = path.join(prdRoot, 'PRD01_example.md');
  await fs.writeFile(prdPath, body, 'utf8');
  return { root, prdRoot, prdPath };
}

async function testDefaultWorkflow() {
  const { prdRoot } = await makeFixture();
  const workflow = await loadPrdWorkflow(prdRoot);
  assert.deepEqual(workflow.statuses.map((status) => status.id), ['draft', 'todo', 'inprogress', 'devdone', 'done']);
  assert.equal(workflow.fallbackStatus, 'draft');

  const dashboard = await buildDashboardData(prdRoot);
  assert.equal(dashboard.prdWorkflow.preset, 'skillgrid-default');
  assert.equal(dashboard.prds[0].status, 'draft');
}

async function testCustomWorkflow() {
  const config = {
    prdWorkflow: {
      source: 'custom',
      fallbackStatus: 'backlog',
      statuses: [
        { id: 'backlog', label: 'Backlog' },
        { id: 'selected', label: 'Selected' },
        { id: 'doing', label: 'Doing' },
        { id: 'review', label: 'Review' },
        { id: 'released', label: 'Released' },
      ],
      phaseStatusMap: {
        plan: 'backlog',
        apply: 'doing',
        validate: 'review',
        finish: 'released',
      },
    },
  };
  const { prdRoot, prdPath } = await makeFixture(config, `---
title: Custom PRD
status: doing
---
# Custom PRD
`);
  const dashboard = await buildDashboardData(prdRoot);
  assert.deepEqual(dashboard.prdWorkflow.statuses.map((status) => status.id), ['backlog', 'selected', 'doing', 'review', 'released']);
  assert.equal(dashboard.prds[0].status, 'doing');

  await fs.access(prdPath);
  await updateStatus(prdRoot, 'PRD01_example.md', { status: 'review' }, dashboard.prds[0].revision);
  const updated = await buildDashboardData(prdRoot);
  assert.equal(updated.prds[0].status, 'review');
}

async function testFallbackForUnknownStatus() {
  const workflow = normalizePrdWorkflow({
    prdWorkflow: {
      fallbackStatus: 'backlog',
      statuses: ['backlog', 'doing'],
    },
  });
  assert.equal(effectiveStatusString('unknown', workflow), 'backlog');

  const prd = extractPrdFromRaw(Buffer.from(`---
title: Unknown
status: unknown
---
# Unknown
`), 'PRD01_unknown.md', { prdWorkflow: workflow });
  assert.equal(prd.status, 'backlog');
}

async function testInvalidConfig() {
  assert.throws(
    () => normalizePrdWorkflow({
      prdWorkflow: {
        fallbackStatus: 'missing',
        statuses: ['backlog', 'doing'],
      },
    }),
    /fallbackStatus "missing" is not listed/
  );
  assert.throws(
    () => normalizePrdWorkflow({
      prdWorkflow: {
        statuses: ['backlog', 'backlog'],
      },
    }),
    /duplicate status id/
  );
}

async function testRejectsPatchOutsideWorkflow() {
  const { prdRoot } = await makeFixture({
    prdWorkflow: {
      fallbackStatus: 'backlog',
      statuses: ['backlog', 'doing'],
    },
  });
  const dashboard = await buildDashboardData(prdRoot);
  await assert.rejects(
    () => updateStatus(prdRoot, 'PRD01_example.md', { status: 'done' }, dashboard.prds[0].revision),
    /invalid_status/
  );
}

await testDefaultWorkflow();
await testCustomWorkflow();
await testFallbackForUnknownStatus();
await testInvalidConfig();
await testRejectsPatchOutsideWorkflow();

console.log('skillgrid-ui workflow tests passed');
