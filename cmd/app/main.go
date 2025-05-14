package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin/v2"

	"github.com/alecthomas/hasansino/github2karakeep/internal/github"
	"github.com/alecthomas/hasansino/github2karakeep/internal/karakeep"
)

const (
	DefaultReposPerPage     = "100"
	RequestTimeout          = "10s"
	DefaultKarakeepListName = "github2karakeep"
	DefaultUpdateInterval   = "24h"
	DefaultExportLimit      = "10"
	DefaultTagName          = "github2karakeep"
	DefaultExtractTopics    = "false"
)

var wg sync.WaitGroup

func main() {
	timeout := kingpin.Flag("timeout", "timeout for the requests").
		Envar("TIMEOUT").Default(RequestTimeout).Duration()
	ghUser := kingpin.Flag("gh-user", "github username").
		Envar("GH_USERNAME").Required().String()
	ghToken := kingpin.Flag("gh-token", "github personal access token").
		Envar("GH_TOKEN").Required().String()
	ghPerPage := kingpin.Flag("gh-per-page", "number of repos per page").
		Envar("GH_PER_PAGE").Default(DefaultReposPerPage).Int()
	ghExtractTopics := kingpin.Flag("gh-extract-topics", "extract topics from repo description as tags").
		Envar("GH_EXTRACT_TOPICS").Default(DefaultExtractTopics).Bool()
	kkHost := kingpin.Flag("kk-host", "karakeep host").
		Envar("KK_HOST").Required().String()
	kkToken := kingpin.Flag("kk-token", "karakeep token").
		Envar("KK_TOKEN").Required().String()
	kkList := kingpin.Flag("kk-list", "karakeep list").
		Envar("KK_LIST").Default(DefaultKarakeepListName).String()
	updateInterval := kingpin.Flag("update-interval", "update interval").
		Envar("UPDATE_INTERVAL").Default(DefaultUpdateInterval).Duration()
	exportLimit := kingpin.Flag("export-limit", "export limit").
		Envar("EXPORT_LIMIT").Default(DefaultExportLimit).Int()
	defaultTag := kingpin.Flag("default-tag", "default tag for bookmark").
		Envar("DEFAULT_TAG").Default(DefaultTagName).String()

	kingpin.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	ghService := github.New(*timeout, *ghToken, *ghPerPage)
	kkService := karakeep.New(*timeout, *kkHost, *kkToken, *defaultTag)

	wg.Add(1)
	go Run(ctx, *updateInterval, *exportLimit, ghService, *ghUser, kkService, *kkList, *ghExtractTopics)

	sys := make(chan os.Signal, 1)
	signal.Notify(sys, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	shutdown(<-sys, cancel)
}

func Run(
	ctx context.Context,
	updateInterval time.Duration,
	exportLimit int,
	ghService *github.Service,
	ghUser string,
	kkService *karakeep.Service,
	kkList string,
	ghExtractTopics bool,
) {
	ticker := time.NewTicker(updateInterval)
	defer wg.Done()
	defer ticker.Stop()

	// always run at start
	err := run(ctx, exportLimit, ghService, ghUser, kkService, kkList, ghExtractTopics)
	if err != nil {
		log.Printf("Failed to execute exporter: %s\n", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := run(ctx, exportLimit, ghService, ghUser, kkService, kkList, ghExtractTopics)
			if err != nil {
				log.Printf("Failed to execute exporter: %s\n", err)
			}
		}
	}
}

func run(
	ctx context.Context,
	exportLimit int,
	ghService *github.Service,
	ghUser string,
	kkService *karakeep.Service,
	kkList string,
	ghExtractTopics bool,
) error {

	log.Printf("Starting exporter...")

	// --- Retrieve starred repos ---

	log.Printf("Retrieving starred repos...")

	allRepos, err := ghService.GetStarredRepos(ctx, ghUser)
	if err != nil {
		return fmt.Errorf("failed to retrieve starred repos: %w", err)
	}

	log.Printf("Total starred repos: %d\n", len(allRepos))

	// --- Retrieve karakeep lists ---
	lists, err := kkService.GetAllLists(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve karakeep lists: %w", err)
	}

	// --- Check if default list exists ---
	listID := ""
	for _, list := range lists {
		if list.Name == kkList {
			listID = list.ID
			log.Printf("Found karakeep list %s with id %s\n", kkList, listID)
			break
		}
	}

	if len(listID) == 0 {
		log.Printf("Karakeep list %s not found. Creating...\n", kkList)
		list, err := kkService.CreateList(ctx, kkList)
		if err != nil {
			return fmt.Errorf("failed to create list: %w", err)
		}
		listID = list.ID
	}

	// --- Create / Update bookmarks ---

	log.Printf("Creating bookmarks...")

	var counter int
	for i := range allRepos {
		repo := allRepos[i]
		if repo.Repository.FullName == nil {
			log.Printf("Repo id %v is invalid - missing full_name field.\n", repo.Repository.ID)
			continue
		}
		if repo.Repository.HTMLURL == nil {
			log.Printf("Repo id %v is invalid - missing html_url feild.\n", repo.Repository.ID)
			continue
		}
		var repoDesc string
		if repo.Repository.Description != nil {
			repoDesc = *repo.Repository.Description
		}
		bookmark, err := kkService.CreateBookmark(
			ctx,
			*repo.Repository.FullName,
			*repo.Repository.HTMLURL,
			repoDesc,
		)
		if err != nil {
			return fmt.Errorf("failed to create bookmark: %w", err)
		}
		err = kkService.AddBookmarkToList(ctx, bookmark.ID, listID)
		if err != nil {
			return fmt.Errorf("failed to attach bookmark to list: %w", err)
		}

		var ghTopics []string
		if ghExtractTopics {
			ghTopics = repo.Repository.Topics
		}

		err = kkService.AddTagsToBookmark(ctx, bookmark.ID, ghTopics)
		if err != nil {
			return fmt.Errorf("failed to attach tags to bookmark: %w", err)
		}

		counter++

		if counter%100 == 0 {
			log.Printf("Processed %d bookmarks.\n", counter)
		}

		if exportLimit > 0 && counter == exportLimit {
			break
		}
	}

	log.Printf("Exporter finished...")

	return nil
}

func shutdown(
	_ os.Signal,
	cancel context.CancelFunc,
) {
	cancel()
	wg.Wait()
	os.Exit(0)
}
