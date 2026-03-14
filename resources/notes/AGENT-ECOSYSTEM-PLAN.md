# AI Agent Ecosystem Implementation Plan
## Complete Study Vault Enhancement System

**Created:** 2026-03-08  
**Location:** `/home/med/Documents/bac/resources/notes/`  
**Vault:** Study-Vault (Math, Physics, Chemistry, Natural Sciences)

---
## Problem Statement

Build a comprehensive, autonomous AI agent system that continuously monitors, analyzes, and improves study vaults with full integration of:
- **Ironclaw** (Telegram control)
- **OpenCode/Kilocode** (code analysis)
- **Kiro CLI** (orchestration, hooks, MCP servers)
- **Cloudflare Workers** (stateful agents, AI processing)

The system must handle local + cloud hybrid architecture, process PDFs/images with OCR, discover web resources, build knowledge graphs, and provide multi-modal interaction.

---

## Requirements Summary

### Core Capabilities
- ✅ Complete autonomous study assistant (content enhancement, knowledge graphs, study optimization, resource discovery)
- ✅ Hybrid automation (minor changes auto-applied, major changes need approval)
- ✅ Local MCP servers + Cloudflare Workers for heavy processing
- ✅ Full interaction flexibility (chat, scheduled reports, on-demand)
- ✅ Process everything: local files + web + study progress + generate problems + create diagrams
- ✅ Specialized agent team: Content Scanner, Resource Finder, Study Coach, Link Builder
- ✅ Ironclaw manages all Telegram bot interactions
- ✅ Complete image processing pipeline (OCR, diagram recognition, auto-annotation, problem solving)

### Technology Stack
- Cloudflare Workers + Durable Objects (stateful agents)
- Cloudflare Workers AI (image analysis, OCR, content generation)
- Local MCP servers (vault operations, file processing)
- Ironclaw/OpenClaw (Telegram gateway)
- Agents SDK (orchestration, state management, scheduling)
- Kiro CLI (hooks, MCP integration, agent management)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         USER INTERFACES                          │
├─────────────────────────────────────────────────────────────────┤
│  Telegram (Ironclaw) │ Kiro CLI │ WebSocket Chat │ Scheduled    │
└──────────┬───────────┴──────┬───┴────────┬───────┴──────────────┘
           │                  │            │
           ▼                  ▼            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    COORDINATOR AGENT (DO)                        │
│  - Orchestrates sub-agents                                       │
│  - Manages approval queue                                        │
│  - Handles scheduling                                            │
│  - Maintains knowledge graph                                     │
└──────────┬──────────────────────────────────────────────────────┘
           │
           ├──────────┬──────────┬──────────┬──────────┐
           ▼          ▼          ▼          ▼          ▼
    ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
    │ Content  │ │ Resource │ │  Study   │ │   Link   │
    │ Scanner  │ │  Finder  │ │  Coach   │ │ Builder  │
    │  Agent   │ │  Agent   │ │  Agent   │ │  Agent   │
    └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘
         │            │            │            │
         ├────────────┴────────────┴────────────┤
         ▼                                      ▼
┌─────────────────────┐              ┌─────────────────────┐
│   Workers AI        │              │   Vault MCP Server  │
│  - OCR/Vision       │              │  - File operations  │
│  - Text generation  │              │  - Search/index     │
│  - Embeddings       │              │  - Link management  │
└─────────────────────┘              └──────────┬──────────┘
                                                │
                                                ▼
                                     ┌─────────────────────┐
                                     │    Study Vault      │
                                     │  - Math, Physics    │
                                     │  - Chemistry, Bio   │
                                     │  - 1200+ problems   │
                                     └─────────────────────┘
