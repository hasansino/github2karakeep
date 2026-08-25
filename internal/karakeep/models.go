package karakeep

import "strings"

type ErrorResponse struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func (e *ErrorResponse) Contains(s string) bool {
	return strings.Contains(e.Error, s)
}

type List struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListsResponse struct {
	Lists []List `json:"lists"`
}

type CreateListRequest struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Bookmark struct {
	ID string `json:"id"`
}

const BookmarkTypeLink = "link"

const BookmarkSourceImport = "import"

type CreateBookmarkRequest struct {
	Type     string `json:"type"`
	URL      string `json:"url"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Source   string `json:"source,omitempty"`
	ImportID string `json:"importSessionId,omitempty"`
}

type AddBookmarkToListRequest struct {
	ListID     string `json:"listId"`
	BookmarkID string `json:"bookmarkId"`
}

type AddTagsToBookmarkRequest struct {
	Tags []AddTagsToBookmarkRequestItem `json:"tags"`
}

type AddTagsToBookmarkRequestItem struct {
	TagName string `json:"tagName"`
}
