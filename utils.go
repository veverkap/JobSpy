package jobspy

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

// Logger is the package-level logger
var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// SetLogLevel sets the logging level
func SetLogLevel(verbose int) {
	var level slog.Level
	switch verbose {
	case 0:
		level = slog.LevelError
	case 1:
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}
	
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

// ProxyRotator manages proxy rotation
type ProxyRotator struct {
	proxies []string
	current int
}

// NewProxyRotator creates a new proxy rotator
func NewProxyRotator(proxies []string) *ProxyRotator {
	return &ProxyRotator{
		proxies: proxies,
		current: 0,
	}
}

// Next returns the next proxy in rotation
func (pr *ProxyRotator) Next() string {
	if len(pr.proxies) == 0 {
		return ""
	}
	proxy := pr.proxies[pr.current]
	pr.current = (pr.current + 1) % len(pr.proxies)
	return proxy
}

// HTTPClient wraps http.Client with additional functionality
type HTTPClient struct {
	client       *http.Client
	proxyRotator *ProxyRotator
	userAgent    string
}

// NewHTTPClient creates a new HTTP client with proxy support
func NewHTTPClient(proxies []string, timeout time.Duration, userAgent string) *HTTPClient {
	client := &http.Client{
		Timeout: timeout,
	}
	
	var proxyRotator *ProxyRotator
	if len(proxies) > 0 {
		proxyRotator = NewProxyRotator(proxies)
	}
	
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	}

	return &HTTPClient{
		client:       client,
		proxyRotator: proxyRotator,
		userAgent:    userAgent,
	}
}

// Get performs a GET request with proxy rotation
func (hc *HTTPClient) Get(ctx context.Context, urlStr string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	// Set default headers
	req.Header.Set("User-Agent", hc.userAgent)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set proxy if available
	if hc.proxyRotator != nil {
		proxy := hc.proxyRotator.Next()
		if proxy != "" {
			proxyURL, err := url.Parse(proxy)
			if err == nil {
				hc.client.Transport = &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				}
			}
		}
	}

	return hc.client.Do(req)
}

// Post performs a POST request
func (hc *HTTPClient) Post(ctx context.Context, urlStr string, headers map[string]string, body string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", hc.userAgent)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return hc.client.Do(req)
}

// MarkdownConverter converts HTML to markdown
func MarkdownConverter(html string) string {
	if html == "" {
		return ""
	}
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(html)
	if err != nil {
		Logger.Warn("Failed to convert HTML to markdown", "error", err)
		return html
	}
	return strings.TrimSpace(markdown)
}

// PlainConverter converts HTML to plain text
func PlainConverter(html string) string {
	if html == "" {
		return ""
	}
	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		Logger.Warn("Failed to parse HTML", "error", err)
		return html
	}
	
	text := doc.Text()
	// Replace multiple whitespace with single space
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

// ExtractEmailsFromText extracts email addresses from text
func ExtractEmailsFromText(text string) []string {
	if text == "" {
		return nil
	}
	
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := emailRegex.FindAllString(text, -1)
	if len(emails) == 0 {
		return nil
	}
	return emails
}

// MapStringToSite converts a string to a Site enum
func MapStringToSite(s string) (Site, error) {
	switch strings.ToLower(s) {
	case "linkedin":
		return SiteLinkedIn, nil
	case "indeed":
		return SiteIndeed, nil
	case "zip_recruiter", "ziprecruiter":
		return SiteZipRecruiter, nil
	case "glassdoor":
		return SiteGlassdoor, nil
	case "google":
		return SiteGoogle, nil
	case "bayt":
		return SiteBayt, nil
	case "naukri":
		return SiteNaukri, nil
	case "bdjobs":
		return SiteBDJobs, nil
	default:
		return "", fmt.Errorf("unknown site: %s", s)
	}
}

// CountryFromString converts a string to a Country enum
func CountryFromString(s string) *Country {
	s = strings.ToLower(s)
	switch s {
	case "usa", "us", "united states":
		c := CountryUSA
		return &c
	case "canada", "ca":
		c := CountryCanada
		return &c
	case "uk", "united kingdom", "gb":
		c := CountryUK
		return &c
	case "germany", "de":
		c := CountryGermany
		return &c
	case "france", "fr":
		c := CountryFrance
		return &c
	default:
		return nil
	}
}

// FormatProxy formats a proxy string
func FormatProxy(proxy string) string {
	if !strings.HasPrefix(proxy, "http://") && !strings.HasPrefix(proxy, "https://") && !strings.HasPrefix(proxy, "socks5://") {
		return "http://" + proxy
	}
	return proxy
}