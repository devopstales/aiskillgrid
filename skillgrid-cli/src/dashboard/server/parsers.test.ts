import { describe, expect, it } from "vitest";
import { firstParagraph, parseMarkdownMetadata, parseTaskStats } from "./parsers.js";

describe("parseMarkdownMetadata", () => {
  it("reads common PRD fields from markdown", () => {
    const metadata = parseMarkdownMetadata(`# Launch Dashboard

Status: In Progress
Change ID: build-dashboard
Owner: platform
Agent: cursor
Preview: .skillgrid/preview/dashboard.html
External: https://example.com

Body text
`);

    expect(metadata).toMatchObject({
      title: "Launch Dashboard",
      status: "In Progress",
      changeId: "build-dashboard",
      owner: "platform",
      agent: "cursor",
      previewUrl: ".skillgrid/preview/dashboard.html",
      externalUrl: "https://example.com"
    });
  });

  it("drops placeholder Spec / change metadata", () => {
    const metadata = parseMarkdownMetadata(`# PRD10

**Spec / change:** _OpenSpec change not created in this repo yet — align with \`openspec/changes/\` when added._
`);

    expect(metadata.specPath).toBe("");
  });
});

describe("firstParagraph", () => {
  it("skips Skillgrid PRD metadata blocks", () => {
    expect(
      firstParagraph(`# PRD10: kubelogin Interactive OIDC Login & IdP CA Management

**File:** \`.skillgrid/prd/PRD10_kubelogin_interactive_login.md\`  
**Spec / change:** _OpenSpec change not created in this repo yet — align with \`openspec/changes/\` when added._  
**Session context:** \`.skillgrid/tasks/context_kubelogin-interactive-login.md\`  
**Status:** draft

> **Imported from:** \`project/prd/kubelogin-interactive-login.md\`

---

### PRD: kubelogin Interactive OIDC Login & IdP CA Management

**OpenSpec change**: \`openspec/changes/kubelogin-interactive-login/\`  
**Status**: Proposed (see OpenSpec for canonical state)

#### Problem / Why

Today KubeDash focuses on generating kubeconfigs that rely on auth-provider=oidc.
`)
    ).toBe("Today KubeDash focuses on generating kubeconfigs that rely on auth-provider=oidc.");
  });
});

describe("parseTaskStats", () => {
  it("counts checkboxes and workflow tags", () => {
    expect(
      parseTaskStats(`
- [x] Create proposal [AFK]
- [ ] Review architecture [HITL]
- [ ] Fix blocked deploy [blocked]
`)
    ).toEqual({
      total: 3,
      completed: 1,
      hitl: 1,
      afk: 1,
      blocked: 1
    });
  });
});
