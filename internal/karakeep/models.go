package karakeep

type List struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Bookmark struct {
	ID string `json:"id"`
}

type ListsResponse struct {
	Lists []List `json:"lists"`
}

type CreateListRequest struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type CreateBookmarkRequest struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

type AddBookmarkToListRequest struct {
	ListID     string `json:"listId"`
	BookmarkID string `json:"bookmarkId"`
}
