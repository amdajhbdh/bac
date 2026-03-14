# Sync Guide

## Overview

Since we use local PostgreSQL (no cloud DB), sync between devices is manual.

## Sync Methods

### Method 1: SQL Export/Import

#### Export (PC 1)
```bash
pg_dump -U med bac_db > ~/bac/db/exports/backup_$(date +%Y%m%d).sql
jj commit -m "DB backup $(date +%Y%m%d)"
```

#### Import (PC 2)
```bash
# Copy the SQL file via USB/cloud
psql -U med bac_db < ~/Downloads/backup_20260223.sql
```

### Method 2: jj Git Export/Import

#### Export
```bash
jj git export
# This creates a git repository you can push to remote
git remote add origin <your-git-repo>
git push -u origin main
```

#### Import
```bash
git clone <your-git-repo> ~/bac-new
cd ~/bac-new
jj git import
```

### Method 3: rsync

```bash
# On PC 1
rsync -avz ~/bac/db/exports/ user@pc2:/home/user/bac/db/exports/
```

## What to Sync

| Item | File | Method |
|------|------|--------|
| Database | db/exports/*.sql | pg_dump |
| Config | config/topics.yaml | jj |
| Scripts | daemons/*.nu | jj |
| Resources | resources/* | rsync/jj |
| PDFs | db/pdfs/ | rsync |

## What NOT to Sync

- `db/postgres/` - PostgreSQL data directory (binary)
- `node_modules/` - npm packages
- `target/` - Rust build output
- `.env` - Credentials

## Daily Sync Workflow

```bash
# Morning: Pull latest
git pull
jj git import

# After changes: Push
pg_dump bac_db > db/exports/backup_$(date +%Y%m%d).sql
jj commit -m "Daily backup"
jj git export
git push
```

## Sharing with Students

```bash
# Export shareable package
tar -czvf bac-share.tar.gz \
    db/pdfs/ \
    db/exports/ \
    config/ \
    resources/

# Extract on student PC
tar -xzvf bac-share.tar.gz
```

## Progress Tracking

All sync actions tracked in jj:
```bash
jj log
jj bookmark add sync-2026-02-23
```
