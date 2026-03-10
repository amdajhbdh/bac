# Vault MCP Server Guide

## Overview
The Vault MCP Server provides programmatic access to your study vault through the Model Context Protocol.

## Installation

### 1. Configure Kiro CLI
Add to your Kiro CLI config (`~/.config/kiro/config.json`):
```json
{
  "mcpServers": {
    "vault": {
      "command": "python3",
      "args": ["/home/med/Documents/bac/resources/notes/vault-mcp-server.py"],
      "env": {
        "VAULT_PATH": "/home/med/Documents/bac/resources/notes/Study-Vault"
      }
    }
  }
}
```

### 2. Test Connection
```bash
python3 /home/med/Documents/bac/resources/notes/vault-mcp-server.py test
```

## Available Methods

### vault/search
Search notes with full-text search
```json
{
  "method": "vault/search",
  "params": {
    "query": "complex numbers",
    "limit": 10,
    "filters": {
      "subject": "Math",
      "type": "concept"
    }
  }
}
```

### vault/get
Get note content and metadata
```json
{
  "method": "vault/get",
  "params": {
    "path": "01-Concepts/Math/Complex-Numbers/Complex Numbers - Basics.md"
  }
}
```

### vault/related
Find related notes via links
```json
{
  "method": "vault/related",
  "params": {
    "path": "01-Concepts/Math/Complex-Numbers/Complex Numbers - Basics.md",
    "limit": 5
  }
}
```

### vault/stats
Get vault statistics
```json
{
  "method": "vault/stats",
  "params": {}
}
```

### vault/orphans
Find notes with no links
```json
{
  "method": "vault/orphans",
  "params": {}
}
```

### vault/formulas
Search mathematical formulas
```json
{
  "method": "vault/formulas",
  "params": {
    "query": "integral",
    "limit": 10
  }
}
```

### vault/rebuild
Rebuild search index
```json
{
  "method": "vault/rebuild",
  "params": {}
}
```

## Integration with Smart Connections

The MCP server works alongside Obsidian's Smart Connections plugin:

1. **Smart Connections**: Provides semantic similarity using embeddings
2. **MCP Server**: Provides structured queries and metadata access
3. **Together**: Agents can use both for comprehensive understanding

## Agent Usage

Agents can now:
- Search vault content programmatically
- Understand note relationships
- Access metadata and statistics
- Find orphaned content
- Search formulas
- Rebuild indexes automatically

## Example Agent Workflow

```python
# Agent searches for related concepts
response = mcp_call("vault/search", {"query": "derivatives", "limit": 5})

# Agent gets full note content
note = mcp_call("vault/get", {"path": response['result'][0]['path']})

# Agent finds related notes
related = mcp_call("vault/related", {"path": note['path']})

# Agent can now understand the full context
```

## Performance

- Search: <100ms
- Get note: <10ms
- Related notes: <50ms
- Stats: <20ms

## Next Steps

1. Install Smart Connections plugin in Obsidian
2. Generate embeddings for all notes
3. Agents will have both:
   - Structured access (MCP)
   - Semantic understanding (Smart Connections)

---

*Part of Vault Agent Ecosystem*
