// Package wallapop uses the web app's private API anonymously: no account, no token.
package wallapop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	DefaultBaseURL   = "https://api.wallapop.com"
	DefaultWebURL    = "https://es.wallapop.com"
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/152.0.0.0 Safari/537.36"
	// The web build number: the API ignores it, but a real request carries it.
	DefaultAppVersion = "826680"
)

type Client struct {
	BaseURL string
	// WebURL serves a listing's page, which carries the details the API leaves out.
	WebURL     string
	Scheme     SignScheme
	UserAgent  string
	AppVersion string
	HTTP       *http.Client
}

func New() *Client {
	return &Client{
		BaseURL:    DefaultBaseURL,
		WebURL:     DefaultWebURL,
		Scheme:     SchemeNone,
		UserAgent:  DefaultUserAgent,
		AppVersion: DefaultAppVersion,
		HTTP:       &http.Client{Timeout: 30 * time.Second},
	}
}

type apiError struct {
	Status int
	Path   string
	Body   string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("wallapop: %s answered %d: %s", e.Path, e.Status, e.Body)
}

// deviceos is sent twice because the web app sends both spellings.
func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es,en-US;q=0.9")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Origin", "https://es.wallapop.com")
	req.Header.Set("Referer", "https://es.wallapop.com/")
	req.Header.Set("deviceos", "0")
	req.Header.Set("X-DeviceOS", "0")
	req.Header.Set("X-AppVersion", c.AppVersion)
}

// public sends no bearer on purpose, so watching searches is not tied to an account.
func (c *Client) public(ctx context.Context, path string, query url.Values, out any) error {
	const method = "GET"
	target := c.BaseURL + path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, target, nil)
	if err != nil {
		return err
	}

	c.setCommonHeaders(req)
	if c.Scheme != SchemeNone {
		ts := time.Now().UnixMilli()
		signature, err := Sign(c.Scheme, method, path, ts)
		if err != nil {
			return err
		}
		req.Header.Set("X-Signature", signature)
		req.Header.Set("Timestamp", strconv.FormatInt(ts, 10))
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("calling %s: %w", path, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	if resp.StatusCode >= 400 {
		return &apiError{Status: resp.StatusCode, Path: path, Body: snippet(raw)}
	}

	if out == nil || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decoding %s: %w (body: %s)", path, err, snippet(raw))
	}
	return nil
}

func snippet(b []byte) string {
	s := string(bytes.TrimSpace(b))
	if len(s) > 300 {
		s = s[:300] + "…"
	}
	return s
}
