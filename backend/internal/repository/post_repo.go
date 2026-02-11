package repository

import (
	"database/sql"
	"grub-exchange/internal/models"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) Create(authorID, stockUserID int, content string) (*models.StockPost, error) {
	var post models.StockPost
	err := r.db.QueryRow(
		`INSERT INTO stock_posts (author_id, stock_user_id, content)
		 VALUES ($1, $2, $3)
		 RETURNING id, author_id, stock_user_id, content, likes, dislikes, created_at`,
		authorID, stockUserID, content,
	).Scan(&post.ID, &post.AuthorID, &post.StockUserID, &post.Content, &post.Likes, &post.Dislikes, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetByStock returns posts for a stock, with the requesting user's vote status.
func (r *PostRepo) GetByStock(stockUserID, requestingUserID, limit int) ([]models.StockPost, error) {
	rows, err := r.db.Query(
		`SELECT p.id, p.author_id, u.username, su.ticker, p.content, p.likes, p.dislikes, p.created_at,
		        COALESCE(v.vote_type, 0)
		 FROM stock_posts p
		 JOIN users u ON p.author_id = u.id
		 JOIN users su ON p.stock_user_id = su.id
		 LEFT JOIN post_votes v ON v.post_id = p.id AND v.user_id = $3
		 WHERE p.stock_user_id = $1
		 ORDER BY p.created_at DESC
		 LIMIT $2`,
		stockUserID, limit, requestingUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.StockPost
	for rows.Next() {
		var p models.StockPost
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.AuthorUsername, &p.StockTicker,
			&p.Content, &p.Likes, &p.Dislikes, &p.CreatedAt, &p.UserVote); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// Vote inserts or updates a vote, and updates the post's like/dislike counts atomically.
func (r *PostRepo) Vote(postID, userID, voteType int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Get existing vote if any
	var existingVote int
	err = tx.QueryRow(
		`SELECT vote_type FROM post_votes WHERE post_id = $1 AND user_id = $2`,
		postID, userID,
	).Scan(&existingVote)

	if err == sql.ErrNoRows {
		// New vote
		_, err = tx.Exec(
			`INSERT INTO post_votes (post_id, user_id, vote_type) VALUES ($1, $2, $3)`,
			postID, userID, voteType,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
		if voteType == 1 {
			_, err = tx.Exec(`UPDATE stock_posts SET likes = likes + 1 WHERE id = $1`, postID)
		} else {
			_, err = tx.Exec(`UPDATE stock_posts SET dislikes = dislikes + 1 WHERE id = $1`, postID)
		}
	} else if err != nil {
		tx.Rollback()
		return err
	} else if existingVote == voteType {
		// Same vote — remove it (toggle off)
		_, err = tx.Exec(`DELETE FROM post_votes WHERE post_id = $1 AND user_id = $2`, postID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if voteType == 1 {
			_, err = tx.Exec(`UPDATE stock_posts SET likes = GREATEST(likes - 1, 0) WHERE id = $1`, postID)
		} else {
			_, err = tx.Exec(`UPDATE stock_posts SET dislikes = GREATEST(dislikes - 1, 0) WHERE id = $1`, postID)
		}
	} else {
		// Switching vote
		_, err = tx.Exec(`UPDATE post_votes SET vote_type = $1 WHERE post_id = $2 AND user_id = $3`, voteType, postID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if voteType == 1 {
			// Was dislike, now like
			_, err = tx.Exec(`UPDATE stock_posts SET likes = likes + 1, dislikes = GREATEST(dislikes - 1, 0) WHERE id = $1`, postID)
		} else {
			// Was like, now dislike
			_, err = tx.Exec(`UPDATE stock_posts SET dislikes = dislikes + 1, likes = GREATEST(likes - 1, 0) WHERE id = $1`, postID)
		}
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetSentimentForStock returns the net sentiment (likes - dislikes) from the top 10
// most-engaged posts for a stock. Used by the market maker.
func (r *PostRepo) GetSentimentForStock(stockUserID int) (int, error) {
	var netSentiment sql.NullInt64
	err := r.db.QueryRow(
		`SELECT SUM(likes - dislikes) FROM (
			SELECT likes, dislikes FROM stock_posts
			WHERE stock_user_id = $1
			ORDER BY (likes + dislikes) DESC
			LIMIT 10
		) top_posts`,
		stockUserID,
	).Scan(&netSentiment)
	if err != nil {
		return 0, err
	}
	if netSentiment.Valid {
		return int(netSentiment.Int64), nil
	}
	return 0, nil
}

// GetAllSentiments returns sentiment scores for all stocks that have posts.
// Used by the market maker to batch-load sentiments.
func (r *PostRepo) GetAllSentiments() (map[int]int, error) {
	rows, err := r.db.Query(
		`SELECT stock_user_id, SUM(net) FROM (
			SELECT stock_user_id, (likes - dislikes) AS net,
			       ROW_NUMBER() OVER (PARTITION BY stock_user_id ORDER BY (likes + dislikes) DESC) AS rn
			FROM stock_posts
		) ranked
		WHERE rn <= 10
		GROUP BY stock_user_id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sentiments := make(map[int]int)
	for rows.Next() {
		var stockID, net int
		if err := rows.Scan(&stockID, &net); err != nil {
			return nil, err
		}
		sentiments[stockID] = net
	}
	return sentiments, nil
}
