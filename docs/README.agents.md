# 🤖 Custom Agents

Custom agents for NAV's GitHub Copilot ecosystem, specialized for Norwegian public sector development patterns.

### How to Install

Agents are `.agent.md` files placed in your repo's `.github/agents/` directory.

| Editor          | Install Method                                                                                    |
| --------------- | ------------------------------------------------------------------------------------------------- |
| **VS Code**     | Click the one-click install button below, or copy the file to `.github/agents/`                   |
| **JetBrains**   | Not supported for Copilot Chat. Works with Copilot coding agent when the file exists in the repo. |
| **Copilot CLI** | Not supported. Use instructions files instead.                                                    |
| **GitHub.com**  | Works with Copilot coding agent when the file exists in the repo.                                 |

> Custom agents (`.agent.md`) are currently a VS Code feature for Copilot Chat. Other editors support them only through the Copilot coding agent. See [support matrix](https://docs.github.com/en/copilot/reference/custom-instructions-support) for details.

**Manual install (any editor):**

```bash
# From your project root
mkdir -p .github/agents
curl -sO --output-dir .github/agents \
  https://raw.githubusercontent.com/navikt/copilot/main/.github/agents/<filename>.agent.md
```

**To use in VS Code:** Type `@agent-name` in Copilot Chat after installing.

## Available Agents

| Agent                                                                                                      | Description                                                                                        | VS Code                                                                                                                                                                                                                                                                                                                                |
| ---------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Nais Agent**<br/>[`@nais-agent`](../.github/agents/nais.agent.md)                                        | NAIS platform deployments, Kubernetes configs, `.nais/app.yaml` manifests, PostgreSQL/Kafka setup. | [![Install](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://min-copilot.ansatt.nav.no/install/agent?url=vscode%3Achat-agent%2Finstall%3Furl%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fnavikt%2Fcopilot%2Fmain%2F.github%2Fagents%2Fnais.agent.md)              |
| **Auth Agent**<br/>[`@auth-agent`](../.github/agents/auth.agent.md)                                        | Azure AD, TokenX, ID-porten, Maskinporten authentication and JWT validation.                       | [![Install](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://min-copilot.ansatt.nav.no/install/agent?url=vscode%3Achat-agent%2Finstall%3Furl%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fnavikt%2Fcopilot%2Fmain%2F.github%2Fagents%2Fauth.agent.md)              |
| **Kafka Agent**<br/>[`@kafka-agent`](../.github/agents/kafka.agent.md)                                     | Rapids & Rivers event handling, Kafka consumers, event schemas, TestRapid.                         | [![Install](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://min-copilot.ansatt.nav.no/install/agent?url=vscode%3Achat-agent%2Finstall%3Furl%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fnavikt%2Fcopilot%2Fmain%2F.github%2Fagents%2Fkafka.agent.md)             |
| **Aksel Agent**<br/>[`@aksel-agent`](../.github/agents/aksel.agent.md)                                     | Aksel Design System, Tailwind-to-Aksel conversion, responsive layouts, mobile-first.               | [![Install](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://min-copilot.ansatt.nav.no/install/agent?url=vscode%3Achat-agent%2Finstall%3Furl%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fnavikt%2Fcopilot%2Fmain%2F.github%2Fagents%2Faksel.agent.md)             |
| **Observability Agent**<br/>[`@observability-agent`](../.github/agents/observability.agent.md)             | Prometheus metrics, Grafana Loki, Tempo tracing, health endpoints, DORA metrics.                   | [![Install](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://min-copilot.ansatt.nav.no/install/agent?url=vscode%3Achat-agent%2Finstall%3Furl%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fnavikt%2Fcopilot%2Fmain%2F.github%2Fagents%2Fobservability.agent.md)     |
| **Security Champion Agent**<br/>[`@security-champion-agent`](../.github/agents/security-champion.agent.md) | NAV security standards, network policies, secrets management, GDPR compliance.                     | [![Install](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://min-copilot.ansatt.nav.no/install/agent?url=vscode%3Achat-agent%2Finstall%3Furl%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fnavikt%2Fcopilot%2Fmain%2F.github%2Fagents%2Fsecurity-champion.agent.md) |
| **Research Agent**<br/>[`@research-agent`](../.github/agents/research.agent.md)                            | Codebase research, issue investigation, pattern analysis, dependency tracing.                      | [![Install](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://min-copilot.ansatt.nav.no/install/agent?url=vscode%3Achat-agent%2Finstall%3Furl%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fnavikt%2Fcopilot%2Fmain%2F.github%2Fagents%2Fresearch.agent.md)          |
| **Forfatter**<br/>[`@forfatter`](../.github/agents/forfatter.agent.md)                                    | Norsk teknisk redaktør: klarspråk, AI-markører, anglismer, fagtermer, mikrotekst.                | [![Install](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://min-copilot.ansatt.nav.no/install/agent?url=vscode%3Achat-agent%2Finstall%3Furl%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fnavikt%2Fcopilot%2Fmain%2F.github%2Fagents%2Fforfatter.agent.md)         |

## Creating Custom Agents

When creating new agents for NAV projects:

1. **Follow NAV Standards**: Align with NAV development principles (Team First, Essential Complexity, DORA Metrics)
2. **Include Context**: Reference NAV tech stack (Kotlin/Ktor, Next.js, NAIS)
3. **Security First**: Always consider security implications and NAV security policies
4. **Norwegian Language**: Support Norwegian text and number formatting where applicable
5. **Platform Integration**: Ensure compatibility with NAIS deployment patterns

## Agent Guidelines

- Agents should be self-contained and focused on specific domains
- Include clear examples and common use cases
- Reference relevant NAV documentation and standards
- Support both local development and NAIS-deployed scenarios
- Consider mobile-first design for frontend-related agents
