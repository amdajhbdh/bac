# Multi-Protocol Sync System Guide

**Status:** ✅ Operational  
**Methods:** Croc (P2P), Git/jj, FTP (optional)

---

## Quick Start

### First Time Setup

**On Local Machine:**
```bash
cd /home/med/Documents/bac/resources/notes/Study-Vault

# Send vault to Cloud Shell
python3 ../sync-vault.py send
# Note the code (e.g., "code-word-word")
```

**On Google Cloud Shell:**
```bash
# Download and run setup script
curl -O https://your-server/cloudshell-setup.sh
bash cloudshell-setup.sh
# Enter the Croc code when prompted
```

---

## Sync Methods

### Method 1: Croc (P2P) - Fastest ⚡

**Send from local:**
```bash
python3 sync-vault.py send
# Outputs: "Code is: word-word-word"
```

**Receive on Cloud Shell:**
```bash
croc --yes word-word-word
tar xzf vault-*.tar.gz
```

**Advantages:**
- Very fast (P2P transfer)
- Works across networks
- No server needed
- Encrypted

**Use when:**
- Initial setup
- Large changes (many PDFs)
- Network is good

---

### Method 2: Git/jj - Most Reliable 🔄

**Push from local:**
```bash
cd Study-Vault
jj describe -m "Your changes"
jj git push
```

**Pull on Cloud Shell:**
```bash
cd ~/Study-Vault
jj git fetch
jj rebase
```

**Or use sync script:**
```bash
python3 sync-vault.py push  # Local
python3 sync-vault.py pull  # Cloud Shell
```

**Advantages:**
- Version history
- Conflict resolution
- Works everywhere
- Automatic

**Use when:**
- Regular syncs
- Small changes (notes only)
- Want history

**Note:** PDFs are excluded (too large). Use Croc for PDFs.

---

### Method 3: FTP/SFTP - Backup Option 📁

**Setup FTP server (local):**
```bash
# Install vsftpd
sudo apt install vsftpd

# Configure
sudo nano /etc/vsftpd.conf
# Set: write_enable=YES

sudo systemctl restart vsftpd
```

**Upload from local:**
```bash
tar czf vault.tar.gz Study-Vault/
lftp -u username,password ftp://your-server
put vault.tar.gz
```

**Download on Cloud Shell:**
```bash
wget ftp://your-server/vault.tar.gz
tar xzf vault.tar.gz
```

**Use when:**
- Croc fails
- Git not available
- Need reliable backup

---

## Sync Strategies

### Strategy 1: Hybrid (Recommended)

**For notes (markdown):** Use Git/jj
- Fast
- Version controlled
- Always synced

**For PDFs/assets:** Use Croc
- Only when needed
- Much faster than Git
- No repo bloat

**Workflow:**
```bash
# Daily: Sync notes via Git
jj describe -m "Daily progress"
jj git push

# Weekly: Sync PDFs via Croc (if changed)
python3 sync-vault.py send
```

---

### Strategy 2: Croc-Only

**Use when:**
- No Git remote configured
- Want simplicity
- Don't need history

**Workflow:**
```bash
# Anytime you want to sync
python3 sync-vault.py send
# On Cloud Shell: croc --yes <code>
```

---

### Strategy 3: Git-Only

**Use when:**
- PDFs already on Cloud Shell
- Only editing notes
- Want full history

**Workflow:**
```bash
# Auto-sync every 30 minutes
jj describe -m "Auto-sync $(date)"
jj git push
```

---

## Automation

### Auto-sync on Save

**Create systemd service (local):**
```bash
cat > ~/.config/systemd/user/vault-sync.service << 'SERVICE'
[Unit]
Description=Vault Auto-Sync

[Service]
ExecStart=/usr/bin/python3 /home/med/Documents/bac/resources/notes/sync-vault.py push
Restart=on-failure

[Install]
WantedBy=default.target
SERVICE

# Enable
systemctl --user enable vault-sync.timer
```

**Or use cron:**
```bash
# Edit crontab
crontab -e

# Add line (sync every 30 min)
*/30 * * * * cd /home/med/Documents/bac/resources/notes/Study-Vault && python3 ../sync-vault.py push
```

---

## Troubleshooting

### Croc: "Connection refused"
- Check firewall
- Try different relay: `croc --relay relay.example.com send`
- Use custom code: `croc send --code my-secret-code`

### Git: "Push rejected"
- Pull first: `jj git fetch && jj rebase`
- Check remote: `jj git remote -v`
- Force push (careful): `jj git push --force`

### FTP: "Permission denied"
- Check FTP user permissions
- Verify write_enable=YES in config
- Check firewall (port 21)

### "Files too large"
- PDFs excluded by .gitignore (correct)
- Use Croc for PDFs
- Or use Git LFS: `git lfs track "*.pdf"`

---

## Performance

**Croc (P2P):**
- Speed: ~10-50 MB/s (depends on network)
- Time for full vault: ~30 seconds

**Git/jj (notes only):**
- Speed: ~1-5 MB/s
- Time for notes: ~2 seconds

**FTP:**
- Speed: ~5-20 MB/s
- Time for full vault: ~1 minute

---

## Security

**Croc:**
- End-to-end encrypted
- Temporary codes
- No data stored on relay

**Git:**
- Use SSH keys
- HTTPS with tokens
- Private repositories

**FTP:**
- Use SFTP (encrypted)
- Or FTP over TLS
- Strong passwords

---

## Cloud Shell Specifics

**Persistent Disk:**
- Cloud Shell has 5GB persistent $HOME
- Vault fits easily (~200MB with PDFs)
- Survives session restarts

**Session Timeout:**
- Sessions end after 20 min inactive
- Vault persists in $HOME
- Just `cd ~/Study-Vault` to resume

**Network:**
- Fast internet (Google network)
- Croc works great
- Git push/pull very fast

---

## Commands Reference

**Sync Script:**
```bash
python3 sync-vault.py push      # Push to cloud
python3 sync-vault.py pull      # Pull from cloud
python3 sync-vault.py send      # Send via Croc
python3 sync-vault.py receive <code>  # Receive via Croc
python3 sync-vault.py status    # Show sync status
```

**jj Commands:**
```bash
jj describe -m "message"  # Commit changes
jj git push               # Push to remote
jj git fetch              # Fetch from remote
jj rebase                 # Update working copy
jj log                    # View history
jj diff                   # See changes
```

**Croc Commands:**
```bash
croc send file.tar.gz           # Send file
croc send --code mycode file    # Custom code
croc --yes code-word-word       # Receive (auto-accept)
croc --relay relay.com send     # Custom relay
```

---

## Next Steps

1. ✅ Setup complete
2. Choose sync strategy
3. Test sync (local ↔ Cloud Shell)
4. Setup automation (optional)
5. Start studying!

---

*Part of Vault Agent Ecosystem*
*Multi-Protocol Sync System v1.0*
