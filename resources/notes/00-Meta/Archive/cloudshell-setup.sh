#!/bin/bash
# Google Cloud Shell Setup Script
# Run this in Cloud Shell to setup your study vault

set -e

echo "=== Study Vault - Cloud Shell Setup ==="
echo ""

# Install dependencies
echo "📦 Installing dependencies..."
sudo apt-get update -qq
sudo apt-get install -y python3-pip sqlite3 tesseract-ocr

# Install jj
if ! command -v jj &> /dev/null; then
    echo "Installing jj..."
    curl -fsSL https://github.com/martinvonz/jj/releases/latest/download/jj-v0.23.0-x86_64-unknown-linux-musl.tar.gz | tar xz
    sudo mv jj /usr/local/bin/
fi

# Install croc
if ! command -v croc &> /dev/null; then
    echo "Installing croc..."
    curl https://getcroc.schollz.com | bash
fi

# Clone vault (if not exists)
if [ ! -d "$HOME/Study-Vault" ]; then
    echo "📥 Cloning vault..."
    cd $HOME
    
    # Option 1: Via Croc (if you have the code)
    read -p "Do you have a Croc code? (y/n): " has_code
    if [ "$has_code" = "y" ]; then
        read -p "Enter Croc code: " croc_code
        croc --yes $croc_code
        tar xzf vault-*.tar.gz
        mv Study-Vault $HOME/
    else
        # Option 2: Via Git (if configured)
        read -p "Enter Git repository URL (or press Enter to skip): " git_url
        if [ ! -z "$git_url" ]; then
            jj git clone $git_url Study-Vault
        else
            echo "⚠️  No sync method provided. Please sync manually."
            exit 1
        fi
    fi
fi

cd $HOME/Study-Vault

# Rebuild index
echo "🔨 Building search index..."
python3 /tmp/vault_indexer.py 2>/dev/null || echo "Index will be built on first use"

# Create aliases
echo "⚙️  Creating aliases..."
cat >> $HOME/.bashrc << 'ALIASES'

# Study Vault aliases
alias vault='cd $HOME/Study-Vault'
alias vault-search='python3 $HOME/Study-Vault/vault-cli.py search'
alias vault-stats='python3 $HOME/Study-Vault/vault-cli.py stats'
alias vault-sync='python3 $HOME/sync-vault.py'
ALIASES

source $HOME/.bashrc

echo ""
echo "✅ Setup complete!"
echo ""
echo "Quick commands:"
echo "  vault          - Go to vault directory"
echo "  vault-search   - Search notes"
echo "  vault-stats    - View statistics"
echo "  vault-sync     - Sync with local machine"
echo ""
echo "Start studying: cd ~/Study-Vault"
