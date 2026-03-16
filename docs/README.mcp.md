# 🔌 MCP Servers

Nav-approved Model Context Protocol (MCP) servers available through the [Nav MCP Registry](https://mcp-registry.nav.no).

MCP servers extend GitHub Copilot with tools that connect to external services — GitHub APIs, design tools, and internal Nav systems.

## How to Install

### VS Code

**Option 1 — MCP Registry (recommended):**

1. Open the Extensions panel (`Cmd+Shift+X`)
2. Click the filter icon → select **MCP Registry**
3. Search for the server name and click **Install**

**Option 2 — Manual configuration:**

Add to `.vscode/mcp.json` in your repo (shared with team) or VS Code `settings.json` (personal):

```json
{
  "servers": {
    "server-name": {
      "type": "http",
      "url": "https://server-url/mcp"
    }
  }
}
```

### JetBrains

Add to `.idea/mcp.json` in your project:

```json
{
  "servers": {
    "server-name": {
      "type": "http",
      "url": "https://server-url/mcp"
    }
  }
}
```

### Copilot CLI

```bash
gh copilot mcp add --type http server-name https://server-url/mcp
```

Or edit `~/.config/github-copilot/mcp.json` directly.

> See the [GitHub MCP docs](https://docs.github.com/en/copilot/customizing-copilot/extending-copilot-chat-with-mcp) for full setup instructions per editor.

## Available MCP Servers

| Server                    | Description                                                                                                 | URL                                        |
| ------------------------- | ----------------------------------------------------------------------------------------------------------- | ------------------------------------------ |
| **GitHub MCP**            | GitHub repositories, issues, pull requests, and code search.                                                | `https://api.githubcopilot.com/mcp/`       |
| **Nav Copilot Discovery** | Discover Nav Copilot customizations, assess agent readiness, generate AGENTS.md. Requires Nav GitHub OAuth. | `https://mcp-onboarding.intern.nav.no/mcp` |
| **Figma MCP**             | Bring Figma design context into your coding workflow for design-to-code generation.                         | `https://mcp.figma.com/mcp`                |

### GitHub MCP

The official GitHub MCP server — available to all Copilot users. No additional setup needed in VS Code (built-in).

```json
{
  "servers": {
    "github": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/"
    }
  }
}
```

### Nav Copilot Discovery (mcp-onboarding)

Internal Nav MCP server for discovering customizations and generating AGENTS.md files. Requires GitHub OAuth through your `@navikt` organization membership.

```json
{
  "servers": {
    "nav-copilot": {
      "type": "http",
      "url": "https://mcp-onboarding.intern.nav.no/mcp"
    }
  }
}
```

**Available tools:** `hello_world`, `greet`, `whoami`, `echo`, `get_time`, `list_instructions`, `list_prompts`, `list_agents`, `list_skills`, `search_customizations`, `suggest_customizations`, `get_installation_guide`, `check_agent_readiness`, `generate_agents_md`, `generate_setup_steps`, `team_readiness`

### Figma MCP

Bring Figma designs into your coding workflow. Requires a Figma account.

```json
{
  "servers": {
    "figma": {
      "type": "http",
      "url": "https://mcp.figma.com/mcp"
    }
  }
}
```

## MCP Registry API

The Nav MCP Registry is available at `https://mcp-registry.nav.no` and implements the [MCP Registry v0.1 specification](https://modelcontextprotocol.io).

| Endpoint                                   | Description                             |
| ------------------------------------------ | --------------------------------------- |
| `GET /v0.1/servers`                        | List all approved MCP servers           |
| `GET /v0.1/servers/{name}/versions/latest` | Get latest version of a specific server |

Server names use reverse-DNS format with URL encoding: `io.github.navikt%2Fgithub-mcp`
