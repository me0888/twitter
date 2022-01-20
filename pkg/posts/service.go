package posts

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

func (s *Service) CreateTweet(ctx context.Context, id int64, content string) (models.Tweet, error) {
	var post models.Tweet

	err := s.pool.QueryRow(ctx,
		`INSERT INTO tweets (user_id, content) VALUES ($1, $2) RETURNING id, content, created_at, updated_at;`,
		id, content).Scan(&post.ID, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return post, fmt.Errorf("Error insert : %v", err)
	}

	return post, nil
}

func (s *Service) GetTweet(ctx context.Context, tweetID string) (models.Tweet, error) {
	var p models.Tweet
	err := s.pool.QueryRow(ctx,
		`SELECT id, content, likes_count, comments_count, retweets_count, created_at, updated_at
		FROM tweets WHERE id = $1;`, tweetID).
		Scan(&p.ID, &p.Content, &p.LikesCount, &p.CommentsCount, &p.RetweetsCount, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, fmt.Errorf("Error select post : %v", err)
	}
	return p, nil
}

func (s *Service) TweetLike(ctx context.Context, userID int64, tweetID string) (models.LikeResponse, error) {
	var response models.LikeResponse

	if err := s.pool.QueryRow(ctx, `SELECT EXISTS (
            SELECT 1 from tweet_likes WHERE user_id = $1 AND tweet_id = $2)
    `, userID, tweetID).Scan(&response.Liked); err != nil {
		return response, fmt.Errorf("Error query select tweet like : %v", err)
	}
	if response.Liked {

		if _, err := s.pool.Exec(ctx, "DELETE FROM tweet_likes WHERE user_id = $1 AND tweet_id = $2", userID, tweetID); err != nil {
			return response, fmt.Errorf("Error query delete tweet like: %v", err)
		}

		if err := s.pool.QueryRow(ctx, "UPDATE tweets SET likes_count = likes_count - 1 WHERE id = $1 RETURNING likes_count", tweetID).
			Scan(&response.LikesCount); err != nil {
			return response, fmt.Errorf("Error update tweet likes count: %v", err)
		}
	} else {

		_, err := s.pool.Exec(ctx, "INSERT INTO tweet_likes (user_id, tweet_id) VALUES ($1, $2)", userID, tweetID)

		if err != nil {
			return response, fmt.Errorf("Error insert tweet like: %v", err)
		}

		if err := s.pool.QueryRow(ctx, "UPDATE tweets SET likes_count = likes_count + 1 WHERE id = $1 RETURNING likes_count", tweetID).
			Scan(&response.LikesCount); err != nil {
			return response, fmt.Errorf("Error update tweet likes count: %v", err)
		}

	}

	response.Liked = !response.Liked
	return response, nil
}

func (s *Service) TweetRetweet(ctx context.Context, userID int64, tweetID string) (models.RetweetResponse, error) {
	var response models.RetweetResponse
	var ownPost bool

	if err := s.pool.QueryRow(ctx, `SELECT EXISTS (
            SELECT 1 from tweet_retweets WHERE user_id = $1 AND tweet_id = $2)
    `, userID, tweetID).Scan(&response.Retweeted); err != nil {
		return response, fmt.Errorf("Error query select tweet retweet : %v", err)
	}

	if err := s.pool.QueryRow(ctx, `SELECT EXISTS (
            SELECT 1 from tweets WHERE user_id = $1 AND id = $2)
    `, userID, tweetID).Scan(&ownPost); err != nil {
		return response, fmt.Errorf("Error query select tweet : %v", err)
	}

	if ownPost {
		return response, fmt.Errorf("You cant`t retweet you own tweet")
	}

	if response.Retweeted {

		if _, err := s.pool.Exec(ctx, "DELETE FROM tweet_retweets WHERE user_id = $1 AND tweet_id = $2", userID, tweetID); err != nil {
			return response, fmt.Errorf("Error query delete retweets : %v", err)
		}

		if err := s.pool.QueryRow(ctx, "UPDATE tweets SET retweets_count = retweets_count - 1 WHERE id = $1 RETURNING retweets_count", tweetID).
			Scan(&response.RetweesCount); err != nil {
			return response, fmt.Errorf("Error update tweet retweets count: %v", err)
		}
	} else {

		_, err := s.pool.Exec(ctx, "INSERT INTO tweet_retweets (user_id, tweet_id) VALUES ($1, $2)", userID, tweetID)

		if err != nil {
			return response, fmt.Errorf("Error insert retweets like: %v", err)
		}

		if err := s.pool.QueryRow(ctx, "UPDATE tweets SET retweets_count = retweets_count + 1 WHERE id = $1 RETURNING retweets_count", tweetID).
			Scan(&response.RetweesCount); err != nil {
			return response, fmt.Errorf("Error update tweet retweets count: %v", err)
		}

	}

	response.Retweeted = !response.Retweeted
	return response, nil
}

