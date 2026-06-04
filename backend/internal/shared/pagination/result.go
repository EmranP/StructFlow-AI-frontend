package pagination

type Result[T any] struct {
	Data []T `json:"data"`

	Page int `json:"page"`

	Limit int `json:"limit"`

	Total int64 `json:"total"`

	Pages int `json:"pages"`
}
