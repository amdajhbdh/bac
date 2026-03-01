## Why

The current online research module opens AI service websites (DeepSeek, Grok, Claude, ChatGPT) in Chrome via Playwright but doesn't actually interact with them - it just opens the browser and takes a screenshot. This leaves the user to manually type problems and copy answers. Completing the Playwright automation will enable the agent to automatically fill prompts, submit questions, and extract AI-generated solutions.

## What Changes

- Add auto-fill functionality to type problems into AI service chat interfaces
- Add submit functionality to send prompts and wait for responses
- Add response extraction to scrape solutions from the AI chat outputs
- Add proper session management with authenticated Chrome profiles
- Add retry logic and error handling for robust automation
- Support multiple AI services: DeepSeek, Grok, Claude, ChatGPT

## Capabilities

### New Capabilities
- **ai-auto-solve**: Automatically send problems to AI services and extract solutions
- **browser-automation**: Core Playwright interaction primitives (fill, click, wait, extract)
- **session-management**: Handle authenticated browser sessions across services

### Modified Capabilities
- None - this is a net new capability

## Impact

- **Code**: `src/agent/internal/online/` - Enhanced Playwright automation
- **Dependencies**: Playwright CLI already installed
- **Config**: May need service-specific selectors in configuration
