package online

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type ServiceSelectors struct {
	InputSelector     string
	SubmitSelector    string
	ResponseSelector  string
	LoadingSelector   string
	AuthRequired      string
	RateLimitSelector string
}

var serviceSelectors = map[string]ServiceSelectors{
	"deepseek": {
		InputSelector:     "textarea[placeholder*='message'], textarea[name='prompt'], div[contenteditable='true']",
		SubmitSelector:    "button[type='submit'], button:has-text('Send'), button:has-text('→')",
		ResponseSelector:  "div[data-message-type='assistant'], div[class*='response'], div[class*='message']",
		LoadingSelector:   "div[class*='loading'], div[data-state='loading']",
		AuthRequired:      "input[type='email'], a[href*='login']",
		RateLimitSelector: "div[class*='rate'], div[class*='limit']",
	},
	"grok": {
		InputSelector:     "textarea[name='prompt'], textarea[placeholder*='Ask'], div[contenteditable='true']",
		SubmitSelector:    "button[type='submit'], button:has-text('Send'), button:has-text('↑')",
		ResponseSelector:  "div[data-testid='assistant-message'], div[class*='response']",
		LoadingSelector:   "div[class*='generating'], div[data-state='streaming']",
		AuthRequired:      "button:has-text('Sign in'), a[href*='login']",
		RateLimitSelector: "div[class*='rate'], div[class*='quota']",
	},
	"claude": {
		InputSelector:     "textarea[name='text'], textarea[placeholder*='Message], div[contenteditable='true']",
		SubmitSelector:    "button[type='submit'], button:has-text('Send'), button:has-text('↑')",
		ResponseSelector:  "div[data-testid='assistant-response'], div[class*='ClaudeMessage']",
		LoadingSelector:   "div[class*='thinking'], div[data-state='thinking']",
		AuthRequired:      "button:has-text('Log in'), a[href*='auth']",
		RateLimitSelector: "div[class*='rate'], div[class*='limit']",
	},
	"chatgpt": {
		InputSelector:     "textarea[name='prompt'], textarea[placeholder*='Send a message], div[contenteditable='true']",
		SubmitSelector:    "button[type='submit'], button[data-testid='send-button']",
		ResponseSelector:  "div[data-message-author='assistant'], div[class*='markdown']",
		LoadingSelector:   "div[data-state='generating'], button[disabled]",
		AuthRequired:      "button:has-text('Log in'), a[href*='login']",
		RateLimitSelector: "div[class*='rate'], div[class*='limit']",
	},
}

type PlaywrightClient struct {
	service    string
	context    playwright.BrowserContext
	page       playwright.Page
	profileDir string
	pw         *playwright.Playwright
}

func NewPlaywrightClient(service string) *PlaywrightClient {
	authDir := filepath.Join(os.Getenv("HOME"), ".bac-agent", "auth")
	os.MkdirAll(authDir, 0755)

	return &PlaywrightClient{
		service:    service,
		profileDir: filepath.Join(authDir, fmt.Sprintf("chrome-%s", service)),
	}
}

func (p *PlaywrightClient) Initialize() error {
	slog.Info("initializing playwright", "service", p.service)

	os.MkdirAll(p.profileDir, 0755)

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("failed to start playwright: %w", err)
	}
	p.pw = pw

	browserContext, err := pw.Chromium.LaunchPersistentContext(
		p.profileDir,
		playwright.BrowserTypeLaunchPersistentContextOptions{
			Headless: playwright.Bool(false),
			Args:     []string{"--start-maximized"},
		},
	)
	if err != nil {
		pw.Stop()
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	p.context = browserContext
	pages := browserContext.Pages()
	if len(pages) > 0 {
		p.page = pages[0]
	} else {
		p.page, _ = browserContext.NewPage()
	}

	slog.Info("playwright initialized", "service", p.service)
	return nil
}

func (p *PlaywrightClient) Open(ctx context.Context, url string) error {
	if p.page == nil {
		if err := p.Initialize(); err != nil {
			return err
		}
	}

	slog.Info("opening URL", "service", p.service, "url", url)
	_, err := p.page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(30000),
	})
	if err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}

	time.Sleep(2 * time.Second)
	slog.Info("page opened", "service", p.service)

	return nil
}

func (p *PlaywrightClient) Close() error {
	slog.Info("closing browser", "service", p.service)

	if p.context != nil {
		if err := p.SaveState(); err != nil {
			slog.Warn("failed to save state", "service", p.service, "err", err)
		}
		p.context.Close()
	}
	if p.pw != nil {
		p.pw.Stop()
	}
	return nil
}

func (p *PlaywrightClient) SaveState() error {
	statePath := filepath.Join(p.profileDir, "state.json")
	state, err := p.context.StorageState()
	if err != nil {
		return fmt.Errorf("failed to get storage state: %w", err)
	}

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	slog.Info("saved browser state", "service", p.service, "path", statePath)
	return nil
}

