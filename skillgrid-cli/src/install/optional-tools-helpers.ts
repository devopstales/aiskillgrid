import type { OptionalToolId } from "./types.js";

export function toolIsSelected(selected: OptionalToolId[], id: OptionalToolId): boolean {
  return selected.includes(id);
}
