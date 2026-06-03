package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/sho0pi/tickli/internal/types"
)

const (
	baseURL     = "https://api.ticktick.com/open/v1"
	authURL     = "https://ticktick.com/oauth/authorize"
	scope       = "tasks:write tasks:read"
	redirectURL = "http://localhost:8080"
)

// tokenURL is a var (not const) so tests can point it at a mock server.
var tokenURL = "https://ticktick.com/oauth/token"

type Client struct {
	http *resty.Client
}

func NewClient(token string) *Client {
	client := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Authorization", "Bearer "+token)

	return &Client{http: client}
}

func GetAuthURL(clientID string) string {
	return fmt.Sprintf("%s?scope=%s&client_id=%s&state=state&redirect_uri=%s&response_type=code",
		authURL, scope, clientID, redirectURL)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func GetAccessToken(clientID, clientSecret, code string) (*TokenResponse, error) {
	client := resty.New()

	resp, err := client.R().
		SetBasicAuth(clientID, clientSecret).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":   "authorization_code",
			"code":         code,
			"redirect_uri": redirectURL,
		}).
		Post(tokenURL)

	if err != nil {
		return nil, errors.Wrap(err, "requesting access token")
	}

	var result TokenResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, errors.Wrap(err, "parsing response")
	}

	return &result, nil
}

func RefreshAccessToken(clientID, clientSecret, refreshToken string) (*TokenResponse, error) {
	client := resty.New()

	resp, err := client.R().
		SetBasicAuth(clientID, clientSecret).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
		}).
		Post(tokenURL)

	if err != nil {
		return nil, errors.Wrap(err, "refreshing access token")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("refresh token rejected (status %d): %s", resp.StatusCode(), resp.String())
	}

	var result TokenResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, errors.Wrap(err, "parsing refresh response")
	}

	return &result, nil
}

func (c *Client) ListProjects() ([]types.Project, error) {
	var projects []types.Project
	resp, err := c.http.R().
		SetResult(&projects).
		Get("/project")

	if err != nil {
		return nil, errors.Wrap(err, "listing projects")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to list projects: %s", resp.String())
	}

	// Adds the default InboxProject - not appears by default
	projects = append(projects, types.InboxProject)

	return projects, nil
}

func (c *Client) GetTask(projectID string, taskID string) (*types.Task, error) {
	var task types.Task
	resp, err := c.http.R().
		SetResult(&task).
		Get(fmt.Sprintf("/project/%s/task/%s", projectID, taskID))

	if err != nil {
		return nil, errors.Wrap(err, "requesting task")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to list tasks: %s", resp.String())
	}

	return &task, nil
}

func (c *Client) ListTasks(projectID string) ([]types.Task, error) {
	var projectData struct {
		Tasks []types.Task `json:"tasks"`
	}
	resp, err := c.http.R().
		SetResult(&projectData).
		Get(fmt.Sprintf("/project/%s/data", projectID))

	if err != nil {
		return nil, errors.Wrap(err, "listing tasks")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to list tasks: %s", resp.String())
	}

	return projectData.Tasks, nil
}

func (c *Client) CreateTask(task *types.Task) (*types.Task, error) {
	if task == nil {
		return nil, errors.New("task cannot be nil")
	}

	resp, err := c.http.R().
		SetBody(task).
		SetResult(task).
		Post("/task")

	if err != nil {
		return nil, errors.Wrap(err, "creating task")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to create task: %s", resp.String())
	}

	return task, nil
}

func (c *Client) UpdateTask(task *types.Task) (*types.Task, error) {
	if task == nil {
		return nil, errors.New("task cannot be nil")
	}

	resp, err := c.http.R().
		SetBody(task).
		SetResult(task).
		Post(fmt.Sprintf("/task/%s", task.ID))

	if err != nil {
		return nil, errors.Wrap(err, "updating task")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to update task: %s", resp.String())
	}

	return task, nil
}

func (c *Client) DeleteTask(projectID, taskID string) error {
	resp, err := c.http.R().
		Delete(fmt.Sprintf("/project/%s/task/%s", projectID, taskID))

	if err != nil {
		return errors.Wrap(err, "deleting task")
	}
	if resp.IsError() {
		return fmt.Errorf("failed to delete task: %s", resp.String())
	}

	return nil
}

func (c *Client) CompleteTask(projectID, taskID string) error {
	resp, err := c.http.R().
		Post(fmt.Sprintf("/project/%s/task/%s/complete", projectID, taskID))

	if err != nil {
		return errors.Wrap(err, "completing task")
	}
	if resp.IsError() {
		return fmt.Errorf("failed to complete task: %s", resp.String())
	}

	return nil
}
