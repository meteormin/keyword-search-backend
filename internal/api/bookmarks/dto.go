package bookmarks

type CreateBookmark struct {
	Subject     string `json:"subject" validate:"required"`
	Description string `json:"description" validate:"required"`
	Url         string `json:"url" validate:"required"`
	Publish     string `json:"publish" validate:"required"`
	UserId      string `json:"user_id" validate:"required"`
}

type UpdateBookmark struct {
	Subject     string `json:"subject" validate:"required"`
	Description string `json:"description" validate:"required"`
	Url         string `json:"url" validate:"required"`
	Publish     string `json:"publish" validate:"required"`
}

type BookmarkResponse struct {
}
