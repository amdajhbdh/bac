# KiloClaw Configuration for Study Vault
# Update your KiloClaw at: https://claw.kilosessions.ai/config

## Current Setup

KiloClaw URL: https://claw.kilosessions.ai

## Tools to Add

Add these tools in your KiloClaw config:

```json
{
  "tools": {
    "vault-search": {
      "command": "curl",
      "args": ["{{vault_api}}/vault?query={{query}}"],
      "description": "Search the study vault"
    },
    "vault-read": {
      "command": "curl", 
      "args": ["{{vault_api}}/vault/{{path}}"],
      "description": "Read a vault note"
    },
    "get-study-progress": {
      "command": "curl",
      "args": ["{{vault_api}}/study/progress"],
      "description": "Get study progress"
    },
    "get-next-review": {
      "command": "curl",
      "args": ["{{vault_api}}/study/next"],
      "description": "Get cards due for review"
    }
  }
}
```

## Environment

Set these in KiloClaw settings:

```
vault_api=https://vault-api.onrender.com
```

## Quick Test

From Telegram:
- `/search [query]` - Search vault
- `/progress` - See study progress
- `/review` - Start review session

## Current Status

✅ KiloClaw: https://claw.kilosessions.ai
⏳ Vault API: Deploy with `just deploy-render`
⏳ Cloudflare Workers: Deploy with `just deploy-cf`