```

---

## Component Breakdown

### 1. Ironclaw Telegram Gateway
- Handles all Telegram Bot API interactions
- Routes commands to Coordinator Agent webhook
- Manages file uploads (PDFs, images)
- Sends reports and notifications

### 2. Coordinator Agent (Durable Object)
- Main orchestrator for all sub-agents
- Maintains persistent state across restarts
- Handles scheduling (daily/hourly tasks)
- Manages approval queue for changes
- Builds and maintains knowledge graph
- Integrates with Vault MCP Server

### 3. Content Scanner Agent
- PDF text extraction and parsing
- Image OCR using Workers AI vision models
- Concept extraction using LLM
- Diagram recognition and annotation
- Stores processed content via Vault MCP

### 4. Resource Finder Agent
- Web search integration (YouTube, educational sites)
- Resource quality scoring and filtering
- Categorization (videos, websites, tools)
- Auto-generates Resource Links entries

### 5. Study Coach Agent
- Progress tracking with Durable Object storage
- Study path recommendation algorithm
- Quiz/practice problem generator
- Spaced repetition scheduler

### 6. Link Builder Agent
- Graph analysis of vault structure
- Automatic cross-reference detection
- MOC (Map of Content) generator
- Missing link identifier

### 7. Vault MCP Server (Local)
- File operations (read, write, search)
- Markdown manipulation
- Link creation and management
- File watching for real-time updates

### 8. Workers AI Integration
- Image analysis and OCR
- Text generation for content
- Embeddings for semantic search

---

## Kiro CLI Integration

### .kiro/ Directory Structure

```
.kiro/
├── agents/
│   ├── coordinator.json           # Main orchestrator config
│   ├── vault-assistant.json       # Interactive vault agent
│   ├── scanner.json                # Content scanning agent
│   └── study-coach.json            # Study guidance agent
│
├── mcp-servers/
│   ├── vault-server/
│   │   ├── package.json
│   │   ├── index.js                # Vault MCP server
│   │   └── tools/
│   │       ├── read-note.js
│   │       ├── write-note.js
│   │       ├── search-content.js
│   │       └── create-link.js
│   │
│   ├── image-processor/
│   │   ├── package.json
│   │   └── index.js                # Image OCR MCP server
│   │
│   └── web-search/
│       ├── package.json
│       └── index.js                # Resource finder MCP server
│
├── hooks/
│   ├── validate-write.sh           # Pre-write validation
│   ├── post-scan.sh                # Post-scan processing
│   ├── daily-report.sh             # Daily summary generator
│   └── sync-vault.sh               # Git sync after changes
│
├── skills/
│   ├── vault-management/
│   │   └── SKILL.md                # Vault operations skill
│   │
│   ├── study-optimization/
│   │   └── SKILL.md                # Study path optimization
│   │
│   └── resource-curation/
│       └── SKILL.md                # Resource discovery skill
│
├── prompts/
│   ├── coordinator-prompt.txt      # Coordinator system prompt
│   ├── scanner-prompt.txt          # Scanner instructions
│   └── coach-prompt.txt            # Study coach personality
│
├── config/
│   ├── cloudflare.toml             # Wrangler config
│   ├── telegram.json               # Ironclaw config
│   └── vault-paths.json            # Vault locations
│
├── scripts/
│   ├── deploy-agents.sh            # Deploy all agents
│   ├── setup-mcp.sh                # Install MCP servers
│   └── test-integration.sh         # Integration tests
│
└── extensions/
    ├── vault-enhancer/
    │   ├── extension.json          # Extension manifest
    │   ├── index.js                # Main entry point
    │   └── tools/                  # Custom tools
    │
    ├── agent-ui/
    │   ├── extension.json
    │   ├── web/                    # Web interface
    │   └── cli/                    # CLI enhancements
    │
    ├── integrations/
    │   ├── anki-sync/              # Anki flashcard sync
    │   ├── notion-export/          # Notion export
    │   └── github-backup/          # GitHub backup
    │
    └── registry.json               # Extension registry
```

---

## Kiro Extension System

### Extension Architecture

Since Kiro doesn't have native plugin support, we'll build an extension system using:
- **MCP Servers** as the extension runtime
- **Hooks** for lifecycle management
- **Agent Configs** for extension activation
- **Custom Tools** for extension functionality

### Extension Manifest Format

```json
{
  "name": "vault-enhancer",
  "version": "1.0.0",
  "description": "Enhanced vault operations with AI",
  "author": "Your Name",
  "type": "mcp-extension",
  "capabilities": [
    "vault-enhancement",
    "agent-tools",
    "ui-components",
    "integrations"
  ],
  "mcpServer": {
    "command": "node",
    "args": ["extensions/vault-enhancer/index.js"]
  },
  "tools": [
    "smart_search",
    "auto_summarize",
    "generate_quiz",
    "create_flashcards"
  ],
  "hooks": {
    "install": "extensions/vault-enhancer/install.sh",
    "uninstall": "extensions/vault-enhancer/uninstall.sh"
  },
  "dependencies": {
    "@cloudflare/workers-ai": "^1.0.0",
    "obsidian-api": "^2.0.0"
  },
  "config": {
    "vaultPath": "/home/med/Documents/bac/resources/notes/Study-Vault",
    "aiModel": "@cf/meta/llama-3-8b-instruct"
  }
}
```

### Extension Types

**1. Vault Enhancement Extensions**
- Smart search with semantic understanding
- Auto-summarization of notes
- Quiz generation from content
- Flashcard creation
- Diagram auto-generation
- LaTeX equation rendering

**2. Agent Capability Extensions**
- Custom tools for specialized tasks
- New MCP servers for external services
- Enhanced reasoning capabilities
- Multi-modal processing (audio, video)

**3. UI/UX Extensions**
- Web dashboard for agent monitoring
- CLI enhancements (better formatting, colors)
- Progress bars and status indicators
- Interactive prompts and menus
- Keyboard shortcuts and hotkeys

**4. Integration Extensions**
- Anki flashcard sync
- Notion export/import
- GitHub backup automation
- Google Drive sync
- Obsidian plugin bridge
- Telegram bot enhancements
- Discord integration
- Slack notifications

### Extension Registry

```json
{
  "extensions": [
    {
      "id": "vault-enhancer",
      "enabled": true,
      "autoload": true,
      "priority": 1
    },
    {
      "id": "agent-ui",
      "enabled": true,
      "autoload": true,
      "priority": 2
    },
    {
      "id": "anki-sync",
      "enabled": false,
      "autoload": false,
      "priority": 3
    }
  ]
}
```

### Extension Loader

Create `.kiro/extensions/loader.js`:

```javascript
// Extension loader for Kiro
const fs = require('fs');
const path = require('path');

