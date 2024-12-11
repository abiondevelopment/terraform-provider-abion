// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
	"terraform-provider-abion/internal/utils"
	"time"

	log "github.com/sirupsen/logrus"
)

const apiKeyHeader = "X-API-KEY"

// Client the Abion API client.
type Client struct {
	apiKey     string
	baseURL    *url.URL
	HTTPClient *http.Client
}

type ApiClient interface {
	GetZones(ctx context.Context, page *Pagination) (*APIResponse[[]Zone], error)
	GetZone(ctx context.Context, name string) (*APIResponse[*Zone], error)
	PatchZone(ctx context.Context, name string, patch ZoneRequest) (*APIResponse[*Zone], error)
}

// NewAbionClient Creates a new Client.
func NewAbionClient(host string, apiKey string, timeout int) (*Client, error) {
	baseURL, err := url.Parse(host)

	if err != nil {
		return nil, err
	}

	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}, nil
}

// GetZone Returns the full information on a single zone.
func (c *Client) GetZone(ctx context.Context, name string) (*APIResponse[*Zone], error) {
	endpoint := c.baseURL.JoinPath("v1", "zones", name)

	req, err := newJSONRequest(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, err
	}

	results := &APIResponse[*Zone]{}

	if err := c.do(req, results); err != nil {
		return nil, fmt.Errorf("could not get zone %s: %w", name, err)
	}

	return results, nil
}

// PatchZone Updates a zone by patching it according to JSON Merge Patch format (RFC 7396).
func (c *Client) PatchZone(ctx context.Context, name string, patch ZoneRequest) (*APIResponse[*Zone], error) {

	ctx = tflog.SetField(ctx, "key", c.apiKey)
	ctx = tflog.SetField(ctx, "url", c.baseURL)
	tflog.Debug(ctx, "Sending patch request")

	endpoint := c.baseURL.JoinPath("v1", "zones", name)

	req, err := newJSONRequest(ctx, http.MethodPatch, endpoint, patch)
	if err != nil {
		return nil, err
	}

	results := &APIResponse[*Zone]{}

	if err := c.do(req, results); err != nil {
		return nil, fmt.Errorf("could not update zone %s: %w", name, err)
	}

	return results, nil
}

func (c *Client) do(req *http.Request, result any) error {
	req.Header.Set(apiKeyHeader, c.apiKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return parseError(resp)
	}

	if result == nil {
		return nil
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body %w", err)
	}

	err = json.Unmarshal(raw, result)
	if err != nil {
		return fmt.Errorf("error unmarshalling response %w", err)
	}

	return nil
}

func newJSONRequest(ctx context.Context, method string, endpoint *url.URL, payload any) (*http.Request, error) {
	buf := new(bytes.Buffer)

	if payload != nil {
		err := json.NewEncoder(buf).Encode(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to create request JSON body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func parseError(resp *http.Response) error {
	raw, _ := io.ReadAll(resp.Body)

	zResp := &APIResponse[any]{}
	err := json.Unmarshal(raw, zResp)
	if err != nil {

		err2 := tryParseHtmlError(resp, raw)
		if err2 != nil {
			return err2
		}

		log.Errorf("error parsing error %s", err)
		return err
	}

	return zResp.Error
}

func tryParseHtmlError(resp *http.Response, raw []byte) error {

	// Handle special whitelist error
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		// Parse the HTML response to extract the title
		doc, err2 := html.Parse(strings.NewReader(string(raw)))
		if err2 != nil {
			log.Printf("Error parsing HTML: %s", err2)
			return err2
		}

		// Traverse the HTML tree to find the <title> tag
		var title string
		var traverseTitle func(*html.Node)
		traverseTitle = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "title" {
				// Extract the title content
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						title = c.Data
						break
					}
				}
			}
			// Continue traversing the child nodes
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				traverseTitle(c)
			}
		}

		// Start traversing from the root node
		traverseTitle(doc)

		if title != "" {
			// Return the error message with the title
			return fmt.Errorf("API error: %s", title)
		}
	}
	return nil
}

func CreateRecordPatchRequest(zoneName string, subDomainOrRoot string, recordType utils.RecordType, data []Record) ZoneRequest {
	records := make(map[string]map[string][]Record)

	if records[subDomainOrRoot] == nil {
		records[subDomainOrRoot] = make(map[string][]Record)
	}
	records[subDomainOrRoot][recordType.String()] = data

	patchRequest := ZoneRequest{
		Data: Zone{
			Type: "zone",
			ID:   zoneName,
			Attributes: Attributes{
				Records: records,
			},
		},
	}
	return patchRequest
}