func (p *PlaywrightClient) LoadState() error {
	statePath := filepath.Join(p.profileDir, "state.json")
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		slog.Info("no saved state found", "service", p.service)
		return nil
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		return fmt.Errorf("failed to read state file: %w", err)
	}

	var state playwright.StorageState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	slog.Info("loaded browser state", "service", p.service, "cookies", len(state.Cookies))
	return nil
}

func (p *PlaywrightClient) TypeText(text string, selector string) error {
	slog.Info("typing text", "service", p.service, "selector", selector)

	element, err := p.page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		return fmt.Errorf("element not found: %w", err)
	}

	if err := element.Fill(text); err != nil {
		return fmt.Errorf("failed to fill: %w", err)
	}

	return nil
}

func (p *PlaywrightClient) Click(selector string) error {
	slog.Info("clicking element", "service", p.service, "selector", selector)

	element, err := p.page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(5000),
	})
	if err != nil {
		return fmt.Errorf("element not found: %w", err)
	}

	if err := element.Click(); err != nil {
		return fmt.Errorf("failed to click: %w", err)
	}

	return nil
}

func (p *PlaywrightClient) Press(key string) error {
	slog.Info("pressing key", "service", p.service, "key", key)

	keyboard := p.page.Keyboard()
	if keyboard == nil {
		return fmt.Errorf("keyboard not available")
	}

	if err := keyboard.Press(key); err != nil {
		return fmt.Errorf("failed to press key: %w", err)
	}

	return nil
}

func (p *PlaywrightClient) WaitForSelector(selector string, timeout time.Duration) error {
	slog.Info("waiting for selector", "service", p.service, "selector", selector, "timeout", timeout)

	_, err := p.page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(timeout.Seconds() * 1000),
	})
	if err != nil {
		return fmt.Errorf("timeout waiting for selector: %w", err)
	}

	return nil
}

func (p *PlaywrightClient) TakeSnapshot() (string, error) {
	slog.Info("taking snapshot", "service", p.service)

	snapshotPath := filepath.Join(os.TempDir(), fmt.Sprintf("snapshot-%s-%d.png",
		p.service, time.Now().Unix()))

	_, err := p.page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(snapshotPath),
	})
	if err != nil {
		return "", fmt.Errorf("failed to screenshot: %w", err)
	}

	return snapshotPath, nil
}

func (p *PlaywrightClient) GetPageHTML() (string, error) {
	slog.Info("getting page HTML", "service", p.service)

	html, err := p.page.Content()
	if err != nil {
		return "", fmt.Errorf("failed to get HTML: %w", err)
	}

	return html, nil
}

func (p *PlaywrightClient) GetElementText(selector string) (string, error) {
	slog.Info("getting element text", "service", p.service, "selector", selector)

	element, err := p.page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(5000),
	})
	if err != nil {
		return "", fmt.Errorf("element not found: %w", err)
	}

	text, err := element.TextContent()
	if err != nil {
		return "", fmt.Errorf("failed to get text: %w", err)
	}

	return strings.TrimSpace(text), nil
}

func (p *PlaywrightClient) IsAuthenticated() bool {
	selectors := p.GetServiceSelectors()
	if selectors.AuthRequired == "" {
		return true
	}

	element, err := p.page.QuerySelector(selectors.AuthRequired)
	if err != nil {
		return true
	}

	return element == nil
}

func (p *PlaywrightClient) HasRateLimit() bool {
	selectors := p.GetServiceSelectors()
	if selectors.RateLimitSelector == "" {
		return false
	}

	element, err := p.page.QuerySelector(selectors.RateLimitSelector)
	if err != nil {
		return false
	}

	return element != nil
}

func (p *PlaywrightClient) GetServiceSelectors() ServiceSelectors {
	if sel, ok := serviceSelectors[p.service]; ok {
		return sel
	}
	return ServiceSelectors{}
}

func (p *PlaywrightClient) WaitForResponse(timeout time.Duration) error {
	slog.Info("waiting for response", "service", p.service, "timeout", timeout)

	selectors := p.GetServiceSelectors()
	responseSel := selectors.ResponseSelector

	if responseSel == "" {
		time.Sleep(5 * time.Second)
		return nil
	}

	_, err := p.page.WaitForSelector(responseSel, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(timeout.Seconds() * 1000),
		State:   playwright.WaitForSelectorStateVisible,
	})

	if err != nil {
		return fmt.Errorf("timeout waiting for response: %w", err)
	}

	if selectors.LoadingSelector != "" {
		_, err = p.page.WaitForSelector(selectors.LoadingSelector, playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateHidden,
			Timeout: playwright.Float(timeout.Seconds() * 1000),
		})
	}

	return nil
}

func (p *PlaywrightClient) GetResponseText() (string, error) {
	selectors := p.GetServiceSelectors()
	responseSel := selectors.ResponseSelector

	if responseSel == "" {
		return p.GetPageHTML()
	}

	element, err := p.page.QuerySelector(responseSel)
	if err != nil || element == nil {
		return p.GetPageHTML()
	}

	text, err := element.TextContent()
	if err != nil {
		return "", fmt.Errorf("failed to get response text: %w", err)
	}

	return strings.TrimSpace(text), nil
}
