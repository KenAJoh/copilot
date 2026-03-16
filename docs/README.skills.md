# 🎯 Agent Skills

Agent Skills are self-contained folders with instructions and bundled resources that enhance AI capabilities for specialized Nav development tasks.

Based on the [Agent Skills specification](https://agentskills.io/specification), each skill contains a `SKILL.md` file with detailed instructions that agents load on-demand.

Skills differ from other primitives by supporting bundled assets (scripts, code samples, reference data) that agents can utilize when performing specialized tasks.

### How to Install

Skills are folders placed in your repo's `.github/skills/` directory.

| Editor          | Install Method                                                                                 |
| --------------- | ---------------------------------------------------------------------------------------------- |
| **VS Code**     | Copy the skill folder to `.github/skills/` in your repo. Agents discover skills automatically. |
| **JetBrains**   | Copy the skill folder to `.github/skills/`. Works with Copilot coding agent.                   |
| **Copilot CLI** | Not supported.                                                                                 |
| **GitHub.com**  | Works with Copilot coding agent when the folder exists in the repo.                            |

> Skills are currently a VS Code feature for Copilot Chat. There is no one-click install — copy the full folder to your repo.

**Manual install:**

```bash
# From your project root — install a single skill
mkdir -p .github/skills
# Clone the repo and copy the skill folder
git clone --depth 1 --filter=blob:none --sparse https://github.com/navikt/copilot.git /tmp/nav-copilot
cd /tmp/nav-copilot && git sparse-checkout set .github/skills/<skill-name>
cp -r .github/skills/<skill-name> /path/to/your/repo/.github/skills/
rm -rf /tmp/nav-copilot
```

**When to use skills vs instructions:** Skills are ideal for complex workflows that need bundled resources (templates, scripts, reference data). For simple coding guidelines, use instructions instead.

## Available Skills

| Name                    | Description                                                                                      | Location                                                                                |
| ----------------------- | ------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------- |
| **ai-news-research**    | Research and write monthly AI coding news summaries in Norwegian with working source links       | [`.github/skills/ai-news-research/`](../.github/skills/ai-news-research/SKILL.md)       |
| **aksel-spacing**       | Responsive layout patterns using Aksel spacing tokens with Box, VStack, HStack, and HGrid        | [`.github/skills/aksel-spacing/`](../.github/skills/aksel-spacing/SKILL.md)             |
| **flyway-migration**    | Database migration patterns using Flyway with versioned SQL scripts                              | [`.github/skills/flyway-migration/`](../.github/skills/flyway-migration/SKILL.md)       |
| **kotlin-app-config**   | Sealed class configuration pattern for Kotlin applications with environment-specific settings    | [`.github/skills/kotlin-app-config/`](../.github/skills/kotlin-app-config/SKILL.md)     |
| **observability-setup** | Setting up Prometheus metrics, OpenTelemetry tracing, and health endpoints for Nais applications | [`.github/skills/observability-setup/`](../.github/skills/observability-setup/SKILL.md) |
| **security-review**     | Pre-commit/PR security checks — use when about to commit, push, or open a pull request           | [`.github/skills/security-review/`](../.github/skills/security-review/SKILL.md)         |
| **tokenx-auth**         | Service-to-service authentication using TokenX token exchange in Nais                            | [`.github/skills/tokenx-auth/`](../.github/skills/tokenx-auth/SKILL.md)                 |
| **web-design-reviewer** | Visual inspection of websites to identify and fix responsive, accessibility, and layout issues   | [`.github/skills/web-design-reviewer/`](../.github/skills/web-design-reviewer/SKILL.md) |

## Creating Nav Skills

When creating agent skills for Nav projects:

1. **Follow Specification**: Adhere to the [Agent Skills specification](https://agentskills.io/specification)
2. **Bundle Resources**: Include templates, scripts, and reference data
3. **Nav Context**: Include Nav-specific patterns and configurations
4. **Self-Contained**: Skills should be independent and reusable
5. **Progressive Disclosure**: Load only when needed for specific tasks

## Skill Structure

```
.github/skills/
└── skill-name/
    ├── SKILL.md              # Main instruction file
    ├── templates/            # Code templates
    ├── scripts/              # Helper scripts
    ├── examples/             # Example implementations
    └── reference/            # Reference documentation
```

## Best Practices

- Keep skills focused on specific domains
- Include practical examples from Nav projects
- Provide clear usage instructions
- Bundle only necessary resources
- Test skills in various Nav contexts
- Document dependencies and requirements
