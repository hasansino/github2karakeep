package karakeep

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const DefaultListIcon = "🔸"

type Service struct {
	httpClient *http.Client
	host       string
	token      string
	defTag     string
}

func New(timeout time.Duration, host string, token string, defTag string) *Service {
	return &Service{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		host:   strings.Trim(host, "/"),
		token:  token,
		defTag: defTag,
	}
}

func (s *Service) GetAllLists(ctx context.Context) ([]List, error) {
	reqUrl := fmt.Sprintf("%s/api/v1/lists", s.host)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.token)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", res.Status)
	}

	response := new(ListsResponse)
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response.Lists, nil
}

func (s *Service) CreateList(ctx context.Context, name string) (*List, error) {
	payload := CreateListRequest{
		Name: name,
		Icon: DefaultListIcon,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	reqUrl := fmt.Sprintf("%s/api/v1/lists", s.host)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.token)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(res.Body)

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("http error: %s", res.Status)
	}

	response := new(List)
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) CreateBookmark(ctx context.Context, title string, url string, desc string) (*Bookmark, error) {
	payload := CreateBookmarkRequest{
		Type:    BookmarkTypeLink,
		Title:   title,
		URL:     url,
		Summary: desc,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	reqUrl := fmt.Sprintf("%s/api/v1/bookmarks", s.host)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.token)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(res.Body)

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("http error: %s", res.Status)
	}

	response := new(Bookmark)
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) AddBookmarkToList(ctx context.Context, bookmarkID string, listID string) error {
	reqUrl := fmt.Sprintf("%s/api/v1/lists/%s/bookmarks/%s", s.host, listID, bookmarkID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.token)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(res.Body)

	// should be idempotent instead
	// @see https://github.com/karakeep-app/karakeep/issues/1402
	if res.StatusCode != http.StatusNoContent {
		errorResp := new(ErrorResponse)
		if err := json.NewDecoder(res.Body).Decode(errorResp); err != nil {
			return err
		}
		if !errorResp.Contains("already in the list") {
			return fmt.Errorf("http error: %s", res.Status)
		}
	}

	return nil
}

func (s *Service) AddTagsToBookmark(ctx context.Context, bookmarkID string, tags []string) error {
	payload := AddTagsToBookmarkRequest{
		Tags: make([]AddTagsToBookmarkRequestItem, 0),
	}
	if len(s.defTag) > 0 {
		payload.Tags = append(payload.Tags, AddTagsToBookmarkRequestItem{s.defTag})
	}
	for i := range tags {
		payload.Tags = append(payload.Tags, AddTagsToBookmarkRequestItem{tags[i]})
	}
	if len(payload.Tags) == 0 {
		return nil
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	reqUrl := fmt.Sprintf("%s/api/v1/bookmarks/%s/tags", s.host, bookmarkID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.token)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http error: %s", res.Status)
	}

	return nil
}
