package github

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/v72/github"
)

type Service struct {
	ghClient *github.Client
	perPage  int
	timeout  time.Duration
}

func New(timeout time.Duration, ghToken string, perPage int) *Service {
	return &Service{
		ghClient: github.NewClient(&http.Client{Timeout: timeout}).WithAuthToken(ghToken),
		perPage:  perPage,
		timeout:  timeout,
	}
}

func (s *Service) GetStarredRepos(ctx context.Context, user string) ([]*github.StarredRepository, error) {
	opt := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			PerPage: s.perPage,
		},
	}

	allRepos := make([]*github.StarredRepository, 0)

	var page int
	for {
		reqCtx, cancel := context.WithTimeout(ctx, s.timeout)
		repos, resp, err := s.ghClient.Activity.ListStarred(reqCtx, user, opt)
		cancel()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve starred: %w", err)
		}

		log.Printf("Retrieved page %d with %d repos\n", page, len(repos))

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

		page++
	}

	return allRepos, nil
}
