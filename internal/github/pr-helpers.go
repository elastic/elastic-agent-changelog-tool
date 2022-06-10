package github

import (
	"context"
	"os/exec"
	"strings"
)

func GetLatestCommitHash(fileName string) (string, error) {
	response, err := exec.Command("git", "log", "--diff-filter=A", "--format=%H", "changelog/fragments/"+fileName).Output()
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(string(response), "\n", ""), nil
}

func FillEmptyPRField(commitHash string, c *Client) ([]int, error) {
	pr, err := FindPR(context.Background(), c, "elastic", "beats", commitHash)
	if err != nil {
		return []int{}, err
	}

	var prIDs []int

	for _, item := range pr.Items {
		prIDs = append(prIDs, item.PullRequestID)
	}

	return prIDs, nil
}

func FindOriginalPR(linkedPR int, c *Client) (int, error) {
	pr, _, err := c.PullRequests.Get(context.Background(), "elastic", "beats", linkedPR)
	if err != nil {
		return 0, err
	}

	prID, err := TestStrategies(pr, &BackportPRNumber{}, &PRNumber{})
	if err != nil {
		return 0, err
	}

	return prID, nil
}
