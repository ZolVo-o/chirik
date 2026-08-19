package repository

import (
	"chirik/internal/models"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
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
		db.Close()
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

func (r *Repository) SearchUsers(query string, excludeID int) ([]*models.User, error) {
	rows, err := r.db.Query(`SELECT id, username, name, bio, created_at
		FROM users WHERE id != ? AND (username LIKE ? OR name LIKE ?)
		ORDER BY username LIMIT 20`, excludeID, "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Name, &user.Bio, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, rows.Err()
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

func (r *Repository) LikeTweet(tweetID, userID int) (bool, error) {
	if exists, err := r.tweetExists(tweetID); err != nil || !exists {
		return false, err
	}

	result, err := r.db.Exec(`INSERT OR IGNORE INTO tweet_likes (tweet_id, user_id, created_at) VALUES (?, ?, ?)`, tweetID, userID, time.Now())
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	if err != nil || count == 0 {
		return false, err
	}

	_, err = r.db.Exec(`UPDATE tweets SET likes = likes + 1 WHERE id = ?`, tweetID)
	return true, err
}

func (r *Repository) tweetExists(id int) (bool, error) {
	var exists int
	err := r.db.QueryRow(`SELECT 1 FROM tweets WHERE id = ?`, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *Repository) ViewTweet(tweetID, userID int) (bool, error) {
	if exists, err := r.tweetExists(tweetID); err != nil || !exists {
		return false, err
	}

	result, err := r.db.Exec(`INSERT OR IGNORE INTO tweet_views (tweet_id, user_id, created_at) VALUES (?, ?, ?)`, tweetID, userID, time.Now())
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	if err != nil || count == 0 {
		return false, err
	}

	_, err = r.db.Exec(`UPDATE tweets SET views = views + 1 WHERE id = ?`, tweetID)
	return true, err
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
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
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
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ============ MESSAGES ============

func (r *Repository) CreateConversation(userID, otherUserID int) (*models.Conversation, error) {
	var conversationID int
	err := r.db.QueryRow(`SELECT c.id FROM conversations c
		JOIN conversation_members m1 ON m1.conversation_id = c.id AND m1.user_id = ?
		JOIN conversation_members m2 ON m2.conversation_id = c.id AND m2.user_id = ?
		WHERE (SELECT COUNT(*) FROM conversation_members cm WHERE cm.conversation_id = c.id) = 2
		LIMIT 1`, userID, otherUserID).Scan(&conversationID)
	if err == nil {
		return r.GetConversation(conversationID, userID)
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	result, err := tx.Exec(`INSERT INTO conversations (updated_at) VALUES (?)`, time.Now())
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	conversationID64, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, memberID := range []int{userID, otherUserID} {
		if _, err := tx.Exec(`INSERT INTO conversation_members (conversation_id, user_id) VALUES (?, ?)`, conversationID64, memberID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetConversation(int(conversationID64), userID)
}

func (r *Repository) GetConversations(userID int) ([]*models.Conversation, error) {
	rows, err := r.db.Query(`SELECT c.id, u.id, u.username, u.name, u.bio, u.created_at,
		m.id, m.sender_id, m.content, m.created_at, m.updated_at, m.deleted, c.updated_at
		FROM conversations c
		JOIN conversation_members mine ON mine.conversation_id = c.id AND mine.user_id = ?
		JOIN conversation_members other ON other.conversation_id = c.id AND other.user_id != ?
		JOIN users u ON u.id = other.user_id
		LEFT JOIN messages m ON m.id = (SELECT id FROM messages WHERE conversation_id = c.id ORDER BY id DESC LIMIT 1)
		ORDER BY c.updated_at DESC`, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*models.Conversation
	for rows.Next() {
		conversation, message, err := scanConversation(rows)
		if err != nil {
			return nil, err
		}
		conversation.LastMessage = message
		conversations = append(conversations, conversation)
	}
	return conversations, rows.Err()
}

func (r *Repository) GetConversation(id, userID int) (*models.Conversation, error) {
	row := r.db.QueryRow(`SELECT c.id, u.id, u.username, u.name, u.bio, u.created_at,
		m.id, m.sender_id, m.content, m.created_at, m.updated_at, m.deleted, c.updated_at
		FROM conversations c
		JOIN conversation_members mine ON mine.conversation_id = c.id AND mine.user_id = ?
		JOIN conversation_members other ON other.conversation_id = c.id AND other.user_id != ?
		JOIN users u ON u.id = other.user_id
		LEFT JOIN messages m ON m.id = (SELECT id FROM messages WHERE conversation_id = c.id ORDER BY id DESC LIMIT 1)
		WHERE c.id = ?`, userID, userID, id)
	conversation, message, err := scanConversation(row)
	if err != nil {
		return nil, err
	}
	conversation.LastMessage = message
	return conversation, nil
}

func (r *Repository) IsConversationMember(conversationID, userID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM conversation_members WHERE conversation_id = ? AND user_id = ?`, conversationID, userID).Scan(&count)
	return count > 0, err
}

func (r *Repository) GetMessages(conversationID int) ([]*models.Message, error) {
	rows, err := r.db.Query(`SELECT m.id, m.conversation_id, m.sender_id, u.username, m.content, m.created_at, m.updated_at, m.deleted
		FROM messages m JOIN users u ON u.id = m.sender_id WHERE m.conversation_id = ? ORDER BY m.id ASC`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []*models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.ConversationID, &message.SenderID, &message.SenderUsername, &message.Content, &message.CreatedAt, &message.UpdatedAt, &message.Deleted); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, rows.Err()
}

func (r *Repository) CreateMessage(message *models.Message) error {
	now := time.Now()
	result, err := r.db.Exec(`INSERT INTO messages (conversation_id, sender_id, content, created_at, updated_at, deleted) VALUES (?, ?, ?, ?, ?, 0)`, message.ConversationID, message.SenderID, message.Content, now, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	if _, err := r.db.Exec(`UPDATE conversations SET updated_at = ? WHERE id = ?`, now, message.ConversationID); err != nil {
		return err
	}
	message.ID = int(id)
	message.CreatedAt = now
	message.UpdatedAt = now
	return nil
}

func (r *Repository) UpdateMessage(messageID, userID int, content string) (*models.Message, error) {
	now := time.Now()
	result, err := r.db.Exec(`UPDATE messages SET content = ?, updated_at = ? WHERE id = ? AND sender_id = ? AND deleted = 0`, content, now, messageID, userID)
	if err != nil {
		return nil, err
	}
	count, err := result.RowsAffected()
	if err != nil || count == 0 {
		return nil, sql.ErrNoRows
	}
	var message models.Message
	err = r.db.QueryRow(`SELECT m.id, m.conversation_id, m.sender_id, u.username, m.content, m.created_at, m.updated_at, m.deleted
		FROM messages m JOIN users u ON u.id = m.sender_id WHERE m.id = ?`, messageID).Scan(&message.ID, &message.ConversationID, &message.SenderID, &message.SenderUsername, &message.Content, &message.CreatedAt, &message.UpdatedAt, &message.Deleted)
	return &message, err
}

func (r *Repository) DeleteMessage(messageID, userID int) (*models.Message, error) {
	now := time.Now()
	result, err := r.db.Exec(`UPDATE messages SET content = '', updated_at = ?, deleted = 1 WHERE id = ? AND sender_id = ? AND deleted = 0`, now, messageID, userID)
	if err != nil {
		return nil, err
	}
	count, err := result.RowsAffected()
	if err != nil || count == 0 {
		return nil, sql.ErrNoRows
	}
	var message models.Message
	err = r.db.QueryRow(`SELECT m.id, m.conversation_id, m.sender_id, u.username, m.content, m.created_at, m.updated_at, m.deleted
		FROM messages m JOIN users u ON u.id = m.sender_id WHERE m.id = ?`, messageID).Scan(&message.ID, &message.ConversationID, &message.SenderID, &message.SenderUsername, &message.Content, &message.CreatedAt, &message.UpdatedAt, &message.Deleted)
	return &message, err
}

type rowScanner interface{ Scan(...any) error }

func scanConversation(row rowScanner) (*models.Conversation, *models.Message, error) {
	var conversation models.Conversation
	var user models.User
	var messageID, senderID, deleted sql.NullInt64
	var content sql.NullString
	var messageCreated, messageUpdated sql.NullTime
	if err := row.Scan(&conversation.ID, &user.ID, &user.Username, &user.Name, &user.Bio, &user.CreatedAt, &messageID, &senderID, &content, &messageCreated, &messageUpdated, &deleted, &conversation.UpdatedAt); err != nil {
		return nil, nil, err
	}
	conversation.OtherUser = &user
	if messageID.Valid {
		conversation.LastMessage = &models.Message{ID: int(messageID.Int64), ConversationID: conversation.ID, SenderID: int(senderID.Int64), Content: content.String, CreatedAt: messageCreated.Time, UpdatedAt: messageUpdated.Time, Deleted: deleted.Int64 == 1}
	}
	return &conversation, conversation.LastMessage, nil
}

// ============ MIGRATIONS ============

func createTables(db *sql.DB) error {
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return err
	}

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
		`CREATE TABLE IF NOT EXISTS tweet_likes (
			tweet_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (tweet_id) REFERENCES tweets(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(tweet_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS tweet_views (
			tweet_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (tweet_id) REFERENCES tweets(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(tweet_id, user_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tweets_user_id ON tweets(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tweets_created_at ON tweets(created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows(follower_id)`,
		`CREATE INDEX IF NOT EXISTS idx_follows_following ON follows(following_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tweet_likes_user ON tweet_likes(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tweet_views_user ON tweet_views(user_id)`,
		`CREATE TABLE IF NOT EXISTS conversations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			updated_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS conversation_members (
			conversation_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			PRIMARY KEY (conversation_id, user_id),
			FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			conversation_id INTEGER NOT NULL,
			sender_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id, id)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
