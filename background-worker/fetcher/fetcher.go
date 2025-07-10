package fetcher

type RepositoryFetcher interface {
	CanHandle(repoUrl string) bool
	Fetch(repoUrl string, destPath string) error
}