class ExtensionManager {
  constructor(extensionsDir) {
    this.extensionsDir = extensionsDir;
    this.registry = this.loadRegistry();
    this.extensions = new Map();
  }

  loadRegistry() {
    const registryPath = path.join(this.extensionsDir, 'registry.json');
    return JSON.parse(fs.readFileSync(registryPath, 'utf8'));
  }

  async loadExtension(extensionId) {
    const manifest = this.loadManifest(extensionId);
    const extension = require(path.join(this.extensionsDir, extensionId, 'index.js'));
    
    // Initialize extension
    await extension.init(manifest.config);
    
    this.extensions.set(extensionId, {
      manifest,
      instance: extension
    });
    
    return extension;
  }

  loadManifest(extensionId) {
    const manifestPath = path.join(this.extensionsDir, extensionId, 'extension.json');
    return JSON.parse(fs.readFileSync(manifestPath, 'utf8'));
  }

  async loadAll() {
    const enabled = this.registry.extensions
      .filter(ext => ext.enabled && ext.autoload)
      .sort((a, b) => a.priority - b.priority);
    
    for (const ext of enabled) {
      await this.loadExtension(ext.id);
    }
  }

  getExtension(extensionId) {
    return this.extensions.get(extensionId);
  }

  listExtensions() {
    return Array.from(this.extensions.keys());
  }
}

module.exports = ExtensionManager;
```

---

## Extension Examples

### Extension 1: Vault Enhancer

**Purpose:** AI-powered vault operations

**Tools:**
- `@vault-enhancer/smart_search` - Semantic search across vault
- `@vault-enhancer/auto_summarize` - Summarize long notes
- `@vault-enhancer/generate_quiz` - Create practice problems
- `@vault-enhancer/create_flashcards` - Generate Anki cards
- `@vault-enhancer/suggest_links` - Find related notes
- `@vault-enhancer/extract_concepts` - Identify key concepts

**Implementation:** `.kiro/extensions/vault-enhancer/index.js`

```javascript
const { WorkersAI } = require('@cloudflare/workers-ai');

class VaultEnhancer {
  async init(config) {
    this.vaultPath = config.vaultPath;
    this.ai = new WorkersAI({ apiKey: process.env.CF_API_KEY });
  }

  async smart_search({ query, limit = 10 }) {
    // Semantic search using embeddings
    const embedding = await this.ai.run('@cf/baai/bge-base-en-v1.5', {
      text: query
    });
    
    // Search vault with embeddings
    const results = await this.searchByEmbedding(embedding, limit);
    return results;
  }

  async auto_summarize({ notePath, maxLength = 200 }) {
    const content = fs.readFileSync(notePath, 'utf8');
    
    const summary = await this.ai.run('@cf/meta/llama-3-8b-instruct', {
      prompt: `Summarize this note in ${maxLength} words:\n\n${content}`
    });
    
    return summary;
  }

  async generate_quiz({ topic, difficulty = 'medium', count = 5 }) {
    const notes = await this.findNotesByTopic(topic);
    const content = notes.map(n => n.content).join('\n\n');
    
    const quiz = await this.ai.run('@cf/meta/llama-3-8b-instruct', {
      prompt: `Generate ${count} ${difficulty} quiz questions about:\n\n${content}`
    });
    
    return quiz;
  }

  async create_flashcards({ notePath }) {
    const content = fs.readFileSync(notePath, 'utf8');
    
    const flashcards = await this.ai.run('@cf/meta/llama-3-8b-instruct', {
      prompt: `Create Anki flashcards from this content:\n\n${content}\n\nFormat: Q: question\nA: answer`
    });
    
    return this.parseFlashcards(flashcards);
  }
}

module.exports = new VaultEnhancer();
```

### Extension 2: Agent UI Dashboard

**Purpose:** Web interface for monitoring agents

**Features:**
- Real-time agent status
- Task queue visualization
- Progress tracking
- Manual triggers
- Approval interface

**Implementation:** `.kiro/extensions/agent-ui/index.js`

```javascript
const express = require('express');
const WebSocket = require('ws');

class AgentUI {
  async init(config) {
    this.app = express();
    this.wss = new WebSocket.Server({ port: 8080 });
    
    this.setupRoutes();
    this.setupWebSocket();
    
    this.app.listen(3000, () => {
      console.log('Agent UI running at http://localhost:3000');
    });
  }

  setupRoutes() {
    this.app.get('/api/agents', (req, res) => {
      res.json(this.getAgentStatus());
    });
    
    this.app.post('/api/trigger/:agent/:action', (req, res) => {
      this.triggerAgent(req.params.agent, req.params.action);
      res.json({ success: true });
    });
    
    this.app.get('/api/queue', (req, res) => {
      res.json(this.getTaskQueue());
    });
  }

  setupWebSocket() {
    this.wss.on('connection', (ws) => {
      ws.on('message', (message) => {
        const data = JSON.parse(message);
        this.handleCommand(data, ws);
      });
    });
  }
}

