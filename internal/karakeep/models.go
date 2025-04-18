package karakeep

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

type CreateBookmarkRequest struct {
	Type    string `json:"type"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
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
