package score

import (
	"log"
	"strings"

	githubactivity "github.com/oss-health/background-worker/internal/github_activity"
)

func ScoreCommunity(metrics *githubactivity.RepositoryMetrics, since int) int {
	hasForum := communityForum(*metrics)
	score := calculateCommunityScore(metrics.Stars, metrics.ActiveContributors, metrics.IssueActivityRatio, hasForum)
	log.Printf("Stars: %d, Active Contributors: %d, Issue Activity Ratio: %.2f, Has Forum: %t", metrics.Stars, metrics.ActiveContributors, metrics.IssueActivityRatio, hasForum)
	log.Printf("Community Score: %d", score)
	return score
}

func calculateCommunityScore(stars int, activeContributors int, issueActivitityRatio float64, hasForum bool) int {
	score := 0

	if stars > 1000 {
		score += 2
	}

	if activeContributors > 5 {
		score += 4
	}

	if issueActivitityRatio > 0.5 {
		score += 2
	}

	if hasForum {
		score += 2
	}

	if score > 10 {
		score = 10
	}

	return score
}


func communityForum(metrics githubactivity.RepositoryMetrics) bool {
	forumPatterns := []string{"discord", "slack"}

	for _, pattern := range forumPatterns {
		if strings.Contains(strings.ToLower(metrics.Readme), pattern) {
			return true
		}
	}

	return false
}
