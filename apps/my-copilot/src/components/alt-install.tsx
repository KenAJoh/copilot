"use client";

import { useState } from "react";
import { CopyButton } from "@navikt/ds-react";

const INSTALL_SCRIPT_COMMAND =
  "curl -fsSL https://raw.githubusercontent.com/navikt/copilot/main/scripts/install.sh | bash";

export function AltInstall() {
  const [open, setOpen] = useState(false);

  return (
    <div style={{ marginTop: "0.5rem" }}>
      <button
        onClick={() => setOpen(!open)}
        style={{
          background: "none",
          border: "none",
          padding: 0,
          cursor: "pointer",
          fontSize: "0.8125rem",
          color: "#64748b",
          textDecoration: "underline",
          textDecorationStyle: "dotted",
          textUnderlineOffset: "2px",
        }}
        aria-expanded={open}
      >
        {open ? "Skjul" : "Ikke Homebrew? Linux / CI →"}
      </button>
      {open && (
        <div
          className="rounded-lg overflow-hidden border border-gray-200 shadow-sm flex items-center justify-between mt-2"
          style={{ background: "#f1f5f9" }}
        >
          <code className="font-mono whitespace-nowrap flex-1 p-3" style={{ fontSize: "0.75rem", color: "#334155" }}>
            {INSTALL_SCRIPT_COMMAND}
          </code>
          <div className="shrink-0 pr-3">
            <CopyButton copyText={INSTALL_SCRIPT_COMMAND} size="xsmall" />
          </div>
        </div>
      )}
    </div>
  );
}
