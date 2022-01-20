package users

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/me0888/twitter/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

var ErrForbiddenFollow = errors.New("you can not follow yourself")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInternal = errors.New("internal error")

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) Save(ctx context.Context, item *models.UserInput) (*models.User_resp, error) {
	var resp models.User_resp
	var err error
	err = s.pool.QueryRow(ctx, `INSERT INTO users (email, username, password) VALUES ($1, $2, $3) 
	RETURNING email, username;
		`, item.Email, item.Username, item.Password).
		Scan(&resp.Email, &resp.Username)

	if err != nil {
		return nil, fmt.Errorf("Error insert user: %v", err)
	}

	return &resp, nil
}

func (s *Service) Update(ctx context.Context, item *models.UserInput, id int64) (*models.User_resp, error) {
	var resp models.User_resp
	var err error
	err = s.pool.QueryRow(ctx, `UPDATE users SET email=$1, username=$2, password=$3 WHERE id=$4 
	RETURNING email, username;
		`, item.Email, item.Username, item.Password, id).
		Scan(&resp.Email, &resp.Username)

	if err != nil {
		return nil, fmt.Errorf("Error update user: %v", err)
	}

	return &resp, nil
}

func (s *Service) UpdateAvatar(ctx context.Context, avatar string, id int64) (string, error) {
	var resp models.UserProfile
	var err error
	err = s.pool.QueryRow(ctx, `UPDATE users SET avatar=$1 WHERE id=$2 
	RETURNING avatar;
		`, avatar, id).
		Scan(&resp.Avatar)

	if err != nil {
		return "", fmt.Errorf("Error update avatar: %v", err)
	}

	return resp.Avatar, nil
}

func (s *Service) GetAvatar(ctx context.Context, id int64) (string, error) {
	var avatar string
	err := s.pool.QueryRow(ctx, `SELECT avatar  FROM users	WHERE id=$1 
		`, id).
		Scan(&avatar)
	if err != nil {
		return "", fmt.Errorf("Error query select avatar: %v", err)
	}

	return avatar, nil
}

func (s *Service) Follow(ctx context.Context, followerID int64, username string) (models.FollowResponse, error) {
	var response models.FollowResponse
	var followeeID int64
	err := s.pool.QueryRow(ctx, `SELECT id FROM users where username = $1;
		`, username).Scan(&followeeID)
	if err != nil {
		return response, fmt.Errorf("Error query select: %v", err)
	}

	if followeeID == followerID {
		return response, ErrForbiddenFollow
	}

	err = s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM follows WHERE follower_id = $1 AND followee_id = $2);
		`, followerID, followeeID).Scan(&response.Following)
	if err != nil {
		return response, fmt.Errorf("Error query select: %v", err)
	}

	if response.Following {
		_, err = s.pool.Exec(ctx, `DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2;`, followerID, followeeID)
		if err != nil {
			return response, fmt.Errorf("Error delete follow : %v", err)
		}

		_, err = s.pool.Exec(ctx, `UPDATE users SET followees_count = followees_count - 1 WHERE id = $1`, followerID)
		if err != nil {
			return response, fmt.Errorf("Error update followees count : %v", err)
		}

		err = s.pool.QueryRow(ctx, `UPDATE users SET followers_count = followers_count - 1 WHERE id = $1 RETURNING followers_count`, followeeID).Scan(&response.FollowersCount)
		if err != nil {
			return response, fmt.Errorf("Error update followers count : %v", err)
		}
	} else {
		_, err = s.pool.Exec(ctx, `INSERT INTO follows (follower_id, followee_id) VALUES ($1, $2);`, followerID, followeeID)
		if err != nil {
			return response, fmt.Errorf("Error insert follow: %v", err)
		}

		_, err = s.pool.Exec(ctx, `UPDATE users SET followees_count = followees_count + 1 where id = $1`, followerID)
		if err != nil {
			return response, fmt.Errorf("Error update follower followees count: %v", err)
		}

		err = s.pool.QueryRow(ctx, `UPDATE users SET followers_count = followers_count + 1 where id = $1 RETURNING followers_count`, followeeID).Scan(&response.FollowersCount)
		if err != nil {
			return response, fmt.Errorf("Error update followee followers count: %v", err)
		}

	}
	response.Following = !response.Following
	return response, nil
}

func (s *Service) Users(ctx context.Context, search string) ([]models.UserProfile, error) {

	rows, err := s.pool.Query(ctx, `
		SELECT id, email, username, followers_count, followees_count
		FROM users
		WHERE username ILIKE '%'|| $1 ||'%'
		ORDER BY username ASC
		`, search)
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

func (s *Service) User(ctx context.Context, id int64) (models.UserProfile, error) {
	var u models.UserProfile
	err := s.pool.QueryRow(ctx, `
		SELECT id, email, username, avatar, followers_count, followees_count 
		FROM users
		WHERE id=$1 
		ORDER BY username ASC
		`, id).
		Scan(&u.ID, &u.Email, &u.UserName, &u.Avatar, &u.FollowersCount, &u.FolloweesCount)
	if err != nil {
		return u, fmt.Errorf("Error query select: %v", err)
	}

	return u, nil
}

func (s *Service) Followers(ctx context.Context, username string) ([]models.UserProfile, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, email, username, followers_count, followees_count
		FROM follows, users
		WHERE follows.followee_id = (SELECT id FROM users WHERE username = $1) 
		AND users.id=follows.follower_id
		ORDER BY username ASC
		`, username)
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

func (s *Service) Followees(ctx context.Context, username string) ([]models.UserProfile, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, email, username, followers_count, followees_count
		FROM follows, users
		WHERE follows.follower_id = (SELECT id FROM users WHERE username = $1) 
		AND users.id=follows.followee_id
		ORDER BY username ASC
		`, username)
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

func (s *Service) Token(ctx context.Context, email string, password string) (token string, err error) {
	var hash string
	var id int64

	err = s.pool.QueryRow(ctx, `SELECT id, password FROM users WHERE email =$1`, email).Scan(&id, &hash)
	if err != nil {
		return "", fmt.Errorf("Error query select : %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}

	token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO users_tokens (token, user_id) VALUES ($1, $2)`, token, id)

	if err != nil {
		return "", fmt.Errorf("Error query insert users tokens : %v", err)
	}

	return token, nil
}

func (s *Service) IDByToken(ctx context.Context, token string) (id int64, err error) {
	var expire bool
	err = s.pool.QueryRow(ctx, `
	SELECT user_id, now()>expire as expire FROM users_tokens WHERE token =$1;`,
		token).Scan(&id, &expire)

	if err != nil {
		return 0, fmt.Errorf("Error query select token: %v", err)
	}

	return id, nil
}
