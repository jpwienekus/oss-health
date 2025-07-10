package fetcher

import (
	"errors"
	"testing"
)

type mockFetcher struct {
	canHandle bool
	fetchErr  error
	called    *bool
}

func (m mockFetcher) CanHandle(url string) bool {
	return m.canHandle
}

func (m mockFetcher) Fetch(url string, dest string) error {
	if m.called != nil {
		*m.called = true
	}
	return m.fetchErr
}

func TestFetcherChain_Fetch_Success(t *testing.T) {
	called := false
	mock := mockFetcher{
		canHandle: true,
		fetchErr:  nil,
		called:    &called,
	}

	chain := &FetcherChain{
		fetchers: []RepositoryFetcher{mock},
	}

	err := chain.Fetch("git@github.com:user/repo.git", "/tmp/repo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("expected fetcher to be called, but it wasn't")
	}
}

func TestFetcherChain_Fetch_NoHandler(t *testing.T) {
	mock := mockFetcher{
		canHandle: false,
		fetchErr:  nil,
	}

	chain := &FetcherChain{
		fetchers: []RepositoryFetcher{mock},
	}

	err := chain.Fetch("not-a-git-url", "/tmp/repo")
	if err == nil || err.Error() != "no fetcher found for repository URL: not-a-git-url" {
		t.Fatalf("expected specific error, got %v", err)
	}
}

func TestFetcherChain_Fetch_ErrorFromFetcher(t *testing.T) {
	called := false
	mock := mockFetcher{
		canHandle: true,
		fetchErr:  errors.New("clone failed"),
		called:    &called,
	}

	chain := &FetcherChain{
		fetchers: []RepositoryFetcher{mock},
	}

	err := chain.Fetch("git@github.com:user/repo.git", "/tmp/repo")
	if err == nil || err.Error() != "clone failed" {
		t.Fatalf("expected fetch error, got %v", err)
	}
	if !called {
		t.Fatal("expected fetcher to be called, but it wasn't")
	}
}
