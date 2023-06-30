// Package npm implements functions for iteracting with the npm registry.
// See here: https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md
package npm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BaseURL = "https://registry.npmjs.org"
)

type Client struct {
	httpClient *http.Client
}

type ClientOption func(*Client)

// WithHttpClient overrides the default http client
func WithHttpClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// NewClient creates a new API client for interacting with the NPM registry API
func NewClient(opts ...ClientOption) *Client {
	c := &Client{http.DefaultClient}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

type RegistryMeta struct {
	DBName             string `json:"db_name"`
	DocCount           int    `json:"doc_count"`
	DocDelCount        int    `json:"doc_del_count"`
	UpdateSeq          int    `json:"update_seq"`
	PurgeSeq           int    `json:"purge_seq"`
	CompactRunning     bool   `json:"compact_running"`
	DiskSize           int    `json:"disk_size"`
	DataSize           int    `json:"data_size"`
	InstanceStartTime  string `json:"instance_start_time"`
	DiskFormatVersion  int    `json:"disk_format_version"`
	CommittedUpdateSeq int    `json:"committed_update_seq"`
}

// GetMeta makes an HTTP request to https://registry.npmjs.org/ and returns the JSON response
func (c *Client) GetMeta(ctx context.Context) (*RegistryMeta, *http.Response, error) {
	path := fmt.Sprintf("%s/", BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, res, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res, err
	}

	var r RegistryMeta
	if err = json.Unmarshal(body, &r); err != nil {
		return nil, res, err
	}

	return &r, res, nil

}

// Package represents an NPM package as returned by the registry API
// See https://github.com/npm/registry/blob/master/docs/responses/package-metadata.md
type Package struct {
	ID          string              `json:"_id"`
	Rev         string              `json:"_rev"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	DistTags    map[string]string   `json:"dist-tags"`
	Versions    map[string]*Version `json:"versions"`
	Time        struct {
		Created  string `json:"created"`
		Modified string `json:"modified"`
	} `json:"time"`
	Author struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		URL   string `json:"url"`
	}
	Repository struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	}
	Readme string `json:"readme"`
}

// Version represents a package version as returned by the NPM registry
// https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md#version
type Version struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Homepage   string `json:"homepage"`
	Repository struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	}
	Dependencies    map[string]any `json:"dependencies"`
	DevDependencies map[string]any `json:"devDependencies"`
	Scripts         map[string]any `json:"scripts"`
	Author          struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		URL   string `json:"url"`
	}
	License        string `json:"license"`
	Readme         string `json:"readme"`
	ReadmeFilename string `json:"readmeFilename"`
	ID             string `json:"_id"`
	Description    string `json:"description"`
	Dist           struct {
		SHASum  string `json:"shasum"`
		Tarball string `json:"tarball"`
	} `json:"dist"`
	NPMVersion  string           `json:"_npmVersion"`
	NPMUser     any              `json:"_npmUser"`
	Maintainers []map[string]any `json:"maintainers"`
}

// GetPackage makes an HTTP request to https://registry.npmjs.org/<<packageName>>
func (c *Client) GetPackage(ctx context.Context, packageName string) (*Package, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", BaseURL, packageName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, res, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res, err
	}

	var p Package
	if err = json.Unmarshal(body, &p); err != nil {
		return nil, res, err
	}

	return &p, res, nil
}

// GetPackageVersion makes an HTTP request to https://registry.npmjs.org/<<packageName>>/<<version>>
func (c *Client) GetPackageVersion(ctx context.Context, packageName, version string) (*Package, *http.Response, error) {
	path := fmt.Sprintf("%s/%s/%s", BaseURL, packageName, version)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, res, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res, err
	}

	var p Package
	if err = json.Unmarshal(body, &p); err != nil {
		return nil, res, err
	}

	return &p, res, nil
}
