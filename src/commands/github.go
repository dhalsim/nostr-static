package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func TriggerGitHubAction(actionName string) error {
	if actionName != "deploy" {
		return fmt.Errorf("invalid action name: %s", actionName)
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Get required environment variables
	token := os.Getenv("GITHUB_TOKEN")
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")
	runID := os.Getenv("GITHUB_WORKFLOW_RUN_ID")

	if token == "" || owner == "" || repo == "" || runID == "" {
		return fmt.Errorf("missing required environment variables")
	}

	// Create the API request
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%s/approve", owner, repo, runID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error: %s - %s", resp.Status, string(body))
	}

	return nil
}
