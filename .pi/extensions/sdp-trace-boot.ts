import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";

export default function (pi: ExtensionAPI) {
  pi.on("session_start", async (event, ctx) => {
    // Only on fresh startup, new session, resume, or fork — skip reload
    if (event.reason === "reload") return;

    // Verify we are in the sdp-trace repo by checking for AGENTS.md
    const fs = await import("node:fs/promises");
    const path = await import("node:path");
    const agentsFile = path.join(ctx.cwd, "AGENTS.md");
    try {
      await fs.access(agentsFile);
    } catch {
      // Not in sdp-trace repo; do nothing
      return;
    }

    // Skip non-interactive modes (print, json, rpc) where follow-up delivery conflicts
    if (!ctx.hasUI) return;

    // Queue the boot reminder as a follow-up so it runs after startup settles
    pi.sendUserMessage("/sdp-trace-boot", { deliverAs: "followUp" });
  });
}
