package models

import "time"

type UserInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	//Avatar   string `json:"avatar"`
}

type Tweet struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"-"`
	Content       string    `json:"content"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	RetweetsCount int       `json:"retweets_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Comment struct {
	ID         int64     `json:"id"`
	Content    string    `json:"content"`
	LikesCount int       `json:"likes_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserProfile struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	UserName       string `json:"username"`
	Avatar         string `json:"avatar"`
	FollowersCount int64  `json:"followers_count"`
	FolloweesCount int64  `json:"followees_count"`
}

type User_resp struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type FollowResponse struct {
	Following      bool `json:"following"`
	FollowersCount int  `json:"followers_count"`
}

type LikeResponse struct {
	Liked      bool `json:"liked"`
	LikesCount int  `json:"likes_count"`
}

type RetweetResponse struct {
	Retweeted    bool `json:"retweeted"`
	RetweesCount int  `json:"retweets_count"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
}

type CreatePostInput struct {
	Content string `json:"content"`
}

type CreateCommentInput struct {
	Content string `json:"content"`
}
