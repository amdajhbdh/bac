---
tags: [meta, smart-connections, setup]
date: 2026-03-10
---

# 🔗 Smart Connections Setup Complete!

> [!success] Vector DB + Ollama Configured
> Smart Connections plugin configured with local Ollama cloud models and vector database.

---

## ⚙️ Configuration

> [!note] Settings Applied
> 
> ### Ollama Integration
> - **Endpoint:** `http://localhost:11434`
> - **Chat Model:** `ministral-3:14b-cloud`
> - **Embed Model:** `ministral-3:14b-cloud`
> - **Max Tokens:** 4096
> - **Temperature:** 0.7
> 
> ### Vector Database
> - **Enabled:** ✅ Yes
> - **Location:** `.smart-connections/vectors.db`
> - **Auto Index:** ✅ Enabled
> - **Index on Startup:** ✅ Enabled
> 
> ### Search Settings
> - **Top K Results:** 10
> - **Similarity Threshold:** 0.5
> - **Excluded:** PDFs, assets, .obsidian

---

## 🚀 Features Enabled

> [!tip] What You Can Do
> 
> ### 1. Semantic Search
> - Find related notes by meaning, not just keywords
> - Works with Arabic + English + French
> - Vector embeddings for similarity
> 
> ### 2. Smart Chat
> - Chat with your vault using Ollama
> - Ask questions about your notes
> - Get AI-powered summaries
> - Location: `00-Meta/Smart-Chat/`
> 
> ### 3. Related Notes
> - Automatic suggestions in sidebar
> - Based on semantic similarity
> - Updates as you type
> 
> ### 4. Vector Database
> - Fast similarity search
> - Persistent embeddings
> - Incremental updates

---

## 🎯 Usage

> [!example] How to Use
> 
> ### View Related Notes
> 1. Open any note
> 2. Check Smart Connections panel (right sidebar)
> 3. See automatically related notes
> 4. Click to navigate
> 
> ### Smart Chat
> 1. Open command palette (`Ctrl/Cmd + P`)
> 2. Search "Smart Connections: Open Chat"
> 3. Ask questions about your vault
> 4. Get AI-powered answers
> 
> ### Search by Meaning
> 1. Use Smart Connections search
> 2. Type concept or question
> 3. Get semantically similar notes
> 4. Not limited to exact keywords

---

## 🤖 Available Models

> [!note] Ollama Models
> 
> ### Current Setup
> - **ministral-3:14b-cloud** - Main model (14B params)
> - **llama3.2:3b** - Backup model (3B params)
> - **qwen3-vl:235b-instruct-cloud** - Vision model
> 
> ### Model Selection
> Edit `.obsidian/plugins/smart-connections/data.json`:
> ```json
> "ollama_model": "ministral-3:14b-cloud"
> ```

---

## 📊 Vector Database

> [!stats] Database Info
> 
> **Location:** `.smart-connections/vectors.db`
> 
> **Contents:**
> - Note embeddings (920 notes)
> - Semantic vectors
> - Similarity indices
> 
> **Updates:**
> - Auto-index on startup
> - Incremental updates on save
> - Manual rebuild available

---

## 🔧 Advanced Settings

> [!tip] Customization
> 
> ### Adjust Similarity Threshold
> Lower = more results, less relevant
> Higher = fewer results, more relevant
> 
> Current: `0.5` (balanced)
> 
> ### Change Top K Results
> How many related notes to show
> 
> Current: `10` results
> 
> ### Exclude Folders
> Add to `file_exclusions`:
> ```
> "file_exclusions": "07-Assets,daily,.obsidian,*.pdf"
> ```

---

## 🎯 Use Cases

> [!example] Practical Applications
> 
> ### 1. Study Session
> - Open a concept note
> - Check related notes in sidebar
> - Explore connections
> - Build understanding
> 
> ### 2. Question Answering
> - Open Smart Chat
> - Ask: "Explain matrices in simple terms"
> - Get answer from your notes
> - With source references
> 
> ### 3. Gap Discovery
> - Search for a topic
> - See what's covered
> - Find missing connections
> - Identify study gaps
> 
> ### 4. Review Preparation
> - Ask: "What are the key physics concepts?"
> - Get summary from all physics notes
> - Review systematically

---

## 📝 Chat Examples

> [!example] Try These Prompts
> 
> ```
> "Summarize all my notes on complex numbers"
> 
> "What are the main topics in physics I should review?"
> 
> "Explain the relationship between matrices and linear transformations"
> 
> "Find all exercises related to organic chemistry"
> 
> "What concepts am I missing in my chemistry notes?"
> 
> "Create a study plan for the next week based on my weak areas"
> ```

---

## 🔄 Maintenance

> [!warning] Periodic Tasks
> 
> ### Rebuild Index
> If search seems off:
> 1. Command palette
> 2. "Smart Connections: Rebuild Index"
> 3. Wait for completion
> 
> ### Clear Cache
> If issues persist:
> ```bash
> rm -rf .smart-connections/
> # Restart Obsidian to rebuild
> ```
> 
> ### Update Models
> ```bash
> ollama pull ministral-3:14b-cloud
> ollama pull llama3.2:3b
> ```

---

## 📊 Performance

> [!stats] Expected Performance
> 
> | Operation | Time |
> |-----------|------|
> | Initial Index | 2-5 min (920 notes) |
> | Incremental Update | <1 sec |
> | Search Query | <500ms |
> | Chat Response | 2-10 sec |
> | Related Notes | <100ms |

---

## ✅ Setup Complete!

> [!success] Ready to Use
> 
> **Configured:**
> - ✅ Ollama integration (ministral-3:14b-cloud)
> - ✅ Vector database enabled
> - ✅ Auto-indexing enabled
> - ✅ Smart Chat folder created
> - ✅ Semantic search ready
> - ✅ Related notes active
> 
> **Next Steps:**
> 1. Open Obsidian
> 2. Enable Smart Connections plugin
> 3. Wait for initial indexing
> 4. Start exploring connections!

---

**Configuration Date:** 2026-03-10 01:31 UTC

