package repository

import (
	"database/sql"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"chirik/internal/models"
)

type Repository struct {
	db *sql.DB
}

func New(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

// ============ USERS ============

func (r *Repository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password, name, bio, created_at)
	          VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Name, user.Bio, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	user.CreatedAt = now
	return nil
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password, name, bio, created_at FROM users WHERE email = ?`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Name, &user.Bio, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, username, email, password, name, bio, created_at FROM users WHERE id = ?`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Name, &user.Bio, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetAllUsers() ([]*models.User, error) {
	query := `SELECT id, username, name, bio, created_at FROM users ORDER BY username`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Name, &user.Bio, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *Repository) UpdateUser(user *models.User) error {
	query := `UPDATE users SET name = ?, bio = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Name, user.Bio, user.ID)
	return err
}

// ============ TWEETS ============

func (r *Repository) CreateTweet(tweet *models.Tweet) error {
	query := `INSERT INTO tweets (user_id, username, content, likes, views, created_at)
	          VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, tweet.UserID, tweet.Username, tweet.Content, 0, 0, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	tweet.ID = int(id)
	tweet.Likes = 0
	tweet.Views = 0
	tweet.CreatedAt = now
	return nil
}

func (r *Repository) GetAllTweets() ([]*models.Tweet, error) {
	query := `SELECT id, user_id, username, content, likes, views, created_at 
	          FROM tweets ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []*models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Username,
			&tweet.Content, &tweet.Likes, &tweet.Views, &tweet.CreatedAt)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}
	return tweets, nil
}

func (r *Repository) GetTweetsByUser(userID int) ([]*models.Tweet, error) {
	query := `SELECT id, user_id, username, content, likes, views, created_at 
	          FROM tweets WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []*models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Username,
			&tweet.Content, &tweet.Likes, &tweet.Views, &tweet.CreatedAt)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}
	return tweets, nil
}

func (r *Repository) GetTweetByID(id int) (*models.Tweet, error) {
	query := `SELECT id, user_id, username, content, likes, views, created_at 
	          FROM tweets WHERE id = ?`

	var tweet models.Tweet
	err := r.db.QueryRow(query, id).Scan(
		&tweet.ID, &tweet.UserID, &tweet.Username,
		&tweet.Content, &tweet.Likes, &tweet.Views, &tweet.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tweet, nil
}

func (r *Repository) LikeTweet(id int) error {
	query := `UPDATE tweets SET likes = likes + 1 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository) IncrementViews(id int) error {
	query := `UPDATE tweets SET views = views + 1 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// ============ FOLLOWS ============

func (r *Repository) Follow(followerID, followingID int) error {
	query := `INSERT INTO follows (follower_id, following_id, created_at) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, followerID, followingID, time.Now())
	return err
}

func (r *Repository) Unfollow(followerID, followingID int) error {
	query := `DELETE FROM follows WHERE follower_id = ? AND following_id = ?`
	_, err := r.db.Exec(query, followerID, followingID)
	return err
}

func (r *Repository) IsFollowing(followerID, followingID int) (bool, error) {
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = ? AND following_id = ?`
	var count int
	err := r.db.QueryRow(query, followerID, followingID).Scan(&count)
	return count > 0, err
}

func (r *Repository) GetFollowing(userID int) ([]int, error) {
	query := `SELECT following_id FROM follows WHERE follower_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repository) GetFollowers(userID int) ([]int, error) {
	query := `SELECT follower_id FROM follows WHERE following_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

// ============ MIGRATIONS ============

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			name TEXT NOT NULL,
			bio TEXT,
			created_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tweets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			username TEXT NOT NULL,
			content TEXT NOT NULL,
			likes INTEGER DEFAULT 0,
			views INTEGER DEFAULT 0,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS follows (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			follower_id INTEGER NOT NULL,
			following_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(follower_id, following_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tweets_user_id ON tweets(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tweets_created_at ON tweets(created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows(follower_id)`,
		`CREATE INDEX IF NOT EXISTS idx_follows_following ON follows(following_id)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
