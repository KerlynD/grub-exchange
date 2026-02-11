package models

type StockPost struct {
	ID             int    `json:"id"`
	AuthorID       int    `json:"author_id"`
	AuthorUsername string `json:"author_username"`
	StockTicker    string `json:"stock_ticker"`
	StockUserID    int    `json:"-"`
	Content        string `json:"content"`
	Likes          int    `json:"likes"`
	Dislikes       int    `json:"dislikes"`
	UserVote       int    `json:"user_vote"` // 1 = liked, -1 = disliked, 0 = none
	CreatedAt      string `json:"created_at"`
}

type PostVote struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	VoteType  int    `json:"vote_type"` // 1 or -1
	CreatedAt string `json:"created_at"`
}

// SentimentScore represents the aggregated sentiment for a stock from its top posts
type SentimentScore struct {
	StockUserID int
	NetLikes    int // sum of (likes - dislikes) across top 10 posts
}
