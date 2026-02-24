# BAC Hyper-Swarm 2026

## What is This?

**BAC Hyper-Swarm** is a smart study assistant designed specifically for students preparing for the **BAC C exam** (Terminale C in Mauritania). It uses artificial intelligence to help you learn faster and more effectively.

Think of it as having a personal tutor that:
- Finds educational videos and articles for you
- Creates practice quizzes from your study materials
- Explains difficult concepts in simple terms
- Helps you track your progress
- Works in Arabic, French, and English

---

## What Can It Do?

### 📚 Smart Study
- **Ask questions** - Get instant answers to any topic
- **Generate quizzes** - Practice with auto-created tests
- **Flashcards** - Review key concepts quickly
- **Audio summaries** - Listen to study materials on the go

### 🔍 Find Content
- **YouTube videos** - Automatically download educational videos
- **Web articles** - Fetch relevant articles from education sites
- **PDFs** - Process and analyze your documents

### 📝 Handwriting Help
- **OCR** - Convert handwritten notes to digital text
- **Save notes** - Store for later review

### 📊 Track Progress
- **Study stats** - See how much you've studied
- **Weak areas** - Identify topics that need more practice

---

## How to Use (3 Simple Steps)

### Step 1: Open Terminal
Open your terminal/command prompt and navigate to the BAC folder:
```bash
cd /path/to/bac
```

### Step 2: Run Commands
Here are the most common commands:

| What you want to do | Command |
|---------------------|---------|
| Ask a question | `bac ask "what is a derivative"` |
| Create a quiz | `bac quiz -t mathematics` |
| Study flashcards | `bac flashcards -t physics` |
| Listen to audio | `bac audio -t chemistry` |
| Check progress | `bac stats` |
| List topics | `bac config list-topics` |

### Step 3: Learn!
That's it! The system will help you study using AI.

---

## Need Help?

- **List all commands**: Run `bac.nu` without arguments
- **Get help**: `bac.nu --help`
- **Database status**: `bac db status`
- **Check daemons**: `bac daemon status`

---

## Example Workflows

### Morning Study Session
```bash
# 1. Review your flashcards
bac flashcards -t mathematics

# 2. Test yourself with a quiz
bac quiz -t "integrals"

# 3. Listen to a summary while getting ready
bac audio -t physics
```

### After Class
```bash
# 1. Ask questions about what you learned
bac ask "explain photosynthesis"

# 2. Check your weak areas
bac stats

# 3. Add a new topic you're studying
bac config add-topic "cell biology"
```

---

## Requirements

- **Nushell** - The command shell this program runs in
- **Internet connection** - For searching content and AI features
- **NotebookLM account** - For advanced AI features (optional)

---

## Troubleshooting

**Command not found?**
Make sure you're in the correct folder and Nushell is installed.

**Need help with a specific command?**
Just type the command without arguments to see usage instructions.

**Questions not answering?**
Check your internet connection and try again.

---

## What's Inside?

- `bac.nu` - The main program
- `daemons/` - Background services that run automatically
- `config/` - Your settings and topics
- `resources/` - Downloaded videos, articles, and notes
- `docs/` - More detailed documentation

---

**Status**: ✅ System Complete - Ready for BAC 2026!

Good luck with your exams! 🎓
