# Daemons Documentation

## Overview

Daemon agents run continuously to automate resource discovery and processing.

## Daemon Types

### YouTube Daemon
- **Schedule**: Every 6 hours
- **Tool**: yt-dlp
- **Function**: Search & download educational YouTube videos
- **Sources**: Topics from `config/topics.yaml`

### Web Crawler Daemon
- **Schedule**: Every 12 hours
- **Tool**: playwright
- **Function**: Scrape Mauritanian education sites
- **Targets**:
  - minedu.mr
  - bacmr.net
  - mesrs.mr
  - men.gov.ma

### NLM Research Daemon
- **Schedule**: Daily
- **Tool**: nlm research
- **Function**: Auto-discover new sources
- **Queries**: From topics list

### Audio Daemon
- **Schedule**: Daily
- **Tool**: nlm audio
- **Function**: Generate podcast summaries
- **Output**: resources/audio/

## Configuration

### Systemd Timer Setup

```bash
# YouTube Daemon
systemctl --user enable youtube-daemon.timer
systemctl --user start youtube-daemon.timer

# Web Crawler
systemctl --user enable webcrawler-daemon.timer
systemctl --user start webcrawler-daemon.timer

# NLM Research
systemctl --user enable nlm-research-daemon.timer
systemctl --user start nlm-research-daemon.timer
```

### Manual Run

```bash
# Run once
nu daemons/youtube-daemon.nu

# Run in background
tmux new -s bac-daemons
nu daemons/youtube-daemon.nu
```

## Output Locations

| Daemon | Output Directory |
|--------|-----------------|
| YouTube | resources/youtube/ |
| Web Crawler | resources/web/ |
| NLM Research | Added to notebook |
| Audio | resources/audio/ |

## Monitoring

```bash
# Check timer status
systemctl --user list-timers

# View logs
journalctl --user -u youtube-daemon -f
```

## Error Handling

- Network failure → Retry next cycle
- Download failure → Log and skip
- NLM unavailable → Use local Ollama fallback
