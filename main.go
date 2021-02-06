package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pbar1/atlantis-go"
)

type (
	GitHubResponse struct {
		Items []Item `json:"items"`
	}

	Item struct {
		Number int `json:"number"`
	}
)

const (
	issuesEndpoint = "https://api.github.com/search/issues"
	helpMsg        = "Atlantis custom run command to check if a GitHub PR is approved by a CODEOWNER"
	notApprovedMsg = `Error: Pull request must be approved by a CODEOWNER

Note: If you are sure that this pull request is approved by a CODEOWNER, try running
"atlantis apply" again. This process checks for approval by querying the GitHub API,
which has delayed consistency with what is shown in the pull request UI. Wait a few
seconds and then try again.`
)

var version = "unknown"

func main() {
	if len(os.Args) > 0 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Fprintln(os.Stderr, helpMsg)
			os.Exit(0)
		}
		if os.Args[1] == "-v" || os.Args[1] == "-V" || os.Args[1] == "--version" {
			fmt.Fprintln(os.Stderr, version)
			os.Exit(0)
		}
	}

	step, err := atlantis.NewRunStep()
	if err != nil {
		log.Fatal(err)
	}

	reqURL := fmt.Sprintf("%s?q=%s/%s+review:approved+is:open+is:pr", issuesEndpoint, step.BaseRepoOwner, step.BaseRepoName)
	token := os.Getenv("ATLANTIS_GH_TOKEN")
	if token == "" {
		log.Fatal("environment variable ATLANTIS_GH_TOKEN must be set")
	}
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "token "+token)

	fmt.Println("Checking pull request for CODEOWNER approval")

	// the GitHub Issues API seems to have delayed consistency; if a PR was just approved
	// by a code owner and the check still fails, retrying it immediately usually works
	for range []struct{}{{}, {}} {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var ghResp GitHubResponse
		if err := json.Unmarshal(body, &ghResp); err != nil {
			log.Fatal(err)
		}
		for _, pr := range ghResp.Items {
			if pr.Number == step.PullNum {
				fmt.Println("Pull request is approved by a CODEOWNER")
				os.Exit(0)
			}
		}
	}

	fmt.Println(notApprovedMsg)
	os.Exit(1)
}
