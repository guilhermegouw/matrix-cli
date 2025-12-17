# Provider System Documentation

This document provides comprehensive documentation for the provider system introduced in [PR #2](https://github.com/guilhermegouw/matrix-cli/pull/2), including the configuration system, quality tooling, and development workflow.

## Table of Contents

- [Overview](#overview)
- [Provider System Architecture](#provider-system-architecture)
  - [Two-Tier Model System](#two-tier-model-system)
  - [Provider Builder](#provider-builder)
  - [Catwalk Integration](#catwalk-integration)
  - [Fantasy Integration](#fantasy-integration)
- [Configuration System](#configuration-system)
  - [Configuration Loading](#configuration-loading)
  - [Configuration Structure](#configuration-structure)
  - [Environment Variable Resolution](#environment-variable-resolution)
  - [Provider Configuration](#provider-configuration)
- [Quality Tooling](#quality-tooling)
  - [GitHub Actions CI](#github-actions-ci)
  - [golangci-lint Configuration](#golangci-lint-configuration)
  - [Taskfile](#taskfile)
- [Test Coverage](#test-coverage)

---

## Overview

The provider system enables Matrix CLI to interact with multiple LLM providers (OpenAI, Anthropic, and OpenAI-compatible APIs). It follows a two-tier model architecture (large/small) for optimizing cost and performance across different task complexities.

Key features:
- **Multi-provider support**: OpenAI, Anthropic, and OpenAI-compatible providers
- **Two-tier model system**: Large models for complex tasks, small models for simpler tasks
- **Catwalk integration**: Provider metadata and model information from Charm's catwalk service
- **Fantasy integration**: LLM orchestration through Charm's fantasy library
- **Environment variable resolution**: Secure configuration with `$VAR` and `${VAR}` syntax
- **Provider caching**: Performance optimization with embedded fallback

---

## Provider System Architecture

### Two-Tier Model System

The system uses two model tiers following the pattern established by Crush CLI:

| Tier | Purpose | Use Case |
|------|---------|----------|
| `large` | Complex tasks requiring full reasoning | Code generation, analysis, complex refactoring |
| `small` | Simpler, faster tasks | Quick queries, simple transformations, summarization |

**Implementation**: `internal/config/config.go:13-21`

```go
type SelectedModelType string

const (
    SelectedModelTypeLarge SelectedModelType = "large"
    SelectedModelTypeSmall SelectedModelType = "small"
)
```

If only a large model is configured, it automatically falls back to use the large model for small tasks as well.

### Provider Builder

The `Builder` struct (`internal/provider/provider.go`) creates and caches fantasy providers:

```go
type Builder struct {
    cfg   *config.Config
    cache map[string]fantasy.Provider
    debug bool
}
```

**Key methods**:
- `NewBuilder(cfg)`: Creates a new builder from configuration
- `BuildModels(ctx)`: Creates large and small models from configuration
- `buildModel(ctx, modelCfg)`: Builds a single model with provider and catwalk metadata
- `getOrBuildProvider(providerCfg, modelCfg)`: Returns cached provider or builds new one

**Provider caching**: Providers are cached by ID to avoid redundant instantiation when the same provider is used for both tiers.

### Catwalk Integration

Catwalk provides provider metadata including:
- Available models per provider
- Default large/small model IDs
- API endpoints
- Model capabilities and pricing

**Loading hierarchy** (`internal/config/providers.go:27-58`):
1. Fetch from Catwalk API (`https://catwalk.charm.sh`)
2. Fall back to local cache (24-hour TTL)
3. Fall back to embedded provider data

**Cache location**: `$XDG_DATA_HOME/matrix/providers.json`

**Manual update**:
```go
UpdateProviders(cfg, source) // source: "embedded", URL, or file path
```

### Fantasy Integration

Fantasy is Charm's LLM orchestration library providing a unified interface across providers.

**Supported providers**:
- `anthropic`: Native Anthropic API
- `openai`: Native OpenAI API
- `openai-compat`: OpenAI-compatible APIs (Ollama, vLLM, etc.)

**Special handling**:
- **Anthropic thinking mode**: Automatically adds `anthropic-beta: interleaved-thinking-2025-05-14` header when `think: true`
- **OAuth tokens**: Detects `Bearer ` prefix and handles authorization header correctly

---

## Configuration System

### Configuration Loading

Configuration is loaded from multiple sources with cascading precedence (`internal/config/load.go`):

1. **Global config**: `$XDG_CONFIG_HOME/matrix/matrix.json`
2. **Project config**: `matrix.json` or `.matrix.json` (searched upward from cwd)

Project configuration takes precedence over global configuration.

**Loading process**:
1. Load global config (if exists)
2. Find and load project config (searching parent directories)
3. Merge configs (project overrides global)
4. Apply defaults
5. Load provider metadata from catwalk
6. Configure providers (merge user config with catwalk metadata)
7. Configure default model selections

### Configuration Structure

**Top-level config** (`internal/config/config.go:75-86`):

```json
{
  "models": {
    "large": { ... },
    "small": { ... }
  },
  "providers": {
    "anthropic": { ... },
    "openai": { ... }
  },
  "options": {
    "debug": false,
    "data_directory": "",
    "context_paths": []
  }
}
```

**Model selection** (`SelectedModel`):

| Field | Type | Description |
|-------|------|-------------|
| `model` | string | Model ID (e.g., "claude-sonnet-4-20250514") |
| `provider` | string | Provider ID matching a key in providers |
| `think` | bool | Enable thinking mode (Anthropic) |
| `reasoning_effort` | string | Reasoning effort (OpenAI) |
| `temperature` | float64 | Sampling temperature (0-1) |
| `top_p` | float64 | Nucleus sampling parameter |
| `top_k` | int64 | Top-k sampling parameter |
| `max_tokens` | int64 | Maximum response tokens |
| `frequency_penalty` | float64 | Reduces repetition |
| `presence_penalty` | float64 | Increases topic diversity |
| `provider_options` | map | Additional provider-specific options |

**Provider configuration** (`ProviderConfig`):

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique provider identifier |
| `name` | string | Human-readable display name |
| `type` | string | Provider type (openai, anthropic, openai-compat) |
| `api_key` | string | Authentication key (supports env vars) |
| `base_url` | string | Custom API endpoint |
| `disable` | bool | Disable this provider |
| `extra_headers` | map | Additional HTTP headers |
| `models` | array | Available models (from catwalk or user) |
| `provider_options` | map | Additional provider-specific options |

### Environment Variable Resolution

The resolver (`internal/config/resolve.go`) expands environment variables in configuration values:

**Supported syntax**:
- `$VAR` - Simple variable reference
- `${VAR}` - Braced variable reference

**Example configuration**:
```json
{
  "providers": {
    "anthropic": {
      "api_key": "$ANTHROPIC_API_KEY"
    },
    "openai": {
      "api_key": "${OPENAI_API_KEY}",
      "base_url": "${OPENAI_BASE_URL}"
    }
  }
}
```

**Behavior**:
- Returns error if referenced variable is not set
- Providers with unresolvable API keys are skipped (not fatal)
- Base URLs fall back to catwalk defaults if not set

### Provider Configuration

Example complete configuration:

```json
{
  "models": {
    "large": {
      "model": "claude-sonnet-4-20250514",
      "provider": "anthropic",
      "think": true,
      "max_tokens": 8192
    },
    "small": {
      "model": "claude-3-5-haiku-latest",
      "provider": "anthropic"
    }
  },
  "providers": {
    "anthropic": {
      "api_key": "$ANTHROPIC_API_KEY"
    },
    "openai": {
      "api_key": "$OPENAI_API_KEY"
    }
  },
  "options": {
    "debug": false
  }
}
```

---

## Quality Tooling

### GitHub Actions CI

The CI workflow (`.github/workflows/ci.yml`) runs on push/PR to main:

| Job | Description |
|-----|-------------|
| **lint** | Runs golangci-lint with 5-minute timeout |
| **test** | Runs tests with race detection and coverage |
| **build** | Builds all packages |
| **security** | Runs govulncheck for vulnerability scanning |

**Features**:
- Go 1.25 support
- Race condition detection
- Code coverage with Codecov integration
- Security vulnerability scanning

### golangci-lint Configuration

The `.golangci.yml` enables 26 linters organized by category:

**Default linters**:
- `errcheck`: Unchecked errors
- `govet`: Suspicious constructs
- `ineffassign`: Unused assignments
- `staticcheck`: Static analysis
- `unused`: Unused code detection

**Additional linters**:
- `bodyclose`: HTTP response body closure
- `gosec`: Security issues
- `gocritic`: Opinionated checks
- `gocyclo`: Cyclomatic complexity (max 15)
- `dupl`: Code duplication (threshold 100)
- `errorlint`: Error wrapping issues
- `exhaustive`: Enum switch exhaustiveness
- And more...

**Exclusions**:
- Test files: Relaxed rules for dupl, gocyclo, gosec, errcheck
- Main files: Allow init functions

### Taskfile

The `Taskfile.yaml` provides common development commands:

**Development**:
```bash
task build          # Build binary with version info
task run            # Run application
task install        # Install to $GOPATH/bin
```

**Quality**:
```bash
task fmt            # Format code (gofmt + goimports)
task lint           # Run golangci-lint
task lint:fix       # Run linters with auto-fix
task vet            # Run go vet
```

**Testing**:
```bash
task test           # Run tests with race detection
task test:coverage  # Run tests with coverage report
task test:short     # Run short tests only
```

**Other**:
```bash
task security       # Run govulncheck
task deps           # Download and tidy dependencies
task deps:update    # Update all dependencies
task clean          # Remove build artifacts
task check          # Run all checks (fmt, lint, vet, test)
task ci             # Run CI checks locally
```

**Build information**: Version, commit, and build date are injected via ldflags:
```go
-X github.com/guilhermegouw/matrix-cli/cmd.Version={{.VERSION}}
-X github.com/guilhermegouw/matrix-cli/cmd.Commit={{.COMMIT}}
-X github.com/guilhermegouw/matrix-cli/cmd.BuildDate={{.BUILD_DATE}}
```

---

## Test Coverage

The provider system has comprehensive test coverage:

| Package | Coverage | Key Test Areas |
|---------|----------|----------------|
| `internal/config` | 88.4% | Config loading, merging, resolution, provider configuration |
| `internal/provider` | 98.8% | Builder, model creation, provider caching, tier selection |

**Test files**:
- `internal/config/config_test.go`: Config struct methods
- `internal/config/load_test.go`: Configuration loading and merging
- `internal/config/providers_test.go`: Provider metadata loading and caching
- `internal/config/resolve_test.go`: Environment variable resolution
- `internal/provider/provider_test.go`: Provider building and model creation
- `internal/provider/tier_test.go`: Tier selection and validation

**Running tests**:
```bash
task test           # All tests
task test:coverage  # With coverage report
go test -v ./internal/config/...    # Config package only
go test -v ./internal/provider/...  # Provider package only
```

---

## CLI Commands

### Root Command

```bash
matrix              # Shows help
matrix --help       # Detailed help
```

### Version Command

```bash
matrix version
# Output:
# matrix v0.1.0
#   commit: abc1234
#   built:  2025-12-17T12:00:00Z
```

---

## File Structure

```
matrix-cli/
├── cmd/
│   ├── root.go           # Root cobra command
│   └── version.go        # Version command with build info
├── internal/
│   ├── config/
│   │   ├── config.go     # Config structures and types
│   │   ├── load.go       # Configuration loading logic
│   │   ├── providers.go  # Catwalk provider integration
│   │   └── resolve.go    # Environment variable resolver
│   └── provider/
│       ├── provider.go   # Provider builder and model creation
│       └── tier.go       # Tier selection utilities
├── .github/
│   └── workflows/
│       └── ci.yml        # GitHub Actions workflow
├── .golangci.yml         # Linter configuration
├── Taskfile.yaml         # Development tasks
└── main.go               # Application entry point
```

---

## Dependencies

Key external dependencies:

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/adrg/xdg` | XDG Base Directory support |
| `github.com/charmbracelet/catwalk` | Provider metadata |
| `charm.land/fantasy` | LLM orchestration |
| `charm.land/fantasy/providers/anthropic` | Anthropic provider |
| `charm.land/fantasy/providers/openai` | OpenAI provider |
