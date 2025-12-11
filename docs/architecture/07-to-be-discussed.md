# Components To Be Discussed

These components require team discussion before finalizing the architecture. Below is a brief description of how each works in the original CDD agent.

---

## Conversation Management

### CDD Behavior

- **Message history tracking** - Stores user messages, assistant responses, and tool results
- **Context window management** - Prunes old messages when approaching token limits
- **Conversation compaction** - When history gets too long:
  - Keeps first 2 messages (initial context)
  - Keeps last 8 messages (recent context)
  - Summarizes middle messages into a condensed form
- **No persistence** - Conversations are session-only
- **Export via `/save`** - Users can manually export conversation to file

### Open Questions

- Keep compaction strategy or change ratios?
- Add persistence across sessions?
- Auto-save history to `~/.matrix/history/`?

---

## The Operator (General Agent)

### CDD Behavior

The main chat agent that handles general conversations and tool execution.

- **Agentic loop**: Message → LLM → Tool use → Loop until done
- **Full tool access** - Can use all tools
- **Streaming responses** - Token-by-token output
- **Background process management** - Can run and monitor background tasks
- **Context injection** - Receives project context in system prompt
- **Slash command handling** - Delegates to specialized agents when user types `/oracle`, etc.

### Open Questions

- Operator personality/persona in Matrix theme?
- Any behavioral changes from CDD?

---

## The Oracle (Socrates Agent)

### CDD Behavior

Requirements gathering through Socratic dialogue.

- **Lifecycle**: `initialize()` → `process()` loop → `finalize()`
- **Read-only tools** - Can explore codebase but not modify
- **State tracking** - Maintains `gathered_info` dict with discovered requirements
- **Conversation phases**: problem_discovery → user_analysis → requirements → edge_cases → wrap_up
- **Codebase exploration** - Can read files on-demand to understand context
- **Conversation compaction** - Summarizes old exchanges when history grows
- **Output** - Produces a spec file (YAML or markdown)

### Open Questions

- Oracle personality in Matrix theme?
- Same phases or different structure?
- Spec file format (YAML vs markdown)?

---

## The Architect (Planner Agent)

### CDD Behavior

Generates step-by-step implementation plans.

- **No tools** - Pure LLM reasoning, no tool execution
- **Reads spec** - Takes output from Oracle as input
- **Context gathering** - Loads project context, codebase structure, relevant files
- **JSON output** - Produces structured plan with steps, dependencies, complexity
- **Validation** - Checks plan against actual codebase
- **Fallback** - Heuristic plan if LLM fails

**Plan structure:**
```json
{
  "overview": "High-level approach",
  "steps": [
    {
      "number": 1,
      "title": "Step title",
      "description": "What to do",
      "complexity": "simple|medium|complex",
      "dependencies": [step numbers],
      "files_affected": [file paths]
    }
  ],
  "total_complexity": "medium",
  "risks": [list of risks]
}
```

### Open Questions

- Architect personality in Matrix theme?
- Plan format changes?
- Interactive plan refinement?

---

## Seraph (Executor Agent)

### CDD Behavior

Executes implementation plans step-by-step.

- **Full tool access** - Can read, write, execute anything
- **Reads spec + plan** - Takes output from Oracle and Architect
- **Dependency tracking** - Only executes steps when dependencies are met
- **Execution state** - JSON persistence for resume capability
- **YOLO mode** - Auto-continue on success, stop on failure
- **Step commands**: skip, retry, status, continue
- **Code generation** - LLM generates complete file contents for each step

**Code response format:**
````
```python:path/to/file.py
# Complete file content here
```
````

### Open Questions

- Seraph personality in Matrix theme?
- Execution modes (normal vs YOLO)?
- Resume behavior?

---

## Session & Slash Commands

### CDD Behavior

Orchestrates agent switching and command routing.

**ChatSession:**
- Manages current active agent
- Handles switching between Operator and specialized agents
- Maintains session state

**Slash Commands:**
- `/socrates {ticket}` → Enter Oracle agent
- `/plan {ticket}` → Enter Architect agent
- `/exec {ticket}` → Enter Seraph agent
- `/init` → Initialize project structure
- `/help` → Show available commands
- `/clear` → Clear conversation history
- `/save` → Export conversation

**Command routing:**
- Extensible router pattern
- Async command execution
- Error handling with user feedback

### Open Questions

- Command names for Matrix theme?
  - `/oracle` instead of `/socrates`?
  - `/architect` instead of `/plan`?
  - `/seraph` instead of `/exec`?
- Additional commands?
- Exit behavior when leaving specialized agents?

---

## File-Based Agent Communication

### CDD Behavior

Agents communicate via files, not direct calls:

```
specs/tickets/{slug}/
├── spec.yaml           # Oracle output
├── plan.md             # Architect output
└── execution-state.json # Seraph state
```

### Open Questions

- Same directory structure?
- File naming conventions?
- Location (project root vs `~/.matrix/`)?
