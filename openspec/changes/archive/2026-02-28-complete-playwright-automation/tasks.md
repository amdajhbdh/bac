## 1. Browser Automation Primitives

- [x] 1.1 Create playwright client wrapper in `src/agent/internal/online/playwright.go`
- [x] 1.2 Implement `typeText()` function to type into input field
- [x] 1.3 Implement `clickSubmit()` function to submit prompt
- [x] 1.4 Implement `waitForResponse()` with configurable timeout
- [x] 1.5 Implement `takeSnapshot()` to capture page DOM state
- [x] 1.6 Add error handling for missing elements

## 2. Service-Specific Selectors

- [x] 2.1 Define DeepSeek selectors (input, submit button, response area)
- [x] 2.2 Define Grok selectors (input, submit button, response area)
- [x] 2.3 Define Claude selectors (input, submit button, response area)
- [x] 2.4 Define ChatGPT selectors (input, submit button, response area)
- [x] 2.5 Create selector configuration in `config/services.yaml`

## 3. Session Management

- [x] 3.1 Update `openWithChrome()` to support session reuse
- [x] 3.2 Implement session health check (detect auth expiry)
- [x] 3.3 Add session cleanup/reset functionality
- [x] 3.4 Configure persistent profiles per service

## 4. AI Auto-Solve Implementation

- [x] 4.1 Implement `autoSolve()` main function in online.go
- [x] 4.2 Add multi-service fallback logic (DeepSeek → Grok → Claude → ChatGPT)
- [x] 4.3 Implement response extraction from snapshot
- [x] 4.4 Add retry logic with exponential backoff
- [x] 4.5 Implement graceful degradation on failure

## 5. Testing & Integration

- [ ] 5.1 Test automation with each AI service
- [ ] 5.2 Verify session persistence works correctly
- [ ] 5.3 Test error handling (rate limits, timeouts)
- [x] 5.4 Update agent to call auto-solve when local models fail
- [x] 5.5 Build and verify agent works end-to-end
