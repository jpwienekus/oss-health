package score

import (
	"log"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	githubactivity "github.com/oss-health/background-worker/internal/github_activity"
)

func ScoreCadence(metrics *githubactivity.RepositoryMetrics, since int) int {
	daysSinceRelease := -1

	if len(metrics.ReleaseDates) > 0 {
		mostRecentRelease := metrics.ReleaseDates[0]
		daysSinceRelease = int(time.Since(mostRecentRelease).Hours() / 24)
	}

	usesSemver, consistenTagging := analyzeSemver(metrics)

	score := calculateCadenceScore(daysSinceRelease, metrics.ReleasesPastYear, usesSemver, consistenTagging, since)
	log.Printf("Days since last release: %d, Releases Per Year: %d, Uses Semver: %t, Consistent Tagging: %t", daysSinceRelease, metrics.ReleasesPastYear, usesSemver, consistenTagging)
	log.Printf("Cadence Score: %d", score)

	return score
}

func calculateCadenceScore(daysSinceLastRelease int, releasesPastYear int, usesSemver bool, consistentTagging bool, since int) int {
	score := 0

	if daysSinceLastRelease > -1 && daysSinceLastRelease < since {
		score += 3
	}

	if releasesPastYear > 4 {
		score += 4
	}

	if usesSemver {
		score += 2
	}

	if consistentTagging {
		score += 1
	}

	if score > 10 {
		score = 10
	}

	return score
}

func analyzeSemver(metrics *githubactivity.RepositoryMetrics) (bool, bool) {
	validCount := 0

	semverTags := []*semver.Version{}
	tagYears := map[int]bool{}
	malformedCount := 0

	for _, tag := range metrics.Tags {
		versionPart := extractSemver(tag.Name)

		v, err := semver.NewVersion(versionPart)
		if err == nil {
			validCount++

			semverTags = append(semverTags, v)
			tagYears[tag.Target.Tag.Tagger.Date.Year()] = true
		} else {
			malformedCount++
		}
	}
	log.Print(tagYears)

	totalTags := len(metrics.Tags)
	usesSemver := float64(validCount)/float64(totalTags) >= 0.7

	if totalTags == 0 {
		return usesSemver, false
	}

	enoughSemver := len(semverTags) >= 5
	lowMalformedRatio := float64(malformedCount)/float64(totalTags) <= 0.2

	currentYear := time.Now().Year()
	taggedThisYear := tagYears[currentYear]
	taggedLastYear := tagYears[currentYear-1]
	activeRecentYears := taggedThisYear || taggedLastYear
	validSemver := enoughSemver && lowMalformedRatio && activeRecentYears

	return usesSemver, validSemver
}

func extractSemver(tag string) string {
	parts := strings.Split(tag, "@")
	if len(parts) != 2 {
		return tag
	}
	return parts[1]
}
