
# Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.