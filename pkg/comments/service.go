package comments

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/me0888/twitter/pkg/models"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) CreateComment(ctx context.Context, userID int64, tweetID, content string) (models.Comment, error) {
	var comment models.Comment

	err := s.pool.QueryRow(ctx, `INSERT INTO comments (tweet_id, user_id, likes_count, content) VALUES($1, $2, $3, $4) 
								RETURNING id, likes_count, content, created_at, updated_at`, tweetID, userID, 0, content).
		Scan(&comment.ID, &comment.LikesCount, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return comment, fmt.Errorf("Error incert comment: %v", err)
	}

	if _, err = s.pool.Exec(ctx, "UPDATE tweets SET comments_count = comments_count + 1 where id = $1", tweetID); err != nil {
		return comment, fmt.Errorf("Error update tweet comments count: %v", err)
	}

	return comment, nil
}

func (s *Service) GetComments(ctx context.Context, tweetID string) ([]models.Comment, error) {

	rows, err := s.pool.Query(ctx, `
	SELECT id, content, likes_count, created_at, updated_at
	FROM comments
	WHERE comments.tweet_id = $1
	ORDER BY created_at DESC`, tweetID)
	if err != nil {
		return nil, fmt.Errorf("Error query comments: %v", err)
	}
	defer rows.Close()
	cc := make([]models.Comment, 0)
	for rows.Next() {
		var c models.Comment
		if err = rows.Scan(&c.ID, &c.Content, &c.LikesCount, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Error scan comment: %v", err)
		}
		cc = append(cc, c)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterate comment rows: %v", err)
	}
	return cc, nil
}

func (s *Service) GetComment(ctx context.Context, commentID string) (models.Comment, error) {
	var comment models.Comment
	err := s.pool.QueryRow(ctx, `
	SELECT id, content, likes_count, created_at, updated_at
	FROM comments
	WHERE id = $1
	ORDER BY created_at DESC`, commentID).Scan(&comment.ID, &comment.Content, &comment.LikesCount, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return comment, fmt.Errorf("Error query select comments: %v", err)
	}

	return comment, nil
}

func (s *Service) UpdateComment(ctx context.Context, id int64, comment models.Comment) (models.Tweet, error) {
	var resp models.Tweet
	var response models.LikeResponse

	if err := s.pool.QueryRow(ctx, `
	SELECT EXISTS (SELECT 1 FROM comments WHERE id = $1 AND user_id = $2)`, comment.ID, id).
		Scan(&response.Liked); err != nil {
		return resp, fmt.Errorf("Error query select comment : %v", err)
	}

	if response.Liked {
		if _, err := s.pool.Exec(ctx, "UPDATE comments SET content=$2, updated_at=$3 WHERE id = $1 ",
			comment.ID, comment.Content, comment.UpdatedAt); err != nil {
			return resp, fmt.Errorf("Error update comments: %v", err)
		}
	} else {
		return resp, fmt.Errorf("Error update comments")
	}

	return resp, nil
}

func (s *Service) DeleteComment(ctx context.Context, id int64, commentID string) (models.Comment, error) {
	var resp models.Comment
	var response models.LikeResponse
	var tweetID int64

	if err := s.pool.QueryRow(ctx, `
	SELECT EXISTS (SELECT 1 FROM comments WHERE id = $1 AND user_id = $2)`, commentID, id).
		Scan(&response.Liked); err != nil {
		return resp, fmt.Errorf("Error query select comment : %v", err)
	}

	if response.Liked {
		err := s.pool.QueryRow(ctx, `SELECT tweet_id FROM comments WHERE id = $1`,
			commentID).Scan(&tweetID)
		if err != nil {
			return resp, fmt.Errorf("Error select  tweet_id : %v", err)
		}

		err = s.pool.QueryRow(ctx, `DELETE FROM comments WHERE id = $1 
		RETURNING id, content, likes_count, created_at, updated_at`,
			commentID).Scan(&resp.ID, &resp.Content, &resp.LikesCount, &resp.CreatedAt, &resp.UpdatedAt)
		if err != nil {
			return resp, fmt.Errorf("Error delete comment: %v", err)
		}

		_, err = s.pool.Exec(ctx, "UPDATE tweets SET comments_count = comments_count -1 where id = $1", tweetID)
		if err != nil {
			return resp, fmt.Errorf("Error update tweet comments count: %v", err)
		}

	} else {
		return resp, fmt.Errorf("Error delete comment")
	}

	return resp, nil
}

func (s *Service) CommentLike(ctx context.Context, userID int64, commentID string) (models.LikeResponse, error) {
	var response models.LikeResponse

	if err := s.pool.QueryRow(ctx, `SELECT EXISTS (
            SELECT 1 FROM comment_likes WHERE user_id = $1 AND comment_id = $2
        )
    `, userID, commentID).Scan(&response.Liked); err != nil {
		return response, fmt.Errorf("Error query select comment like : %v", err)
	}

	if response.Liked {

		if _, err := s.pool.Exec(ctx, "DELETE FROM comment_likes WHERE user_id = $1 AND comment_id = $2", userID, commentID); err != nil {
			return response, fmt.Errorf("Error query delete comment like: %v", err)
		}

		if err := s.pool.QueryRow(ctx, "UPDATE comments SET likes_count = likes_count - 1 WHERE id = $1 RETURNING likes_count", commentID).Scan(&response.LikesCount); err != nil {
			return response, fmt.Errorf("Error update comment likes count: %v", err)
		}
	} else {

		_, err := s.pool.Exec(ctx, "INSERT INTO comment_likes (user_id, comment_id) VALUES ($1, $2)", userID, commentID)

		if err != nil {
			return response, fmt.Errorf("Error insert comment like: %v", err)
		}

		if err := s.pool.QueryRow(ctx, "UPDATE comments SET likes_count = likes_count + 1 WHERE id = $1 RETURNING likes_count", commentID).Scan(&response.LikesCount); err != nil {
			return response, fmt.Errorf("Error update comments likes count: %v", err)
		}

	}

	response.Liked = !response.Liked
	return response, nil
}

func (s *Service) GetCommetsLikedUsers(ctx context.Context, commentID string) ([]models.UserProfile, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, email, username, followers_count, followees_count
		FROM comment_likes, users
		WHERE comment_likes.comment_id = $1 
		AND users.id=comment_likes.user_id
		ORDER BY username ASC
		`, commentID)
	if err != nil {
		return nil, fmt.Errorf("Error query select : %v", err)
	}
	defer rows.Close()
	uu := make([]models.UserProfile, 0)
	for rows.Next() {
		var u models.UserProfile

		if err = rows.Scan(&u.ID, &u.Email, &u.UserName, &u.FollowersCount, &u.FolloweesCount); err != nil {
			return nil, fmt.Errorf("Error scan user: %v", err)
		}

		uu = append(uu, u)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterate user rows: %v", err)
	}
	return uu, nil
}
