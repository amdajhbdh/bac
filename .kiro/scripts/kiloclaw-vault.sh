#!/bin/bash
# KiloClaw Vault Integration
# Add this to your KiloClaw config

KILO_CLAW_URL="https://claw.kilosessions.ai"

echo "=== KiloClaw Vault Integration ==="
echo ""
echo "KiloClaw: $KILO_CLAW_URL"
echo ""

# Test KiloClaw
echo "Testing KiloClaw..."
curl -sf "$KILO_CLAW_URL" >/dev/null 2>&1 && echo "✅ KiloClaw: Online" || echo "⚠️  KiloClaw: Check URL"

# Show integration info
echo ""
echo "To connect KiloClaw to your vault:"
echo "1. Go to: $KILO_CLAW_URL/config"
echo "2. Add vault_api env var: https://vault-api.onrender.com"
echo "3. Add HTTP tools for vault access"
echo ""
echo "Deploy vault API first: just deploy-render"
