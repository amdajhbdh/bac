#!/bin/bash
# BAC Study - Termux Setup Script
# Run this on Termux (Android)

echo "📚 BAC Study - Termux Setup"
echo "=============================="

# Update and install dependencies
echo "[1/5] Installing dependencies..."
pkg update -y
pkg install -y python curl git

# Install Python packages
echo "[2/5] Installing Python packages..."
pip install --break-system-packages pyyaml requests typer rich

# Clone or update BAC repo (if needed)
# cd ~/Documents/bac 2>/dev/null || mkdir -p ~/Documents/bac

# Test connection
echo "[3/5] Testing Worker connection..."
STATUS=$(curl -s https://bac-api.amdajhbdh.workers.dev/rag/status)
if echo "$STATUS" | grep -q "total_vectors"; then
	echo "✅ Worker connected"
else
	echo "⚠️ Worker unreachable"
fi

# Test commands
echo "[4/5] Testing CLI..."
python3 ~/Documents/bac/scripts/bac-typer.py status

echo "[5/5] Quick aliases - add to ~/.bashrc:"
echo ""
echo "# BAC Study aliases"
echo "alias bac='python3 ~/Documents/bac/scripts/bac-typer.py'"
echo "alias bac-q='python3 ~/Documents/bac/scripts/bac-typer.py query'"
echo "alias bac-s='python3 ~/Documents/bac/scripts/bac-typer.py solve'"
echo "alias bac-st='python3 ~/Documents/bac/scripts/bac-typer.py status'"
echo ""
echo "=============================="
echo "✅ Setup complete!"
echo ""
echo "Usage:"
echo "  bac query \"question\"    # Query vault"
echo "  bac solve \"question\"    # Solve problem"
echo "  bac status              # Check system"
echo "  bac practice            # Get practice question"
