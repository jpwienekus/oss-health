package fetcher

import "fmt"

type FetcherChain struct {
	fetchers []RepositoryFetcher
}

func NewFetcherChain() *FetcherChain {
	return &FetcherChain{
		fetchers: []RepositoryFetcher{
			GithubFetcher{},
		},
	}
}

func (fc *FetcherChain) Fetch(repoUrl string, destPath string) error {
	for _, fetcher := range fc.fetchers {
		if fetcher.CanHandle(repoUrl) {
			return fetcher.Fetch(repoUrl, destPath)
		}
	}

	return fmt.Errorf("no fetcher found for repository URL: %s", repoUrl)
}
