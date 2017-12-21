package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"strings"
)

func (c *Client) SearchIssues(terms []string) (*IssuesSearchResult, error) {
	req, err := c.newRequest("/search/issues", "GET", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", strings.Join(terms, " "))
	req.URL.RawQuery = q.Encode()

	var result IssuesSearchResult
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetIssues(username, repository string) ([]*Issue, error) {
	req, err := c.newRequest(path.Join("repos", username, repository, "issues"), "GET", nil)
	if err != nil {
		return nil, err
	}
	var issues []*Issue
	if err := c.do(req, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

func (c *Client) GetIssue(username, repository string, number int) (*Issue, error) {
	req, err := c.newRequest(path.Join("repos", username, repository, "issues", fmt.Sprint(number)), "GET", nil)
	if err != nil {
		return nil, err
	}
	var out Issue
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateIssue(username, repository, title, body string) (*Issue, error) {
	issue := struct {
		Body  string `json:"body"`
		Title string `json:"title"`
	}{
		Body:  body,
		Title: title,
	}
	encoded, err := json.Marshal(issue)
	if err != nil {
		return nil, err
	}
	b := bytes.NewReader(encoded)
	req, err := c.newRequest(path.Join("repos", username, repository, "issues"), "POST", b)
	if err != nil {
		return nil, err
	}
	var out Issue
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) EditIssue(username, repository string, number int, title, body string) (*Issue, error) {
	issue := struct {
		Body  string `json:"body"`
		Title string `json:"title"`
	}{
		Body:  body,
		Title: title,
	}
	encoded, err := json.Marshal(issue)
	if err != nil {
		return nil, err
	}
	b := bytes.NewReader(encoded)
	req, err := c.newRequest(path.Join("repos", username, repository, "issues", fmt.Sprint(number)), "PATCH", b)
	if err != nil {
		return nil, err
	}
	var out Issue
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CloseIssue(username, repository string, number int) (*Issue, error) {
	issue := struct {
		State string `json:"state"`
	}{
		State: "closed",
	}
	encoded, err := json.Marshal(issue)
	if err != nil {
		return nil, err
	}
	b := bytes.NewReader(encoded)
	req, err := c.newRequest(path.Join("repos", username, repository, "issues", fmt.Sprint(number)), "PATCH", b)
	if err != nil {
		return nil, err
	}
	var out Issue
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