func (s *Service) TweetLikes(ctx context.Context, tweetID string) ([]models.UserProfile, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, email, username, followers_count, followees_count
		FROM tweet_likes, users
		WHERE tweet_likes.tweet_id = $1 
		AND users.id=tweet_likes.user_id
		ORDER BY username ASC
		`, tweetID)
	if err != nil {
		return nil, fmt.Errorf("Error query select: %v", err)
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

func (s *Service) TweetRetweetedUsers(ctx context.Context, tweetID string) ([]models.UserProfile, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, email, username, followers_count, followees_count
		FROM tweet_retweets, users
		WHERE tweet_retweets.tweet_id = $1 
		AND users.id=tweet_retweets.user_id
		ORDER BY username ASC
		`, tweetID)
	if err != nil {
		return nil, fmt.Errorf("Error query select: %v", err)
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

func (s *Service) GetTweets(ctx context.Context, username string) ([]models.Tweet, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, content, likes_count, comments_count, created_at
		FROM tweets
		WHERE user_id = (SELECT id FROM users WHERE username = $1) 
		ORDER BY created_at DESC
		`, username)
	if err != nil {
		return nil, fmt.Errorf("Error query select: %v", err)
	}

	defer rows.Close()

	pp := make([]models.Tweet, 0)
	for rows.Next() {
		var p models.Tweet
		if err = rows.Scan(&p.ID, &p.Content, &p.LikesCount, &p.CommentsCount, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("Error scan user: %v", err)
		}

		pp = append(pp, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterate user rows: %v", err)
	}
	return pp, nil
}

func (s *Service) UpdateTweet(ctx context.Context, id int64, tweet models.Tweet) (models.Tweet, error) {
	var resp models.Tweet
	var exsist bool

	if err := s.pool.QueryRow(ctx, `
	SELECT EXISTS (SELECT 1 FROM tweets WHERE id = $1 AND user_id = $2)`, tweet.ID, id).
		Scan(&exsist); err != nil {
		return resp, fmt.Errorf("Error query select tweets : %v", err)
	}

	if exsist {
		if _, err := s.pool.Exec(ctx, "UPDATE tweets SET content=$2, updated_at=$3 WHERE id = $1 ",
			tweet.ID, tweet.Content, tweet.UpdatedAt); err != nil {
			return resp, fmt.Errorf("Error update tweet: %v", err)
		}
	} else {
		return resp, fmt.Errorf("Error update tweet")
	}

	return resp, nil
}

func (s *Service) DeleteTweet(ctx context.Context, id int64, tweetID string) (models.Tweet, error) {
	var resp models.Tweet
	var exsist bool

	if err := s.pool.QueryRow(ctx, `
	SELECT EXISTS (SELECT 1 FROM tweets WHERE id = $1 AND user_id = $2)`, tweetID, id).
		Scan(&exsist); err != nil {
		return resp, fmt.Errorf("Error query select : %v", err)
	}

	if exsist {
		if err := s.pool.QueryRow(ctx, `DELETE FROM tweets WHERE id = $1 
		RETURNING id, content, likes_count, comments_count, created_at, updated_at`,
			tweetID).Scan(&resp.ID, &resp.Content, &resp.LikesCount, &resp.CommentsCount, &resp.CreatedAt, &resp.UpdatedAt); err != nil {
			return resp, fmt.Errorf("Error delete tweet: %v", err)
		}
	} else {
		return resp, fmt.Errorf("Error delete tweet")
	}

	return resp, nil
}

func (s *Service) ReadTweets(ctx context.Context, id int64) ([]models.Tweet, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT tweets.id, content, likes_count, comments_count, retweets_count, created_at, updated_at
		FROM follows, users, tweets
		WHERE follows.follower_id = $1 AND users.id=follows.followee_id AND tweets.user_id=followee_id 
		UNION
		SELECT tweets.id, content, likes_count, comments_count, retweets_count, created_at, updated_at
		FROM tweets, tweet_retweets, users, follows
		WHERE tweet_retweets.user_id =(SELECT followee_id WHERE follower_id = $1) 
		AND tweets.id=tweet_retweets.tweet_id AND users.id=tweet_retweets.user_id
		ORDER BY updated_at ASC
		`, id)

	if err != nil {
		return nil, fmt.Errorf("Error query select : %v", err)
	}

	defer rows.Close()

	pp := make([]models.Tweet, 0)
	for rows.Next() {
		var p models.Tweet
		if err = rows.Scan(&p.ID, &p.Content, &p.LikesCount, &p.CommentsCount, &p.RetweetsCount, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Error scan user: %v", err)
		}
		pp = append(pp, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterate user rows: %v", err)
	}
	return pp, nil
}
