package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

var DefaultEndpoint = "https://api.github.com/"

type Client struct {
	baseURL    *url.URL
	token      string
	httpClient *http.Client
}

func NewClient(token string) (*Client, error) {
	u, err := url.ParseRequestURI(DefaultEndpoint)
	if err != nil {
		return nil, err
	}
	return &Client{u, token, http.DefaultClient}, nil
}

func (c *Client) newRequest(spath, method string, body io.Reader) (*http.Request, error) {
	u := *c.baseURL
	u.Path = path.Join(u.Path, spath)
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.token)
	return req, nil
}

func (c *Client) do(req *http.Request, r interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
	}()
	if resp.StatusCode > 299 {
		return fmt.Errorf("failed to request: %v", resp.Status)
	}
	if r != nil {
		if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
			return err
		}
	}
	return nil
}
