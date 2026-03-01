# Skill: Jujutsu (jj) Workflow

## Purpose

Jujutsu (jj) version control for BAC Unified - NEVER use git directly.

## When to use

- Any version control operation
- Creating branches/changes
- Committing and pushing
- Viewing history

## MANDATORY: Never Use Git Directly

❌ **FORBIDDEN**: `git commit`, `git push`, `git rebase -i`
✅ **REQUIRED**: `jj commit`, `jj push`, `jj rebase`

## Core Commands

### Starting Work

```bash
# Create new change (like branch)
jj new -m "Add new feature"

# Or from specific commit
jj new v1.0.0 -m "Feature from v1"
```

### Making Changes

```bash
# Edit files normally (jj auto-tracks)

# See what changed
jj diff

# See status
jj status
```

### Committing

```bash
# Describe the change (sets message)
jj describe -m "Add user authentication"

# Commit
jj commit
```

### Viewing History

```bash
# View log
jj log

# View recent
jj log -r -5

# View change details
jj log --patch
```

### Pushing & Pulling

```bash
# Push to remote
jj push

# Pull from remote
jj pull

# Push with new branch name
jj push -r @ -b feature-name
```

### Undo & Redo

```bash
# Undo last operation
jj undo

# Undo specific operation
jj operation undo <op-id>
```

### Rebasing

```bash
# Rebase current change onto another
jj rebase -r @ -d main

# Rebase onto stack
jj rebase -r feature@1 -d feature
```

### Squashing

```bash
# Squash current with previous
jj squash

# Squash into specific commit
jj squash -r <commit>
```

## Workflow Example

```bash
# 1. Start new feature
jj new -m "Add math solver"

# 2. Make code changes
vim src/agent/internal/solver/solver.go

# 3. Run tests
cd src/agent && go test -v ./...

# 4. Format
go fmt ./...

# 5. Describe commit
jj describe -m "Add quadratic equation solver"

# 6. Commit
jj commit

# 7. Push
jj push
```

## Conflicts

```bash
# See conflicts
jj log --conflicts

# Resolve (edit files, then)
jj resolve

# Abandon conflict
jj abandon
```

## Git Interop

```bash
# Import from git
jj git import

# Export to git
jj git export

# Set git remote
jj git remote add origin https://github.com/bac-unified/bac.git
```

## Common Issues

### "Working copy has conflicts"

```bash
# See what's conflicted
jj status

# Resolve manually, then
jj resolve
```

### "The working copy is locked"

```bash
# Remove lock
jj workspace lock --unlock
```

## Installation

```bash
# macOS
brew install jj

# Linux
cargo install jj

# Windows
winget install jj
```

## Anti-Patterns

- ❌ Using `git commit`
- ❌ Using `git push`
- ❌ Using `git rebase -i`
- ❌ Using `git reset --hard`
- ❌ Using `git checkout`