module.exports = new AgentUI();
```

### Extension 3: Anki Sync Integration

**Purpose:** Sync flashcards to Anki

**Tools:**
- `@anki-sync/export_deck` - Export vault notes as Anki deck
- `@anki-sync/sync_cards` - Sync flashcards to AnkiConnect
- `@anki-sync/import_progress` - Import study progress from Anki

**Implementation:** `.kiro/extensions/integrations/anki-sync/index.js`

```javascript
const axios = require('axios');

class AnkiSync {
  async init(config) {
    this.ankiConnectUrl = config.ankiConnectUrl || 'http://localhost:8765';
    this.deckName = config.deckName || 'Study Vault';
  }

  async export_deck({ topic }) {
    const flashcards = await this.generateFlashcards(topic);
    
    for (const card of flashcards) {
      await this.addNote({
        deckName: this.deckName,
        modelName: 'Basic',
        fields: {
          Front: card.question,
          Back: card.answer
        },
        tags: [topic, 'vault-generated']
      });
    }
    
    return { exported: flashcards.length };
  }

  async addNote(note) {
    const response = await axios.post(this.ankiConnectUrl, {
      action: 'addNote',
      version: 6,
      params: { note }
    });
    
    return response.data;
  }
}

module.exports = new AnkiSync();
```

### Extension 4: GitHub Backup

**Purpose:** Automated vault backup to GitHub

**Tools:**
- `@github-backup/commit` - Commit vault changes
- `@github-backup/push` - Push to remote
- `@github-backup/auto_backup` - Scheduled backups

**Implementation:** `.kiro/extensions/integrations/github-backup/index.js`

```javascript
const simpleGit = require('simple-git');

class GitHubBackup {
  async init(config) {
    this.vaultPath = config.vaultPath;
    this.git = simpleGit(this.vaultPath);
  }

  async commit({ message = 'Auto-backup from agents' }) {
    await this.git.add('.');
    await this.git.commit(message);
    return { success: true };
  }

  async push() {
    await this.git.push('origin', 'main');
    return { success: true };
  }

  async auto_backup() {
    const status = await this.git.status();
    
    if (status.files.length > 0) {
      await this.commit({ message: `Auto-backup: ${status.files.length} files changed` });
      await this.push();
      return { backed_up: status.files.length };
    }
    
    return { backed_up: 0 };
  }
}

module.exports = new GitHubBackup();
```

### Extension 5: Obsidian Plugin Bridge

**Purpose:** Control Obsidian from agents

**Tools:**
- `@obsidian/open_note` - Open note in Obsidian
- `@obsidian/create_diagram` - Create Excalidraw diagram
- `@obsidian/run_dataview` - Execute Dataview query
- `@obsidian/update_graph` - Refresh graph view

**Implementation:** `.kiro/extensions/integrations/obsidian-bridge/index.js`

```javascript
const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

class ObsidianBridge {
  async init(config) {
    this.vaultPath = config.vaultPath;
  }

  async open_note({ notePath }) {
    const obsidianUri = `obsidian://open?vault=Study-Vault&file=${encodeURIComponent(notePath)}`;
    await execPromise(`xdg-open "${obsidianUri}"`);
    return { success: true };
  }

  async create_diagram({ name, content }) {
    const diagramPath = `${this.vaultPath}/Assets/Diagrams/${name}.excalidraw`;
    fs.writeFileSync(diagramPath, JSON.stringify(content));
    return { path: diagramPath };
  }

  async run_dataview({ query }) {
    // Execute Dataview query via Obsidian API
    // This requires Obsidian Local REST API plugin
    const response = await axios.post('http://localhost:27123/vault/dataview', {
      query
    });
    return response.data;
  }
}

module.exports = new ObsidianBridge();
```

---

## Task 16: Build Extension System

**Objective:** Create extensible architecture for custom Kiro extensions

**Implementation:**
1. Create extension loader in `.kiro/extensions/loader.js`
2. Define extension manifest format
3. Build extension registry
4. Create CLI commands for extension management
5. Implement extension lifecycle (install, enable, disable, uninstall)

**Extension Manager CLI:**

```bash
# List extensions
kiro-ext list

# Install extension
kiro-ext install vault-enhancer

# Enable extension
kiro-ext enable vault-enhancer

# Disable extension
kiro-ext disable anki-sync

# Uninstall extension
kiro-ext uninstall notion-export

