# Orkflow

**YAML-driven orchestration engine for multi-agent AI systems.**

Define your agents, wire their collaboration patterns, and run complex workflows — all without writing orchestration code.

## Core Principle

> *Configuration defines collaboration. Execution is automatic.*

## Features

| Feature | Description |
|---------|-------------|
| **Declarative YAML** | Define agents, roles, goals, and workflows in simple YAML |
| **Sequential Execution** | Chain agents in order with automatic context passing |
| **Parallel Execution** | Run agents concurrently with aggregation support |
| **Shared Memory** | Agents publish/subscribe to data via `outputs`/`requires` |
| **Multi-Provider** | OpenAI, Gemini, Anthropic, Ollama, and any OpenAI-compatible API |
| **Built-in Tools** | `calc`, `file`, `script` tools for agent capabilities |
| **MCP Support** | Connect external tool servers (filesystem, databases, etc.) |
| **Session Persistence** | Automatic session saving and continuation |
| **Execution Logs** | Detailed file-based logging with `--log` flag |

## Quick Start

```bash
# Build
go build -o orka cmd/orka/main.go

# Run a workflow
./orka run examples/sequential-workflow.yaml

# With logging enabled
./orka run examples/parallel-workflow.yaml --log

# Validate a workflow
./orka validate examples/tool-enabled-workflow.yaml

# List sessions
./orka sessions list

# View session with workflow graph
./orka sessions show <session-id> --workflow
```

## YAML Examples

### Sequential Workflow
```yaml
models:
  gpt4:
    provider: openai
    model: gpt-4o-mini

agents:
  - id: researcher
    role: Research Assistant
    goal: Research electric vehicles
    model: gpt4
    outputs:
      - research_notes

  - id: writer
    role: Content Writer
    goal: Write a summary using the research
    model: gpt4
    requires:
      - research_notes

workflow:
  type: sequential
  steps:
    - agent: researcher
    - agent: writer
```

### Parallel Workflow
```yaml
agents:
  - id: backend
    role: Backend Engineer
    goal: Design API
    outputs: [api_design]

  - id: frontend
    role: Frontend Engineer
    goal: Design UI
    outputs: [ui_design]

  - id: reviewer
    role: Tech Lead
    goal: Review both designs
    requires: [api_design, ui_design]

workflow:
  type: parallel
  branches: [backend, frontend]
  then:
    agent: reviewer
```

### Tool-Enabled Workflow
```yaml
agents:
  - id: analyst
    role: Data Analyst
    goal: Calculate metrics and inspect files
    tools:
      - calc   # Math expressions
      - file   # Filesystem operations
      - script # Tengo scripts
```

### MCP Integration
```yaml
mcp_servers:
  filesystem:
    command: npx
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]

agents:
  - id: dev
    role: Developer
    goal: List files in /tmp
    toolsets:
      - filesystem
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `orka run <file.yaml>` | Execute a workflow |
| `orka run <file.yaml> --log` | Execute with file logging |
| `orka run <file.yaml> --continue` | Continue last session |
| `orka validate <file.yaml>` | Validate workflow syntax |
| `orka sessions list` | List all sessions |
| `orka sessions show <id>` | Show session details |
| `orka sessions show <id> --workflow` | Show workflow visualization |

## Environment Variables

| Variable | Provider |
|----------|----------|
| `OPENAI_API_KEY` | OpenAI |
| `GEMINI_API_KEY` | Google Gemini |
| `ANTHROPIC_API_KEY` | Anthropic Claude |

## Project Structure

```
Orkflow/
├── cmd/orka/           # CLI entrypoint
├── internal/
│   ├── agent/          # LLM clients (OpenAI, Gemini, Ollama, etc.)
│   ├── cli/            # CLI commands
│   ├── engine/         # Workflow executor
│   ├── logging/        # Execution logger
│   ├── mcp/            # MCP client and tool adapter
│   ├── memory/         # Session and shared memory
│   ├── parser/         # YAML parser
│   └── tools/          # Built-in tools (calc, file, script)
├── pkg/types/          # Shared types
└── examples/           # Example workflows
```

## License

MIT
