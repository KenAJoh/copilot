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
          className="rounded-lg p-4 overflow-x-auto flex items-center gap-3 mt-2"
          style={{ background: "#1e1e1e" }}
        >
          <code
            className="font-mono whitespace-nowrap flex-1"
            style={{ fontSize: "0.8rem", color: "#d4d4d4" }}
          >
            {INSTALL_SCRIPT_COMMAND}
          </code>
          <CopyButton copyText={INSTALL_SCRIPT_COMMAND} size="xsmall" style={{ color: "white" }} />
        </div>
      )}
    </div>
  );
}
