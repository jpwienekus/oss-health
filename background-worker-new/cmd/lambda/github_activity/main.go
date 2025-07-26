package main

import (
	"context"
	// "encoding/json"
	"log"

	// "os"

	"github.com/joho/godotenv"
	githubactivity "github.com/oss-health/background-worker/internal/github_activity"
	"github.com/oss-health/background-worker/internal/github_activity/score"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// connectionString := os.Getenv("DATABASE_URL")
	ctx := context.Background()

	token := ""
	monthsBack := 1
	test := githubactivity.NewGithubClient(token, ctx)

	repos := []struct {
		owner string
		name  string
	}{
		{"facebook", "react"},
		{"vuejs", "vue"},
		{"vercel", "next.js"},
		{"sveltejs", "svelte"},
	}

	for _, repo := range repos {
		data, err := test.FetchRepositoryData(ctx, repo.owner, repo.name, monthsBack)
		if err != nil {
			log.Fatalf("Failed to fetch github activity for %s/%s: %v", repo.owner, repo.name, err)
		}

		log.Print("===============")
		log.Print(repo.name)
		score.ScoreMaintenance(data, monthsBack)
		score.ScoreCommunity(data, monthsBack)
		score.ScoreCadence(data, monthsBack)
		log.Print("===============")
	}

	// jsonBytes, err := json.MarshalIndent(data, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Failed to marshal JSON: %v", err)
	// }
	//
	// log.Println(string(jsonBytes))
}

// func testScoring(metrics *githubactivity.RepositoryMetrics, since int) error {
// 	lastCommitAgeInDays := float64(-1)
// 	commitsPerMonth := float64(-1)
//
// 	if len(metrics.Commits) > 0 {
// 		lastCommitAgeInDays = float64(time.Since(metrics.Commits[0]).Hours() / 24)
// 		commitsPerMonth = float64(len(metrics.Commits)) / float64(since)
// 	}
//
// 	issueClosureRate := float64(metrics.ClosedIssues) / float64(metrics.ClosedIssues+metrics.OpenIssues)
// 	avgPRMergeTimeInDays := float64(-1)
//
// 	if len(metrics.PullRequests) > 0 {
// 		var totalPRMergeDays float64
// 		for _, pr := range metrics.PullRequests {
// 			totalPRMergeDays += float64(pr.MergedAt.Sub(pr.CreatedAt).Hours() / 24)
// 		}
// 		avgPRMergeTimeInDays = totalPRMergeDays / float64(len(metrics.PullRequests))
// 	}
//
// 	log.Printf("Last Commit Age (days): %.2f, Commits Per Month: %.2f, Issue Closure Rate: %.2f, Average PR Merge Time (days): %.2f", lastCommitAgeInDays, commitsPerMonth, issueClosureRate, avgPRMergeTimeInDays)
// 	log.Printf("Maintenance Score: %d", score.ScoreMaintenance(lastCommitAgeInDays, commitsPerMonth, issueClosureRate, avgPRMergeTimeInDays))
//
// 	return nil
// }
