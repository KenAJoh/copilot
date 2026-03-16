# 📦 Copilot Collections

Collections are curated sets of agents, prompts, instructions, and skills organized around specific Nav development domains or workflows.

Collections provide a convenient way to install related Copilot customizations as a complete package, ensuring consistent tooling and practices across Nav teams.

### How to Use Collections

**To Install:**
- Download the complete collection folder
- Copy to your workspace's `.github/` directory
- Collections automatically integrate with GitHub Copilot

**What's Included:**
- **Agents**: Specialized AI assistants for domain-specific tasks
- **Prompts**: Reusable templates for common workflows
- **Instructions**: Technology and pattern-specific guidelines
- **Skills**: Bundled resources and templates

## Available Collections

| Name        | Description                        | Components | Location               |
| ----------- | ---------------------------------- | ---------- | ---------------------- |
| Coming Soon | Nav collections will be added here |            | `.github/collections/` |

## Potential Nav Collections

### NAIS Platform Collection
Complete set of tools for NAIS platform development:
- **Agents**: Platform deployment, configuration management
- **Prompts**: Manifest generation, troubleshooting
- **Instructions**: NAIS best practices, security patterns
- **Skills**: Deployment templates, monitoring setup

### Full-Stack Nav Application Collection
Everything needed for a complete Nav application:
- **Agents**: Frontend (Aksel/Next.js), Backend (Kotlin/Ktor), Database
- **Prompts**: Component creation, API design, database migrations
- **Instructions**: Coding standards for all layers
- **Skills**: Application templates, testing patterns

### Auth & Security Collection
Authentication and security patterns for Nav:
- **Agents**: Azure AD setup, TokenX configuration
- **Prompts**: Security validation, access control
- **Instructions**: Security best practices
- **Skills**: Auth configuration templates

### Observability Collection
Complete observability stack setup:
- **Agents**: Metrics, logging, tracing
- **Prompts**: Dashboard creation, alert rules
- **Instructions**: Observability patterns
- **Skills**: Instrumentation templates

## Collection Structure

```
.github/collections/
└── collection-name/
    ├── README.md             # Collection documentation
    ├── agents/              # Collection-specific agents
    ├── prompts/             # Collection-specific prompts
    ├── instructions/        # Collection-specific instructions
    └── skills/              # Collection-specific skills
```

## Creating Collections

When creating collections for Nav:

1. **Cohesive Theme**: Focus on specific domain or workflow
2. **Complete Package**: Include all necessary components
3. **Nav Alignment**: Follow Nav development principles
4. **Documentation**: Provide clear setup and usage guides
5. **Dependencies**: Document required MCP servers or tools
6. **Team Input**: Gather feedback from Nav teams

## Installation

To install a collection:

```bash
# Navigate to your project
cd /path/to/your/project

# Copy the collection
cp -r /path/to/copilot/docs/collections/collection-name .github/
```

Or use the provided installation tasks in VS Code.

## Best Practices

- Use collections for comprehensive domain coverage
- Keep collections focused and maintainable
- Update collections based on Nav team feedback
- Version collections for stability
- Document collection dependencies
- Provide examples and getting started guides
