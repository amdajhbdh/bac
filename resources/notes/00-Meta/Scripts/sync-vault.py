#!/usr/bin/env python3
"""
Multi-Protocol Vault Sync System
Supports: Croc (P2P), FTP/SFTP, Git/jj
"""
import subprocess
import json
import time
from pathlib import Path
from datetime import datetime

VAULT = Path("/home/med/Documents/bac/resources/notes/Study-Vault")
CONFIG = VAULT / ".sync-config.json"

class SyncOrchestrator:
    def __init__(self):
        self.vault = VAULT
        self.config = self.load_config()
    
    def load_config(self):
        if CONFIG.exists():
            return json.loads(CONFIG.read_text())
        return {
            "methods": ["croc", "git"],
            "croc_code": None,
            "git_remote": None,
            "last_sync": None,
            "sync_history": []
        }
    
    def save_config(self):
        CONFIG.write_text(json.dumps(self.config, indent=2))
    
    def run_cmd(self, cmd, cwd=None):
        """Run shell command"""
        result = subprocess.run(
            cmd, shell=True, cwd=cwd or self.vault,
            capture_output=True, text=True
        )
        return result.returncode == 0, result.stdout, result.stderr
    
    def sync_croc_send(self, code=None):
        """Send vault via Croc"""
        print("📡 Syncing via Croc (P2P)...")
        
        # Create tarball
        tarball = f"/tmp/vault-{int(time.time())}.tar.gz"
        success, out, err = self.run_cmd(
            f"tar czf {tarball} --exclude='.jj' --exclude='.git' .",
            cwd=self.vault
        )
        
        if not success:
            return False, "Failed to create tarball"
        
        # Send via croc
        if code:
            cmd = f"croc send --code {code} {tarball}"
        else:
            cmd = f"croc send {tarball}"
        
        success, out, err = self.run_cmd(cmd)
        
        # Extract code from output
        if "Code is:" in out:
            code = out.split("Code is:")[1].strip().split()[0]
            self.config['croc_code'] = code
            self.save_config()
        
        return success, out
    
    def sync_croc_receive(self, code):
        """Receive vault via Croc"""
        print(f"📡 Receiving via Croc (code: {code})...")
        
        tarball = f"/tmp/vault-received-{int(time.time())}.tar.gz"
        success, out, err = self.run_cmd(f"croc --yes {code}")
        
        if success:
            # Extract
            self.run_cmd(f"tar xzf {tarball} -C {self.vault}")
            return True, "Received and extracted"
        
        return False, err
    
    def sync_git_push(self):
        """Push to git remote"""
        print("🔄 Syncing via Git...")
        
        # Commit with jj
        timestamp = datetime.now().isoformat()
        success, out, err = self.run_cmd(
            f'jj commit -m "Auto-sync: {timestamp}"'
        )
        
        if not success:
            return False, "Failed to commit"
        
        # Push
        success, out, err = self.run_cmd("jj git push")
        
        return success, out if success else err
    
    def sync_git_pull(self):
        """Pull from git remote"""
        print("🔄 Pulling from Git...")
        
        success, out, err = self.run_cmd("jj git fetch")
        if not success:
            return False, err
        
        success, out, err = self.run_cmd("jj rebase")
        return success, out if success else err
    
    def sync_auto(self, direction="push"):
        """Auto-select best sync method"""
        methods = self.config.get('methods', ['git'])
        
        for method in methods:
            if method == "croc" and direction == "push":
                success, msg = self.sync_croc_send(self.config.get('croc_code'))
                if success:
                    self.log_sync("croc", direction, True)
                    return True, msg
            
            elif method == "git":
                if direction == "push":
                    success, msg = self.sync_git_push()
                else:
                    success, msg = self.sync_git_pull()
                
                if success:
                    self.log_sync("git", direction, True)
                    return True, msg
        
        return False, "All sync methods failed"
    
    def log_sync(self, method, direction, success):
        """Log sync operation"""
        self.config['last_sync'] = datetime.now().isoformat()
        self.config['sync_history'].append({
            "timestamp": datetime.now().isoformat(),
            "method": method,
            "direction": direction,
            "success": success
        })
        
        # Keep last 50 entries
        self.config['sync_history'] = self.config['sync_history'][-50:]
        self.save_config()
    
    def status(self):
        """Show sync status"""
        print("\n📊 Sync Status\n")
        print(f"Vault: {self.vault}")
        print(f"Last sync: {self.config.get('last_sync', 'Never')}")
        print(f"\nConfigured methods: {', '.join(self.config.get('methods', []))}")
        
        if self.config.get('croc_code'):
            print(f"Croc code: {self.config['croc_code']}")
        
        print(f"\nRecent syncs:")
        for entry in self.config.get('sync_history', [])[-5:]:
            status = "✓" if entry['success'] else "✗"
            print(f"  {status} {entry['timestamp'][:19]} - {entry['method']} ({entry['direction']})")

if __name__ == "__main__":
    import sys
    
    sync = SyncOrchestrator()
    
    if len(sys.argv) < 2:
        print("Usage: sync-vault.py <command> [args]")
        print("\nCommands:")
        print("  push              - Push vault to cloud")
        print("  pull              - Pull vault from cloud")
        print("  send [code]       - Send via Croc")
        print("  receive <code>    - Receive via Croc")
        print("  status            - Show sync status")
        sys.exit(1)
    
    cmd = sys.argv[1]
    
    if cmd == "push":
        success, msg = sync.sync_auto("push")
        print(msg)
        sys.exit(0 if success else 1)
    
    elif cmd == "pull":
        success, msg = sync.sync_auto("pull")
        print(msg)
        sys.exit(0 if success else 1)
    
    elif cmd == "send":
        code = sys.argv[2] if len(sys.argv) > 2 else None
        success, msg = sync.sync_croc_send(code)
        print(msg)
        sys.exit(0 if success else 1)
    
    elif cmd == "receive":
        if len(sys.argv) < 3:
            print("Error: code required")
            sys.exit(1)
        success, msg = sync.sync_croc_receive(sys.argv[2])
        print(msg)
        sys.exit(0 if success else 1)
    
    elif cmd == "status":
        sync.status()
    
    else:
        print(f"Unknown command: {cmd}")
        sys.exit(1)
