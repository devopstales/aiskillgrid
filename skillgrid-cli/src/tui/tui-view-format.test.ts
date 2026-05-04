import { describe, expect, it } from "vitest";
import type { BoardIssue, BoardLane } from "../dashboard/shared/types.js";
import { groupIssuesByLane } from "./tui-view-format.js";

describe("groupIssuesByLane", () => {
  const lanes: BoardLane[] = [
    { id: "a", label: "A" },
    { id: "b", label: "B" }
  ];

  it("places issues in configured lanes; unknown lane uses first lane", () => {
    const issues: BoardIssue[] = [
      { id: "1", key: "X-1", title: "t1", lane: "a", status: "s", source: "prd", path: "p", taskStats: { total: 0, completed: 0, hitl: 0, afk: 0, blocked: 0 }, subtasks: [] },
      { id: "2", key: "X-2", title: "t2", lane: "unknown", status: "s", source: "prd", path: "p", taskStats: { total: 0, completed: 0, hitl: 0, afk: 0, blocked: 0 }, subtasks: [] },
      { id: "3", key: "X-3", title: "t3", lane: "b", status: "s", source: "prd", path: "p", taskStats: { total: 0, completed: 0, hitl: 0, afk: 0, blocked: 0 }, subtasks: [] }
    ];
    const m = groupIssuesByLane(issues, lanes);
    expect(m.get("a")!.map((i) => i.id).sort()).toEqual(["1", "2"]);
    expect(m.get("b")!.map((i) => i.id)).toEqual(["3"]);
  });
});