# Create new extension
kiro-ext create my-extension
```

**Kiro Integration:**
- Extensions load as MCP servers
- Tools prefixed with `@extension-name/tool-name`
- Hooks for extension lifecycle
- Agent configs reference extension tools

**Test:** Install vault-enhancer, verify tools available  
**Demo:** Use `@vault-enhancer/smart_search` to find related notes

---

## Task 17: Build Core Extensions

**Objective:** Implement 5 essential extensions

**Extensions to Build:**
1. **Vault Enhancer** - AI-powered vault operations
2. **Agent UI** - Web dashboard for monitoring
3. **Anki Sync** - Flashcard integration
4. **GitHub Backup** - Automated backups
5. **Obsidian Bridge** - Control Obsidian from agents

**Test:** All 5 extensions installed and functional  
**Demo:** Complete workflow using all extensions

---

## Updated Task Breakdown

### **Task 1: Set up Vault MCP Server**

**Objective:** Create local MCP server for vault file operations

**Implementation:**
1. Create Node.js MCP server in `.kiro/mcp-servers/vault-server/`
2. Implement tools:
   - `read_note` - Read markdown files
   - `write_note` - Create/update notes
   - `list_topics` - List vault structure
   - `search_content` - Full-text search
   - `create_link` - Add cross-references
3. Add file watching for real-time updates
4. Configure in agent JSON:
```json
{
  "mcpServers": {
    "vault": {
      "command": "node",
      "args": [".kiro/mcp-servers/vault-server/index.js"],
      "env": {
        "VAULT_PATH": "/home/med/Documents/bac/resources/notes/Study-Vault"
      }
    }
  }
}
```

**Kiro Integration:**
- Use `/mcp` to verify server loaded
- Add to `vault-assistant.json` agent config
- Set `allowedTools: ["@vault/*"]` for auto-approval

**Test:** `kiro-cli chat` → `/mcp` → verify vault server initialized  
**Demo:** Ask "What topics are in my vault?" → agent uses `@vault/list_topics`

---

### **Task 2: Build Coordinator Agent on Cloudflare**

**Objective:** Deploy main orchestrator agent with state management

**Implementation:**
1. Create Cloudflare Worker with Agents SDK
2. Extend `Agent` class with:
   - State management for sub-agent tracking
   - RPC methods: `triggerScan`, `getReport`, `queryVault`, `addContent`
   - Scheduling for daily/hourly tasks
   - MCP client integration
3. Deploy with wrangler:
```bash
cd .kiro/agents/coordinator
wrangler deploy
```

**Kiro Integration:**
- Create `.kiro/agents/coordinator.json`:
```json
{
  "name": "coordinator",
  "description": "Main orchestrator for vault agents",
  "prompt": "file://.kiro/prompts/coordinator-prompt.txt",
  "tools": ["fs_read", "execute_bash"],
  "mcpServers": {
    "vault": {
      "command": "node",
      "args": [".kiro/mcp-servers/vault-server/index.js"]
    }
  },
  "hooks": {
    "agentSpawn": [{
      "command": "date",
      "description": "Log coordinator start time"
    }]
  },
  "keyboardShortcut": "ctrl+shift+c",
  "welcomeMessage": "Coordinator agent ready. Managing vault operations."
}
```

**Test:** Deploy agent, verify state persistence across restarts  
**Demo:** Call RPC methods via wrangler, see agent respond with vault data

---

### **Task 3: Implement Content Scanner Agent**

**Objective:** PDF/image processing with OCR and concept extraction

**Implementation:**
1. Create Durable Object for scanner with queue system
2. Integrate Workers AI for image OCR:
```typescript
const result = await env.AI.run('@cf/meta/llama-3.2-11b-vision-instruct', {
  image: imageData,
  prompt: 'Extract all text and equations from this image'
});
```
3. Add PDF text extraction
4. Implement concept extraction using LLM
5. Store via Vault MCP

**Kiro Integration:**
- Create `.kiro/agents/scanner.json`:
```json
{
  "name": "scanner",
  "description": "Scans PDFs and images, extracts content",
  "tools": ["fs_read"],
  "allowedTools": ["@vault/write_note", "@vault/create_link"],
  "mcpServers": {
    "vault": { "command": "node", "args": [".kiro/mcp-servers/vault-server/index.js"] },
    "image": { "command": "node", "args": [".kiro/mcp-servers/image-processor/index.js"] }
  },
  "hooks": {
    "postToolUse": [{
      "matcher": "@image/ocr",
      "command": ".kiro/hooks/post-scan.sh",
      "description": "Process OCR results"
    }]
  }
}
```

**Test:** Upload PDF/image, verify OCR extraction and note creation  
**Demo:** Scan textbook page → show generated markdown note in vault

---

### **Task 4: Build Resource Finder Agent**

**Objective:** Web search and resource curation

**Implementation:**
1. Create MCP server for web search in `.kiro/mcp-servers/web-search/`
2. Implement tools:
   - `search_youtube` - Find educational videos
   - `search_khan_academy` - Find Khan Academy content
   - `search_educational_sites` - Search curated sites
3. Add resource quality scoring
4. Auto-generate Resource Links entries

**Kiro Integration:**
- Create `.kiro/agents/resource-finder.json`:
```json
{
  "name": "resource-finder",
  "description": "Finds and curates educational resources",
  "tools": ["fs_read"],
  "allowedTools": ["@web/*", "@vault/write_note"],
  "mcpServers": {
    "web": { "command": "node", "args": [".kiro/mcp-servers/web-search/index.js"] },
    "vault": { "command": "node", "args": [".kiro/mcp-servers/vault-server/index.js"] }
  },
  "resources": [
    "file://Study-Vault/Resource Links.md"
  ]
}
```

**Test:** Search for "complex numbers tutorial", verify curated results  
**Demo:** Agent finds 5 relevant videos and adds them to vault

---

### **Task 5: Create Study Coach Agent**

**Objective:** Progress tracking and personalized recommendations

**Implementation:**
1. Create Durable Object for progress storage
2. Implement study path recommendation algorithm
3. Build quiz generator using Workers AI
4. Create spaced repetition scheduler

**Kiro Integration:**
- Create `.kiro/agents/study-coach.json`:
```json
{
  "name": "study-coach",
  "description": "Tracks progress and provides study guidance",
  "prompt": "file://.kiro/prompts/coach-prompt.txt",
  "tools": ["fs_read"],
  "allowedTools": ["@vault/read_note", "@vault/search_content"],
  "mcpServers": {
    "vault": { "command": "node", "args": [".kiro/mcp-servers/vault-server/index.js"] }
  },
  "hooks": {
    "userPromptSubmit": [{
      "command": ".kiro/hooks/track-progress.sh",
      "description": "Log study session"
    }]
  },
  "keyboardShortcut": "ctrl+shift+s",
  "welcomeMessage": "What would you like to study today?"
}
```

**Test:** Track completion of 3 topics, verify recommendations update  
**Demo:** Agent suggests next topic and generates 5 practice problems

---

### **Task 6: Develop Link Builder Agent**

**Objective:** Knowledge graph construction and cross-referencing

**Implementation:**
1. Implement graph analysis of vault structure
2. Add automatic cross-reference detection
3. Create MOC generator
4. Build missing link identifier

**Kiro Integration:**
- Create `.kiro/agents/link-builder.json`:
```json
{
  "name": "link-builder",
  "description": "Builds knowledge graph and cross-references",
  "tools": ["fs_read", "grep"],
  "allowedTools": ["@vault/*"],
  "mcpServers": {
    "vault": { "command": "node", "args": [".kiro/mcp-servers/vault-server/index.js"] }
  },
  "hooks": {
    "stop": [{
      "command": ".kiro/hooks/update-graph.sh",
      "description": "Update Obsidian graph.json"
    }]
  }
}
```

**Test:** Analyze vault, verify link suggestions are relevant  
**Demo:** Agent creates 10 new cross-references between related notes

---

### **Task 7: Integrate Ironclaw Telegram Gateway**

**Objective:** Enable Telegram control of all agents

**Implementation:**
1. Install OpenClaw/Ironclaw
2. Configure Telegram Bot API token
3. Set up webhook routing to Coordinator Agent
4. Implement command handlers:
   - `/scan` - Trigger content scanner
   - `/report` - Get daily summary
   - `/query <question>` - Ask vault question
   - `/add` - Upload file for processing
5. Add file upload handling

**Kiro Integration:**
- Create `.kiro/config/telegram.json`:
```json
{
  "botToken": "$TELEGRAM_BOT_TOKEN",
  "webhookUrl": "https://coordinator.your-domain.workers.dev/telegram",
  "allowedUsers": ["your-telegram-id"],
  "commands": {
    "scan": "triggerScan",
    "report": "getReport",
    "query": "queryVault",
    "add": "addContent"
  }
}
```

**Test:** Send Telegram message, verify agent responds  
**Demo:** Full conversation via Telegram - ask question, get vault answer

---

### **Task 8: Add Chat Interface to Coordinator**

**Objective:** Real-time WebSocket chat with context management

**Implementation:**
1. Implement WebSocket connection in Coordinator
2. Add conversational context management
3. Integrate with all sub-agents
4. Create response formatting for different channels

**Kiro Integration:**
- Use Kiro CLI as primary chat interface
- Switch agents with keyboard shortcuts:
  - `Ctrl+Shift+C` - Coordinator
  - `Ctrl+Shift+S` - Study Coach
  - `Ctrl+Shift+V` - Vault Assistant
- Use `/agent` to switch between specialized agents

**Test:** Chat session maintains context across 5 messages  
**Demo:** Ask "What should I study next?" and get personalized answer

---

### **Task 9: Build Scheduled Reporting System**

**Objective:** Automated daily summaries and progress reports

**Implementation:**
1. Add daily summary generation in Coordinator
2. Implement progress reports with statistics
3. Create suggestion compilation from all agents
4. Configure delivery via Telegram

**Kiro Integration:**
- Create `.kiro/hooks/daily-report.sh`:
```bash
#!/bin/bash
# Generate daily study report
REPORT=$(curl -X POST https://coordinator.workers.dev/rpc/getReport)
echo "$REPORT" | jq -r '.summary'
```
- Add to coordinator agent:
```json
{
  "hooks": {
    "agentSpawn": [{
      "command": ".kiro/hooks/daily-report.sh",
      "description": "Generate daily report",
      "cache_ttl_seconds": 86400
    }]
  }
}
```

**Test:** Trigger scheduled task, verify report generated  
**Demo:** Receive daily report with vault stats and recommendations

---

### **Task 10: Create On-Demand Trigger System**

**Objective:** Manual agent triggers via CLI and Telegram

**Implementation:**
1. Implement CLI commands in Kiro
2. Add Telegram commands for specific actions
3. Build batch processing queue
4. Create status tracking for long-running tasks

**Kiro Integration:**
- Create custom slash commands (if supported) or use hooks
- Use `execute_bash` tool to trigger agent actions:
```bash
kiro-cli chat --agent scanner --prompt "Scan all PDFs in Assets/PDFs/"
```

**Test:** Trigger bulk PDF scan, monitor progress  
**Demo:** Process 10 PDFs via CLI command, see progress updates

---

### **Task 11: Implement Hybrid Approval Workflow**

**Objective:** Auto-apply minor changes, require approval for major changes

**Implementation:**
1. Add change classification logic
2. Create approval queue in Coordinator state
3. Build review interface via Telegram
4. Implement auto-apply rules:
   - Auto: formatting, links, tags
   - Manual: content changes, deletions

**Kiro Integration:**
- Use `preToolUse` hooks for validation:
```json
{
  "hooks": {
    "preToolUse": [{
      "matcher": "@vault/write_note",
      "command": ".kiro/hooks/validate-write.sh",
      "description": "Validate note changes"
    }]
  }
}
```
- Hook script returns exit code 2 to block major changes

**Test:** Agent suggests major change, verify approval required  
**Demo:** Review and approve/reject 3 suggested changes via Telegram

---

### **Task 12: Add Image Analysis Pipeline**

**Objective:** Complete image processing with OCR, diagram recognition, annotation

**Implementation:**
1. Integrate Workers AI vision models
2. Implement equation extraction
3. Add auto-annotation with explanations
4. Create Excalidraw diagram generator

**Kiro Integration:**
- Create `.kiro/mcp-servers/image-processor/`:
```javascript
// MCP server for image processing
export const tools = {
  ocr_image: async ({ imagePath }) => {
    const result = await workersAI.run('@cf/meta/llama-3.2-11b-vision-instruct', {
      image: fs.readFileSync(imagePath),
      prompt: 'Extract text, equations, and describe diagrams'
    });
    return result;
  },
  recognize_diagram: async ({ imagePath }) => {
    // Diagram recognition logic
  }
};
```

**Test:** Upload physics diagram, verify recognition and annotation  
**Demo:** Image of equation → extracted LaTeX + explanation + linked concepts

---

### **Task 13: Build Knowledge Graph Visualization**

**Objective:** Interactive graph of vault connections

**Implementation:**
1. Implement graph data structure in Coordinator
2. Add Obsidian graph.json generator
3. Create interactive visualization endpoint
4. Build graph query API

**Kiro Integration:**
- Use `stop` hook to update graph after changes:
```json
{
  "hooks": {
    "stop": [{
      "command": ".kiro/hooks/update-graph.sh",
      "description": "Regenerate knowledge graph"
    }]
  }
}
```

**Test:** Generate graph from vault, verify connections accurate  
**Demo:** View knowledge graph showing all topic relationships

---

### **Task 14: Integrate OpenCode/Kilocode**

**Objective:** Code analysis for programming notes

**Implementation:**
1. Connect OpenCode for code analysis
2. Add Kilocode for advanced processing
3. Implement code example generation
4. Create code snippet extraction

**Kiro Integration:**
- Add as MCP servers if available
- Or use as external tools via `execute_bash`
- Create skill for code-related operations

**Test:** Analyze Python code in notes, verify suggestions  
**Demo:** Agent generates code examples for algorithms topic

---

### **Task 15: Create Comprehensive .kiro/ Structure**

**Objective:** Organize all components in .kiro/ directory

**Implementation:**
1. Organize agent configs in `.kiro/agents/`
2. Set up MCP servers in `.kiro/mcp-servers/`
3. Create hooks in `.kiro/hooks/`
4. Add skills in `.kiro/skills/`
5. Document setup and usage

**Kiro Integration:**
- Create `.kiro/README.md` with setup instructions
- Add `.kiro/scripts/setup.sh` for automated installation
- Create `.kiro/scripts/deploy.sh` for Cloudflare deployment
- Add `.kiro/config/` for all configuration files

**Structure:**
```
.kiro/
├── README.md                       # Setup guide
├── agents/                         # Agent configurations
├── mcp-servers/                    # Local MCP servers
├── hooks/                          # Hook scripts
├── skills/                         # Skill definitions
├── prompts/                        # System prompts
├── config/                         # Configuration files
└── scripts/                        # Deployment scripts
```

**Test:** Fresh install from .kiro/ directory works  
**Demo:** Complete system running from organized .kiro/ structure

---

## Kiro-Specific Features

### Hooks Configuration

**agentSpawn Hooks:**
- Log agent start time
- Check vault status
- Display available topics

**userPromptSubmit Hooks:**
- Track study sessions
- Log user questions
- Update progress metrics

**preToolUse Hooks:**
- Validate file writes
- Check for conflicts
- Require approval for major changes

**postToolUse Hooks:**
- Process OCR results
- Update knowledge graph
- Sync with git

**stop Hooks:**
- Generate summaries
- Update statistics
- Commit changes

### MCP Server Integration

**Vault Server Tools:**
- `@vault/read_note` - Read markdown files
- `@vault/write_note` - Create/update notes
- `@vault/list_topics` - List vault structure
- `@vault/search_content` - Full-text search
- `@vault/create_link` - Add cross-references

**Image Processor Tools:**
- `@image/ocr` - Extract text from images
- `@image/recognize_diagram` - Identify diagrams
- `@image/extract_equations` - Extract LaTeX

**Web Search Tools:**
- `@web/search_youtube` - Find videos
- `@web/search_educational` - Find resources
- `@web/quality_score` - Rate resources

### Agent Keyboard Shortcuts

- `Ctrl+Shift+C` - Coordinator Agent
- `Ctrl+Shift+S` - Study Coach Agent
- `Ctrl+Shift+V` - Vault Assistant Agent
- `Ctrl+Shift+L` - Link Builder Agent

### Skills

**Vault Management Skill:**
```markdown
---
name: vault-management
description: Operations for managing study vault including reading, writing, searching, and linking notes. Use when user asks about vault content or wants to modify notes.
---
```

**Study Optimization Skill:**
```markdown
---
name: study-optimization
description: Analyzes study progress and recommends optimal learning paths. Use when user asks "what should I study" or needs study guidance.
---
```

**Resource Curation Skill:**
```markdown
---
name: resource-curation
description: Finds and evaluates educational resources from web. Use when user needs videos, websites, or practice problems for a topic.
---
```

---

## Deployment Steps

### 1. Local Setup
```bash
# Create .kiro structure
mkdir -p .kiro/{agents,mcp-servers,hooks,skills,prompts,config,scripts}

# Install MCP servers
cd .kiro/mcp-servers/vault-server && npm install
cd ../image-processor && npm install
cd ../web-search && npm install

# Make hooks executable
chmod +x .kiro/hooks/*.sh

# Configure agents
cp agent-templates/*.json .kiro/agents/
```

### 2. Cloudflare Deployment
```bash
# Deploy Coordinator Agent
cd .kiro/agents/coordinator
wrangler deploy

# Deploy sub-agents
cd ../scanner && wrangler deploy
cd ../resource-finder && wrangler deploy
cd ../study-coach && wrangler deploy
cd ../link-builder && wrangler deploy
```

### 3. Ironclaw Setup
```bash
# Install OpenClaw/Ironclaw
npm install -g openclaw

# Configure Telegram bot
openclaw config set telegram.token $TELEGRAM_BOT_TOKEN
openclaw config set telegram.webhook https://coordinator.workers.dev/telegram

# Start gateway
openclaw start
```

### 4. Kiro CLI Configuration
```bash
# Verify MCP servers
kiro-cli chat
/mcp

# Test agents
/agent
# Select vault-assistant

# Test hooks
/hooks
```

---

## Testing Plan

### Unit Tests
- [ ] Vault MCP server tools work correctly
- [ ] Image OCR extracts text accurately
- [ ] Resource finder returns quality results
- [ ] Study coach generates valid recommendations
- [ ] Link builder creates correct cross-references

### Integration Tests
- [ ] Coordinator orchestrates sub-agents correctly
- [ ] Telegram commands trigger correct actions
- [ ] Hooks execute at proper trigger points
- [ ] Approval workflow blocks/allows correctly
- [ ] Knowledge graph updates after changes

### End-to-End Tests
- [ ] Upload PDF via Telegram → OCR → note created in vault
- [ ] Ask "what should I study" → coach analyzes progress → suggests topic
- [ ] Request resources for topic → finder searches web → adds to vault
- [ ] Scan vault → link builder finds connections → creates cross-references
- [ ] Daily report generated → sent via Telegram

---

## Success Metrics

### Functionality
- ✅ All 4 specialized agents deployed and operational
- ✅ Vault MCP server handles 100+ operations/day
- ✅ Image OCR accuracy > 95%
- ✅ Resource finder adds 10+ quality resources/day
- ✅ Study coach provides personalized recommendations
- ✅ Link builder creates 20+ cross-references/day

### Performance
- ✅ Agent response time < 2 seconds
- ✅ OCR processing < 5 seconds per image
- ✅ Web search < 3 seconds per query
- ✅ Knowledge graph updates < 1 second

### User Experience
- ✅ Telegram bot responds within 1 second
- ✅ Kiro CLI agents switch instantly
- ✅ Approval workflow clear and intuitive
- ✅ Daily reports comprehensive and actionable

---

## Future Enhancements

### Phase 2
- Voice input via Telegram
- Mobile app integration
- Collaborative study features
- Gamification (points, achievements)

### Phase 3
- Multi-vault support
- Team collaboration
- Advanced analytics dashboard
- AI tutor with video explanations

### Phase 4
- Integration with learning platforms (Coursera, edX)
- Automated exam preparation
- Peer study matching
- Content marketplace

---

## Documentation

### User Guide
- Getting started with vault agents
- Using Telegram commands
- Kiro CLI keyboard shortcuts
- Approval workflow guide

### Developer Guide
- Adding new MCP tools
- Creating custom hooks
- Extending agents
- Deploying to Cloudflare

### API Reference
- Coordinator RPC methods
- Vault MCP server API
- Image processor API
- Web search API

---

## Support & Maintenance

### Monitoring
- Cloudflare Workers analytics
- MCP server health checks
- Telegram bot uptime
- Daily error reports

### Updates
- Weekly dependency updates
- Monthly feature releases
- Quarterly major versions
- Continuous security patches

### Backup
- Daily vault backups to git
- Weekly Cloudflare state exports
- Monthly full system snapshots

---

**Status:** Ready for implementation  
**Next Step:** Begin Task 1 - Set up Vault MCP Server  
**Estimated Timeline:** 2-3 weeks for full implementation
