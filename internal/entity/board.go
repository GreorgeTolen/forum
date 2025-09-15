package entity

type Board struct {
	ID          int    `json:"id,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}
