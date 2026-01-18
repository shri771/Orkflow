```
Orkflow/
├── cmd/
│   └── orka/                    # CLI entry point
│       └── main.go
│
├── internal/
│   ├── cli/                     # Cobra commands
│   │   ├── root.go              # Root command + global flags
│   │   ├── run.go               # `orka run workflow.yaml`
│   │   ├── validate.go          # `orka validate workflow.yaml`
│   │   └── sessions.go          # `orka sessions list/show`
│   │
│   ├── engine/                  # Core orchestration engine
│   │   ├── executor.go          # Runs sequential/parallel workflows
│   │   └── state.go             # Execution state tracking
│   │
│   ├── parser/                  # YAML parsing & validation
│   │   └── yaml.go              # Parser + validator
│   │
│   ├── agent/                   # LLM client implementations
│   │   ├── agent.go             # Runner + prompt builder
│   │   ├── llm.go               # Client factory
│   │   ├── openai.go            # OpenAI client
│   │   ├── gemini.go            # Google Gemini client
│   │   ├── anthropic.go         # Anthropic Claude client
│   │   ├── ollama.go            # Ollama (local) client
│   │   ├── generic.go           # OpenAI-compatible providers
│   │   └── context.go           # Shared context manager
│   │
│   ├── tools/                   # Built-in tool system
│   │   ├── registry.go          # Tool registry
│   │   ├── executor.go          # Tool call parser/executor
│   │   ├── calc.go              # Math expression evaluator
│   │   ├── file.go              # Filesystem operations
│   │   └── script.go            # Tengo script executor
│   │
│   ├── mcp/                     # Model Context Protocol
│   │   ├── client.go            # MCP server connection manager
│   │   └── tools.go             # MCP-to-registry adapter
│   │
│   ├── memory/                  # Session & shared memory
│   │   ├── session.go           # Persistent session storage
│   │   └── shared.go            # Inter-agent shared memory
│   │
│   ├── logging/                 # Execution logging
│   │   └── logger.go            # File-based logger
│   │
│   └── vectorstore/             # Vector storage (optional)
│       └── chromem.go           # ChromaDB integration
│
├── pkg/
│   └── types/                   # Shared domain types
│       ├── agent.go             # Agent struct
│       └── config.go            # WorkflowConfig struct
│
├── examples/                    # Example YAML workflows
│   ├── sequential-workflow.yaml
│   ├── parallel-workflow.yaml
│   ├── tool-enabled-workflow.yaml
│   ├── tool-1-calc.yaml
│   ├── tool-2-file.yaml
│   └── tool-3-script.yaml
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```
