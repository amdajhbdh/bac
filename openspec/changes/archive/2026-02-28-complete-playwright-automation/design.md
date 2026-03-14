## Context

The BAC Unified agent needs to leverage online AI services (DeepSeek, Grok, Claude, ChatGPT) for research when local Ollama models fail or produce low-confidence results. Currently, the online module only opens these services in a browser but requires manual interaction.

**Current State:**
- Uses playwright-cli with Chrome persistent profiles
- Opens service URLs in separate browser sessions
- Takes screenshots but doesn't interact with pages
- User must manually type problems and copy solutions

**Constraints:**
- Must use existing authenticated Chrome sessions (user has sessions in Chrome)
- Must work with playwright-cli (not raw Playwright)
- Must handle rate limiting and captchas gracefully
- Must extract structured solutions from chat interfaces

## Goals / Non-Goals

**Goals:**
- Automatically type problems into AI service chat inputs
- Submit prompts and wait for AI responses
- Extract solution text from chat outputs
- Handle authentication via persistent Chrome profiles
- Support all four AI services: DeepSeek, Grok, Claude, ChatGPT

**Non-Goals:**
- Handle account authentication (assumes existing sessions)
- Bypass CAPTCHAs or rate limits (will gracefully fail)
- Extract complex formatted content (math, images)
- Real-time streaming of AI responses

## Decisions

### 1. Use playwright-cli commands over CDP
**Decision:** Use playwright-cli's CLI commands (type, click, snapshot) rather than raw CDP
**Rationale:** Simpler integration, already available, matches existing code pattern
**Alternative considered:** Direct CDP connection - more complex, requires separate setup

### 2. Service-specific selector strategies
**Decision:** Define selectors per service in configuration rather than a generic approach
**Rationale:** Each AI service has different DOM structures; service-specific is more maintainable
**Alternative considered:** Generic selectors with heuristics - less reliable

### 3. Snapshot-based response extraction
**Decision:** Use playwright-cli snapshot to get page state, parse for response text
**Rationale:** Captures current DOM state, works without complex wait logic
**Alternative considered:** Wait for specific elements - fragile across service updates

### 4. Graceful degradation
**Decision:** If automation fails, fall back to opening browser for manual interaction
**Rationale:** User can still interact manually; agent continues with partial success
**Alternative considered:** Hard failure - poor UX when services change

## Risks / Trade-offs

- **[Risk] Service DOM changes** → [Mitigation] Version selectors, fail gracefully, allow manual override
- **[Risk] Rate limiting from AI services** → [Mitigation] Exponential backoff, user-configurable delays
- **[Risk] Session expiry** → [Mitigation] Detect auth issues, prompt user to re-authenticate
- **[Risk] Slow automation** → [Mitigation] Timeout configuration, async operation

## Open Questions

- Should we cache responses to avoid repeated API calls?
- How to handle multi-step conversations vs. single prompt?
- Need to determine optimal wait times per service
