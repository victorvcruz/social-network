package entities

type Comment struct {
	ID        string
	AccountID string
	PostID    string
	Content   string
	CreatedAt string
	UpdatedAt string
	Removed   bool
}
