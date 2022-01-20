DROP TABLE IF EXISTS users CASCADE; 
DROP TABLE IF EXISTS users_tokens CASCADE; 
DROP TABLE IF EXISTS tweets CASCADE;  
DROP TABLE IF EXISTS tweet_likes CASCADE;  
DROP TABLE IF EXISTS tweet_retweets CASCADE; 
DROP TABLE IF EXISTS comments CASCADE;
DROP TABLE IF EXISTS comment_likes CASCADE;  
DROP TABLE IF EXISTS follows CASCADE;  

CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    avatar TEXT DEFAULT '',
    followers_count BIGINT NOT NULL DEFAULT 0 CHECK (followers_count >= 0),
    followees_count BIGINT NOT NULL DEFAULT 0 CHECK (followees_count >= 0)
);

CREATE TABLE IF NOT EXISTS users_tokens (
   token    TEXT NOT NULL UNIQUE, 
   user_id BIGINT NOT NULL REFERENCES users,
   expire   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '24 hour',
   created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
 );

 CREATE TABLE IF NOT EXISTS follows (
  follower_id INT NOT NULL,
  followee_id INT NOT NULL,
  PRIMARY KEY (follower_id, followee_id)
);

CREATE TABLE IF NOT EXISTS tweets (
   id SERIAL NOT NULL PRIMARY KEY,
   user_id INT NOT NULL REFERENCES users,
   content TEXT NOT NULL,
   likes_count INT NOT NULL DEFAULT 0 CHECK (likes_count >= 0),
   comments_count INT NOT NULL DEFAULT 0 CHECK (comments_count >= 0),
   retweets_count INT NOT NULL DEFAULT 0 CHECK (comments_count >= 0),
   created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tweet_likes (
   user_id INT NOT NULL REFERENCES users,
   tweet_id INT NOT NULL REFERENCES tweets,
   PRIMARY KEY (user_id, tweet_id)
);

CREATE TABLE IF NOT EXISTS tweet_retweets (
   user_id INT NOT NULL REFERENCES users,
   tweet_id INT NOT NULL REFERENCES tweets,
   PRIMARY KEY (user_id, tweet_id)
);

CREATE TABLE IF NOT EXISTS comments (
   id SERIAL NOT NULL PRIMARY KEY,
   user_id INT NOT NULL REFERENCES users,
   tweet_id INT NOT NULL REFERENCES tweets,
   content TEXT NOT NULL,
   likes_count INT NOT NULL DEFAULT 0 CHECK (likes_count >= 0),
   created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS comment_likes (
   user_id INT NOT NULL REFERENCES users,
   comment_id INT NOT NULL REFERENCES comments,
   PRIMARY KEY (user_id, comment_id)
);